package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

// resolveRefs recursively resolves $ref entries for schemas and parameters
func resolveRefs(obj interface{}, components map[string]interface{}) interface{} {
	switch node := obj.(type) {
	case map[string]interface{}:
		// Handle $ref
		if refVal, ok := node["$ref"].(string); ok {
			if compsSchemas, found := components["schemas"].(map[string]interface{}); found &&
				len(refVal) > len("#/components/schemas/") &&
				refVal[:len("#/components/schemas/")] == "#/components/schemas/" {
				schemaName := refVal[len("#/components/schemas/"):]
				if resolved, exists := compsSchemas[schemaName]; exists {
					return resolveRefs(resolved, components)
				}
			} else if compsParams, found := components["parameters"].(map[string]interface{}); found &&
				len(refVal) > len("#/components/parameters/") &&
				refVal[:len("#/components/parameters/")] == "#/components/parameters/" {
				paramName := refVal[len("#/components/parameters/"):]
				if resolved, exists := compsParams[paramName]; exists {
					return resolveRefs(resolved, components)
				}
			}
		}
		// Recurse into map
		for k, v := range node {
			node[k] = resolveRefs(v, components)
		}
		return node

	case []interface{}:
		for i, v := range node {
			node[i] = resolveRefs(v, components)
		}
		return node

	default:
		// Primitive types
		return obj
	}
}

// splitOpenAPIByTags reads an OpenAPI YAML, splits paths by tag, and writes per-tag YAML files.
func splitOpenAPIByTags(inputFile, outputDir string) error {
	// Read input
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Parse YAML into generic map
	var spec map[string]interface{}
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Prepare base info
	baseInfo := map[string]interface{}{}
	for _, key := range []string{"openapi", "info", "servers", "security"} {
		if v, ok := spec[key]; ok {
			baseInfo[key] = v
		}
	}

	// Get tags list and components
	allTags, _ := spec["tags"].([]interface{})
	components, _ := spec["components"].(map[string]interface{})

	// Group paths by tag
	pathsByTag := map[string]map[string]map[string]interface{}{}
	if paths, ok := spec["paths"].(map[string]interface{}); ok {
		for path, methodsRaw := range paths {
			if methods, ok := methodsRaw.(map[string]interface{}); ok {
				for method, detailsRaw := range methods {
					if details, ok := detailsRaw.(map[string]interface{}); ok {
						tagsRaw, _ := details["tags"].([]interface{})
						for _, tagVal := range tagsRaw {
							if tagName, ok := tagVal.(string); ok {
								if _, exists := pathsByTag[tagName]; !exists {
									pathsByTag[tagName] = map[string]map[string]interface{}{}
								}
								if _, exists := pathsByTag[tagName][path]; !exists {
									pathsByTag[tagName][path] = map[string]interface{}{}
								}
								pathsByTag[tagName][path][method] = details
							}
						}
					}
				}
			}
		}
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write per-tag YAML
	for tagName, tagPaths := range pathsByTag {
		// Build output spec
		out := map[string]interface{}{}
		for k, v := range baseInfo {
			out[k] = v
		}
		// Select only the tag object matching tagName
		selectedTags := []interface{}{}
		for _, t := range allTags {
			if tagObj, ok := t.(map[string]interface{}); ok {
				if name, _ := tagObj["name"].(string); name == tagName {
					selectedTags = append(selectedTags, tagObj)
				}
			}
		}
		out["tags"] = selectedTags
		
		// Resolve and assign paths
		outPaths := map[string]interface{}{}
		for path, methods := range tagPaths {
			outPaths[path] = map[string]interface{}{}
			for method, details := range methods {
				resolved := resolveRefs(details, components)
				outPaths[path].(map[string]interface{})[method] = resolved
			}
		}
		out["paths"] = outPaths

		// Marshal and write file
		bytes, err := yaml.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML for tag %s: %w", tagName, err)
		}
		outFile := filepath.Join(outputDir, fmt.Sprintf("%s.yml", tagName))
		if err := ioutil.WriteFile(outFile, bytes, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", outFile, err)
		}
		fmt.Printf("Created %s\n", outFile)
	}

	return nil
}

func main() {
	input := flag.String("input", "authapi.yml", "Path to the input OpenAPI YAML file")
	output := flag.String("output", "auth_yaml_divided", "Directory for output per-tag YAML files")
	flag.Parse()

	if err := splitOpenAPIByTags(*input, *output); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
