package snapshot

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// FallbackMechanism は、メソッド呼び出しの後方互換性を保証するフォールバック機構を提供します
// 要件 4.1, 4.2: 既存のメソッド呼び出しパターンを維持し、未知のメソッドパターンに対する安全な処理を実装
type FallbackMechanism struct {
	logger       *SnapshotLogger
	errorHandler *ReflectionErrorHandler
	strategies   []FallbackStrategy
}

// FallbackStrategy は、フォールバック戦略を定義するインターフェースです
type FallbackStrategy interface {
	// CanHandle は、この戦略が指定されたメソッドを処理できるかどうかを判定します
	CanHandle(methodName string, methodType reflect.Type, parameters interface{}) bool

	// Execute は、メソッドを実行します
	Execute(method reflect.Value, methodName string, methodType reflect.Type, ctx context.Context, parameters interface{}) ([]reflect.Value, error)

	// GetStrategyName は、戦略の名前を返します
	GetStrategyName() string
}

// NewFallbackMechanism は、新しいフallbackMechanismを作成します
func NewFallbackMechanism(logger *SnapshotLogger, errorHandler *ReflectionErrorHandler) *FallbackMechanism {
	fm := &FallbackMechanism{
		logger:       logger,
		errorHandler: errorHandler,
		strategies:   make([]FallbackStrategy, 0),
	}

	// デフォルトの戦略を登録（優先度順）
	fm.registerDefaultStrategies()

	return fm
}

// registerDefaultStrategies は、デフォルトのフォールバック戦略を登録します
func (fm *FallbackMechanism) registerDefaultStrategies() {
	// 1. Enhanced method calling strategy (最優先)
	fm.AddStrategy(&EnhancedMethodStrategy{
		logger:       fm.logger,
		errorHandler: fm.errorHandler,
	})

	// 2. Legacy method calling strategy (既存パターン)
	fm.AddStrategy(&LegacyMethodStrategy{
		logger:       fm.logger,
		errorHandler: fm.errorHandler,
	})

	// 3. Safe method calling strategy (未知パターン用)
	fm.AddStrategy(&SafeMethodStrategy{
		logger:       fm.logger,
		errorHandler: fm.errorHandler,
	})

	// 4. Emergency fallback strategy (最後の手段)
	fm.AddStrategy(&EmergencyFallbackStrategy{
		logger:       fm.logger,
		errorHandler: fm.errorHandler,
	})
}

// AddStrategy は、新しいフォールバック戦略を追加します
func (fm *FallbackMechanism) AddStrategy(strategy FallbackStrategy) {
	fm.strategies = append(fm.strategies, strategy)
	fm.logger.LogDebug(fmt.Sprintf("Added fallback strategy: %s", strategy.GetStrategyName()))
}

// ExecuteWithFallback は、フォールバック機構を使用してメソッドを実行します
// 要件 4.1, 4.2: 既存のメソッド呼び出しパターンを維持し、未知のメソッドパターンに対する安全な処理
func (fm *FallbackMechanism) ExecuteWithFallback(method reflect.Value, methodName string, methodType reflect.Type, ctx context.Context, parameters interface{}) ([]reflect.Value, error) {
	fm.logger.LogDebug(fmt.Sprintf("Starting fallback execution for method: %s", methodName))

	var lastError error
	var attemptedStrategies []string

	// 各戦略を順番に試行
	for i, strategy := range fm.strategies {
		strategyName := strategy.GetStrategyName()
		fm.logger.LogDebug(fmt.Sprintf("Trying strategy %d/%d: %s for method %s",
			i+1, len(fm.strategies), strategyName, methodName))

		// 戦略がこのメソッドを処理できるかチェック
		if !strategy.CanHandle(methodName, methodType, parameters) {
			fm.logger.LogDebug(fmt.Sprintf("Strategy %s cannot handle method %s", strategyName, methodName))
			continue
		}

		// 戦略を実行
		startTime := time.Now()
		results, err := strategy.Execute(method, methodName, methodType, ctx, parameters)
		duration := time.Since(startTime)

		attemptedStrategies = append(attemptedStrategies, strategyName)

		if err == nil {
			fm.logger.LogInfo(fmt.Sprintf("Strategy %s succeeded for method %s (duration: %v)",
				strategyName, methodName, duration))
			return results, nil
		}

		// エラーをログに記録
		fm.logger.LogDebug(fmt.Sprintf("Strategy %s failed for method %s: %v (duration: %v)",
			strategyName, methodName, err, duration))
		lastError = err
	}

	// すべての戦略が失敗した場合
	fm.logger.LogError(fmt.Sprintf("All fallback strategies failed for method %s", methodName), lastError)
	return nil, fmt.Errorf("all fallback strategies failed for method %s (attempted: %v): %w",
		methodName, attemptedStrategies, lastError)
}

