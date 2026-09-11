package snapshot

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// CompatibilityChecker は、既存テストケースが影響を受けないことを保証する機能を提供します
// 要件 4.1, 4.2: 既存テストケースが影響を受けないことを保証する機能
type CompatibilityChecker struct {
	logger            *SnapshotLogger
	fallbackMechanism *FallbackMechanism
	client            interface{}
}

// NewCompatibilityChecker は、新しいCompatibilityCheckerを作成します
func NewCompatibilityChecker(logger *SnapshotLogger, fallbackMechanism *FallbackMechanism, client interface{}) *CompatibilityChecker {
	return &CompatibilityChecker{
		logger:            logger,
		fallbackMechanism: fallbackMechanism,
		client:            client,
	}
}

// CompatibilityCheckResult は、互換性チェックの結果を表します
type CompatibilityCheckResult struct {
	TotalStories            int
	CompatibleStories       int
	IncompatibleStories     int
	CompatibilityPercentage float64
	StoryResults            map[string]*StoryCompatibilityResult
	Summary                 *CompatibilitySummary
	Timestamp               time.Time
}

// StoryCompatibilityResult は、個別のストーリーの互換性結果を表します
type StoryCompatibilityResult struct {
	StoryName          string
	IsCompatible       bool
	TotalSteps         int
	CompatibleSteps    int
	IncompatibleSteps  int
	StepResults        map[string]*StepCompatibilityResult
	RequiredStrategies []string
	Issues             []string
}

// StepCompatibilityResult は、個別のステップの互換性結果を表します
type StepCompatibilityResult struct {
	StepName         string
	MethodName       string
	IsCompatible     bool
	RequiredStrategy string
	Issues           []string
	MethodSignature  string
}

// CompatibilitySummary は、互換性チェックの要約を表します
type CompatibilitySummary struct {
	OverallCompatibility float64
	MethodCompatibility  float64
	StrategyUsage        map[string]int
	CommonIssues         []string
	RecommendedActions   []string
}

// CheckStoriesCompatibility は、ストーリーリストの互換性をチェックします
// 要件 4.1, 4.2: 既存テストケースが影響を受けないことを保証する機能
func (c *CompatibilityChecker) CheckStoriesCompatibility(stories []testlib.Story) *CompatibilityCheckResult {
	c.logger.LogInfo(fmt.Sprintf("Starting compatibility check for %d stories", len(stories)))
	startTime := time.Now()

	result := &CompatibilityCheckResult{
		TotalStories: len(stories),
		StoryResults: make(map[string]*StoryCompatibilityResult),
		Timestamp:    startTime,
	}

	// 各ストーリーの互換性をチェック
	for _, story := range stories {
		storyResult := c.checkStoryCompatibility(story)
		result.StoryResults[story.Name] = storyResult

		if storyResult.IsCompatible {
			result.CompatibleStories++
		} else {
			result.IncompatibleStories++
		}
	}

	// 互換性パーセンテージを計算
	if result.TotalStories > 0 {
		result.CompatibilityPercentage = float64(result.CompatibleStories) / float64(result.TotalStories) * 100
	}

	// サマリーを生成
	result.Summary = c.generateCompatibilitySummary(result)

	duration := time.Since(startTime)
	c.logger.LogInfo(fmt.Sprintf("Compatibility check completed in %v: %.1f%% compatible (%d/%d stories)",
		duration, result.CompatibilityPercentage, result.CompatibleStories, result.TotalStories))

	return result
}

