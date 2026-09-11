package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// DiagnosticReporter は診断情報を出力します
type DiagnosticReporter struct {
	verbose    bool
	metrics    *EnhancedMappingMetrics
	moduleName string
}

// NewDiagnosticReporter は新しいDiagnosticReporterを作成します
func NewDiagnosticReporter(moduleName string, verbose bool) *DiagnosticReporter {
	return &DiagnosticReporter{
		verbose:    verbose,
		moduleName: moduleName,
		metrics: &EnhancedMappingMetrics{
			TypePatternStats:   make(map[string]int),
			MatchStrategyStats: make(map[string]int),
			FailedMappings:     []DetailedFailedMapping{},
		},
	}
}

// RecordSuccess はマッピング成功を記録します
func (dr *DiagnosticReporter) RecordSuccess(groupName string, fieldCount int) {
	dr.metrics.SuccessfulGroups++
	dr.metrics.MappedFields += fieldCount
}

// RecordFailure はマッピング失敗を記録します
func (dr *DiagnosticReporter) RecordFailure(failed DetailedFailedMapping) {
	dr.metrics.FailedMappings = append(dr.metrics.FailedMappings, failed)
}

// RecordTypePattern は型パターンの使用を記録します
func (dr *DiagnosticReporter) RecordTypePattern(pattern string) {
	dr.metrics.TypePatternStats[pattern]++
}

// RecordMatchStrategy はマッチング戦略の使用を記録します
func (dr *DiagnosticReporter) RecordMatchStrategy(strategy string) {
	dr.metrics.MatchStrategyStats[strategy]++
}

// UpdateTotals は総数を更新します
func (dr *DiagnosticReporter) UpdateTotals(totalGroups, totalFields int) {
	dr.metrics.TotalGroups = totalGroups
	dr.metrics.TotalFields = totalFields
}

// GetMetrics はメトリクスを取得します
func (dr *DiagnosticReporter) GetMetrics() *EnhancedMappingMetrics {
	return dr.metrics
}

// ReportMappingResults はマッピング結果を報告します
func (dr *DiagnosticReporter) ReportMappingResults() {
	fmt.Println("\n=== マッピング診断レポート ===")
	fmt.Printf("\nモジュール: %s\n", dr.moduleName)
	fmt.Printf("総パラメータグループ: %d\n", dr.metrics.TotalGroups)

	successRate := 0.0
	if dr.metrics.TotalGroups > 0 {
		successRate = float64(dr.metrics.SuccessfulGroups) / float64(dr.metrics.TotalGroups) * 100
	}
	fmt.Printf("成功したグループ: %d (%.1f%%)\n", dr.metrics.SuccessfulGroups, successRate)

	fmt.Printf("総フィールド: %d\n", dr.metrics.TotalFields)

	fieldSuccessRate := 0.0
	if dr.metrics.TotalFields > 0 {
		fieldSuccessRate = float64(dr.metrics.MappedFields) / float64(dr.metrics.TotalFields) * 100
	}
	fmt.Printf("マッピング成功: %d (%.1f%%)\n", dr.metrics.MappedFields, fieldSuccessRate)

	// 失敗したマッピングの詳細
	if len(dr.metrics.FailedMappings) > 0 {
		dr.reportFailedMappings()
	}

	// パターン使用統計
	if len(dr.metrics.TypePatternStats) > 0 {
		dr.reportPatternStats()
	}

	// マッチング戦略統計
	if len(dr.metrics.MatchStrategyStats) > 0 {
		dr.reportStrategyStats()
	}

	// 改善提案
	suggestions := dr.GenerateSuggestions()
	if len(suggestions) > 0 {
		dr.reportSuggestions(suggestions)
	}

	fmt.Println("\n================================")
}

// reportFailedMappings は失敗したマッピングを報告します
func (dr *DiagnosticReporter) reportFailedMappings() {
	fmt.Println("\n--- 失敗したマッピング ---")

	// 最初の10件のみ表示（詳細モードでは全件）
	displayCount := len(dr.metrics.FailedMappings)
	if !dr.verbose && displayCount > 10 {
		displayCount = 10
	}

	for i := 0; i < displayCount; i++ {
		failed := dr.metrics.FailedMappings[i]
		fmt.Printf("\n%d. モジュール: %s, グループ: %s, 操作: %s, フィールド: %s\n",
			i+1, failed.ModuleName, failed.GroupName, failed.Operation, failed.FieldName)

		if len(failed.TriedFilePaths) > 0 {
			fmt.Println("   試行したファイルパス:")
			for _, path := range failed.TriedFilePaths {
				fmt.Printf("     - %s\n", path)
			}
		}

		if len(failed.TriedPatterns) > 0 {
			fmt.Println("   試行したパターン:")
			displayPatterns := len(failed.TriedPatterns)
			if displayPatterns > 5 {
				displayPatterns = 5
			}
			for j := 0; j < displayPatterns; j++ {
				pattern := failed.TriedPatterns[j]
				fmt.Printf("     - %s (カテゴリ: %s, 優先度: %d)\n",
					pattern.Pattern, pattern.Category, pattern.Priority)
			}
			if len(failed.TriedPatterns) > 5 {
				fmt.Printf("     ... 他 %d パターン\n", len(failed.TriedPatterns)-5)
			}
		}

		if len(failed.TriedStrategies) > 0 {
			fmt.Println("   試行したマッチング戦略:")
			for _, strategy := range failed.TriedStrategies {
				fmt.Printf("     - %s\n", strategy)
			}
		}

		if failed.BestMatch != nil {
			fmt.Printf("   最も近かったマッチ: %s (信頼度: %.2f)\n",
				failed.BestMatch.Field.Name, failed.BestMatch.Confidence)
		}

		if failed.Reason != "" {
			fmt.Printf("   理由: %s\n", failed.Reason)
		}

		if failed.Suggestion != "" {
			fmt.Printf("   提案: %s\n", failed.Suggestion)
		}
	}

	if !dr.verbose && len(dr.metrics.FailedMappings) > 10 {
		fmt.Printf("\n... 他 %d 件の失敗（詳細は -verbose フラグで確認）\n",
			len(dr.metrics.FailedMappings)-10)
	}
}

