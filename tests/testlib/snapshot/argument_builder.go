package snapshot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// ArgumentBuilder インターフェースは、メソッド呼び出しのための引数を構築する戦略を定義します
type ArgumentBuilder interface {
	BuildArguments(ctx context.Context, parameters interface{}) ([]reflect.Value, error)
}

// TypedArgumentBuilder インターフェースは、メソッド型を考慮した引数構築をサポートします
type TypedArgumentBuilder interface {
	ArgumentBuilder
	BuildArgumentsWithMethodType(ctx context.Context, parameters interface{}, methodType reflect.Type) ([]reflect.Value, error)
}

// StandardArgumentBuilder は、標準的なメソッド（ctx + parameters）の引数を構築します
// TypedArgumentBuilderインターフェースを実装
type StandardArgumentBuilder struct{}

// BuildArguments は、標準メソッドの引数を構築します
func (s *StandardArgumentBuilder) BuildArguments(ctx context.Context, parameters interface{}) ([]reflect.Value, error) {
	args := []reflect.Value{reflect.ValueOf(ctx)}

	if parameters != nil {
		args = append(args, reflect.ValueOf(parameters))
	}

	return args, nil
}

// BuildArgumentsWithMethodType は、メソッド型を考慮して標準メソッドの引数を構築します
func (s *StandardArgumentBuilder) BuildArgumentsWithMethodType(ctx context.Context, parameters interface{}, methodType reflect.Type) ([]reflect.Value, error) {
	args := []reflect.Value{reflect.ValueOf(ctx)}

	if parameters != nil {
		paramValue := reflect.ValueOf(parameters)

		// If parameters is a struct, extract its fields as separate arguments only when
		// the method expects multiple non-context parameters. Otherwise, pass the struct as-is.
		if paramValue.Kind() == reflect.Struct {
			// Determine the number of non-context, non-variadic parameters
			expectedArgs := methodType.NumIn()
			if methodType.IsVariadic() {
				expectedArgs-- // Exclude variadic argument
			}
			expectedArgs-- // Exclude context argument

			paramType := paramValue.Type()
			targetType := methodType.In(1)
			handledStruct := false

			// If struct type already matches the target type, use it directly
			if paramType.AssignableTo(targetType) {
				args = append(args, paramValue)
				handledStruct = true
			}

			// If pointer to struct matches the target type, pass pointer
			if !handledStruct && paramValue.CanAddr() {
				addrType := paramValue.Addr().Type()
				if addrType.AssignableTo(targetType) {
					args = append(args, paramValue.Addr())
					handledStruct = true
				}
			}

			numFields := paramValue.NumField()

			// If the method expects multiple parameters (excluding context),
			// and the struct fields match that count, expand the struct into fields.
			if !handledStruct && expectedArgs > 1 && numFields == expectedArgs {
				for i := 0; i < numFields; i++ {
					args = append(args, paramValue.Field(i))
				}
				handledStruct = true
			}

			// For single-field structs, attempt to map the field to the expected parameter type.
			if !handledStruct && numFields == 1 {
				field := paramValue.Field(0)
				if field.Type().AssignableTo(targetType) {
					args = append(args, field)
					handledStruct = true
				} else if field.Type().ConvertibleTo(targetType) {
					args = append(args, field.Convert(targetType))
					handledStruct = true
				}
			}

			// Fallback: use the struct as-is (common for request body structs)
			if !handledStruct {
				args = append(args, paramValue)
			}
		} else {
			// For non-struct parameters, use as-is
			args = append(args, paramValue)
		}
	}

	// 可変長引数（RequestEditorFn）がある場合は空のスライスを追加
	if methodType != nil && methodType.NumIn() >= 3 && methodType.IsVariadic() {
		reqEditorsType := methodType.In(methodType.NumIn() - 1) // 最後の引数（可変長引数）
		if reqEditorsType.Kind() == reflect.Slice {
			emptyReqEditors := reflect.MakeSlice(reqEditorsType, 0, 0)
			args = append(args, emptyReqEditors)
		}
	}

	return args, nil
}

// WithBodyArgumentBuilder は、WithBodyメソッド（ctx + contentType + body + reqEditors）の引数を構築します
// TypedArgumentBuilderインターフェースを実装
type WithBodyArgumentBuilder struct{}

