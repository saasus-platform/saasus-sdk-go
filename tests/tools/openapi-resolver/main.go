package main

import (
	"fmt"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	// CLIハンドラーを作成
	handler := NewCLIHandler()

	// コマンドラインフラグをパース
	if err := handler.ParseFlags(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n\n", err)
		handler.PrintUsage()
		return 1
	}

	// CLIコマンドを実行
	if err := handler.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		return 1
	}

	return 0
}