// GetCompatibilityReport は、指定されたメソッドリストの互換性レポートを生成します
func (fm *FallbackMechanism) GetCompatibilityReport(client interface{}, methodNames []string) *FallbackCompatibilityReport {
	report := &FallbackCompatibilityReport{
		TotalMethods:      len(methodNames),
		CompatibleMethods: 0,
		StrategyUsage:     make(map[string]int),
		MethodStrategies:  make(map[string]string),
	}

	clientValue := reflect.ValueOf(client)

	for _, methodName := range methodNames {
		method := clientValue.MethodByName(methodName)
		if !method.IsValid() {
			fm.logger.LogWarning(fmt.Sprintf("Method %s not found on client", methodName))
			continue
		}

		methodType := method.Type()

		// どの戦略がこのメソッドを処理できるかチェック
		for _, strategy := range fm.strategies {
			if strategy.CanHandle(methodName, methodType, nil) {
				strategyName := strategy.GetStrategyName()
				report.CompatibleMethods++
				report.StrategyUsage[strategyName]++
				report.MethodStrategies[methodName] = strategyName
				break
			}
		}
	}

	if report.TotalMethods > 0 {
		report.CompatibilityPercentage = float64(report.CompatibleMethods) / float64(report.TotalMethods) * 100
	}

	return report
}

// FallbackCompatibilityReport は、フォールバック機構の互換性レポートを表します
type FallbackCompatibilityReport struct {
	TotalMethods            int
	CompatibleMethods       int
	CompatibilityPercentage float64
	StrategyUsage           map[string]int
	MethodStrategies        map[string]string
}

// GetSummary は、互換性レポートの要約を返します
func (r *FallbackCompatibilityReport) GetSummary() string {
	return fmt.Sprintf("Fallback Compatibility: %.1f%% (%d/%d methods compatible)",
		r.CompatibilityPercentage, r.CompatibleMethods, r.TotalMethods)
}

// EnhancedMethodStrategy は、拡張されたメソッド呼び出し戦略を実装します
// 要件 1.1, 1.2: メソッドシグネチャ解析を統合し、動的な引数構築ロジックを組み込む
type EnhancedMethodStrategy struct {
	logger       *SnapshotLogger
	errorHandler *ReflectionErrorHandler
}

// CanHandle は、この戦略が指定されたメソッドを処理できるかどうかを判定します
func (s *EnhancedMethodStrategy) CanHandle(methodName string, methodType reflect.Type, parameters interface{}) bool {
	// メソッドシグネチャアナライザーで解析可能なメソッドを処理
	analyzer := NewMethodSignatureAnalyzer(methodType, methodName)
	if err := analyzer.ValidateMethodSignature(); err != nil {
		return false
	}

	// 既知のパターンのメソッドを処理
	strategy := analyzer.AnalyzeSignature()
	return strategy == WithBodyMethod || strategy == StandardMethod || strategy == NoParameterMethod
}

// Execute は、拡張されたメソッド呼び出しを実行します
func (s *EnhancedMethodStrategy) Execute(method reflect.Value, methodName string, methodType reflect.Type, ctx context.Context, parameters interface{}) ([]reflect.Value, error) {
	s.logger.LogDebug(fmt.Sprintf("Executing enhanced method strategy for: %s", methodName))

	// メソッドシグネチャを解析
	analyzer := NewMethodSignatureAnalyzer(methodType, methodName)
	strategy := analyzer.AnalyzeSignature()

	// 適切なArgumentBuilderを選択
	var argumentBuilder ArgumentBuilder
	switch strategy {
	case WithBodyMethod:
		argumentBuilder = &WithBodyArgumentBuilder{}
	case NoParameterMethod:
		argumentBuilder = &NoParameterArgumentBuilder{}
	default:
		argumentBuilder = &StandardArgumentBuilder{}
	}

	// 引数を構築
	var args []reflect.Value
	var err error

	if typedBuilder, ok := argumentBuilder.(TypedArgumentBuilder); ok {
		args, err = typedBuilder.BuildArgumentsWithMethodType(ctx, parameters, methodType)
	} else {
		args, err = argumentBuilder.BuildArguments(ctx, parameters)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to build arguments: %w", err)
	}

	// 引数数の検証
	if len(args) != methodType.NumIn() {
		return nil, s.errorHandler.HandleArgumentMismatch(methodName, methodType, args)
	}

	// メソッドを実行
	return s.executeMethodCall(method, methodName, args)
}

