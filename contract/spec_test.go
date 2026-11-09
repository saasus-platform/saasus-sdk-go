package contract

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// spec is a parsed OpenAPI document, the source of the request inputs.
//
// Values come from the spec rather than being hardcoded, so they are acceptable
// by construction: the Prism mock validates requests against the same document.
// It also keeps the approach language-independent — the only language-specific
// part is binding an operationId to a generated method.
type spec struct {
	root map[string]any
	ops  map[string]*operation
}

type operation struct {
	id string
	// pathParams are ordered as they appear in the path template, which is the
	// order the generated method takes them in.
	pathParams []map[string]any
	// params holds the query and header parameters together, because the
	// generator gathers both into the single Params struct argument.
	params     []map[string]any
	bodySchema map[string]any
	statuses   map[int]bool
	hasSuccess bool
}

const syntheticValue = "synthetic"

var (
	httpMethods = []string{"get", "put", "post", "delete", "patch", "head", "options"}
	pathParamRe = regexp.MustCompile(`\{([^}]+)\}`)
)

func loadSpec(dir, module string) (*spec, error) {
	raw, err := os.ReadFile(filepath.Join(dir, module+".yml"))
	if err != nil {
		return nil, err
	}
	var decoded any
	if err := yaml.Unmarshal(raw, &decoded); err != nil {
		return nil, err
	}
	root, ok := normalize(decoded).(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: not an OpenAPI document", module)
	}
	s := &spec{root: root, ops: map[string]*operation{}}
	s.index()
	return s, nil
}

// normalize turns YAML mappings with non-string keys into map[string]any.
// Unquoted response codes (200:) decode to integer keys, so without this the
// document would carry two different map types.
func normalize(v any) any {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			t[k] = normalize(val)
		}
		return t
	case map[any]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			out[fmt.Sprint(k)] = normalize(val)
		}
		return out
	case []any:
		for i, val := range t {
			t[i] = normalize(val)
		}
		return t
	default:
		return v
	}
}

// --- indexing ---------------------------------------------------------------

func (s *spec) index() {
	paths, _ := s.root["paths"].(map[string]any)
	for path, raw := range paths {
		pathItem, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		for _, method := range httpMethods {
			rawOp, ok := pathItem[method].(map[string]any)
			if !ok {
				continue
			}
			id, ok := rawOp["operationId"].(string)
			if !ok {
				continue
			}
			op := &operation{id: id, statuses: map[int]bool{}}
			s.loadParams(op, path, pathItem, rawOp)
			op.bodySchema = s.bodySchema(rawOp)
			loadStatuses(op, rawOp)
			s.ops[id] = op
		}
	}
}

func (s *spec) loadParams(op *operation, path string, pathItem, rawOp map[string]any) {
	// Path-level parameters apply to every operation of the path item.
	byName := map[string]map[string]any{}
	for _, list := range []any{pathItem["parameters"], rawOp["parameters"]} {
		items, _ := list.([]any)
		for _, raw := range items {
			m, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			p := s.resolve(m)
			name, _ := p["name"].(string)
			switch p["in"] {
			case "path":
				byName[name] = p
			case "query", "header":
				op.params = append(op.params, p)
			}
		}
	}
	for _, match := range pathParamRe.FindAllStringSubmatch(path, -1) {
		if p, ok := byName[match[1]]; ok {
			op.pathParams = append(op.pathParams, p)
		}
	}
}

func (s *spec) bodySchema(rawOp map[string]any) map[string]any {
	body, ok := rawOp["requestBody"].(map[string]any)
	if !ok {
		return nil
	}
	content, _ := s.resolve(body)["content"].(map[string]any)
	jsonContent, _ := content["application/json"].(map[string]any)
	schema, ok := jsonContent["schema"].(map[string]any)
	if !ok {
		return nil
	}
	return s.resolve(schema)
}

func loadStatuses(op *operation, rawOp map[string]any) {
	responses, _ := rawOp["responses"].(map[string]any)
	for key := range responses {
		code, err := strconv.Atoi(key)
		if err != nil {
			continue // "default" and the like carry no status
		}
		op.statuses[code] = true
		if code >= 200 && code < 300 {
			op.hasSuccess = true
		}
	}
}

// --- $ref resolution --------------------------------------------------------

func (s *spec) resolve(node map[string]any) map[string]any {
	for i := 0; i < 32; i++ {
		ref, ok := node["$ref"].(string)
		if !ok {
			return node
		}
		next := s.resolveRef(ref)
		if next == nil {
			return map[string]any{}
		}
		node = next
	}
	return node
}

func (s *spec) resolveRef(ref string) map[string]any {
	if !strings.HasPrefix(ref, "#/") {
		return nil
	}
	var cur any = s.root
	for _, seg := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
		seg = strings.ReplaceAll(strings.ReplaceAll(seg, "~1", "/"), "~0", "~")
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[seg]
	}
	m, _ := cur.(map[string]any)
	return m
}

// --- value synthesis --------------------------------------------------------

// paramsValues returns the query and header parameters as a JSON-shaped object
// keyed by their spec names. The generated Params struct tags every field with
// its spec name, so unmarshalling binds by name.
func (s *spec) paramsValues(op *operation) map[string]any {
	out := make(map[string]any, len(op.params))
	for _, p := range op.params {
		name, _ := p["name"].(string)
		out[name] = s.paramValue(p)
	}
	return out
}