// checkStoryCompatibility は、個別のストーリーの互換性をチェックします
func (c *CompatibilityChecker) checkStoryCompatibility(story testlib.Story) *StoryCompatibilityResult {
	c.logger.LogDebug(fmt.Sprintf("Checking compatibility for story: %s", story.Name))

	result := &StoryCompatibilityResult{
		StoryName:          story.Name,
		TotalSteps:         len(story.Steps),
		StepResults:        make(map[string]*StepCompatibilityResult),
		RequiredStrategies: make([]string, 0),
		Issues:             make([]string, 0),
	}

	clientValue := reflect.ValueOf(c.client)
	strategyUsage := make(map[string]bool)

	// 各ステップの互換性をチェック
	for _, step := range story.Steps {
		stepResult := c.checkStepCompatibility(step, clientValue)
		result.StepResults[step.Name] = stepResult

		if stepResult.IsCompatible {
			result.CompatibleSteps++
		} else {
			result.IncompatibleSteps++
			result.Issues = append(result.Issues, fmt.Sprintf("Step '%s': %s", step.Name, strings.Join(stepResult.Issues, ", ")))
		}

		// 使用される戦略を記録
		if stepResult.RequiredStrategy != "" {
			strategyUsage[stepResult.RequiredStrategy] = true
		}
	}

	// 使用される戦略のリストを作成
	for strategy := range strategyUsage {
		result.RequiredStrategies = append(result.RequiredStrategies, strategy)
	}

	// ストーリー全体の互換性を判定
	result.IsCompatible = result.IncompatibleSteps == 0

	c.logger.LogDebug(fmt.Sprintf("Story '%s' compatibility: %t (%d/%d steps compatible)",
		story.Name, result.IsCompatible, result.CompatibleSteps, result.TotalSteps))

	return result
}

// checkStepCompatibility は、個別のステップの互換性をチェックします
func (c *CompatibilityChecker) checkStepCompatibility(step testlib.Step, clientValue reflect.Value) *StepCompatibilityResult {
	result := &StepCompatibilityResult{
		StepName:   step.Name,
		MethodName: step.ClientMethod,
		Issues:     make([]string, 0),
	}

	// メソッドの存在をチェック
	method := clientValue.MethodByName(step.ClientMethod)
	if !method.IsValid() {
		result.IsCompatible = false
		result.Issues = append(result.Issues, fmt.Sprintf("Method '%s' not found", step.ClientMethod))
		return result
	}

	methodType := method.Type()
	result.MethodSignature = c.getMethodSignatureString(methodType)

	// フォールバック機構で処理可能かチェック
	compatibilityReport := c.fallbackMechanism.GetCompatibilityReport(c.client, []string{step.ClientMethod})

	if compatibilityReport.CompatibleMethods > 0 {
		result.IsCompatible = true
		if strategy, exists := compatibilityReport.MethodStrategies[step.ClientMethod]; exists {
			result.RequiredStrategy = strategy
		}
	} else {
		result.IsCompatible = false
		result.Issues = append(result.Issues, "No compatible fallback strategy found")
	}

	// 特定のパターンに対する追加チェック
	c.performAdditionalCompatibilityChecks(step, methodType, result)

	return result
}

// performAdditionalCompatibilityChecks は、特定のパターンに対する追加の互換性チェックを実行します
func (c *CompatibilityChecker) performAdditionalCompatibilityChecks(step testlib.Step, methodType reflect.Type, result *StepCompatibilityResult) {
	// WithBodyメソッドの特別なチェック
	if IsWithBodyMethod(step.ClientMethod) {
		if err := c.checkWithBodyMethodCompatibility(step, methodType); err != nil {
			result.Issues = append(result.Issues, fmt.Sprintf("WithBody method issue: %v", err))
		}
	}

	// パラメータなしメソッドのチェック
	if IsNoParameterMethod(methodType) && step.Parameters != nil {
		result.Issues = append(result.Issues, "Method expects no parameters but parameters provided")
	}

	// 引数数の基本チェック
	expectedArgs := methodType.NumIn()
	if expectedArgs == 0 {
		result.Issues = append(result.Issues, "Method expects no arguments")
	}

	// 可変長引数メソッドのチェック
	if methodType.IsVariadic() {
		c.logger.LogDebug(fmt.Sprintf("Method %s is variadic", step.ClientMethod))
	}
}

