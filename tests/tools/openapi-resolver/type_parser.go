package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// TypeDefinitionParser parses types.gen.go to extract type definitions
type TypeDefinitionParser struct {
	moduleName string
	typesPath  string
	typeCache  map[string]*TypeDefinition
}

// TypeDefinition represents a complete type definition
type TypeDefinition struct {
	Name   string
	Fields []FieldDefinition
}

// FieldDefinition represents a struct field
type FieldDefinition struct {
	Name     string
	Type     string
	JSONTag  string
	Required bool
	Comment  string
}

// NewTypeDefinitionParser creates a new parser
func NewTypeDefinitionParser(moduleName string) (*TypeDefinitionParser, error) {
	// generated/{module}api/types.gen.go のパスを構築
	// プロジェクトルートからの相対パスを試す
	typesPath := filepath.Join("generated", moduleName+"api", "types.gen.go")

	// ファイルが存在しない場合、../../../ を追加して試す（テスト実行時用）
	if _, err := os.Stat(typesPath); os.IsNotExist(err) {
		typesPath = filepath.Join("../../../generated", moduleName+"api", "types.gen.go")
		if _, err := os.Stat(typesPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("types.gen.go file not found for module: %s", moduleName)
		}
	}

	parser := &TypeDefinitionParser{
		moduleName: moduleName,
		typesPath:  typesPath,
		typeCache:  make(map[string]*TypeDefinition),
	}

	return parser, nil
}

// ParseAllTypes parses all type definitions
func (p *TypeDefinitionParser) ParseAllTypes() (map[string]*TypeDefinition, error) {
	// ファイルを読み込んでASTを解析
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, p.typesPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", p.typesPath, err)
	}

	// すべての型定義を解析
	ast.Inspect(node, func(n ast.Node) bool {
		// 型定義を検索
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// 構造体型のみを処理
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// 型定義を抽出
		typeDef := p.extractTypeDefinition(typeSpec.Name.Name, structType)
		if typeDef != nil {
			p.typeCache[typeDef.Name] = typeDef
		}

		return true
	})

	if len(p.typeCache) == 0 {
		return nil, fmt.Errorf("no type definitions found in %s", p.typesPath)
	}

	return p.typeCache, nil
}

// extractTypeDefinition extracts type definition from AST struct type
func (p *TypeDefinitionParser) extractTypeDefinition(typeName string, structType *ast.StructType) *TypeDefinition {
	typeDef := &TypeDefinition{
		Name:   typeName,
		Fields: []FieldDefinition{},
	}

	// 各フィールドを解析
	for _, field := range structType.Fields.List {
		// フィールド名を取得
		var fieldName string
		if len(field.Names) > 0 {
			fieldName = field.Names[0].Name
		} else {
			// 埋め込みフィールドの場合はスキップ
			continue
		}

		// フィールド型を取得
		fieldType := p.exprToString(field.Type)

		// JSONタグを取得
		jsonTag := ""
		required := false
		if field.Tag != nil {
			tag := strings.Trim(field.Tag.Value, "`")
			jsonTag = p.extractJSONTag(tag)
			// requiredの判定（オプショナルでない場合はrequired）
			required = !strings.Contains(tag, "omitempty")
		}

		// コメントを取得
		comment := ""
		if field.Doc != nil && len(field.Doc.List) > 0 {
			comment = strings.TrimSpace(field.Doc.List[0].Text)
			comment = strings.TrimPrefix(comment, "//")
			comment = strings.TrimSpace(comment)
		}

		typeDef.Fields = append(typeDef.Fields, FieldDefinition{
			Name:     fieldName,
			Type:     fieldType,
			JSONTag:  jsonTag,
			Required: required,
			Comment:  comment,
		})
	}

	return typeDef
}

