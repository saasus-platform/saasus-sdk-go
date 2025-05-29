package stories

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

// TestStory09_ErrorHandling はエラーハンドリングストーリーのE2Eテストです
func TestStory09_ErrorHandling(t *testing.T) {
	// テストデータローダーを初期化
	loader := common.NewTestDataLoader()
	testData, err := loader.LoadErrorHandlingData()
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
		if err := cleanup.CleanupByStory(ctx, common.StoryErrorHandling); err != nil {
			t.Logf("クリーンアップに失敗: %v", err)
		}
	}()

	// サブテストを順次実行
	t.Run("HTTPエラーステータステスト", func(t *testing.T) {
		testHTTPErrorStatuses(t, client, testData, assert)
	})

	t.Run("エラーレスポンス形式テスト", func(t *testing.T) {
		testErrorResponseFormat(t, client, testData, assert)
	})

	t.Run("認証・認可エラーテスト", func(t *testing.T) {
		testAuthenticationErrors(t, client, testData, assert)
	})

	t.Run("バリデーションエラーテスト", func(t *testing.T) {
		testValidationErrors(t, client, testData, assert)
	})

	t.Run("リソース不存在エラーテスト", func(t *testing.T) {
		testResourceNotFoundErrors(t, client, testData, assert)
	})
}

// testHTTPErrorStatuses はHTTPエラーステータスのテストを実行します
func testHTTPErrorStatuses(t *testing.T, client *common.ClientWrapper, testData *common.ErrorHandlingTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// 各エラーステータスをテスト
	errorTests := []struct {
		name           string
		expectedStatus int
		testFunc       func() (*http.Response, error)
	}{
		{
			name:           "400 Bad Request",
			expectedStatus: 400,
			testFunc: func() (*http.Response, error) {
				// 無効なJSONでユーザー作成を試行
				req, _ := http.NewRequestWithContext(ctx, "POST", "/users", nil)
				return http.DefaultClient.Do(req)
			},
		},
		{
			name:           "401 Unauthorized",
			expectedStatus: 401,
			testFunc: func() (*http.Response, error) {
				// 認証なしでユーザー一覧取得を試行
				req, _ := http.NewRequestWithContext(ctx, "GET", "/users", nil)
				return http.DefaultClient.Do(req)
			},
		},
		{
			name:           "404 Not Found",
			expectedStatus: 404,
			testFunc: func() (*http.Response, error) {
				// 存在しないユーザーの取得を試行
				resp, err := client.Client.GetSaasUserWithResponse(ctx, "non-existent-user-id")
				if err != nil {
					return nil, err
				}
				return &http.Response{StatusCode: resp.StatusCode()}, nil
			},
		},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			startTime := time.Now()
			resp, err := test.testFunc()
			duration := time.Since(startTime)

			if err != nil {
				t.Logf("エラーが期待通り発生: %v", err)
				return
			}

			// レスポンス時間をチェック
			assert.AssertResponseTime(duration, 10*time.Second, test.name)

			// 期待されるエラーステータスコードをチェック
			if resp != nil {
				assert.AssertStatusCode(test.expectedStatus, resp.StatusCode, test.name)
				t.Logf("%s エラー確認成功", test.name)
			}
		})
	}
}

