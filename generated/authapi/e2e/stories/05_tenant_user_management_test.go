package stories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

// TestStory05_TenantUserManagement はテナントユーザー管理ストーリーのE2Eテストです
func TestStory05_TenantUserManagement(t *testing.T) {
	// テストデータローダーを初期化
	loader := common.NewTestDataLoader()
	testData, err := loader.LoadTenantUserManagementData()
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
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		if err := cleanup.CleanupByStory(ctx, common.StoryTenantUserManagement); err != nil {
			t.Logf("クリーンアップに失敗: %v", err)
		}
	}()

	// 前提条件の準備
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// テスト用テナントとユーザーを作成
	testTenantID, testUserIDs := setupTenantUserTestData(t, client, ctx)
	if testTenantID == "" || len(testUserIDs) == 0 {
		t.Fatal("テスト用データの準備に失敗")
	}

	// サブテストを順次実行
	t.Run("テナントユーザー管理", func(t *testing.T) {
		testTenantUserManagement(t, client, testData, assert, cleanup, testTenantID, testUserIDs)
	})

	t.Run("全テナントユーザー管理", func(t *testing.T) {
		testAllTenantUserManagement(t, client, testData, assert, cleanup)
	})

	t.Run("役割管理", func(t *testing.T) {
		testRoleManagement(t, client, testData, assert, testTenantID, testUserIDs)
	})
}

// setupTenantUserTestData はテスト用のテナントとユーザーを準備します
func setupTenantUserTestData(t *testing.T, client *common.ClientWrapper, ctx context.Context) (string, []string) {
	// テスト用テナントを作成
	tenantParam := authapi.CreateTenantParam{
		Name:                 "テナントユーザーテスト用テナント",
		BackOfficeStaffEmail: "tenant-user-test@example.com",
		Attributes:           &map[string]interface{}{},
	}

	tenantResp, err := client.Client.CreateTenantWithResponse(ctx, tenantParam)
	if err != nil || tenantResp.JSON201 == nil {
		t.Logf("テスト用テナントの作成に失敗: %v", err)
		return "", nil
	}

	testTenantID := tenantResp.JSON201.Id

	// テスト用ユーザーを作成
	var testUserIDs []string
	for i := 0; i < 2; i++ {
		userParam := authapi.CreateSaasUserParam{
			Email: fmt.Sprintf("tenant-user-test-%d@example.com", i+1),
			Attributes: &map[string]interface{}{
				"test": "value",
			},
		}

		userResp, err := client.Client.CreateSaasUserWithResponse(ctx, userParam)
		if err != nil || userResp.JSON201 == nil {
			t.Logf("テスト用ユーザー%dの作成に失敗: %v", i+1, err)
			continue
		}

		testUserIDs = append(testUserIDs, userResp.JSON201.Id)
	}

	return testTenantID, testUserIDs
}