// BuildArguments は、WithBodyメソッドの引数を構築します
// createUpdateStripeInfoBodyParams()の戻り値を適切に分解して処理します
func (w *WithBodyArgumentBuilder) BuildArguments(ctx context.Context, parameters interface{}) ([]reflect.Value, error) {
	args := []reflect.Value{reflect.ValueOf(ctx)}

	// parametersがBodyParamsインターフェースを実装しているかチェック
	if bodyParams, ok := parameters.(BodyParams); ok {
		return w.buildFromBodyParams(args, bodyParams)
	}

	// createUpdateStripeInfoBodyParams()の戻り値のような構造体を処理
	if parameters != nil {
		return w.buildFromStructParams(args, parameters)
	}

	return nil, fmt.Errorf("parameters do not contain required ContentType and Body fields for WithBody method")
}

// BuildArgumentsWithMethodType は、メソッド型を考慮してWithBodyメソッドの引数を構築します
func (w *WithBodyArgumentBuilder) BuildArgumentsWithMethodType(ctx context.Context, parameters interface{}, methodType reflect.Type) ([]reflect.Value, error) {
	args := []reflect.Value{reflect.ValueOf(ctx)}

	// parametersがBodyParamsインターフェースを実装しているかチェック
	if bodyParams, ok := parameters.(BodyParams); ok {
		return w.buildFromBodyParamsWithMethodType(args, bodyParams, methodType)
	}

	// createUpdateStripeInfoBodyParams()の戻り値のような構造体を処理
	if parameters != nil {
		return w.buildFromStructParamsWithMethodType(args, parameters, methodType)
	}

	return nil, fmt.Errorf("parameters do not contain required ContentType and Body fields for WithBody method")
}

// buildFromBodyParams は、BodyParamsインターフェースから引数を構築します
func (w *WithBodyArgumentBuilder) buildFromBodyParams(args []reflect.Value, bodyParams BodyParams) ([]reflect.Value, error) {
	// ContentTypeとBodyを個別の引数として抽出
	args = append(args, reflect.ValueOf(bodyParams.GetContentType()))
	args = append(args, reflect.ValueOf(bodyParams.GetBody()))

	// RequestEditorFnスライスを空配列として追加
	args = append(args, w.createEmptyRequestEditors())

	return args, nil
}

// buildFromStructParams は、createUpdateStripeInfoBodyParams()のような構造体から引数を構築します
func (w *WithBodyArgumentBuilder) buildFromStructParams(args []reflect.Value, parameters interface{}) ([]reflect.Value, error) {
	paramValue := reflect.ValueOf(parameters)
	if paramValue.Kind() == reflect.Ptr {
		paramValue = paramValue.Elem()
	}

	if paramValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("parameters must be a struct with ContentType and Body fields")
	}

	// ContentTypeフィールドを抽出
	contentTypeField := paramValue.FieldByName("ContentType")
	if !contentTypeField.IsValid() {
		return nil, fmt.Errorf("parameters struct must have ContentType field")
	}

	// Bodyフィールドを抽出
	bodyField := paramValue.FieldByName("Body")
	if !bodyField.IsValid() {
		return nil, fmt.Errorf("parameters struct must have Body field")
	}

	// ContentTypeとBodyを個別の引数として追加
	args = append(args, contentTypeField)
	args = append(args, bodyField)

	// RequestEditorFnスライスを空配列として追加
	args = append(args, w.createEmptyRequestEditors())

	return args, nil
}

// buildFromStructParamsWithMethodType は、メソッド型を考慮して構造体から引数を構築します
func (w *WithBodyArgumentBuilder) buildFromStructParamsWithMethodType(args []reflect.Value, parameters interface{}, methodType reflect.Type) ([]reflect.Value, error) {
	paramValue := reflect.ValueOf(parameters)
	if paramValue.Kind() == reflect.Ptr {
		paramValue = paramValue.Elem()
	}

	if paramValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("parameters must be a struct with ContentType and Body fields")
	}

	// ContentTypeフィールドを抽出
	contentTypeField := paramValue.FieldByName("ContentType")
	if !contentTypeField.IsValid() {
		return nil, fmt.Errorf("parameters struct must have ContentType field")
	}

	// Bodyフィールドを抽出
	bodyField := paramValue.FieldByName("Body")
	if !bodyField.IsValid() {
		return nil, fmt.Errorf("parameters struct must have Body field")
	}

	// メソッドが期待する引数数を計算（context + 可変長引数を除く）
	expectedArgs := methodType.NumIn()
	if methodType.IsVariadic() {
		expectedArgs-- // 可変長引数を除く
	}
	expectedArgs-- // contextを除く

	// ContentTypeとBodyの前に追加のパラメータがある場合（例：feedbackId）
	if expectedArgs > 2 {
		// 構造体の全フィールドを順番に追加（ContentTypeとBodyを除く）
		paramType := paramValue.Type()
		for i := 0; i < paramValue.NumField(); i++ {
			fieldName := paramType.Field(i).Name
			if fieldName != "ContentType" && fieldName != "Body" {
				args = append(args, paramValue.Field(i))
			}
		}
	}

	// ContentTypeとBodyを個別の引数として追加
	args = append(args, contentTypeField)
	args = append(args, bodyField)

	// メソッド型から正しいRequestEditorFn型を取得して空配列を作成
	if methodType != nil && methodType.NumIn() >= 4 {
		reqEditorsType := methodType.In(methodType.NumIn() - 1) // 最後の引数（可変長引数）
		args = append(args, w.createEmptyRequestEditorsForType(reqEditorsType))
	} else {
		// フォールバック：汎用的な空のスライス
		args = append(args, w.createEmptyRequestEditors())
	}

	return args, nil
}