// checkWithBodyMethodCompatibility は、WithBodyメソッドの互換性をチェックします
func (c *CompatibilityChecker) checkWithBodyMethodCompatibility(step testlib.Step, methodType reflect.Type) error {
	if step.Parameters == nil {
		return fmt.Errorf("WithBody method requires parameters")
	}

	// パラメータがBodyParamsインターフェースを実装しているかチェック
	if _, ok := step.Parameters.(BodyParams); ok {
		return nil // OK
	}

	// 構造体からContentTypeとBodyを抽出できるかチェック
	if err := ValidateWithBodyParams(step.Parameters); err != nil {
		return fmt.Errorf("invalid WithBody parameters: %w", err)
	}

	return nil
}

// generateCompatibilitySummary は、互換性チェックのサマリーを生成します
func (c *CompatibilityChecker) generateCompatibilitySummary(result *CompatibilityCheckResult) *CompatibilitySummary {
	summary := &CompatibilitySummary{
		OverallCompatibility: result.CompatibilityPercentage,
		StrategyUsage:        make(map[string]int),
		CommonIssues:         make([]string, 0),
		RecommendedActions:   make([]string, 0),
	}

	// 戦略使用状況を集計
	totalMethods := 0
	compatibleMethods := 0
	issueCount := make(map[string]int)

	for _, storyResult := range result.StoryResults {
		for _, stepResult := range storyResult.StepResults {
			totalMethods++
			if stepResult.IsCompatible {
				compatibleMethods++
				if stepResult.RequiredStrategy != "" {
					summary.StrategyUsage[stepResult.RequiredStrategy]++
				}
			}

			// 問題を集計
			for _, issue := range stepResult.Issues {
				issueCount[issue]++
			}
		}
	}

	// メソッド互換性を計算
	if totalMethods > 0 {
		summary.MethodCompatibility = float64(compatibleMethods) / float64(totalMethods) * 100
	}

	// 一般的な問題を特定
	for issue, count := range issueCount {
		if count >= 2 { // 2回以上発生した問題を一般的な問題とする
			summary.CommonIssues = append(summary.CommonIssues, fmt.Sprintf("%s (occurs %d times)", issue, count))
		}
	}

	// 推奨アクションを生成
	summary.RecommendedActions = c.generateRecommendedActions(summary, result)

	return summary
}

// generateRecommendedActions は、推奨アクションを生成します
func (c *CompatibilityChecker) generateRecommendedActions(summary *CompatibilitySummary, result *CompatibilityCheckResult) []string {
	actions := make([]string, 0)

	// 互換性が低い場合の推奨アクション
	if summary.OverallCompatibility < 90 {
		actions = append(actions, "Consider updating test cases to use compatible method patterns")
	}

	// 特定の戦略が多用されている場合
	for strategy, count := range summary.StrategyUsage {
		if strategy == "Emergency" && count > 0 {
			actions = append(actions, "Some methods require emergency fallback - consider implementing proper support")
		}
		if strategy == "Safe" && count > 5 {
			actions = append(actions, "Many methods use safe fallback - consider optimizing method signatures")
		}
	}

	// 一般的な問題に基づく推奨アクション
	for _, issue := range summary.CommonIssues {
		if strings.Contains(issue, "WithBody") {
			actions = append(actions, "Update WithBody method parameter handling")
		}
		if strings.Contains(issue, "not found") {
			actions = append(actions, "Verify client method names and availability")
		}
	}

	if len(actions) == 0 {
		actions = append(actions, "No specific actions required - compatibility is good")
	}

	return actions
}

// getMethodSignatureString は、メソッドシグネチャの文字列表現を返します
func (c *CompatibilityChecker) getMethodSignatureString(methodType reflect.Type) string {
	// 引数
	var args []string
	for i := 0; i < methodType.NumIn(); i++ {
		argType := methodType.In(i)
		args = append(args, argType.String())
	}

	// 戻り値
	var returns []string
	for i := 0; i < methodType.NumOut(); i++ {
		returnType := methodType.Out(i)
		returns = append(returns, returnType.String())
	}

	signature := fmt.Sprintf("(%s) (%s)", strings.Join(args, ", "), strings.Join(returns, ", "))

	if methodType.IsVariadic() {
		signature += " [variadic]"
	}

	return signature
}