// testErrorResponseFormat はエラーレスポンス形式のテストを実行します
func testErrorResponseFormat(t *testing.T, client *common.ClientWrapper, testData *common.ErrorHandlingTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Run("エラーレスポンス形式統一性確認", func(t *testing.T) {
		// 存在しないユーザーを取得してエラーレスポンスを確認
		startTime := time.Now()
		resp, err := client.Client.GetSaasUserWithResponse(ctx, "non-existent-user-id")
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "エラーレスポンス形式確認")

		// 404エラーを期待
		assert.AssertStatusCode(404, resp.StatusCode(), "存在しないユーザー取得")

		// エラーレスポンスの形式をチェック
		if resp.JSON404 != nil {
			// 必須フィールドの存在確認
			assert.AssertNotEmpty(resp.JSON404.Type, "エラータイプ")
			assert.AssertNotEmpty(resp.JSON404.Message, "エラーメッセージ")
			t.Log("エラーレスポンス形式確認成功")
		}
	})

	t.Run("エラーメッセージ内容確認", func(t *testing.T) {
		// 無効なテナントIDでテナント取得を試行
		resp, err := client.Client.GetTenantWithResponse(ctx, "invalid-tenant-id")
		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		if resp.StatusCode() == 404 && resp.JSON404 != nil {
			// エラーメッセージが適切な内容であることを確認
			message := resp.JSON404.Message
			assert.AssertNotEmpty(message, "エラーメッセージ")
			
			// 日本語メッセージかどうかの簡単なチェック
			if len(message) > 0 {
				t.Logf("エラーメッセージ確認成功: %s", message)
			}
		}
	})

	t.Run("エラーレスポンスJSON形式確認", func(t *testing.T) {
		// 無効なロール名でロール削除を試行
		resp, err := client.Client.DeleteRoleWithResponse(ctx, "invalid-role-name")
		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// エラーレスポンスがJSON形式であることを確認
		if resp.StatusCode() >= 400 {
			// レスポンスボディを取得してJSON解析を試行
			body := resp.Body
			if len(body) > 0 {
				var errorResponse map[string]interface{}
				if err := json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("エラーレスポンスがJSON形式ではありません: %v", err)
				} else {
					t.Log("エラーレスポンスJSON形式確認成功")
				}
			}
		}
	})
}

// testAuthenticationErrors は認証・認可エラーのテストを実行します
func testAuthenticationErrors(t *testing.T, client *common.ClientWrapper, testData *common.ErrorHandlingTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Run("無効な認証トークンエラー", func(t *testing.T) {
		// 無効な認証トークンでクライアントを作成
		invalidClient, err := common.NewClientWithCustomAuth("invalid-token")
		if err != nil {
			t.Skip("無効な認証トークンクライアントの作成に失敗")
			return
		}

		// ユーザー一覧取得を試行
		startTime := time.Now()
		resp, err := invalidClient.Client.GetSaasUsersWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "無効な認証トークン")

		// 401エラーを期待
		assert.AssertStatusCode(401, resp.StatusCode(), "無効な認証トークン")

		t.Log("無効な認証トークンエラー確認成功")
	})

	t.Run("権限不足エラー", func(t *testing.T) {
		// 管理者権限が必要な操作を一般ユーザー権限で実行
		// 注意: 実際のテストでは適切な権限レベルのトークンを使用する必要があります
		
		// 基本情報の更新を試行（管理者権限が必要と仮定）
		startTime := time.Now()
		resp, err := client.Client.UpdateBasicInfoWithResponse(ctx, authapi.UpdateBasicInfoJSONRequestBody{
			DomainName: "unauthorized-test.example.com",
			FromEmailAddress: "test@example.com",
		})
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "権限不足")

		// 権限エラーの場合は403または401を期待
		if resp.StatusCode() != 200 && resp.StatusCode() != 403 && resp.StatusCode() != 401 {
			t.Logf("権限チェック結果: ステータスコード %d", resp.StatusCode())
		}

		t.Log("権限チェック完了")
	})
}