// GetStrategyName は、戦略の名前を返します
func (s *EnhancedMethodStrategy) GetStrategyName() string {
	return "Enhanced"
}

// executeMethodCall は、実際のメソッド呼び出しを実行します
func (s *EnhancedMethodStrategy) executeMethodCall(method reflect.Value, methodName string, args []reflect.Value) ([]reflect.Value, error) {
	var results []reflect.Value
	var callErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				callErr = fmt.Errorf("panic during enhanced method call %s: %v", methodName, r)
			}
		}()

		methodType := method.Type()
		if methodType.IsVariadic() && len(args) > 0 {
			lastArgIndex := len(args) - 1
			lastArg := args[lastArgIndex]
			variadicType := methodType.In(methodType.NumIn() - 1)

			if lastArg.Kind() == reflect.Slice && lastArg.Type().AssignableTo(variadicType) {
				results = method.CallSlice(args)
				return
			}
		}

		results = method.Call(args)
	}()

	if callErr != nil {
		return nil, callErr
	}

	return results, nil
}

// LegacyMethodStrategy は、既存のメソッド呼び出しパターンを維持する戦略を実装します
// 要件 4.1: 既存のメソッド呼び出しパターンを維持するフォールバック処理
type LegacyMethodStrategy struct {
	logger       *SnapshotLogger
	errorHandler *ReflectionErrorHandler
}

// CanHandle は、この戦略が指定されたメソッドを処理できるかどうかを判定します
func (s *LegacyMethodStrategy) CanHandle(methodName string, methodType reflect.Type, parameters interface{}) bool {
	// 標準的なパターン（ctx + parameters）を処理
	numIn := methodType.NumIn()
	return numIn >= 1 && numIn <= 3 // ctx, parameters, reqEditors
}

// Execute は、レガシーメソッド呼び出しを実行します
func (s *LegacyMethodStrategy) Execute(method reflect.Value, methodName string, methodType reflect.Type, ctx context.Context, parameters interface{}) ([]reflect.Value, error) {
	s.logger.LogDebug(fmt.Sprintf("Executing legacy method strategy for: %s", methodName))

	args := []reflect.Value{reflect.ValueOf(ctx)}

	// パラメータを追加（メソッドが期待する場合）
	if parameters != nil && methodType.NumIn() > 1 {
		args = append(args, reflect.ValueOf(parameters))
	}

	// RequestEditorFnを追加（メソッドが期待する場合）
	if methodType.NumIn() > len(args) {
		reqEditorsType := reflect.TypeOf([]RequestEditorFn{})
		emptyReqEditors := reflect.MakeSlice(reqEditorsType, 0, 0)
		args = append(args, emptyReqEditors)
	}

	// 引数数の検証
	if len(args) != methodType.NumIn() {
		return nil, fmt.Errorf("legacy pattern argument count mismatch: expected %d, got %d", methodType.NumIn(), len(args))
	}

	// メソッドを実行
	return s.executeMethodCall(method, methodName, args)
}

// GetStrategyName は、戦略の名前を返します
func (s *LegacyMethodStrategy) GetStrategyName() string {
	return "Legacy"
}

// executeMethodCall は、実際のメソッド呼び出しを実行します
func (s *LegacyMethodStrategy) executeMethodCall(method reflect.Value, methodName string, args []reflect.Value) ([]reflect.Value, error) {
	var results []reflect.Value
	var callErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				callErr = fmt.Errorf("panic during legacy method call %s: %v", methodName, r)
			}
		}()

		results = method.Call(args)
	}()

	if callErr != nil {
		return nil, callErr
	}

	return results, nil
}

// SafeMethodStrategy は、未知のメソッドパターンに対する安全な処理を実装します
// 要件 4.2: 未知のメソッドパターンに対する安全な処理
type SafeMethodStrategy struct {
	logger       *SnapshotLogger
	errorHandler *ReflectionErrorHandler
}

// CanHandle は、この戦略が指定されたメソッドを処理できるかどうかを判定します
func (s *SafeMethodStrategy) CanHandle(methodName string, methodType reflect.Type, parameters interface{}) bool {
	// すべてのメソッドを処理可能（最後の手段として）
	return true
}

