package e2e

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

var (
	verbose = flag.Bool("verbose", false, "詳細ログを出力")
)

// TestMain はE2Eテストのエントリーポイントです
// テストの実行順序を制御し、セットアップとクリーンアップを管理します
func TestMain(m *testing.M) {
	flag.Parse()

	// 環境変数の確認
	if !checkEnvironmentVariables() {
		log.Println("必要な環境変数が設定されていません。E2Eテストをスキップします。")
		os.Exit(0)
	}

	// グローバルセットアップ
	if err := globalSetup(); err != nil {
		log.Fatalf("グローバルセットアップに失敗: %v", err)
	}

	// テスト実行
	code := m.Run()

	// グローバルクリーンアップ
	if err := globalCleanup(); err != nil {
		log.Printf("グローバルクリーンアップに失敗: %v", err)
	}

	os.Exit(code)
}

// checkEnvironmentVariables は必要な環境変数が設定されているかチェックします
func checkEnvironmentVariables() bool {
	required := []string{
		"SAASUS_SAAS_ID",
		"SAASUS_API_KEY", 
		"SAASUS_SECRET_KEY",
	}

	for _, env := range required {
		if os.Getenv(env) == "" {
			log.Printf("環境変数 %s が設定されていません", env)
			return false
		}
	}

	return true
}

// globalSetup はテスト実行前のグローバルセットアップを行います
func globalSetup() error {
	if *verbose {
		log.Println("=== E2Eテスト グローバルセットアップ開始 ===")
	}

	// 共通クライアントの初期化テスト
	client, err := common.NewAuthenticatedClient()
	if err != nil {
		return fmt.Errorf("認証付きクライアントの作成に失敗: %w", err)
	}

	// 接続テスト
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := common.TestConnection(ctx, client); err != nil {
		return fmt.Errorf("API接続テストに失敗: %w", err)
	}

	if *verbose {
		log.Println("API接続テスト成功")
		log.Println("=== E2Eテスト グローバルセットアップ完了 ===")
	}

	return nil
}

// globalCleanup はテスト実行後のグローバルクリーンアップを行います
func globalCleanup() error {
	if *verbose {
		log.Println("=== E2Eテスト グローバルクリーンアップ開始 ===")
	}

	// 共通クライアントの作成
	client, err := common.NewAuthenticatedClient()
	if err != nil {
		log.Printf("クリーンアップ用クライアントの作成に失敗: %v", err)
		return err
	}

	// テスト用リソースのクリーンアップ
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := common.CleanupTestResources(ctx, client); err != nil {
		log.Printf("テストリソースのクリーンアップに失敗: %v", err)
		return err
	}

	if *verbose {
		log.Println("=== E2Eテスト グローバルクリーンアップ完了 ===")
	}

	return nil
}

// TestE2EStoryExecution はストーリーテストの実行順序を管理します
func TestE2EStoryExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("ショートテストモードではE2Eテストをスキップ")
	}

	// ストーリーテストの実行順序
	// 依存関係を考慮して順次実行
	stories := []struct {
		name string
		test func(*testing.T)
	}{
		{"01_基本セットアップ", testStory01BasicSetup},
		{"02_ユーザー管理", testStory02UserManagement},
		{"03_ロール管理", testStory03RoleManagement},
		{"04_テナント管理", testStory04TenantManagement},
		{"05_テナントユーザー管理", testStory05TenantUserManagement},
		{"06_招待管理", testStory06InvitationManagement},
		{"07_認証フロー", testStory07AuthenticationFlow},
		{"08_シングルテナント管理", testStory08SingleTenantManagement},
		{"09_エラーハンドリング", testStory09ErrorHandling},
		{"10_統合テスト", testStory10IntegrationTest},
	}

	for _, story := range stories {
		t.Run(story.name, func(t *testing.T) {
			if *verbose {
				log.Printf("=== ストーリーテスト開始: %s ===", story.name)
			}

			// ストーリー開始前のセットアップ
			if err := storySetup(t, story.name); err != nil {
				t.Fatalf("ストーリーセットアップに失敗: %v", err)
			}

			// ストーリーテスト実行
			story.test(t)

			// ストーリー終了後のクリーンアップ
			if err := storyCleanup(t, story.name); err != nil {
				t.Errorf("ストーリークリーンアップに失敗: %v", err)
			}

			if *verbose {
				log.Printf("=== ストーリーテスト完了: %s ===", story.name)
			}
		})
	}
}