// PrintCompatibilityReport は、互換性レポートを読みやすい形式で出力します
func (c *CompatibilityChecker) PrintCompatibilityReport(result *CompatibilityCheckResult) {
	c.logger.LogInfo("=== Compatibility Check Report ===")
	c.logger.LogInfo(fmt.Sprintf("Timestamp: %s", result.Timestamp.Format(time.RFC3339)))
	c.logger.LogInfo(fmt.Sprintf("Total Stories: %d", result.TotalStories))
	c.logger.LogInfo(fmt.Sprintf("Compatible Stories: %d (%.1f%%)", result.CompatibleStories, result.CompatibilityPercentage))
	c.logger.LogInfo(fmt.Sprintf("Incompatible Stories: %d", result.IncompatibleStories))

	if result.Summary != nil {
		c.logger.LogInfo(fmt.Sprintf("Method Compatibility: %.1f%%", result.Summary.MethodCompatibility))

		if len(result.Summary.StrategyUsage) > 0 {
			c.logger.LogInfo("Strategy Usage:")
			for strategy, count := range result.Summary.StrategyUsage {
				c.logger.LogInfo(fmt.Sprintf("  - %s: %d methods", strategy, count))
			}
		}

		if len(result.Summary.CommonIssues) > 0 {
			c.logger.LogInfo("Common Issues:")
			for _, issue := range result.Summary.CommonIssues {
				c.logger.LogInfo(fmt.Sprintf("  - %s", issue))
			}
		}

		if len(result.Summary.RecommendedActions) > 0 {
			c.logger.LogInfo("Recommended Actions:")
			for _, action := range result.Summary.RecommendedActions {
				c.logger.LogInfo(fmt.Sprintf("  - %s", action))
			}
		}
	}

	// 非互換ストーリーの詳細
	if result.IncompatibleStories > 0 {
		c.logger.LogInfo("Incompatible Stories Details:")
		for storyName, storyResult := range result.StoryResults {
			if !storyResult.IsCompatible {
				c.logger.LogInfo(fmt.Sprintf("  Story: %s (%d/%d steps compatible)",
					storyName, storyResult.CompatibleSteps, storyResult.TotalSteps))
				for _, issue := range storyResult.Issues {
					c.logger.LogInfo(fmt.Sprintf("    - %s", issue))
				}
			}
		}
	}

	c.logger.LogInfo("=== End Compatibility Report ===")
}

// ValidateTestCaseCompatibility は、テストケースの互換性を検証します
// 要件 4.1, 4.2: 既存テストケースが影響を受けないことを保証する機能
func (c *CompatibilityChecker) ValidateTestCaseCompatibility(stories []testlib.Story, minCompatibilityThreshold float64) error {
	result := c.CheckStoriesCompatibility(stories)

	if result.CompatibilityPercentage < minCompatibilityThreshold {
		return fmt.Errorf("compatibility check failed: %.1f%% < %.1f%% threshold (%d/%d stories compatible)",
			result.CompatibilityPercentage, minCompatibilityThreshold, result.CompatibleStories, result.TotalStories)
	}

	c.logger.LogInfo(fmt.Sprintf("Compatibility validation passed: %.1f%% >= %.1f%% threshold",
		result.CompatibilityPercentage, minCompatibilityThreshold))

	return nil
}

// GetIncompatibleMethods は、非互換メソッドのリストを返します
func (c *CompatibilityChecker) GetIncompatibleMethods(stories []testlib.Story) []string {
	result := c.CheckStoriesCompatibility(stories)
	var incompatibleMethods []string

	for _, storyResult := range result.StoryResults {
		for _, stepResult := range storyResult.StepResults {
			if !stepResult.IsCompatible {
				incompatibleMethods = append(incompatibleMethods, stepResult.MethodName)
			}
		}
	}

	return incompatibleMethods
}

// GetRequiredStrategies は、必要な戦略のリストを返します
func (c *CompatibilityChecker) GetRequiredStrategies(stories []testlib.Story) map[string]int {
	result := c.CheckStoriesCompatibility(stories)
	return result.Summary.StrategyUsage
}