// buildFromBodyParamsWithMethodType は、メソッド型を考慮してBodyParamsから引数を構築します
func (w *WithBodyArgumentBuilder) buildFromBodyParamsWithMethodType(args []reflect.Value, bodyParams BodyParams, methodType reflect.Type) ([]reflect.Value, error) {
	// ContentTypeとBodyを個別の引数として抽出
	args = append(args, reflect.ValueOf(bodyParams.GetContentType()))
	args = append(args, reflect.ValueOf(bodyParams.GetBody()))

	// メソッド型から正しいRequestEditorFn型を取得して空配列を作成
	if methodType != nil && methodType.NumIn() >= 4 {
		reqEditorsType := methodType.In(methodType.NumIn() - 1) // 最後の引数（可変長引数）
		args = append(args, w.createEmptyRequestEditorsForType(reqEditorsType))
	} else {
		// フォールバック：汎用的な空のスライス
		args = append(args, w.createEmptyRequestEditors())
	}

	return args, nil
}

// createEmptyRequestEditors は、空のRequestEditorFnスライスを作成します
func (w *WithBodyArgumentBuilder) createEmptyRequestEditors() reflect.Value {
	// 可変長引数として渡すため、スライスの各要素を個別に展開する必要がある
	// 空のスライスの場合は、何も追加しない
	reqEditorsType := reflect.TypeOf([]RequestEditorFn{})
	return reflect.MakeSlice(reqEditorsType, 0, 0)
}

// createEmptyRequestEditorsForType は、指定された型の空のRequestEditorFnスライスを作成します
func (w *WithBodyArgumentBuilder) createEmptyRequestEditorsForType(targetType reflect.Type) reflect.Value {
	// 可変長引数の型を取得（[]RequestEditorFn）
	if targetType.Kind() == reflect.Slice {
		return reflect.MakeSlice(targetType, 0, 0)
	}
	// フォールバック：汎用的な空のスライス
	reqEditorsType := reflect.TypeOf([]RequestEditorFn{})
	return reflect.MakeSlice(reqEditorsType, 0, 0)
}

// NoParameterArgumentBuilder は、パラメータなしメソッド（ctxのみ）の引数を構築します
// TypedArgumentBuilderインターフェースを実装
type NoParameterArgumentBuilder struct{}

// BuildArguments は、パラメータなしメソッドの引数を構築します
func (n *NoParameterArgumentBuilder) BuildArguments(ctx context.Context, parameters interface{}) ([]reflect.Value, error) {
	args := []reflect.Value{reflect.ValueOf(ctx)}

	// デフォルトではcontextのみを返す
	// メソッド型が必要な場合はBuildArgumentsWithMethodTypeを使用
	return args, nil
}

