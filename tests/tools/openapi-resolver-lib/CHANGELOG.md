# OpenAPI Resolver Library - 変更履歴

## ClientInterfaceParser の追加

### 実装日
2024-10-24

### 概要
`client.gen.go`からメソッド情報を抽出する`ClientInterfaceParser`を追加しました。これにより、生成されたクライアントコードのメソッドシグネチャを解析し、operationIDからメソッド名やパラメータ型を取得できるようになりました。

### 追加された機能

#### ClientInterfaceParser 構造体
```go
type ClientInterfaceParser struct {
    moduleName string
    clientPath string
    methods    map[string]*MethodSignature
}
```

生成された`client.gen.go`ファイルを解析し、`ClientWithResponses`インターフェースからメソッド情報を抽出します。

#### MethodSignature 構造体
```go
type MethodSignature struct {
    MethodName     string
    OperationID    string
    ParameterTypes []ParameterType
    ReturnType     string
}
```

メソッドのシグネチャ情報を保持します。

#### 主要メソッド

1. **NewClientInterfaceParser**
   - モジュール名から`client.gen.go`のパスを自動構築
   - ファイルの存在確認を実行
   - パーサーインスタンスを作成

2. **ParseAllMethods**
   - `ClientWithResponses`インターフェースのすべてのメソッドを解析
   - Go ASTを使用してメソッドシグネチャを抽出
   - メソッド名、パラメータ型、戻り値型を取得

3. **FindMethodByOperationID**
   - operationIDからメソッドシグネチャを検索
   - WithResponse、WithBodyなどのサフィックスを考慮
   - 見つからない場合は利用可能なメソッドリストを提示

### 使用例

```go
// パーサーを作成
parser, err := NewClientInterfaceParser("communication")
if err != nil {
    log.Fatal(err)
}

// すべてのメソッドを解析
methods, err := parser.ParseAllMethods()
if err != nil {
    log.Fatal(err)
}

for _, method := range methods {
    fmt.Printf("メソッド: %s\n", method.MethodName)
    fmt.Printf("  パラメータ: %d個\n", len(method.ParameterTypes))
}

// operationIDから検索
method, err := parser.FindMethodByOperationID("CreateFeedback")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("メソッド名: %s\n", method.MethodName)
for _, param := range method.ParameterTypes {
    fmt.Printf("  - %s: %s\n", param.Name, param.Type)
}
```

### 技術詳細

#### AST解析
Go標準ライブラリの`go/ast`と`go/parser`を使用して、生成されたGoコードを解析します。

```go
fset := token.NewFileSet()
node, err := parser.ParseFile(fset, clientPath, nil, parser.ParseComments)
```

#### インターフェース検索
ASTを走査して`ClientWithResponses`インターフェースを特定し、そのメソッドを抽出します。

```go
ast.Inspect(node, func(n ast.Node) bool {
    typeSpec, ok := n.(*ast.TypeSpec)
    if !ok || typeSpec.Name.Name != "ClientWithResponses" {
        return true
    }
    // メソッドを抽出
})
```

#### 型情報の文字列化
AST式を文字列表現に変換します。ポインタ型、配列型、マップ型などをサポート。

```go
func exprToString(expr ast.Expr) string {
    switch t := expr.(type) {
    case *ast.Ident:
        return t.Name
    case *ast.StarExpr:
        return "*" + exprToString(t.X)
    case *ast.ArrayType:
        return "[]" + exprToString(t.Elt)
    // ...
    }
}
```

### 統合ポイント

この機能は以下のツールと統合されます：

1. **postman-analyzer**: operationIDから正確なメソッド名を取得
2. **test-param-generator**: メソッドシグネチャからパラメータ型を特定
3. **openapi-resolver CLI**: メソッド情報の表示と検証

### エラーハンドリング

- `client.gen.go`が見つからない場合は`ResolverError`を返す
- メソッドが見つからない場合は利用可能なメソッドリストを提示
- AST解析エラーは詳細なエラーメッセージで報告

### 今後の拡張

1. **types.gen.goの解析**: パラメータ型の詳細情報を取得
2. **ドキュメント抽出**: コメントからドキュメントを抽出
3. **バリデーション**: メソッドシグネチャの妥当性検証

---

## SchemaObject型の追加とMediaTypeObjectの型安全性向上

