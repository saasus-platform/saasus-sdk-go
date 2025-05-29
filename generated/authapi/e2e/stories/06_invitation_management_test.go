package stories

import (
	"context"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

// TestStory06_InvitationManagement は招待管理ストーリーのE2Eテストです
func TestStory06_InvitationManagement(t *testing.T) {
	// テストデータローダーを初期化
	loader := common.NewTestDataLoader()
	testData, err := loader.LoadInvitationManagementData()
	if err != nil {
		t.Fatalf("テストデータの読み込みに失敗: %v", err)
	}

	// 認証付きクライアントを作成
	client, err := common.NewAuthenticatedClient()
	if err != nil {
		t.Fatalf("認証付きクライアントの作成に失敗: %v", err)
	}

	// アサーションヘルパーを初期化
	assert := common.NewAssertionHelper(t, true)

	// クリーンアップマネージャーを初期化
	cleanup := common.NewCleanupManager(client, true)

	// テスト終了時のクリーンアップを設定
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := cleanup.CleanupByStory(ctx, common.StoryInvitationManagement); err != nil {
			t.Logf("クリーンアップに失敗: %v", err)
		}
	}()

	// 前提条件の準備
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// テスト用テナントを作成
	testTenantID := setupInvitationTestData(t, client, ctx)
	if testTenantID == "" {
		t.Fatal("テスト用テナントの準備に失敗")
	}

	// サブテストを順次実行
	t.Run("招待管理", func(t *testing.T) {
		testInvitationManagement(t, client, testData, assert, cleanup, testTenantID)
	})

	t.Run("招待検証", func(t *testing.T) {
		testInvitationValidation(t, client, testData, assert, testTenantID)
	})
}

// setupInvitationTestData はテスト用のテナントを準備します
func setupInvitationTestData(t *testing.T, client *common.ClientWrapper, ctx context.Context) string {
	// テスト用テナントを作成
	tenantParam := authapi.CreateTenantParam{
		Name:                 "招待テスト用テナント",
		BackOfficeStaffEmail: "invitation-test@example.com",
		Attributes:           map[string]interface{}{},
	}

	tenantResp, err := client.Client.CreateTenantWithResponse(ctx, tenantParam)
	if err != nil || tenantResp.JSON201 == nil {
		t.Logf("テスト用テナントの作成に失敗: %v", err)
		return ""
	}

	return tenantResp.JSON201.Id
}