// testTenantUserManagement はテナントユーザー管理のテストを実行します
func testTenantUserManagement(t *testing.T, client *common.ClientWrapper, testData *common.TenantUserManagementTestData, assert *common.AssertionHelper, cleanup *common.CleanupManager, testTenantID string, testUserIDs []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var addedUserIDs []string

	t.Run("テナントユーザー一覧取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetTenantUsersWithResponse(ctx, testTenantID)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("テナントユーザー一覧取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "テナントユーザー一覧取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "テナントユーザー一覧取得")

		t.Logf("テナントユーザー一覧取得成功: ユーザー数=%d", len(resp.JSON200.Users))
	})

	t.Run("テナントユーザー追加", func(t *testing.T) {
		for i, userID := range testUserIDs {
			if i >= len(testData.TenantUsers.Create.Params) {
				break
			}

			param := testData.TenantUsers.Create.Params[i]
			t.Run(userID, func(t *testing.T) {
				// 追加パラメータを準備
				createParam := authapi.CreateTenantUserParam{
					UserId: userID,
					Roles:  &param.Roles,
					Envs:   &param.Envs,
				}

				// テナントユーザーを追加
				startTime := time.Now()
				createResp, err := client.Client.CreateTenantUserWithResponse(ctx, testTenantID, createParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("テナントユーザー追加APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 20*time.Second, "テナントユーザー追加")

				// ステータスコードをチェック
				assert.AssertStatusCode(201, createResp.StatusCode(), "テナントユーザー追加")

				addedUserIDs = append(addedUserIDs, userID)

				// リソース追跡に追加
				client.CreateTestResource(
					common.ResourceTypeTenantUser,
					fmt.Sprintf("%s:%s", testTenantID, userID),
					fmt.Sprintf("TenantUser-%s", userID),
					common.StoryTenantUserManagement,
					map[string]interface{}{
						"tenant_id": testTenantID,
						"user_id":   userID,
						"roles":     param.Roles,
						"envs":      param.Envs,
					},
				)

				t.Logf("テナントユーザー追加成功: TenantID=%s, UserID=%s", testTenantID, userID)
			})
		}
	})

	t.Run("テナントユーザー詳細取得", func(t *testing.T) {
		for _, userID := range addedUserIDs {
			t.Run(userID, func(t *testing.T) {
				startTime := time.Now()
				resp, err := client.Client.GetTenantUserWithResponse(ctx, testTenantID, userID)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("テナントユーザー詳細取得APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 10*time.Second, "テナントユーザー詳細取得")

				// ステータスコードをチェック
				assert.AssertStatusCode(200, resp.StatusCode(), "テナントユーザー詳細取得")

				if resp.JSON200 != nil {
					// ユーザー情報をチェック
					assert.AssertEquals(userID, resp.JSON200.Id, "ユーザーID")
				}

				t.Logf("テナントユーザー詳細取得成功: TenantID=%s, UserID=%s", testTenantID, userID)
			})
		}
	})

	t.Run("テナントユーザー情報更新", func(t *testing.T) {
		if len(addedUserIDs) > 0 {
			userID := addedUserIDs[0]
			t.Run(userID, func(t *testing.T) {
				// 更新パラメータを準備
				updateParam := authapi.UpdateTenantUserParam{
					Roles: testData.TenantUsers.Update.Params.Roles,
					Envs:  testData.TenantUsers.Update.Params.Envs,
				}

				// テナントユーザー情報を更新
				startTime := time.Now()
				updateResp, err := client.Client.UpdateTenantUserWithResponse(ctx, testTenantID, userID, updateParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("テナントユーザー情報更新APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 15*time.Second, "テナントユーザー情報更新")

				// ステータスコードをチェック
				assert.AssertStatusCode(200, updateResp.StatusCode(), "テナントユーザー情報更新")

				t.Logf("テナントユーザー情報更新成功: TenantID=%s, UserID=%s", testTenantID, userID)
			})
		}
	})

	t.Run("テナントユーザー削除", func(t *testing.T) {
		for _, userID := range addedUserIDs {
			t.Run(userID, func(t *testing.T) {
				// テナントユーザーを削除
				startTime := time.Now()
				deleteResp, err := client.Client.DeleteTenantUserWithResponse(ctx, testTenantID, userID)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("テナントユーザー削除APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 20*time.Second, "テナントユーザー削除")

				// ステータスコードをチェック
				assert.AssertResourceDeleted(deleteResp.StatusCode(), "テナントユーザー")

				// リソースを削除済みとしてマーク
				resourceID := fmt.Sprintf("%s:%s", testTenantID, userID)
				client.MarkResourceCleaned(resourceID, nil)

				// 削除確認
				confirmResp, err := client.Client.GetTenantUserWithResponse(ctx, testTenantID, userID)
				if err == nil && confirmResp.StatusCode() != 404 {
					t.Errorf("テナントユーザーが削除されていません: ステータスコード %d", confirmResp.StatusCode())
				}

				t.Logf("テナントユーザー削除成功: TenantID=%s, UserID=%s", testTenantID, userID)
			})
		}
	})
}

// testAllTenantUserManagement は全テナントユーザー管理のテストを実行します
func testAllTenantUserManagement(t *testing.T, client *common.ClientWrapper, testData *common.TenantUserManagementTestData, assert *common.AssertionHelper, cleanup *common.CleanupManager) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var createdUserIDs []string

	t.Run("全テナントユーザー一覧取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetAllTenantUsersWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("全テナントユーザー一覧取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "全テナントユーザー一覧取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "全テナントユーザー一覧取得")

		t.Logf("全テナントユーザー一覧取得成功: ユーザー数=%d", len(resp.JSON200.Users))
	})

	t.Run("全テナントユーザー作成", func(t *testing.T) {
		for i, param := range testData.AllTenantUsers.Create.Params {
			t.Run(param.Email, func(t *testing.T) {
				// 作成パラメータを準備
				createParam := authapi.CreateAllTenantUserParam{
					Email:      param.Email,
					Attributes: &param.Attributes,
					Roles:      &param.Roles,
					Envs:       &param.Envs,
				}

				// 全テナントユーザーを作成
				startTime := time.Now()
				createResp, err := client.Client.CreateAllTenantUserWithResponse(ctx, createParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("全テナントユーザー作成APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 25*time.Second, "全テナントユーザー作成")

				// ステータスコードをチェック
				assert.AssertStatusCode(201, createResp.StatusCode(), "全テナントユーザー作成")

				if createResp.JSON201 != nil {
					userID := createResp.JSON201.Id
					createdUserIDs = append(createdUserIDs, userID)

					// リソース追跡に追加
					client.CreateTestResource(
						common.ResourceTypeTenantUser,
						userID,
						param.Email,
						common.StoryTenantUserManagement,
						map[string]interface{}{
							"email":      param.Email,
							"attributes": param.Attributes,
							"roles":      param.Roles,
							"envs":       param.Envs,
						},
					)

					t.Logf("全テナントユーザー作成成功: ID=%s, Email=%s", userID, param.Email)
				}
			})
		}
	})

	t.Run("全テナントユーザー詳細取得", func(t *testing.T) {
		for _, userID := range createdUserIDs {
			t.Run(userID, func(t *testing.T) {
				startTime := time.Now()
				resp, err := client.Client.GetAllTenantUserWithResponse(ctx, userID)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("全テナントユーザー詳細取得APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 10*time.Second, "全テナントユーザー詳細取得")

				// ステータスコードをチェック
				assert.AssertStatusCode(200, resp.StatusCode(), "全テナントユーザー詳細取得")

				if resp.JSON200 != nil {
					// ユーザー情報をチェック
					assert.AssertEquals(userID, resp.JSON200.Id, "ユーザーID")
				}

				t.Logf("全テナントユーザー詳細取得成功: ID=%s", userID)
			})
		}
	})

	t.Run("全テナントユーザー属性更新", func(t *testing.T) {
		if len(createdUserIDs) > 0 {
			userID := createdUserIDs[0]
			t.Run(userID, func(t *testing.T) {
				// 更新パラメータを準備
				updateParam := authapi.UpdateAllTenantUserAttributesParam{
					Attributes: testData.AllTenantUsers.Update.Params.Attributes,
				}

				// 全テナントユーザー属性を更新
				startTime := time.Now()
				updateResp, err := client.Client.UpdateAllTenantUserAttributesWithResponse(ctx, userID, updateParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("全テナントユーザー属性更新APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 15*time.Second, "全テナントユーザー属性更新")

				// ステータスコードをチェック
				assert.AssertStatusCode(200, updateResp.StatusCode(), "全テナントユーザー属性更新")

				t.Logf("全テナントユーザー属性更新成功: ID=%s", userID)
			})
		}
	})

	t.Run("全テナントユーザー削除", func(t *testing.T) {
		for _, userID := range createdUserIDs {
			t.Run(userID, func(t *testing.T) {
				// 全テナントユーザーを削除
				startTime := time.Now()
				deleteResp, err := client.Client.DeleteAllTenantUserWithResponse(ctx, userID)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("全テナントユーザー削除APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 25*time.Second, "全テナントユーザー削除")

				// ステータスコードをチェック
				assert.AssertResourceDeleted(deleteResp.StatusCode(), "全テナントユーザー")

				// リソースを削除済みとしてマーク
				client.MarkResourceCleaned(userID, nil)

				// 削除確認
				confirmResp, err := client.Client.GetAllTenantUserWithResponse(ctx, userID)
				if err == nil && confirmResp.StatusCode() != 404 {
					t.Errorf("全テナントユーザーが削除されていません: ステータスコード %d", confirmResp.StatusCode())
				}

				t.Logf("全テナントユーザー削除成功: ID=%s", userID)
			})
		}
	})
}

// testRoleManagement は役割管理のテストを実行します
func testRoleManagement(t *testing.T, client *common.ClientWrapper, testData *common.TenantUserManagementTestData, assert *common.AssertionHelper, testTenantID string, testUserIDs []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if len(testUserIDs) == 0 {
		t.Skip("テスト用ユーザーが存在しないためスキップ")
		return
	}

	testUserID := testUserIDs[0]

	// まずテナントユーザーとして追加
	createParam := authapi.CreateTenantUserParam{
		UserId: testUserID,
		Roles:  &[]string{},
		Envs:   &[]string{"dev"},
	}

	_, err := client.Client.CreateTenantUserWithResponse(ctx, testTenantID, createParam)
	if err != nil {
		t.Logf("テナントユーザー追加に失敗: %v", err)
		return
	}

	// テスト終了時にテナントユーザーを削除
	defer func() {
		client.Client.DeleteTenantUserWithResponse(ctx, testTenantID, testUserID)
	}()

	t.Run("役割付与", func(t *testing.T) {
		envID := testData.RoleManagement.AttachRole.Params.EnvID
		roleName := testData.RoleManagement.AttachRole.Params.RoleName

		// 役割を付与
		startTime := time.Now()
		attachResp, err := client.Client.AttachRoleWithResponse(ctx, testTenantID, testUserID, authapi.Id(envID), roleName)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("役割付与APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "役割付与")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, attachResp.StatusCode(), "役割付与")

		t.Logf("役割付与成功: TenantID=%s, UserID=%s, EnvID=%s, Role=%s", testTenantID, testUserID, envID, roleName)
	})

	t.Run("役割削除", func(t *testing.T) {
		envID := testData.RoleManagement.DetachRole.Params.EnvID
		roleName := testData.RoleManagement.DetachRole.Params.RoleName

		// 役割を削除
		startTime := time.Now()
		detachResp, err := client.Client.DetachRoleWithResponse(ctx, testTenantID, testUserID, authapi.Id(envID), roleName)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("役割削除APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "役割削除")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, detachResp.StatusCode(), "役割削除")

		t.Logf("役割削除成功: TenantID=%s, UserID=%s, EnvID=%s, Role=%s", testTenantID, testUserID, envID, roleName)
	})
}