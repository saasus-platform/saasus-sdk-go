package main

import (
	"fmt"
	"os"
)

func main() {
	// CLI引数の解析
	cliConfig, err := ParseArgs()
	if err != nil {
		printError(err)
		ShowUsage()
		os.Exit(1)
	}

	// ヘルプ表示
	if cliConfig.ShowHelp {
		ShowUsage()
		os.Exit(0)
	}

	// 設定の作成
	config := &Config{
		ReportPath: cliConfig.ReportPath,
		OutputDir:  cliConfig.OutputDir,
		ModuleName: cliConfig.ModuleName,
		Overwrite:  cliConfig.Overwrite,
		Verbose:    cliConfig.Verbose,
	}

	// メイン処理の実行
	if err := run(config); err != nil {
		printError(err)
		os.Exit(1)
	}
}

// printError はエラーをユーザーフレンドリーな形式で出力します
func printError(err error) {
	fmt.Fprintf(os.Stderr, "\n❌ エラーが発生しました:\n")

	// エラーの種類に応じて適切なメッセージを表示
	switch {
	case IsFileNotFoundError(err):
		fmt.Fprintf(os.Stderr, "  ファイルが見つかりません\n")
		fmt.Fprintf(os.Stderr, "  %v\n", err)
		fmt.Fprintf(os.Stderr, "\n💡 ヒント: ファイルパスが正しいか確認してください\n")

	case IsInvalidJSONError(err):
		fmt.Fprintf(os.Stderr, "  JSONの形式が不正です\n")
		fmt.Fprintf(os.Stderr, "  %v\n", err)
		fmt.Fprintf(os.Stderr, "\n💡 ヒント: JSONファイルの構文を確認してください\n")

	case IsInvalidReportError(err):
		fmt.Fprintf(os.Stderr, "  レポートの構造が不正です\n")
		fmt.Fprintf(os.Stderr, "  %v\n", err)
		fmt.Fprintf(os.Stderr, "\n💡 ヒント: postman-analyzerで生成された正しいレポートファイルを使用してください\n")

	default:
		fmt.Fprintf(os.Stderr, "  %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "\n")
}

// run はメイン処理を実行します
func run(config *Config) error {
	// 詳細ログモードの場合、開始メッセージを表示
	if config.Verbose {
		fmt.Println("=== Test Param Generator ===")
		fmt.Printf("レポートファイル: %s\n", config.ReportPath)
		fmt.Printf("モジュール名: %s\n", config.ModuleName)
		fmt.Printf("出力ディレクトリ: %s\n", config.OutputDir)
		fmt.Println()
	}

	// 1. 設定の検証
	if config.Verbose {
		fmt.Println("ステップ 1/6: 設定を検証しています...")
	}
	if err := config.Validate(); err != nil {
		return WrapError(err, "設定の検証に失敗")
	}

	// 2. 出力ディレクトリの作成
	if config.Verbose {
		fmt.Println("ステップ 2/6: 出力ディレクトリを作成しています...")
	}
	if err := config.EnsureOutputDir(); err != nil {
		return WrapError(err, "出力ディレクトリの作成に失敗")
	}

	// 3. レポートの解析
	if config.Verbose {
		fmt.Println("ステップ 3/6: レポートファイルを解析しています...")
	}
	report, err := ParseReport(config.ReportPath)
	if err != nil {
		return WrapError(err, "レポートの解析に失敗")
	}
	if config.Verbose {
		fmt.Printf("  - コレクション名: %s\n", report.CollectionName)
		fmt.Printf("  - テストフロー数: %d\n", len(report.TestFlow))
		fmt.Printf("  - クライアントメソッド数: %d\n", len(report.ClientMethods))
	}

	// 4. レポートの分析
	if config.Verbose {
		fmt.Println("ステップ 4/7: レポートを分析しています...")
	}
	analysis, err := AnalyzeReport(report, config.ModuleName)
	if err != nil {
		return WrapError(err, "レポートの分析に失敗")
	}
	if config.Verbose {
		fmt.Printf("  - パラメータグループ数: %d\n", len(analysis.ParamGroups))
	}

	// 5. iterationDataのマージ
	if config.Verbose {
		fmt.Println("ステップ 5/7: iterationDataをマージしています...")
	}
	// レポート情報を含むGeneratorを作成（要件11.4）
	generator := NewGeneratorWithReport(config, analysis, report)
	if err := generator.mergeIterationData(); err != nil {
		return WrapError(err, "iterationDataのマージに失敗")
	}

	// 6. パラメータファイルの生成
	if config.Verbose {
		fmt.Println("ステップ 6/7: パラメータファイルを生成しています...")
	}
	if err := generator.GenerateParamsFile(); err != nil {
		return WrapError(err, "パラメータファイルの生成に失敗")
	}

	// 7. ヘルパー関数ファイルの生成
	if config.Verbose {
		fmt.Println("ステップ 7/7: ヘルパー関数ファイルを生成しています...")
	}
	if err := generator.GenerateHelpersFile(); err != nil {
		return WrapError(err, "ヘルパー関数ファイルの生成に失敗")
	}

	// 成功メッセージの表示
	paramsPath, helpersPath := config.GetOutputPaths()
	fmt.Println("\n✓ ファイルの生成が完了しました:")
	fmt.Printf("  - %s\n", paramsPath)
	fmt.Printf("  - %s\n", helpersPath)

	return nil
}
