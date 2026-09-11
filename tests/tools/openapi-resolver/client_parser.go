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

// ClientInterfaceParser parses client.gen.go to extract method information
type ClientInterfaceParser struct {
	moduleName string
	clientPath string
	methods    map[string]*MethodSignature
}

// MethodSignature represents a client method signature
type MethodSignature struct {
	MethodName     string
	OperationID    string
	ParameterTypes []ParameterType
	ReturnType     string
}

// ParameterType represents a method parameter
type ParameterType struct {
	Name string
	Type string
}

// NewClientInterfaceParser creates a new parser
func NewClientInterfaceParser(moduleName string) (*ClientInterfaceParser, error) {
	// generated/{module}api/client.gen.go のパスを構築
	// プロジェクトルートからの相対パスを試す
	clientPath := filepath.Join("generated", moduleName+"api", "client.gen.go")

	// ファイルが存在しない場合、../../../ を追加して試す（テスト実行時用）
	if _, err := os.Stat(clientPath); os.IsNotExist(err) {
		clientPath = filepath.Join("../../../generated", moduleName+"api", "client.gen.go")
		if _, err := os.Stat(clientPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("client.gen.go file not found for module: %s", moduleName)
		}
	}

	parser := &ClientInterfaceParser{
		moduleName: moduleName,
		clientPath: clientPath,
		methods:    make(map[string]*MethodSignature),
	}

	return parser, nil
}

// ParseAllMethods parses all methods in ClientInterface
func (p *ClientInterfaceParser) ParseAllMethods() ([]MethodSignature, error) {
	// ファイルを読み込んでASTを解析
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, p.clientPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", p.clientPath, err)
	}

	// ClientWithResponsesインターフェースを検索
	var methods []MethodSignature
	ast.Inspect(node, func(n ast.Node) bool {
		// インターフェース型を検索
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// ClientInterfaceを特定
		if typeSpec.Name.Name != "ClientInterface" {
			return true
		}

		// インターフェース型を取得
		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		// 各メソッドを解析
		for _, method := range interfaceType.Methods.List {
			if len(method.Names) == 0 {
				continue
			}

			methodName := method.Names[0].Name
			funcType, ok := method.Type.(*ast.FuncType)
			if !ok {
				continue
			}

			// メソッドシグネチャを抽出
			sig := p.extractMethodSignature(methodName, funcType)
			if sig != nil {
				methods = append(methods, *sig)
				p.methods[sig.OperationID] = sig
			}
		}

		return false // ClientWithResponsesを見つけたので探索終了
	})

	if len(methods) == 0 {
		return nil, fmt.Errorf("no methods found in ClientInterface in %s", p.clientPath)
	}

	return methods, nil
}

// extractMethodSignature extracts method signature from AST function type
func (p *ClientInterfaceParser) extractMethodSignature(methodName string, funcType *ast.FuncType) *MethodSignature {
	sig := &MethodSignature{
		MethodName:     methodName,
		OperationID:    methodName, // デフォルトではメソッド名をoperationIDとして使用
		ParameterTypes: []ParameterType{},
	}

	// パラメータを抽出
	if funcType.Params != nil {
		for _, field := range funcType.Params.List {
			paramType := p.exprToString(field.Type)

			// パラメータ名がある場合
			if len(field.Names) > 0 {
				for _, name := range field.Names {
					sig.ParameterTypes = append(sig.ParameterTypes, ParameterType{
						Name: name.Name,
						Type: paramType,
					})
				}
			} else {
				// 無名パラメータの場合
				sig.ParameterTypes = append(sig.ParameterTypes, ParameterType{
					Name: "",
					Type: paramType,
				})
			}
		}
	}

	// 戻り値の型を抽出
	if funcType.Results != nil && len(funcType.Results.List) > 0 {
		sig.ReturnType = p.exprToString(funcType.Results.List[0].Type)
	}

	return sig
}

// exprToString converts AST expression to string representation
func (p *ClientInterfaceParser) exprToString(expr ast.Expr) string {
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

// FindMethodByOperationID finds method name by operationId
func (p *ClientInterfaceParser) FindMethodByOperationID(operationID string) (*MethodSignature, error) {
	// メソッドがまだ解析されていない場合は解析を実行
	if len(p.methods) == 0 {
		if _, err := p.ParseAllMethods(); err != nil {
			return nil, err
		}
	}

	// operationIDで検索
	if method, ok := p.methods[operationID]; ok {
		return method, nil
	}

	// WithResponseサフィックスを除いて検索
	baseOperationID := strings.TrimSuffix(operationID, "WithResponse")
	if method, ok := p.methods[baseOperationID]; ok {
		return method, nil
	}

	// WithBodyサフィックスを除いて検索
	baseOperationID = strings.TrimSuffix(operationID, "WithBody")
	baseOperationID = strings.TrimSuffix(baseOperationID, "WithResponse")
	if method, ok := p.methods[baseOperationID]; ok {
		return method, nil
	}

	// 利用可能なメソッドのリストを作成
	availableMethods := make([]string, 0, len(p.methods))
	for methodName := range p.methods {
		availableMethods = append(availableMethods, methodName)
	}

	return nil, fmt.Errorf("method not found for operationId '%s'. Available methods: %v", operationID, availableMethods)
}
