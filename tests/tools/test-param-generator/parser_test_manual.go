package main

import (
	"fmt"
)

// manualTestParser は手動でパーサーをテストするための関数
func manualTestParser() {
	fmt.Println("=== Parser Manual Test ===")
	fmt.Println()

	// テストケース1: 正常なレポート
	testCase1 := "../postman-analyzer/reports/communication/postman_analysis.json"
	fmt.Printf("テストケース1: 正常なレポート\n")
	fmt.Printf("ファイル: %s\n", testCase1)

	report, err := ParseReport(testCase1)
	if err != nil {
		fmt.Printf("✗ エラー: %v\n\n", err)
	} else {
		fmt.Printf("✓ 解析成功\n")
		fmt.Printf("  Collection Name: %s\n", report.CollectionName)
		fmt.Printf("  Test Flow Steps: %d\n", len(report.TestFlow))
		fmt.Printf("  Client Methods: %d\n\n", len(report.ClientMethods))
	}

	// テストケース2: 存在しないファイル
	testCase2 := "nonexistent.json"
	fmt.Printf("テストケース2: 存在しないファイル\n")
	fmt.Printf("ファイル: %s\n", testCase2)

	_, err = ParseReport(testCase2)
	if err != nil {
		fmt.Printf("✓ 期待通りエラー: %v\n\n", err)
	} else {
		fmt.Printf("✗ エラーが発生すべきでした\n\n")
	}

	// テストケース3: 不正なJSON
	testCase3 := "testdata/invalid.json"
	fmt.Printf("テストケース3: 不正なJSON\n")
	fmt.Printf("ファイル: %s\n", testCase3)

	_, err = ParseReport(testCase3)
	if err != nil {
		fmt.Printf("✓ 期待通りエラー: %v\n\n", err)
	} else {
		fmt.Printf("✗ エラーが発生すべきでした\n\n")
	}

	// テストケース4: 空のレポート
	testCase4 := "testdata/empty_report.json"
	fmt.Printf("テストケース4: 空のレポート\n")
	fmt.Printf("ファイル: %s\n", testCase4)

	_, err = ParseReport(testCase4)
	if err != nil {
		fmt.Printf("✓ 期待通りエラー: %v\n\n", err)
	} else {
		fmt.Printf("✗ エラーが発生すべきでした\n\n")
	}

	// テストケース5: 必須フィールド欠落
	testCase5 := "testdata/missing_fields.json"
	fmt.Printf("テストケース5: 必須フィールド欠落\n")
	fmt.Printf("ファイル: %s\n", testCase5)

	_, err = ParseReport(testCase5)
	if err != nil {
		fmt.Printf("✓ 期待通りエラー: %v\n\n", err)
	} else {
		fmt.Printf("✗ エラーが発生すべきでした\n\n")
	}

	fmt.Println("=== テスト完了 ===")
}