### 変更内容
`MediaTypeObject`構造体の`Schema`フィールドの型を`interface{}`から`SchemaObject`構造体に変更しました。

**変更前:**
```go
type MediaTypeObject struct {
    Schema interface{} `yaml:"schema"`
}
```

**変更後:**
```go
type MediaTypeObject struct {
    Schema SchemaObject `yaml:"schema"`
}

type SchemaObject struct {
    Ref  string `yaml:"$ref"`
    Type string `yaml:"type"`
}
```

### 改善点
1. **型安全性の向上**: `interface{}`から具体的な構造体型に変更することで、コンパイル時の型チェックが可能に
2. **コードの可読性向上**: スキーマオブジェクトの構造が明確になり、どのフィールドが利用可能かが一目瞭然
3. **保守性の向上**: IDEの補完機能が効くようになり、開発効率が向上
4. **バグの早期発見**: 型の不一致がコンパイル時に検出されるため、ランタイムエラーを削減

### 影響範囲
- `extractRequestBodyType`関数で`mediaType.Schema.Ref`に直接アクセス可能
- YAML解析時に自動的に`SchemaObject`構造体にマッピング
- 既存の機能には影響なし（後方互換性あり）

### 使用例
```go
// requestBodyの型名を抽出
mediaType, ok := requestBody.Content["application/json"]
if ok {
    // 型安全にアクセス可能
    ref := mediaType.Schema.Ref
    schemaType := mediaType.Schema.Type
}
```

## TypeDefinitionParserの追加（実装中）

### 実装日
2024-10-24

### 概要
`types.gen.go`から型定義を抽出する`TypeDefinitionParser`を追加しました。これにより、パラメータ型の詳細情報（フィールド名、型、JSONタグ等）を取得できるようになります。

### 追加された機能

#### TypeDefinitionParser 構造体
```go
type TypeDefinitionParser struct {
    moduleName string
    typesPath  string
    typeCache  map[string]*TypeDefinition
}
```

`generated/{module}api/types.gen.go`ファイルを解析し、構造体定義から型情報を抽出します。

#### TypeDefinition 構造体
```go
type TypeDefinition struct {
    Name   string
    Fields []FieldDefinition
}
```

型定義の完全な情報を保持します。

#### FieldDefinition 構造体
```go
type FieldDefinition struct {
    Name     string
    Type     string
    JSONTag  string
    Required bool
    Comment  string
}
```

構造体フィールドの詳細情報を保持します。

### 主要メソッド

1. **NewTypeDefinitionParser**
   - モジュール名から`types.gen.go`のパスを自動構築
   - ファイルの存在確認を実行
   - パーサーインスタンスを作成

2. **GetTypeDefinition**
   - 型名から型定義を取得
   - キャッシュ機構により高速な検索を実現
   - 見つからない場合はエラーを返す

3. **parseAllTypes**（内部メソッド）
   - Go ASTを使用してすべての型定義を解析
   - 構造体定義からフィールド情報を抽出
   - 結果をキャッシュに保存

### 使用例（実装完了後）

```go
// パーサーを作成
parser, err := NewTypeDefinitionParser("communication")
if err != nil {
    log.Fatal(err)
}

// 型定義を取得
typeDef, err := parser.GetTypeDefinition("CreateFeedbackJSONRequestBody")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("型名: %s\n", typeDef.Name)
for _, field := range typeDef.Fields {
    fmt.Printf("  フィールド: %s\n", field.Name)
    fmt.Printf("    型: %s\n", field.Type)
    fmt.Printf("    JSONタグ: %s\n", field.JSONTag)
    fmt.Printf("    必須: %v\n", field.Required)
}
```

### 実装状況

- ✅ 基本構造の定義
- ✅ ファイルパス解決
- ✅ AST解析の実装
- ✅ フィールド情報の抽出
- ✅ キャッシュ機構の実装
- ✅ テストの追加
- ✅ 再帰的型解析のサポート
- ✅ JSONタグ抽出
- ✅ プリミティブ型判定

### 統合ポイント

この機能は以下のツールと統合されます：

1. **postman-analyzer**: 型定義からフィールド情報を取得し、レポートに含める
2. **test-param-generator**: 正確な型情報を使用してパラメータマッピングを実行
3. **openapi-resolver CLI**: 型定義の表示と検証

### 実装された機能の詳細

