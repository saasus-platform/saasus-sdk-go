package snapshot

import (
	"fmt"
	"reflect"
	"strings"
)

// MethodCallStrategy メソッド呼び出し戦略を表す列挙型
type MethodCallStrategy int

const (
	// StandardMethod 標準メソッド（ctx + parameters）
	StandardMethod MethodCallStrategy = iota
	// WithBodyMethod WithBodyメソッド（ctx + contentType + body + reqEditors）
	WithBodyMethod
	// NoParameterMethod パラメータなしメソッド（ctxのみ）
	NoParameterMethod
	// WithResponseMethod レスポンス付きメソッド
	WithResponseMethod
)

// String MethodCallStrategyの文字列表現を返す
func (m MethodCallStrategy) String() string {
	switch m {
	case StandardMethod:
		return "StandardMethod"
	case WithBodyMethod:
		return "WithBodyMethod"
	case NoParameterMethod:
		return "NoParameterMethod"
	case WithResponseMethod:
		return "WithResponseMethod"
	default:
		return "UnknownMethod"
	}
}

// ArgumentType 引数タイプを表す列挙型
type ArgumentType int

const (
	// ContextArg context.Context引数
	ContextArg ArgumentType = iota
	// ParameterArg パラメータ引数
	ParameterArg
	// ContentTypeArg Content-Type引数
	ContentTypeArg
	// BodyArg ボディ引数
	BodyArg
	// RequestEditorsArg RequestEditorFn引数
	RequestEditorsArg
)

// String ArgumentTypeの文字列表現を返す
func (a ArgumentType) String() string {
	switch a {
	case ContextArg:
		return "ContextArg"
	case ParameterArg:
		return "ParameterArg"
	case ContentTypeArg:
		return "ContentTypeArg"
	case BodyArg:
		return "BodyArg"
	case RequestEditorsArg:
		return "RequestEditorsArg"
	default:
		return "UnknownArg"
	}
}

// MethodCallInfo メソッド呼び出し情報を保持する構造体
type MethodCallInfo struct {
	Name          string
	Strategy      MethodCallStrategy
	RequiredArgs  []ArgumentType
	ParameterType reflect.Type
}

// MethodSignatureAnalyzer メソッドシグネチャを解析するアナライザー
type MethodSignatureAnalyzer struct {
	methodType reflect.Type
	methodName string
}

// NewMethodSignatureAnalyzer 新しいMethodSignatureAnalyzerを作成する
func NewMethodSignatureAnalyzer(methodType reflect.Type, methodName string) *MethodSignatureAnalyzer {
	return &MethodSignatureAnalyzer{
		methodType: methodType,
		methodName: methodName,
	}
}

// AnalyzeSignature メソッドシグネチャを解析してメソッド呼び出し戦略を決定する
func (m *MethodSignatureAnalyzer) AnalyzeSignature() MethodCallStrategy {
	if m.methodType == nil {
		return StandardMethod
	}

	numIn := m.methodType.NumIn()

	// WithBodyメソッドの検出
	if m.isWithBodyMethod() {
		if numIn >= 4 { // ctx, contentType, body, reqEditors...
			return WithBodyMethod
		}
	}

	// パラメータなしメソッド（ctx + reqEditors のみ）の検出
	if numIn == 2 && m.hasOnlyRequestEditorsAsSecondParam() {
		return NoParameterMethod
	}

	// パラメータなしメソッドの検出（contextのみ）
	if numIn == 1 {
		return NoParameterMethod
	}

	// 標準メソッドの検出（ctx + parameters）
	if numIn >= 2 {
		return StandardMethod
	}

	// デフォルトは標準メソッド
	return StandardMethod
}

// GetRequiredArguments 必要な引数タイプのリストを取得する
func (m *MethodSignatureAnalyzer) GetRequiredArguments() []ArgumentType {
	strategy := m.AnalyzeSignature()

	switch strategy {
	case WithBodyMethod:
		return []ArgumentType{ContextArg, ContentTypeArg, BodyArg, RequestEditorsArg}
	case NoParameterMethod:
		return []ArgumentType{ContextArg}
	case StandardMethod:
		return []ArgumentType{ContextArg, ParameterArg}
	case WithResponseMethod:
		return []ArgumentType{ContextArg, ParameterArg}
	default:
		return []ArgumentType{ContextArg, ParameterArg}
	}
}

// GetMethodCallInfo メソッド呼び出し情報を取得する
func (m *MethodSignatureAnalyzer) GetMethodCallInfo() MethodCallInfo {
	strategy := m.AnalyzeSignature()
	requiredArgs := m.GetRequiredArguments()

	var parameterType reflect.Type
	if m.methodType != nil && m.methodType.NumIn() > 1 {
		parameterType = m.methodType.In(1) // 2番目の引数（最初はcontext）
	}

	return MethodCallInfo{
		Name:          m.methodName,
		Strategy:      strategy,
		RequiredArgs:  requiredArgs,
		ParameterType: parameterType,
	}
}

