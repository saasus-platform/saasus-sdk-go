package main

// PostmanReport はPostman Analyzerレポートの構造を表します
type PostmanReport struct {
	CollectionName string         `json:"collection_name"`
	TestFlow       []TestStep     `json:"test_flow"`
	Variables      any            `json:"variables"`
	ClientMethods  []ClientMethod `json:"client_methods"`
}

// TestStep はテストフローの1ステップを表します
type TestStep struct {
	Name              string            `json:"name"`
	Method            string            `json:"method"`
	Endpoint          string            `json:"endpoint"`
	Module            string            `json:"module"` // SaaSus module name (auth, pricing, billing, etc.)
	Description       string            `json:"description"`
	Tests             []string          `json:"tests"`
	ExpectedStatus    int               `json:"expected_status"`
	Variables         map[string]string `json:"variables"`
	RequiredVariables []string          `json:"required_variables"`
	RequestBody       map[string]any    `json:"request_body"`
	BodyFields        []string          `json:"body_fields"`

	// 型情報フィールド（要件11.1）
	OperationID         string               `json:"operation_id,omitempty"`          // OpenAPIのoperationId
	ClientMethodName    string               `json:"client_method_name,omitempty"`    // client.gen.goのメソッド名
	ParameterTypeName   string               `json:"parameter_type_name,omitempty"`   // パラメータ型名
	ParameterTypeFields map[string]FieldInfo `json:"parameter_type_fields,omitempty"` // 型のフィールド情報
}

// ClientMethod はクライアントメソッドの情報を表します
type ClientMethod struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	HTTPMethod  string `json:"http_method"`
	Endpoint    string `json:"endpoint"`
	Module      string `json:"module"` // SaaSus module name (auth, pricing, billing, etc.)
	HasBody     bool   `json:"has_body"`
	Description string `json:"description"`
	OperationID string `json:"operation_id"`
}

// ParamGroup はパラメータのグループを表します
type ParamGroup struct {
	Name           string         // グループ名（例: "feedbacks", "comments"）
	CreateParams   map[string]any // 作成用パラメータ
	UpdateParams   map[string]any // 更新用パラメータ
	CreateEndpoint string         // 作成用エンドポイント
	UpdateEndpoint string         // 更新用エンドポイント
	CreateMethod   string         // 作成用HTTPメソッド
	UpdateMethod   string         // 更新用HTTPメソッド
	ModuleCandidates []string     // 当該グループで参照されたモジュール候補
}

// AnalysisResult は分析結果を表します
type AnalysisResult struct {
	ModuleName   string
	ParamGroups  map[string]*ParamGroup
	UtilityFuncs []string // 必要なユーティリティ関数のリスト
	ModuleDependencies []string // 型参照に使用するモジュール候補
}

// TypeInfo は生成された型の情報を保持します
type TypeInfo struct {
	Name   string
	Fields map[string]FieldInfo
}

// FieldInfo はフィールドの情報を保持します
type FieldInfo struct {
	Name     string
	Type     string
	JSONTag  string
	Required bool
}

// CLIConfig はCLIの設定を保持します
type CLIConfig struct {
	ReportPath string // Postman Analyzerレポートのパス
	OutputDir  string // 出力ディレクトリ
	ModuleName string // モジュール名
	Overwrite  bool   // 既存ファイルの上書き
	Verbose    bool   // 詳細ログ出力
	ShowHelp   bool   // ヘルプ表示
}

// Config はツールの設定を保持します
type Config struct {
	ReportPath string
	OutputDir  string
	ModuleName string
	Overwrite  bool
	Verbose    bool
}

// TemplateData はテンプレートに渡すデータを表します
type TemplateData struct {
	PackageName string
	ModuleName  string
	ParamGroups map[string]*ParamGroup
	ImportPaths []string
	StructsCode string // 構造体定義のコード
}

// MappingMetrics はマッピングの成功率を追跡します
type MappingMetrics struct {
	TotalGroups      int
	SuccessfulGroups int
	TotalFields      int
	MappedFields     int
	FailedMappings   []FailedMapping
}

// FailedMapping はマッピング失敗の詳細を保持します
type FailedMapping struct {
	GroupName     string
	Operation     string
	FieldName     string
	TriedPatterns []string
	Reason        string
}

// EnhancedTypeInfo は拡張された型情報を保持します
type EnhancedTypeInfo struct {
	Name          string
	Fields        map[string]EnhancedFieldInfo
	IsRequestBody bool   // JSONRequestBody型かどうか
	IsParams      bool   // Params型かどうか
	HTTPMethod    string // 推測されるHTTPメソッド（Create, Update等）
	ResourceName  string // 推測されるリソース名
}

// EnhancedFieldInfo は拡張されたフィールド情報を保持します
type EnhancedFieldInfo struct {
	Name         string
	Type         string
	JSONTag      string
	Required     bool
	IsPointer    bool
	IsArray      bool
	IsNested     bool
	NestedType   string
	DefaultValue any
	Aliases      []string // フィールド名の別名候補
}

// TypeMetadata は型名から推測される情報を保持します
type TypeMetadata struct {
	HTTPMethod   string // "Create", "Update", "Get", "Delete"
	ResourceName string // "Feedback", "Comment", "User"
	TypeCategory string // "RequestBody", "Params", "Response"
}

// TypeNamePattern は型名パターンと優先度を保持します
type TypeNamePattern struct {
	Pattern  string
	Priority int    // 優先度（低いほど優先）
	Category string // パターンのカテゴリ（"exact", "standard", "variant", "fuzzy"）
}

// FieldMatchResult はマッチング結果を保持します
type FieldMatchResult struct {
	Field           EnhancedFieldInfo
	MatchStrategy   string  // 使用されたマッチング戦略
	Confidence      float64 // マッチの信頼度（0.0-1.0）
	TriedStrategies []string
}

// EnhancedMappingMetrics は拡張されたマッピングメトリクスを保持します
type EnhancedMappingMetrics struct {
	TotalGroups        int
	SuccessfulGroups   int
	TotalFields        int
	MappedFields       int
	FailedMappings     []DetailedFailedMapping
	TypePatternStats   map[string]int // 各パターンの使用回数
	MatchStrategyStats map[string]int // 各マッチング戦略の使用回数
}

// DetailedFailedMapping は詳細なマッピング失敗情報を保持します
type DetailedFailedMapping struct {
	ModuleName      string // モジュール名
	GroupName       string
	Operation       string
	FieldName       string
	TriedPatterns   []TypeNamePattern
	TriedStrategies []string
	TriedFilePaths  []string          // 試行したファイルパス
	BestMatch       *FieldMatchResult // 最も近かったマッチ
	Reason          string
	Suggestion      string // 改善提案
}

// TypeCacheEntry はキャッシュエントリを表します
type TypeCacheEntry struct {
	TypeInfo  map[string]EnhancedTypeInfo
	Timestamp int64 // キャッシュ作成時刻（Unix時間）
}

// PatternCacheEntry はパターン生成結果のキャッシュエントリを表します
type PatternCacheEntry struct {
	Patterns  []TypeNamePattern
	Timestamp int64
}

// ModuleMappingStats はモジュール別の統計を保持します
// 要件: 4.1, 4.2, 4.3
type ModuleMappingStats struct {
	ModuleName       string                  // モジュール名
	TotalGroups      int                     // 総グループ数
	SuccessfulGroups int                     // 成功したグループ数
	TotalFields      int                     // 総フィールド数
	MappedFields     int                     // マッピング成功フィールド数
	SuccessRate      float64                 // 成功率
	FailedMappings   []DetailedFailedMapping // 失敗したマッピング
}
