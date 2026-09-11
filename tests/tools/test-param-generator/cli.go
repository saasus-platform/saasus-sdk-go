package main

import (
	"flag"
	"fmt"
)

// ParseArgs はコマンドライン引数を解析します
func ParseArgs() (*CLIConfig, error) {
	config := &CLIConfig{}

	flag.StringVar(&config.ReportPath, "report", "", "Postman Analyzerレポートファイルのパス（必須）")
	flag.StringVar(&config.OutputDir, "output", "", "出力ディレクトリ（デフォルト: tests/e2e/{module}/）")
	flag.StringVar(&config.ModuleName, "module", "", "モジュール名（必須）")
	flag.BoolVar(&config.Overwrite, "overwrite", false, "既存ファイルを確認なしで上書き")
	flag.BoolVar(&config.Verbose, "verbose", false, "詳細ログを出力")
	flag.BoolVar(&config.ShowHelp, "help", false, "使用方法を表示")

	flag.Parse()

	// ヘルプフラグが指定された場合
	if config.ShowHelp {
		return config, nil
	}

	// 必須パラメータの検証
	if config.ReportPath == "" {
		return nil, NewInvalidConfigError("report", "必須パラメータが不足しています: -report を指定してください")
	}

	if config.ModuleName == "" {
		return nil, NewInvalidConfigError("module", "必須パラメータが不足しています: -module を指定してください")
	}

	// 出力ディレクトリのデフォルト値設定
	if config.OutputDir == "" {
		config.OutputDir = defaultOutputDir(config.ModuleName)
	}

	return config, nil
}

// ShowUsage は使用方法を表示します
func ShowUsage() {
	fmt.Println("Test Param Generator - Postman Analyzerレポートからテストパラメータを生成")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  test-param-generator -report <レポートファイル> -module <モジュール名> [オプション]")
	fmt.Println()
	fmt.Println("必須パラメータ:")
	fmt.Println("  -report string")
	fmt.Println("        Postman Analyzerレポートファイルのパス")
	fmt.Println("  -module string")
	fmt.Println("        モジュール名（例: communication, billing, auth）")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -output string")
	fmt.Println("        出力ディレクトリ（デフォルト: tests/e2e/{module}/）")
	fmt.Println("  -overwrite")
	fmt.Println("        既存ファイルを確認なしで上書き")
	fmt.Println("  -verbose")
	fmt.Println("        詳細ログを出力")
	fmt.Println("  -help")
	fmt.Println("        このヘルプメッセージを表示")
	fmt.Println()
	fmt.Println("使用例:")
	fmt.Println("  # 基本的な使用方法")
	fmt.Println("  test-param-generator -report reports/communication/postman_analysis.json -module communication")
	fmt.Println()
	fmt.Println("  # 出力先を指定")
	fmt.Println("  test-param-generator -report reports/billing/postman_analysis.json -module billing -output tests/e2e/custom-billing")
	fmt.Println()
	fmt.Println("  # 上書き確認なし")
	fmt.Println("  test-param-generator -report reports/auth/postman_analysis.json -module auth -overwrite")
	fmt.Println()
	fmt.Println("  # 詳細ログ付き")
	fmt.Println("  test-param-generator -report reports/pricing/postman_analysis.json -module pricing -verbose")
}