// BuildArgumentsWithMethodType は、メソッド型を考慮してパラメータなしメソッドの引数を構築します
func (n *NoParameterArgumentBuilder) BuildArgumentsWithMethodType(ctx context.Context, parameters interface{}, methodType reflect.Type) ([]reflect.Value, error) {
	args := []reflect.Value{reflect.ValueOf(ctx)}

	// メソッドが期待する引数数をチェック
	if methodType != nil {
		expectedArgCount := methodType.NumIn()

		// contextのみを期待するメソッド（引数数1）の場合
		if expectedArgCount == 1 {
			return args, nil
		}

		// context + RequestEditorFnを期待するメソッド（引数数2以上）の場合
		if expectedArgCount >= 2 {
			reqEditorsType := methodType.In(methodType.NumIn() - 1) // 最後の引数（可変長引数）
			if reqEditorsType.Kind() == reflect.Slice {
				emptyReqEditors := reflect.MakeSlice(reqEditorsType, 0, 0)
				args = append(args, emptyReqEditors)
			} else {
				// フォールバック：汎用的な空のスライス
				reqEditorsType := reflect.TypeOf([]RequestEditorFn{})
				emptyReqEditors := reflect.MakeSlice(reqEditorsType, 0, 0)
				args = append(args, emptyReqEditors)
			}
		}
	} else {
		// methodTypeがnilの場合はcontextのみを返す
		return args, nil
	}

	return args, nil
}

// BodyParams インターフェースは、WithBodyメソッドで使用されるパラメータを定義します
type BodyParams interface {
	GetContentType() string
	GetBody() io.Reader
}

// RequestEditorFn は、リクエストエディター関数の型を定義します
// 生成されたコードと同じ型定義を使用
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// ArgumentBuilderFactory は、メソッドタイプに基づいて適切なArgumentBuilderを作成します
type ArgumentBuilderFactory struct {
	builders map[MethodCallStrategy]ArgumentBuilder
}

// NewArgumentBuilderFactory は、新しいArgumentBuilderFactoryを作成します
func NewArgumentBuilderFactory() *ArgumentBuilderFactory {
	return &ArgumentBuilderFactory{
		builders: map[MethodCallStrategy]ArgumentBuilder{
			StandardMethod:    &StandardArgumentBuilder{},
			WithBodyMethod:    &WithBodyArgumentBuilder{},
			NoParameterMethod: &NoParameterArgumentBuilder{},
		},
	}
}

// GetBuilder は、指定された戦略に対応するArgumentBuilderを返します
func (f *ArgumentBuilderFactory) GetBuilder(strategy MethodCallStrategy) (ArgumentBuilder, error) {
	builder, exists := f.builders[strategy]
	if !exists {
		return nil, fmt.Errorf("no argument builder found for strategy: %d", strategy)
	}
	return builder, nil
}

// BuildArgumentsForMethod は、メソッド名とパラメータから適切な引数を構築します
func BuildArgumentsForMethod(methodName string, methodType reflect.Type, ctx context.Context, parameters interface{}) ([]reflect.Value, error) {
	factory := NewArgumentBuilderFactory()

	// 既存のMethodSignatureAnalyzerを使用して戦略を決定
	analyzer := NewMethodSignatureAnalyzer(methodType, methodName)
	strategy := analyzer.AnalyzeSignature()

	builder, err := factory.GetBuilder(strategy)
	if err != nil {
		return nil, fmt.Errorf("failed to get argument builder for method %s: %w", methodName, err)
	}

	return builder.BuildArguments(ctx, parameters)
}

// ExtractBodyParamsFromStruct は、createUpdateStripeInfoBodyParams()のような構造体から
// ContentTypeとBodyを抽出するヘルパー関数です
func ExtractBodyParamsFromStruct(parameters interface{}) (contentType string, body io.Reader, err error) {
	if parameters == nil {
		return "", nil, fmt.Errorf("parameters cannot be nil")
	}

	paramValue := reflect.ValueOf(parameters)
	if paramValue.Kind() == reflect.Ptr {
		paramValue = paramValue.Elem()
	}

	if paramValue.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("parameters must be a struct")
	}

	// ContentTypeフィールドを抽出
	contentTypeField := paramValue.FieldByName("ContentType")
	if !contentTypeField.IsValid() {
		return "", nil, fmt.Errorf("struct must have ContentType field")
	}

	contentTypeStr, ok := contentTypeField.Interface().(string)
	if !ok {
		return "", nil, fmt.Errorf("ContentType field must be a string")
	}

	// Bodyフィールドを抽出
	bodyField := paramValue.FieldByName("Body")
	if !bodyField.IsValid() {
		return "", nil, fmt.Errorf("struct must have Body field")
	}

	bodyReader, ok := bodyField.Interface().(io.Reader)
	if !ok {
		return "", nil, fmt.Errorf("Body field must implement io.Reader")
	}

	return contentTypeStr, bodyReader, nil
}

// ValidateWithBodyParams は、WithBodyメソッドのパラメータが有効かどうかを検証します
func ValidateWithBodyParams(parameters interface{}) error {
	_, _, err := ExtractBodyParamsFromStruct(parameters)
	return err
}