// Execute は、安全なメソッド呼び出しを実行します
func (s *SafeMethodStrategy) Execute(method reflect.Value, methodName string, methodType reflect.Type, ctx context.Context, parameters interface{}) ([]reflect.Value, error) {
	s.logger.LogDebug(fmt.Sprintf("Executing safe method strategy for: %s", methodName))

	numIn := methodType.NumIn()
	args := []reflect.Value{reflect.ValueOf(ctx)}

	// 引数を安全に構築
	if numIn == 1 {
		// contextのみ
		return s.executeMethodCall(method, methodName, args)
	}

	if numIn == 2 && parameters != nil {
		// context + parameters
		args = append(args, reflect.ValueOf(parameters))
		return s.executeMethodCall(method, methodName, args)
	}

	// より複雑なケース：ゼロ値で埋める
	if parameters != nil && numIn > 1 {
		args = append(args, reflect.ValueOf(parameters))
	}

	// 残りの引数をゼロ値で埋める
	for len(args) < numIn {
		paramType := methodType.In(len(args))
		zeroValue := reflect.Zero(paramType)
		args = append(args, zeroValue)
	}

	// 引数数を調整
	if len(args) > numIn {
		args = args[:numIn]
	}

	return s.executeMethodCall(method, methodName, args)
}

// GetStrategyName は、戦略の名前を返します
func (s *SafeMethodStrategy) GetStrategyName() string {
	return "Safe"
}

// executeMethodCall は、実際のメソッド呼び出しを実行します
func (s *SafeMethodStrategy) executeMethodCall(method reflect.Value, methodName string, args []reflect.Value) ([]reflect.Value, error) {
	var results []reflect.Value
	var callErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				callErr = fmt.Errorf("panic during safe method call %s: %v", methodName, r)
			}
		}()

		results = method.Call(args)
	}()

	if callErr != nil {
		return nil, callErr
	}

	return results, nil
}

// EmergencyFallbackStrategy は、最後の手段としての緊急フォールバック戦略を実装します
type EmergencyFallbackStrategy struct {
	logger       *SnapshotLogger
	errorHandler *ReflectionErrorHandler
}

// CanHandle は、この戦略が指定されたメソッドを処理できるかどうかを判定します
func (s *EmergencyFallbackStrategy) CanHandle(methodName string, methodType reflect.Type, parameters interface{}) bool {
	// 緊急時のみ使用（他のすべての戦略が失敗した場合）
	return true
}

// Execute は、緊急フォールバック呼び出しを実行します
func (s *EmergencyFallbackStrategy) Execute(method reflect.Value, methodName string, methodType reflect.Type, ctx context.Context, parameters interface{}) ([]reflect.Value, error) {
	s.logger.LogWarning(fmt.Sprintf("Using emergency fallback strategy for: %s", methodName))

	// 最も基本的なパターンを試行
	args := []reflect.Value{reflect.ValueOf(ctx)}

	// メソッドが1つの引数のみを期待する場合
	if methodType.NumIn() == 1 {
		return s.executeMethodCall(method, methodName, args)
	}

	// 緊急時：エラーを返すのではなく、ダミーの結果を返す
	s.logger.LogError(fmt.Sprintf("Emergency fallback: Cannot safely call method %s", methodName), nil)

	// ダミーの結果を作成（メソッドの戻り値型に基づく）
	numOut := methodType.NumOut()
	results := make([]reflect.Value, numOut)

	for i := 0; i < numOut; i++ {
		returnType := methodType.Out(i)
		results[i] = reflect.Zero(returnType)
	}

	return results, fmt.Errorf("emergency fallback used for method %s: method call was not executed", methodName)
}

// GetStrategyName は、戦略の名前を返します
func (s *EmergencyFallbackStrategy) GetStrategyName() string {
	return "Emergency"
}

// executeMethodCall は、実際のメソッド呼び出しを実行します
func (s *EmergencyFallbackStrategy) executeMethodCall(method reflect.Value, methodName string, args []reflect.Value) ([]reflect.Value, error) {
	var results []reflect.Value
	var callErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				callErr = fmt.Errorf("panic during emergency method call %s: %v", methodName, r)
			}
		}()

		results = method.Call(args)
	}()

	if callErr != nil {
		return nil, callErr
	}

	return results, nil
}

// IsWithBodyMethod は、メソッド名がWithBodyパターンかどうかを判定します
func IsWithBodyMethod(methodName string) bool {
	return strings.Contains(methodName, "WithBody")
}

// IsNoParameterMethod は、メソッドがパラメータを持たないかどうかを判定します
func IsNoParameterMethod(methodType reflect.Type) bool {
	return methodType.NumIn() == 1 // contextのみ
}

// GetMethodComplexity は、メソッドの複雑度を計算します
func GetMethodComplexity(methodType reflect.Type) int {
	complexity := 0

	// 引数数による複雑度
	complexity += methodType.NumIn()

	// 戻り値数による複雑度
	complexity += methodType.NumOut()

	// 可変長引数による複雑度
	if methodType.IsVariadic() {
		complexity += 2
	}

	return complexity
}