func (s *spec) paramValue(p map[string]any) any {
	if example, ok := p["example"]; ok {
		return example
	}
	schema, _ := p["schema"].(map[string]any)
	return s.synth(schema)
}

// bodyVariants returns the labels of the request body's top-level oneOf/anyOf
// members, or nil when the body is neither. Each member becomes its own case
// so that every variant is actually sent.
func (s *spec) bodyVariants(op *operation) []string {
	members := variants(op.bodySchema)
	if members == nil {
		return nil
	}
	discriminator, _ := op.bodySchema["discriminator"].(map[string]any)
	labels := make([]string, 0, len(members))
	for _, member := range members {
		labels = append(labels, discriminatorValue(discriminator, refOf(member)))
	}
	return labels
}

// bodyValue synthesizes the request body, selecting the given top-level
// oneOf/anyOf member.
func (s *spec) bodyValue(op *operation, variantIndex int) any {
	members := variants(op.bodySchema)
	if variantIndex < len(members) {
		return s.synthVariant(op.bodySchema, members[variantIndex])
	}
	return s.synth(op.bodySchema)
}

func (s *spec) synth(schema map[string]any) any {
	if schema == nil {
		return nil
	}
	schema = s.resolve(schema)
	if example, ok := schema["example"]; ok {
		return example
	}
	if enum, ok := schema["enum"].([]any); ok && len(enum) > 0 {
		return enum[0]
	}
	if members := variants(schema); members != nil {
		return s.synthVariant(schema, members[0])
	}
	if members, ok := schema["allOf"].([]any); ok {
		return s.synthObject(s.mergeAllOf(members))
	}
	kind, _ := schema["type"].(string)
	if _, hasProps := schema["properties"].(map[string]any); kind == "" && hasProps {
		kind = "object"
	}
	switch kind {
	case "object":
		return s.synthObject(schema)
	case "array":
		items, _ := schema["items"].(map[string]any)
		if item := s.synth(items); item != nil {
			return []any{item}
		}
		return []any{}
	default:
		return synthScalar(schema)
	}
}

// synthObject fills every declared property, not just the required ones: an
// optional property that is never sent leaves its own serialization, and the
// required fields nested below it, unverified.
func (s *spec) synthObject(schema map[string]any) any {
	out := map[string]any{}
	for name, raw := range mapOf(schema["properties"]) {
		if prop, ok := raw.(map[string]any); ok {
			out[name] = s.synth(prop)
		}
	}
	return out
}

// mergeAllOf composes the member schemas' properties, which is what allOf is
// used for in these specs.
func (s *spec) mergeAllOf(members []any) map[string]any {
	props := map[string]any{}
	for _, raw := range members {
		member, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		part := s.resolve(member)
		for name, prop := range mapOf(part["properties"]) {
			props[name] = prop
		}
	}
	return map[string]any{"properties": props}
}

// synthVariant synthesizes one oneOf member and, when the parent declares a
// discriminator, sets the discriminator property so that the value resolves to
// that member.
func (s *spec) synthVariant(parent map[string]any, member any) any {
	schema, ok := member.(map[string]any)
	if !ok {
		return map[string]any{}
	}
	value := s.synth(schema)
	discriminator, ok := parent["discriminator"].(map[string]any)
	if !ok {
		return value
	}
	object, ok := value.(map[string]any)
	if !ok {
		return value
	}
	if property, ok := discriminator["propertyName"].(string); ok {
		object[property] = discriminatorValue(discriminator, refOf(member))
	}
	return object
}

func synthScalar(schema map[string]any) any {
	switch kind, _ := schema["type"].(string); kind {
	case "integer", "number":
		return 1
	case "boolean":
		return true
	}
	switch format, _ := schema["format"].(string); format {
	case "date":
		return "2024-01-01"
	case "date-time":
		return "2024-01-01T00:00:00Z"
	}
	// A string with no example still has to clear the declared length floor or
	// the mock rejects it. Padding with an alphanumeric also satisfies the
	// character patterns these specs declare.
	value := syntheticValue
	for minLen := intOf(schema["minLength"]); len(value) < minLen; {
		value += "x"
	}
	return value
}

func intOf(v any) int {
	switch t := v.(type) {
	case int:
		return t
	case float64:
		return int(t)
	}
	return 0
}

// --- small helpers ----------------------------------------------------------

// variants returns the schema's oneOf or anyOf members, or nil when it is
// neither.
func variants(schema map[string]any) []any {
	if schema == nil {
		return nil
	}
	members, _ := schema["oneOf"].([]any)
	if len(members) == 0 {
		members, _ = schema["anyOf"].([]any)
	}
	if len(members) == 0 {
		return nil
	}
	return members
}

// discriminatorValue returns the value whose mapping points at the given member
// reference, falling back to the referenced schema name.
func discriminatorValue(discriminator map[string]any, ref string) string {
	for value, target := range mapOf(discriminator["mapping"]) {
		if fmt.Sprint(target) == ref && ref != "" {
			return value
		}
	}
	if i := strings.LastIndex(ref, "/"); i >= 0 {
		return ref[i+1:]
	}
	return syntheticValue
}

func refOf(member any) string {
	schema, ok := member.(map[string]any)
	if !ok {
		return ""
	}
	ref, _ := schema["$ref"].(string)
	return ref
}

func mapOf(v any) map[string]any {
	m, _ := v.(map[string]any)
	return m
}