// isWithBodyMethod WithBodyメソッドかどうかを判定する
func (m *MethodSignatureAnalyzer) isWithBodyMethod() bool {
	// メソッド名に"WithBody"が含まれているかチェック
	if strings.Contains(m.methodName, "WithBody") {
		// 名前にWithBodyが含まれる場合は、シグネチャも確認
		if m.methodType != nil && m.methodType.NumIn() >= 4 {
			return true
		}
		// 引数が不足している場合はWithBodyメソッドではない
		return false
	}

	// メソッドシグネチャからも判定
	if m.methodType != nil && m.methodType.NumIn() >= 4 {
		// 2番目の引数がstring（ContentType）、3番目がio.Reader（Body）の場合
		if m.methodType.NumIn() >= 3 {
			secondArgType := m.methodType.In(1)
			thirdArgType := m.methodType.In(2)

			// ContentTypeがstring型かチェック
			if secondArgType.Kind() == reflect.String {
				// BodyがReaderインターフェースを実装しているかチェック
				if m.implementsReader(thirdArgType) {
					return true
				}
			}
		}
	}

	return false
}

// implementsReader 型がio.Readerインターフェースを実装しているかチェック
func (m *MethodSignatureAnalyzer) implementsReader(t reflect.Type) bool {
	// io.Readerインターフェースの型を取得
	readerType := reflect.TypeOf((*interface {
		Read([]byte) (int, error)
	})(nil)).Elem()

	return t.Implements(readerType)
}

// hasOnlyRequestEditorsAsSecondParam 2番目の引数がRequestEditorFnスライスのみかどうかをチェック
func (m *MethodSignatureAnalyzer) hasOnlyRequestEditorsAsSecondParam() bool {
	if m.methodType == nil || m.methodType.NumIn() != 2 {
		return false
	}

	secondArgType := m.methodType.In(1)

	// スライス型かチェック
	if secondArgType.Kind() != reflect.Slice {
		return false
	}

	// スライスの要素型がRequestEditorFnかチェック
	elemType := secondArgType.Elem()
	elemTypeName := elemType.String()

	// RequestEditorFnの型名をチェック（パッケージ名は異なる可能性がある）
	// より厳密にチェックするため、関数型であることも確認
	if elemType.Kind() == reflect.Func && strings.Contains(elemTypeName, "RequestEditorFn") {
		return true
	}

	// 汎用的な関数型の場合（interface{}として定義されている場合）
	if elemType.Kind() == reflect.Interface && elemTypeName == "interface {}" {
		return true
	}

	return false
}

// ValidateMethodSignature メソッドシグネチャが期待される形式かどうかを検証する
func (m *MethodSignatureAnalyzer) ValidateMethodSignature() error {
	if m.methodType == nil {
		return fmt.Errorf("method type is nil for method %s", m.methodName)
	}

	numIn := m.methodType.NumIn()
	if numIn == 0 {
		return fmt.Errorf("method %s has no input parameters", m.methodName)
	}

	// 最初の引数がcontext.Contextかチェック
	firstArgType := m.methodType.In(0)
	contextTypeName := firstArgType.String()

	// context.Contextの型名をチェック
	if contextTypeName != "context.Context" {
		return fmt.Errorf("method %s first argument is not context.Context, got %s", m.methodName, contextTypeName)
	}

	return nil
}

// GetExpectedArgumentCount 期待される引数数を取得する
func (m *MethodSignatureAnalyzer) GetExpectedArgumentCount() int {
	if m.methodType == nil {
		return 0
	}
	return m.methodType.NumIn()
}

// GetActualArgumentTypes 実際の引数タイプのリストを取得する
func (m *MethodSignatureAnalyzer) GetActualArgumentTypes() []reflect.Type {
	if m.methodType == nil {
		return nil
	}

	var types []reflect.Type
	for i := 0; i < m.methodType.NumIn(); i++ {
		types = append(types, m.methodType.In(i))
	}

	return types
}

// GetReturnTypes 戻り値タイプのリストを取得する
func (m *MethodSignatureAnalyzer) GetReturnTypes() []reflect.Type {
	if m.methodType == nil {
		return nil
	}

	var types []reflect.Type
	for i := 0; i < m.methodType.NumOut(); i++ {
		types = append(types, m.methodType.Out(i))
	}

	return types
}

// IsVariadic メソッドが可変長引数を持つかどうかを判定する
func (m *MethodSignatureAnalyzer) IsVariadic() bool {
	if m.methodType == nil {
		return false
	}
	return m.methodType.IsVariadic()
}

// GetMethodSignatureString メソッドシグネチャの文字列表現を取得する
func (m *MethodSignatureAnalyzer) GetMethodSignatureString() string {
	if m.methodType == nil {
		return fmt.Sprintf("%s(unknown signature)", m.methodName)
	}

	var params []string
	for i := 0; i < m.methodType.NumIn(); i++ {
		paramType := m.methodType.In(i)
		params = append(params, paramType.String())
	}

	var returns []string
	for i := 0; i < m.methodType.NumOut(); i++ {
		returnType := m.methodType.Out(i)
		returns = append(returns, returnType.String())
	}

	paramStr := strings.Join(params, ", ")
	returnStr := strings.Join(returns, ", ")

	if len(returns) > 1 {
		returnStr = "(" + returnStr + ")"
	}

	return fmt.Sprintf("%s(%s) %s", m.methodName, paramStr, returnStr)
}