// extractJSONTag extracts JSON tag from struct tag
func (p *TypeDefinitionParser) extractJSONTag(tag string) string {
	// json:"field_name" または json:"field_name,omitempty" の形式から抽出
	parts := strings.Split(tag, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "json:") {
			jsonPart := strings.TrimPrefix(part, "json:")
			jsonPart = strings.Trim(jsonPart, "\"")
			// カンマで分割して最初の部分を取得
			jsonFields := strings.Split(jsonPart, ",")
			if len(jsonFields) > 0 && jsonFields[0] != "-" {
				return jsonFields[0]
			}
		}
	}
	return ""
}

// exprToString converts AST expression to string representation
func (p *TypeDefinitionParser) exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + p.exprToString(t.X)
	case *ast.SelectorExpr:
		return p.exprToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + p.exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + p.exprToString(t.Key) + "]" + p.exprToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return "..." + p.exprToString(t.Elt)
	default:
		return ""
	}
}

// GetTypeDefinition gets type definition by type name
func (p *TypeDefinitionParser) GetTypeDefinition(typeName string) (*TypeDefinition, error) {
	// キャッシュがまだ構築されていない場合は解析を実行
	if len(p.typeCache) == 0 {
		if _, err := p.ParseAllTypes(); err != nil {
			return nil, err
		}
	}

	// ポインタ型の場合は * を除去
	typeName = strings.TrimPrefix(typeName, "*")

	// 型名で検索
	if typeDef, ok := p.typeCache[typeName]; ok {
		return typeDef, nil
	}

	// 利用可能な型のリストを作成
	availableTypes := make([]string, 0, len(p.typeCache))
	for name := range p.typeCache {
		availableTypes = append(availableTypes, name)
	}

	return nil, fmt.Errorf("type definition not found for type '%s'. Available types: %v", typeName, availableTypes)
}

// GetTypeDefinitionRecursive gets type definition recursively including nested types
func (p *TypeDefinitionParser) GetTypeDefinitionRecursive(typeName string, visited map[string]bool) (*TypeDefinition, error) {
	// 循環参照の検出
	if visited == nil {
		visited = make(map[string]bool)
	}

	// ポインタ型の場合は * を除去
	typeName = strings.TrimPrefix(typeName, "*")

	// 既に訪問済みの型の場合は循環参照
	if visited[typeName] {
		return nil, fmt.Errorf("circular reference detected for type '%s'", typeName)
	}

	// 訪問済みとしてマーク
	visited[typeName] = true

	// 型定義を取得
	typeDef, err := p.GetTypeDefinition(typeName)
	if err != nil {
		return nil, err
	}

	// ネストした型を再帰的に解析
	for i, field := range typeDef.Fields {
		// フィールド型がカスタム型の場合
		fieldType := strings.TrimPrefix(field.Type, "*")
		fieldType = strings.TrimPrefix(fieldType, "[]")

		// マップ型の場合は値の型を抽出
		if strings.HasPrefix(fieldType, "map[") {
			// map[string]CustomType の形式から CustomType を抽出
			parts := strings.Split(fieldType, "]")
			if len(parts) > 1 {
				fieldType = parts[1]
			}
		}

		// プリミティブ型やパッケージ型はスキップ
		if p.isPrimitiveType(fieldType) || strings.Contains(fieldType, ".") {
			continue
		}

		// カスタム型の場合は再帰的に解析
		if _, ok := p.typeCache[fieldType]; ok {
			nestedTypeDef, err := p.GetTypeDefinitionRecursive(fieldType, visited)
			if err != nil {
				// 循環参照の場合はスキップ
				if strings.Contains(err.Error(), "circular reference") {
					continue
				}
				return nil, err
			}

			// ネストした型の情報をコメントに追加
			if nestedTypeDef != nil {
				typeDef.Fields[i].Comment = fmt.Sprintf("%s (nested type with %d fields)", field.Comment, len(nestedTypeDef.Fields))
			}
		}
	}

	return typeDef, nil
}

// isPrimitiveType checks if the type is a primitive Go type
func (p *TypeDefinitionParser) isPrimitiveType(typeName string) bool {
	primitives := []string{
		"string", "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "bool", "byte", "rune",
		"interface{}", "error",
	}

	for _, primitive := range primitives {
		if typeName == primitive {
			return true
		}
	}

	return false
}