#### 1. AST解析
Go標準ライブラリの`go/ast`と`go/parser`を使用して、types.gen.goファイルを解析します。

```go
fset := token.NewFileSet()
node, err := parser.ParseFile(fset, typesPath, nil, parser.ParseComments)
```

#### 2. 型定義の抽出
ASTを走査してすべての構造体定義を抽出し、フィールド情報を取得します。

```go
ast.Inspect(node, func(n ast.Node) bool {
    typeSpec, ok := n.(*ast.TypeSpec)
    if !ok {
        return true
    }
    structType, ok := typeSpec.Type.(*ast.StructType)
    if !ok {
        return true
    }
    // 型定義を抽出
})
```

#### 3. フィールド情報の抽出
各フィールドから以下の情報を抽出します：
- フィールド名
- 型（ポインタ、配列、マップ、カスタム型をサポート）
- JSONタグ
- required判定（omitemptyの有無）
- ドキュメントコメント

#### 4. 再帰的型解析
ネストした型を再帰的に解析し、循環参照を検出します。

```go
func (p *TypeDefinitionParser) GetTypeDefinitionRecursive(typeName string, visited map[string]bool) (*TypeDefinition, error)
```

#### 5. キャッシュ機構
一度解析した型定義をキャッシュし、高速な検索を実現します。

### 今後の拡張予定

1. **バリデーションタグのサポート**: `validate`タグの解析
2. **パフォーマンス最適化**: 大きなファイル（10MB以上）の効率的な処理
3. **型エイリアスのサポート**: `type Alias = OriginalType`形式の解析
4. **ジェネリクスのサポート**: Go 1.18+のジェネリック型の解析

---

## パスパラメータ名の柔軟なマッチング対応

### 問題
以前の実装では、パスパラメータの名前が完全一致しないとマッチングできませんでした：

- OpenAPI仕様: `/tenants/{tenant_id}`
- Postman Collection: `/tenants/:id`

これらは構造的に同じパスですが、パラメータ名が異なるため一致しませんでした。

### 解決策
`PathMatcher` を導入し、パスの構造を比較してマッチングするようにしました。マッチした場合は**OpenAPI仕様のパス（正しいパラメータ名を含む）を返す**ため、後続の処理でOpenAPI仕様のパラメータ名が使用されます。

#### マッチングの仕組み

1. 入力パス（Postman Collection）を正規化: `/tenants/:id` → `/tenants/{id}`
2. OpenAPI仕様の各パスも正規化: `/tenants/{tenant_id}` → `/tenants/{tenant_id}`
3. パスの構造を比較（パラメータ名は無視）:
   - セグメント数が同じか
   - 各セグメントが固定値かパラメータか
   - 固定値の場合は完全一致するか
4. マッチした場合は**OpenAPI仕様のパス**を返す: `/tenants/{tenant_id}`

#### マッチング例

| Postman Collection | OpenAPI仕様 | マッチ結果 | 返されるパス |
|-------------------|------------|----------|------------|
| `/tenants/:id` | `/tenants/{tenant_id}` | ✅ | `/tenants/{tenant_id}` |
| `/tenants/{{id}}` | `/tenants/{tenant_id}` | ✅ | `/tenants/{tenant_id}` |
| `/tenants/{id}` | `/tenants/{tenant_id}` | ✅ | `/tenants/{tenant_id}` |
| `/tenants/:id/users/:uid` | `/tenants/{tenant_id}/users/{user_id}` | ✅ | `/tenants/{tenant_id}/users/{user_id}` |
| `/tenants/:id` | `/users/{id}` | ❌ | (固定セグメントが異なる) |

### 影響範囲
- パラメータ名が異なっても、パスの構造が同じであればマッチするようになりました
- マッチした場合は常にOpenAPI仕様のパス（正しいパラメータ名）が使用されます
- 既存の機能には影響ありません（後方互換性あり）
- Postman AnalyzerがOpenAPI仕様とPostman Collectionのエンドポイントを正しくマッチングできるようになりました

### テスト
`lib_test.go` に以下のテストケースを追加：
- `TestNormalizePath`: 各種パラメータ形式の正規化テスト
- `TestPathsMatch`: パスの構造比較テスト（パラメータ名が異なっても一致することを確認）
- `TestPathMatcher_FindMatchingPath`: 実際のマッチング動作テスト（OpenAPI仕様のパスが返されることを確認）