// reportPatternStats はパターン使用統計を報告します
func (dr *DiagnosticReporter) reportPatternStats() {
	fmt.Println("\n--- パターン使用統計 ---")

	// 使用回数でソート
	type patternStat struct {
		pattern string
		count   int
	}
	stats := make([]patternStat, 0, len(dr.metrics.TypePatternStats))
	for pattern, count := range dr.metrics.TypePatternStats {
		stats = append(stats, patternStat{pattern, count})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].count > stats[j].count
	})

	// 上位10件を表示
	displayCount := len(stats)
	if displayCount > 10 {
		displayCount = 10
	}

	for i := 0; i < displayCount; i++ {
		fmt.Printf("%s: %d回\n", stats[i].pattern, stats[i].count)
	}
}

// reportStrategyStats はマッチング戦略統計を報告します
func (dr *DiagnosticReporter) reportStrategyStats() {
	fmt.Println("\n--- マッチング戦略統計 ---")

	// 使用回数でソート
	type strategyStat struct {
		strategy string
		count    int
	}
	stats := make([]strategyStat, 0, len(dr.metrics.MatchStrategyStats))
	for strategy, count := range dr.metrics.MatchStrategyStats {
		stats = append(stats, strategyStat{strategy, count})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].count > stats[j].count
	})

	for _, stat := range stats {
		fmt.Printf("%s: %d回\n", stat.strategy, stat.count)
	}
}

// reportSuggestions は改善提案を報告します
func (dr *DiagnosticReporter) reportSuggestions(suggestions []string) {
	fmt.Println("\n--- 改善提案 ---")
	for i, suggestion := range suggestions {
		fmt.Printf("%d. %s\n", i+1, suggestion)
	}
}

// ReportFailedMapping は失敗したマッピングを報告します（個別）
func (dr *DiagnosticReporter) ReportFailedMapping(failed DetailedFailedMapping) {
	if dr.verbose {
		fmt.Printf("[VERBOSE] マッピング失敗: グループ=%s, 操作=%s, フィールド=%s, 理由=%s\n",
			failed.GroupName, failed.Operation, failed.FieldName, failed.Reason)
	}
	dr.RecordFailure(failed)
}

// GenerateSuggestions は改善提案を生成します
func (dr *DiagnosticReporter) GenerateSuggestions() []string {
	suggestions := []string{}

	// マッピング成功率が低い場合
	if dr.metrics.TotalGroups > 0 {
		successRate := float64(dr.metrics.SuccessfulGroups) / float64(dr.metrics.TotalGroups) * 100
		if successRate < 50 {
			suggestions = append(suggestions,
				"マッピング成功率が50%未満です。リソース名の正規化ルールを見直してください")
		}
	}

	// 型パターンが見つからない失敗が多い場合
	typeNotFoundCount := 0
	for _, failed := range dr.metrics.FailedMappings {
		if strings.Contains(failed.Reason, "型名パターン") || strings.Contains(failed.Reason, "型が見つかりません") {
			typeNotFoundCount++
		}
	}
	if typeNotFoundCount > 5 {
		suggestions = append(suggestions,
			fmt.Sprintf("%d件の型名パターンマッチング失敗があります。types.gen.goに不足している型定義を追加してください", typeNotFoundCount))
	}

	// フィールドマッチングの失敗が多い場合
	fieldNotFoundCount := 0
	for _, failed := range dr.metrics.FailedMappings {
		if strings.Contains(failed.Reason, "フィールド") {
			fieldNotFoundCount++
		}
	}
	if fieldNotFoundCount > 5 {
		suggestions = append(suggestions,
			fmt.Sprintf("%d件のフィールドマッチング失敗があります。フィールド名の別名ルールを追加してください", fieldNotFoundCount))
	}

	// カスタムマッピングルールの提案
	if len(dr.metrics.FailedMappings) > 10 {
		suggestions = append(suggestions,
			"多数のマッピング失敗があります。カスタムマッピングルールの追加を検討してください")
	}

	// 特定のパターンが全く使用されていない場合
	if len(dr.metrics.TypePatternStats) == 0 && dr.metrics.TotalGroups > 0 {
		suggestions = append(suggestions,
			"型パターンが全く使用されていません。パターン生成ロジックを確認してください")
	}

	return suggestions
}

