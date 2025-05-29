package stories

import (
	"context"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

// TestStory02_UserManagement はユーザー管理ストーリーのE2Eテストです
func TestStory02_UserManagement(t *testing.T) {
	// テストデータローダーを初期化
	loader := common.NewTestDataLoader()
	testData, err := loader.LoadUserManagementData()
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
		if err := cleanup.CleanupByStory(ctx, common.StoryUserManagement); err != nil {
			t.Logf("クリーンアップに失敗: %v", err)
		}
	}()

	// サブテストを順次実行
	t.Run("ユーザー属性管理", func(t *testing.T) {
		testUserAttributeManagement(t, client, testData, assert, cleanup)
	})

	t.Run("SaaSユーザー管理", func(t *testing.T) {
		testSaasUserManagement(t, client, testData, assert, cleanup)
	})

	t.Run("ユーザー情報更新", func(t *testing.T) {
		testUserInfoManagement(t, client, testData, assert)
	})
}

// testUserAttributeManagement はユーザー属性管理のテストを実行します
func testUserAttributeManagement(t *testing.T, client *common.ClientWrapper, testData *common.UserManagementTestData, assert *common.AssertionHelper, cleanup *common.CleanupManager) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var createdAttributes []string

	t.Run("ユーザー属性一覧取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetUserAttributesWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("ユーザー属性一覧取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "ユーザー属性一覧取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "ユーザー属性一覧取得")

		t.Logf("ユーザー属性一覧取得成功: 属性数=%d", len(resp.JSON200.UserAttributes))
	})

	t.Run("ユーザー属性作成", func(t *testing.T) {
		for _, param := range testData.UserAttributes.Create.Params {
			t.Run(param.AttributeName, func(t *testing.T) {
				// 作成パラメータを準備
				createParam := authapi.CreateUserAttributeParam{
					AttributeName: param.AttributeName,
					DisplayName:   param.DisplayName,
				}

				// ユーザー属性を作成
				startTime := time.Now()
				createResp, err := client.Client.CreateUserAttributeWithResponse(ctx, createParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("ユーザー属性作成APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 15*time.Second, "ユーザー属性作成")

				// ステータスコードをチェック
				assert.AssertStatusCode(201, createResp.StatusCode(), "ユーザー属性作成")

				createdAttributes = append(createdAttributes, param.AttributeName)

				// リソース追跡に追加
				client.CreateTestResource(
					common.ResourceTypeUserAttribute,
					param.AttributeName,
					param.DisplayName,
					common.StoryUserManagement,
					map[string]interface{}{
						"display_name": param.DisplayName,
					},
				)

				t.Logf("ユーザー属性作成成功: %s", param.AttributeName)
			})
		}
	})

	t.Run("ユーザー属性削除", func(t *testing.T) {
		for _, attributeName := range createdAttributes {
			t.Run(attributeName, func(t *testing.T) {
				// ユーザー属性を削除
				startTime := time.Now()
				deleteResp, err := client.Client.DeleteUserAttributeWithResponse(ctx, attributeName)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("ユーザー属性削除APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 15*time.Second, "ユーザー属性削除")

				// ステータスコードをチェック
				assert.AssertResourceDeleted(deleteResp.StatusCode(), "ユーザー属性")

				// リソースを削除済みとしてマーク
				client.MarkResourceCleaned(attributeName, nil)

				t.Logf("ユーザー属性削除成功: %s", attributeName)
			})
		}
	})
}

// testSaasUserManagement はSaaSユーザー管理のテストを実行します
func testSaasUserManagement(t *testing.T, client *common.ClientWrapper, testData *common.UserManagementTestData, assert *common.AssertionHelper, cleanup *common.CleanupManager) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var createdUserIDs []string

	t.Run("SaaSユーザー一覧取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetUsersWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("SaaSユーザー一覧取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "SaaSユーザー一覧取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "SaaSユーザー一覧取得")

		t.Logf("SaaSユーザー一覧取得成功: ユーザー数=%d", len(resp.JSON200.Users))
	})

	t.Run("SaaSユーザー作成", func(t *testing.T) {
		for i, param := range testData.SaasUsers.Create.Params {
			t.Run(param.Email, func(t *testing.T) {
				// 作成パラメータを準備
				createParam := authapi.CreateSaasUserParam{
					Email:      param.Email,
					Attributes: &param.Attributes,
				}

				// SaaSユーザーを作成
				startTime := time.Now()
				createResp, err := client.Client.CreateSaasUserWithResponse(ctx, createParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("SaaSユーザー作成APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 20*time.Second, "SaaSユーザー作成")

				// ステータスコードをチェック
				assert.AssertStatusCode(201, createResp.StatusCode(), "SaaSユーザー作成")

				if createResp.JSON201 != nil {
					userID := createResp.JSON201.Id
					createdUserIDs = append(createdUserIDs, userID)

					// リソース追跡に追加
					client.CreateTestResource(
						common.ResourceTypeSaasUser,
						userID,
						param.Email,
						common.StoryUserManagement,
						map[string]interface{}{
							"email":      param.Email,
							"attributes": param.Attributes,
						},
					)

					t.Logf("SaaSユーザー作成成功: ID=%s, Email=%s", userID, param.Email)
				}
			})
		}
	})

	t.Run("SaaSユーザー詳細取得", func(t *testing.T) {
		for _, userID := range createdUserIDs {
			t.Run(userID, func(t *testing.T) {
				startTime := time.Now()
				resp, err := client.Client.GetSaasUserWithResponse(ctx, userID)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("SaaSユーザー詳細取得APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 10*time.Second, "SaaSユーザー詳細取得")

				// ステータスコードをチェック
				assert.AssertStatusCode(200, resp.StatusCode(), "SaaSユーザー詳細取得")

				if resp.JSON200 != nil {
					// ユーザー情報をチェック
					assert.AssertEquals(userID, resp.JSON200.Id, "ユーザーID")
				}

				t.Logf("SaaSユーザー詳細取得成功: ID=%s", userID)
			})
		}
	})

	t.Run("ユーザー属性更新", func(t *testing.T) {
		if len(createdUserIDs) > 0 {
			userID := createdUserIDs[0]
			t.Run(userID, func(t *testing.T) {
				// 更新パラメータを準備
				updateParam := authapi.UpdateSaasUserAttributesParam{
					Attributes: testData.SaasUsers.Update.Params.Attributes,
				}

				// ユーザー属性を更新
				startTime := time.Now()
				updateResp, err := client.Client.UpdateSaasUserAttributesWithResponse(ctx, userID, updateParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("ユーザー属性更新APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 15*time.Second, "ユーザー属性更新")

				// ステータスコードをチェック
				assert.AssertStatusCode(200, updateResp.StatusCode(), "ユーザー属性更新")

				t.Logf("ユーザー属性更新成功: ID=%s", userID)
			})
		}
	})

	t.Run("SaaSユーザー削除", func(t *testing.T) {
		for _, userID := range createdUserIDs {
			t.Run(userID, func(t *testing.T) {
				// SaaSユーザーを削除
				startTime := time.Now()
				deleteResp, err := client.Client.DeleteSaasUserWithResponse(ctx, userID)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("SaaSユーザー削除APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 20*time.Second, "SaaSユーザー削除")

				// ステータスコードをチェック
				assert.AssertResourceDeleted(deleteResp.StatusCode(), "SaaSユーザー")

				// リソースを削除済みとしてマーク
				client.MarkResourceCleaned(userID, nil)

				// 削除確認
				confirmResp, err := client.Client.GetSaasUserWithResponse(ctx, userID)
				if err == nil && confirmResp.StatusCode() != 404 {
					t.Errorf("SaaSユーザーが削除されていません: ステータスコード %d", confirmResp.StatusCode())
				}

				t.Logf("SaaSユーザー削除成功: ID=%s", userID)
			})
		}
	})
}

// testUserInfoManagement はユーザー情報管理のテストを実行します
func testUserInfoManagement(t *testing.T, client *common.ClientWrapper, testData *common.UserManagementTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// テスト用ユーザーを作成
	createParam := authapi.CreateSaasUserParam{
		Email: "test-user-info@example.com",
		Attributes: &map[string]interface{}{
			"test": "value",
		},
	}

	createResp, err := client.Client.CreateSaasUserWithResponse(ctx, createParam)
	if err != nil || createResp.JSON201 == nil {
		t.Skip("テスト用ユーザーの作成に失敗したためスキップ")
		return
	}

	testUserID := createResp.JSON201.Id

	// テスト終了時にユーザーを削除
	defer func() {
		client.Client.DeleteSaasUserWithResponse(ctx, testUserID)
	}()

	t.Run("メールアドレス変更", func(t *testing.T) {
		// 更新パラメータを準備
		updateParam := authapi.UpdateSaasUserEmailParam{
			Email: testData.UserInfo.UpdateEmail.Params.Email,
		}

		// メールアドレスを更新
		startTime := time.Now()
		updateResp, err := client.Client.UpdateSaasUserEmailWithResponse(ctx, testUserID, updateParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("メールアドレス変更APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "メールアドレス変更")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, updateResp.StatusCode(), "メールアドレス変更")

		t.Log("メールアドレス変更成功")
	})

	t.Run("パスワード変更", func(t *testing.T) {
		// 更新パラメータを準備
		updateParam := authapi.UpdateSaasUserPasswordParam{
			Password: testData.UserInfo.UpdatePassword.Params.Password,
		}

		// パスワードを更新
		startTime := time.Now()
		updateResp, err := client.Client.UpdateSaasUserPasswordWithResponse(ctx, testUserID, updateParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("パスワード変更APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "パスワード変更")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, updateResp.StatusCode(), "パスワード変更")

		t.Log("パスワード変更成功")
	})
}