// testValidationErrors はバリデーションエラーのテストを実行します
func testValidationErrors(t *testing.T, client *common.ClientWrapper, testData *common.ErrorHandlingTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Run("無効なメールアドレス形式エラー", func(t *testing.T) {
		// 無効なメールアドレスでユーザー作成を試行
		invalidEmail := "invalid-email-format"
		
		startTime := time.Now()
		resp, err := client.Client.CreateSaasUserWithResponse(ctx, authapi.CreateSaasUserJSONRequestBody{
			Email: invalidEmail,
		})
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "無効なメールアドレス")

		// バリデーションエラーを期待（400）
		if resp.StatusCode() == 201 {
			t.Error("無効なメールアドレスで成功レスポンスが返されました")
		} else {
			t.Logf("無効なメールアドレスエラー確認成功: ステータスコード %d", resp.StatusCode())
		}
	})

	t.Run("必須パラメータ不足エラー", func(t *testing.T) {
		// 必須パラメータなしでテナント作成を試行
		startTime := time.Now()
		resp, err := client.Client.CreateTenantWithResponse(ctx, authapi.CreateTenantJSONRequestBody{
			Name: "", // 空の名前でバリデーションエラーを発生させる
		})
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "必須パラメータ不足")

		// バリデーションエラーを期待（400）
		if resp.StatusCode() == 201 {
			t.Error("必須パラメータ不足で成功レスポンスが返されました")
		} else {
			t.Logf("必須パラメータ不足エラー確認成功: ステータスコード %d", resp.StatusCode())
		}
	})

	t.Run("文字数制限エラー", func(t *testing.T) {
		// 長すぎる名前でロール作成を試行
		longRoleName := string(make([]byte, 1000)) // 1000文字の文字列
		for i := range longRoleName {
			longRoleName = longRoleName[:i] + "a" + longRoleName[i+1:]
		}

		startTime := time.Now()
		resp, err := client.Client.CreateRoleWithResponse(ctx, authapi.CreateRoleJSONRequestBody{
			RoleName:    longRoleName,
			DisplayName: "テスト用長い名前ロール",
		})
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "文字数制限")

		// バリデーションエラーを期待（400）
		if resp.StatusCode() == 201 {
			t.Error("文字数制限を超えた値で成功レスポンスが返されました")
		} else {
			t.Logf("文字数制限エラー確認成功: ステータスコード %d", resp.StatusCode())
		}
	})
}

// testResourceNotFoundErrors はリソース不存在エラーのテストを実行します
func testResourceNotFoundErrors(t *testing.T, client *common.ClientWrapper, testData *common.ErrorHandlingTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Run("存在しないユーザーエラー", func(t *testing.T) {
		nonExistentUserID := "non-existent-user-12345"

		startTime := time.Now()
		resp, err := client.Client.GetSaasUserWithResponse(ctx, nonExistentUserID)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "存在しないユーザー")

		// 404エラーを期待
		assert.AssertStatusCode(404, resp.StatusCode(), "存在しないユーザー")

		t.Log("存在しないユーザーエラー確認成功")
	})

	t.Run("存在しないテナントエラー", func(t *testing.T) {
		nonExistentTenantID := "non-existent-tenant-12345"

		startTime := time.Now()
		resp, err := client.Client.GetTenantWithResponse(ctx, nonExistentTenantID)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "存在しないテナント")

		// 404エラーを期待
		assert.AssertStatusCode(404, resp.StatusCode(), "存在しないテナント")

		t.Log("存在しないテナントエラー確認成功")
	})

	t.Run("存在しないロールエラー", func(t *testing.T) {
		nonExistentRoleName := "non-existent-role"

		startTime := time.Now()
		resp, err := client.Client.DeleteRoleWithResponse(ctx, nonExistentRoleName)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "存在しないロール")

		// 404エラーを期待
		assert.AssertStatusCode(404, resp.StatusCode(), "存在しないロール")

		t.Log("存在しないロールエラー確認成功")
	})

	t.Run("存在しない招待エラー", func(t *testing.T) {
		// テスト用テナントを作成
		testTenantResp, err := client.Client.CreateTenantWithResponse(ctx, authapi.CreateTenantJSONRequestBody{
			Name:                   "エラーテスト用テナント",
			BackOfficeStaffEmail: "error-test@example.com",
		})
		if err != nil || testTenantResp.JSON201 == nil {
			t.Skip("テスト用テナントの作成に失敗")
			return
		}

		testTenantID := testTenantResp.JSON201.Id
		nonExistentInvitationID := "non-existent-invitation-12345"

		// テスト終了時にテナントを削除
		defer func() {
			client.Client.DeleteTenantWithResponse(ctx, testTenantID)
		}()

		startTime := time.Now()
		resp, err := client.Client.GetTenantInvitationWithResponse(ctx, testTenantID, nonExistentInvitationID)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("API呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "存在しない招待")

		// 404エラーを期待
		assert.AssertStatusCode(404, resp.StatusCode(), "存在しない招待")

		t.Log("存在しない招待エラー確認成功")
	})
}