// ExportDiagnosticReport は診断レポートをファイルに出力します
func (dr *DiagnosticReporter) ExportDiagnosticReport(outputPath string) error {
	// JSON形式でエクスポート
	data := map[string]any{
		"module_name":          dr.moduleName,
		"total_groups":         dr.metrics.TotalGroups,
		"successful_groups":    dr.metrics.SuccessfulGroups,
		"total_fields":         dr.metrics.TotalFields,
		"mapped_fields":        dr.metrics.MappedFields,
		"failed_mappings":      dr.metrics.FailedMappings,
		"type_pattern_stats":   dr.metrics.TypePatternStats,
		"match_strategy_stats": dr.metrics.MatchStrategyStats,
		"suggestions":          dr.GenerateSuggestions(),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("診断レポートのJSON変換に失敗しました: %w", err)
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("診断レポートの書き込みに失敗しました: %w", err)
	}

	fmt.Printf("\n診断レポートを出力しました: %s\n", outputPath)
	return nil
}

// GetModuleStats はモジュール別の統計を取得します
// 要件: 4.1, 4.2, 4.3
func (dr *DiagnosticReporter) GetModuleStats() *ModuleMappingStats {
	successRate := 0.0
	if dr.metrics.TotalFields > 0 {
		successRate = float64(dr.metrics.MappedFields) / float64(dr.metrics.TotalFields) * 100
	}

	return &ModuleMappingStats{
		ModuleName:       dr.moduleName,
		TotalGroups:      dr.metrics.TotalGroups,
		SuccessfulGroups: dr.metrics.SuccessfulGroups,
		TotalFields:      dr.metrics.TotalFields,
		MappedFields:     dr.metrics.MappedFields,
		SuccessRate:      successRate,
		FailedMappings:   dr.metrics.FailedMappings,
	}
}

// ReportByModule はモジュール別のレポートを出力します
// 要件: 4.1, 4.2, 4.3
func (dr *DiagnosticReporter) ReportByModule() {
	stats := dr.GetModuleStats()

	fmt.Println("\n=== モジュール別マッピングレポート ===")
	fmt.Printf("\nモジュール: %s\n", stats.ModuleName)
	fmt.Printf("総パラメータグループ: %d\n", stats.TotalGroups)
	fmt.Printf("成功したグループ: %d (%.1f%%)\n", stats.SuccessfulGroups,
		float64(stats.SuccessfulGroups)/float64(stats.TotalGroups)*100)
	fmt.Printf("総フィールド: %d\n", stats.TotalFields)
	fmt.Printf("マッピング成功: %d (%.1f%%)\n", stats.MappedFields, stats.SuccessRate)

	// 失敗パターンの集計
	if len(stats.FailedMappings) > 0 {
		fmt.Println("\n--- 失敗パターンの集計 ---")

		// 失敗理由別の集計
		reasonCounts := make(map[string]int)
		for _, failed := range stats.FailedMappings {
			reasonCounts[failed.Reason]++
		}

		fmt.Println("失敗理由別:")
		for reason, count := range reasonCounts {
			fmt.Printf("  %s: %d件\n", reason, count)
		}

		// グループ別の失敗数
		groupCounts := make(map[string]int)
		for _, failed := range stats.FailedMappings {
			groupCounts[failed.GroupName]++
		}

		fmt.Println("\nグループ別失敗数:")
		type groupCount struct {
			name  string
			count int
		}
		groups := make([]groupCount, 0, len(groupCounts))
		for name, count := range groupCounts {
			groups = append(groups, groupCount{name, count})
		}
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].count > groups[j].count
		})

		displayCount := len(groups)
		if displayCount > 10 {
			displayCount = 10
		}
		for i := 0; i < displayCount; i++ {
			fmt.Printf("  %s: %d件\n", groups[i].name, groups[i].count)
		}
		if len(groups) > 10 {
			fmt.Printf("  ... 他 %d グループ\n", len(groups)-10)
		}
	}

	fmt.Println("\n====================================")
}

// ExportModuleReport はモジュール別レポートをファイルに出力します
// 要件: 4.1, 4.2, 4.3
func (dr *DiagnosticReporter) ExportModuleReport(outputPath string) error {
	stats := dr.GetModuleStats()

	// JSON形式でエクスポート
	jsonData, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("モジュール別レポートのJSON変換に失敗しました: %w", err)
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("モジュール別レポートの書き込みに失敗しました: %w", err)
	}

	fmt.Printf("\nモジュール別レポートを出力しました: %s\n", outputPath)
	return nil
}