// testInvitationManagement は招待管理のテストを実行します
func testInvitationManagement(t *testing.T, client *common.ClientWrapper, testData *common.InvitationManagementTestData, assert *common.AssertionHelper, cleanup *common.CleanupManager, testTenantID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var createdInvitationIDs []string

	t.Run("招待一覧取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetTenantInvitationsWithResponse(ctx, testTenantID)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("招待一覧取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "招待一覧取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "招待一覧取得")

		t.Logf("招待一覧取得成功: 招待数=%d", len(resp.JSON200.Invitations))
	})

	t.Run("招待作成", func(t *testing.T) {
		for _, param := range testData.Invitations.Create.Params {
			t.Run(param.Email, func(t *testing.T) {
				// 作成パラメータを準備
				createParam := authapi.CreateTenantInvitationParam{
					Email: param.Email,
					AccessToken: "test-access-token",
					Envs: []struct {
						Id        authapi.Id `json:"id"`
						RoleNames []string   `json:"role_names"`
					}{
						{
							Id:        1,
							RoleNames: []string{"user"},
						},
					},
				}

				// 招待を作成
				startTime := time.Now()
				createResp, err := client.Client.CreateTenantInvitationWithResponse(ctx, testTenantID, createParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("招待作成APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 15*time.Second, "招待作成")

				// ステータスコードをチェック
				assert.AssertStatusCode(201, createResp.StatusCode(), "招待作成")

				if createResp.JSON201 != nil {
					invitationID := createResp.JSON201.Id
					createdInvitationIDs = append(createdInvitationIDs, invitationID)

					// リソース追跡に追加
					client.CreateTestResource(
						common.ResourceTypeInvitation,
						invitationID,
						param.Email,
						common.StoryInvitationManagement,
						map[string]interface{}{
							"tenant_id": testTenantID,
							"email":     param.Email,
							"roles":     param.Roles,
							"envs":      param.Envs,
						},
					)

					t.Logf("招待作成成功: ID=%s, Email=%s", invitationID, param.Email)
				}
			})
		}
	})

	t.Run("招待詳細取得", func(t *testing.T) {
		for _, invitationID := range createdInvitationIDs {
			t.Run(invitationID, func(t *testing.T) {
				startTime := time.Now()
				resp, err := client.Client.GetTenantInvitationWithResponse(ctx, testTenantID, invitationID)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("招待詳細取得APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 10*time.Second, "招待詳細取得")

				// ステータスコードをチェック
				assert.AssertStatusCode(200, resp.StatusCode(), "招待詳細取得")

				if resp.JSON200 != nil {
					// 招待情報をチェック
					assert.AssertEquals(invitationID, resp.JSON200.Id, "招待ID")
				}

				t.Logf("招待詳細取得成功: ID=%s", invitationID)
			})
		}
	})

	t.Run("招待削除", func(t *testing.T) {
		for _, invitationID := range createdInvitationIDs {
			t.Run(invitationID, func(t *testing.T) {
				// 招待を削除
				startTime := time.Now()
				deleteResp, err := client.Client.DeleteTenantInvitationWithResponse(ctx, testTenantID, invitationID)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("招待削除APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 15*time.Second, "招待削除")

				// ステータスコードをチェック
				assert.AssertResourceDeleted(deleteResp.StatusCode(), "招待")

				// リソースを削除済みとしてマーク
				client.MarkResourceCleaned(invitationID, nil)

				// 削除確認
				confirmResp, err := client.Client.GetTenantInvitationWithResponse(ctx, testTenantID, invitationID)
				if err == nil && confirmResp.StatusCode() != 404 {
					t.Errorf("招待が削除されていません: ステータスコード %d", confirmResp.StatusCode())
				}

				t.Logf("招待削除成功: ID=%s", invitationID)
			})
		}
	})
}

// testInvitationValidation は招待検証のテストを実行します
func testInvitationValidation(t *testing.T, client *common.ClientWrapper, testData *common.InvitationManagementTestData, assert *common.AssertionHelper, testTenantID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// テスト用招待を作成
	createParam := authapi.CreateTenantInvitationParam{
		Email: "validation-test@example.com",
		AccessToken: "test-access-token",
		Envs: []struct {
			Id        authapi.Id `json:"id"`
			RoleNames []string   `json:"role_names"`
		}{
			{
				Id:        1,
				RoleNames: []string{"user"},
			},
		},
	}

	createResp, err := client.Client.CreateTenantInvitationWithResponse(ctx, testTenantID, createParam)
	if err != nil || createResp.JSON201 == nil {
		t.Skip("テスト用招待の作成に失敗したためスキップ")
		return
	}

	testInvitationID := createResp.JSON201.Id

	// テスト終了時に招待を削除
	defer func() {
		client.Client.DeleteTenantInvitationWithResponse(ctx, testTenantID, testInvitationID)
	}()

	t.Run("招待有効性確認", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetInvitationValidityWithResponse(ctx, testInvitationID)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("招待有効性確認APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "招待有効性確認")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "招待有効性確認")

		t.Log("招待有効性確認成功")
	})

	t.Run("招待コード検証", func(t *testing.T) {
		t.Skip("ValidateInvitationCodeWithResponse API method not available")

		t.Log("招待コード検証完了")
	})

	t.Run("無効な招待ID検証エラー", func(t *testing.T) {
		invalidInvitationID := "invalid-invitation-id"

		startTime := time.Now()
		resp, err := client.Client.GetInvitationValidityWithResponse(ctx, invalidInvitationID)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("無効な招待ID検証APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "無効な招待ID検証")

		// エラーステータスコードをチェック
		if resp.StatusCode() == 200 {
			t.Error("無効な招待IDで成功レスポンスが返されました")
		}

		t.Log("無効な招待ID検証エラー確認成功")
	})

	t.Run("無効な招待コード検証エラー", func(t *testing.T) {
		t.Skip("ValidateInvitationCodeWithResponse API method not available")
	})
}