// storySetup は各ストーリーテスト開始前のセットアップを行います
func storySetup(t *testing.T, storyName string) error {
	if *verbose {
		log.Printf("ストーリーセットアップ: %s", storyName)
	}

	// ストーリー固有のセットアップロジックをここに追加
	// 例: 特定のテストデータの準備、前提条件の確認など

	return nil
}

// storyCleanup は各ストーリーテスト終了後のクリーンアップを行います
func storyCleanup(t *testing.T, storyName string) error {
	if *verbose {
		log.Printf("ストーリークリーンアップ: %s", storyName)
	}

	// ストーリー固有のクリーンアップロジックをここに追加
	// 例: ストーリーで作成したリソースの削除

	return nil
}

// testStory01BasicSetup は基本セットアップストーリーのプレースホルダーです
// 実際のテストは stories/01_basic_setup_test.go で実装されます
func testStory01BasicSetup(t *testing.T) {
	t.Log("基本セットアップストーリーテスト（プレースホルダー）")
	// 実際のテストロジックは stories/ ディレクトリで実装
}

// testStory03RoleManagement はロール管理ストーリーのプレースホルダーです
// 実際のテストは stories/03_role_management_test.go で実装されます
func testStory03RoleManagement(t *testing.T) {
	t.Log("ロール管理ストーリーテスト（プレースホルダー）")
	// 実際のテストロジックは stories/ ディレクトリで実装
}

// testStory02UserManagement はユーザー管理ストーリーのプレースホルダーです
// 実際のテストは stories/02_user_management_test.go で実装されます
func testStory02UserManagement(t *testing.T) {
	t.Log("ユーザー管理ストーリーテスト（プレースホルダー）")
	// 実際のテストロジックは stories/ ディレクトリで実装
}

// testStory04TenantManagement はテナント管理ストーリーのプレースホルダーです
// 実際のテストは stories/04_tenant_management_test.go で実装されます
func testStory04TenantManagement(t *testing.T) {
	t.Log("テナント管理ストーリーテスト（プレースホルダー）")
	// 実際のテストロジックは stories/ ディレクトリで実装
}

// testStory05TenantUserManagement はテナントユーザー管理ストーリーのプレースホルダーです
// 実際のテストは stories/05_tenant_user_management_test.go で実装されます
func testStory05TenantUserManagement(t *testing.T) {
	t.Log("テナントユーザー管理ストーリーテスト（プレースホルダー）")
	// 実際のテストロジックは stories/ ディレクトリで実装
}

// testStory06InvitationManagement は招待管理ストーリーのプレースホルダーです
// 実際のテストは stories/06_invitation_management_test.go で実装されます
func testStory06InvitationManagement(t *testing.T) {
	t.Log("招待管理ストーリーテスト（プレースホルダー）")
	// 実際のテストロジックは stories/ ディレクトリで実装
}

// testStory07AuthenticationFlow は認証フローストーリーのプレースホルダーです
// 実際のテストは stories/07_authentication_flow_test.go で実装されます
func testStory07AuthenticationFlow(t *testing.T) {
	t.Log("認証フローストーリーテスト（プレースホルダー）")
	// 実際のテストロジックは stories/ ディレクトリで実装
}

// testStory08SingleTenantManagement はシングルテナント管理ストーリーのプレースホルダーです
// 実際のテストは stories/08_single_tenant_management_test.go で実装されます
func testStory08SingleTenantManagement(t *testing.T) {
	t.Log("シングルテナント管理ストーリーテスト（プレースホルダー）")
	// 実際のテストロジックは stories/ ディレクトリで実装
}

// testStory09ErrorHandling はエラーハンドリングストーリーのプレースホルダーです
// 実際のテストは stories/09_error_handling_test.go で実装されます
func testStory09ErrorHandling(t *testing.T) {
	t.Log("エラーハンドリングストーリーテスト（プレースホルダー）")
	// 実際のテストロジックは stories/ ディレクトリで実装
}

// testStory10IntegrationTest は統合テストストーリーのプレースホルダーです
// 実際のテストは stories/10_integration_test.go で実装されます
func testStory10IntegrationTest(t *testing.T) {
	t.Log("統合テストストーリーテスト（プレースホルダー）")
	// 実際のテストロジックは stories/ ディレクトリで実装
}

// TestHealthCheck はE2Eテスト環境の基本的なヘルスチェックを行います
func TestHealthCheck(t *testing.T) {
	client, err := common.NewAuthenticatedClient()
	if err != nil {
		t.Fatalf("認証付きクライアントの作成に失敗: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := common.TestConnection(ctx, client); err != nil {
		t.Fatalf("API接続テストに失敗: %v", err)
	}

	t.Log("E2Eテスト環境のヘルスチェック成功")
}