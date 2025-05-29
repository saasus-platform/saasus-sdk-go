package stories

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

// TestStory01_BasicSetup は基本セットアップストーリーのE2Eテストです
func TestStory01_BasicSetup(t *testing.T) {
	// テストデータローダーを初期化
	loader := common.NewTestDataLoader()
	testData, err := loader.LoadBasicSetupData()
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
		if err := cleanup.CleanupByStory(ctx, common.StoryBasicSetup); err != nil {
			t.Logf("クリーンアップに失敗: %v", err)
		}
	}()

	// サブテストを順次実行
	t.Run("基本情報管理", func(t *testing.T) {
		testBasicInfoManagement(t, client, testData, assert)
	})

	t.Run("認証情報管理", func(t *testing.T) {
		testAuthInfoManagement(t, client, testData, assert)
	})

	t.Run("環境管理", func(t *testing.T) {
		testEnvManagement(t, client, testData, assert, cleanup)
	})

	t.Run("外部IDプロバイダ管理", func(t *testing.T) {
		testIdentityProviderManagement(t, client, assert)
	})

	t.Run("サインイン設定管理", func(t *testing.T) {
		testSignInSettingsManagement(t, client, assert)
	})

	t.Run("カスタマイズ設定管理", func(t *testing.T) {
		testCustomizeSettingsManagement(t, client, assert)
	})
}

// testBasicInfoManagement は基本情報管理のテストを実行します
func testBasicInfoManagement(t *testing.T, client *common.ClientWrapper, testData *common.BasicSetupTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("基本情報取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetBasicInfoWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("基本情報取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "基本情報取得")

		// レスポンス内容をアサート
		assert.AssertBasicInfoResponse(resp)

		t.Logf("基本情報取得成功: ドメイン=%s", resp.JSON200.DomainName)
	})

	t.Run("基本情報更新", func(t *testing.T) {
		// 現在の設定を取得（復元用）
		originalResp, err := client.Client.GetBasicInfoWithResponse(ctx)
		if err != nil {
			t.Fatalf("現在の基本情報取得に失敗: %v", err)
		}

		// 更新パラメータを準備
		updateParam := authapi.UpdateBasicInfoParam{
			DomainName:        testData.BasicInfo.Update.Params.DomainName,
			FromEmailAddress:  testData.BasicInfo.Update.Params.FromEmailAddress,
			ReplyEmailAddress: testData.BasicInfo.Update.Params.ReplyEmailAddress,
		}

		// 基本情報を更新
		startTime := time.Now()
		updateResp, err := client.Client.UpdateBasicInfoWithResponse(ctx, updateParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("基本情報更新APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "基本情報更新")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, updateResp.StatusCode(), "基本情報更新")

		// 更新後の情報を確認
		updatedResp, err := client.Client.GetBasicInfoWithResponse(ctx)
		if err != nil {
			t.Fatalf("更新後の基本情報取得に失敗: %v", err)
		}

		// 更新が反映されているかチェック
		assert.AssertEquals(testData.BasicInfo.Update.Params.DomainName, 
			updatedResp.JSON200.DomainName, "更新後ドメイン名")
		assert.AssertEquals(testData.BasicInfo.Update.Params.FromEmailAddress, 
			updatedResp.JSON200.FromEmailAddress, "更新後送信元メールアドレス")

		// 元の設定に復元
		if originalResp.JSON200 != nil {
			restoreParam := authapi.UpdateBasicInfoParam{
				DomainName:        originalResp.JSON200.DomainName,
				FromEmailAddress:  originalResp.JSON200.FromEmailAddress,
				ReplyEmailAddress: &originalResp.JSON200.ReplyEmailAddress,
			}
			_, err = client.Client.UpdateBasicInfoWithResponse(ctx, restoreParam)
			if err != nil {
				t.Logf("基本情報の復元に失敗: %v", err)
			}
		}

		t.Log("基本情報更新成功")
	})
}

// testAuthInfoManagement は認証情報管理のテストを実行します
func testAuthInfoManagement(t *testing.T, client *common.ClientWrapper, testData *common.BasicSetupTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("認証情報取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetAuthInfoWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "認証情報取得")

		// レスポンス内容をアサート
		assert.AssertAuthInfoResponse(resp)

		t.Logf("認証情報取得成功: コールバックURL=%s", resp.JSON200.CallbackUrl)
	})

	t.Run("認証情報更新", func(t *testing.T) {
		// 現在の設定を取得（復元用）
		originalResp, err := client.Client.GetAuthInfoWithResponse(ctx)
		if err != nil {
			t.Fatalf("現在の認証情報取得に失敗: %v", err)
		}

		// 更新パラメータを準備
		updateParam := authapi.UpdateAuthInfoParam{
			CallbackUrl: testData.AuthInfo.Update.Params.CallbackUrl,
		}

		// 認証情報を更新
		startTime := time.Now()
		updateResp, err := client.Client.UpdateAuthInfoWithResponse(ctx, updateParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報更新APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "認証情報更新")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, updateResp.StatusCode(), "認証情報更新")

		// 更新後の情報を確認
		updatedResp, err := client.Client.GetAuthInfoWithResponse(ctx)
		if err != nil {
			t.Fatalf("更新後の認証情報取得に失敗: %v", err)
		}

		// 更新が反映されているかチェック
		assert.AssertEquals(testData.AuthInfo.Update.Params.CallbackUrl, 
			updatedResp.JSON200.CallbackUrl, "更新後コールバックURL")

		// 元の設定に復元
		if originalResp.JSON200 != nil {
			restoreParam := authapi.UpdateAuthInfoParam{
				CallbackUrl: originalResp.JSON200.CallbackUrl,
			}
			_, err = client.Client.UpdateAuthInfoWithResponse(ctx, restoreParam)
			if err != nil {
				t.Logf("認証情報の復元に失敗: %v", err)
			}
		}

		t.Log("認証情報更新成功")
	})
}

// testEnvManagement は環境管理のテストを実行します
func testEnvManagement(t *testing.T, client *common.ClientWrapper, testData *common.BasicSetupTestData, assert *common.AssertionHelper, cleanup *common.CleanupManager) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var createdEnvID string

	t.Run("環境一覧取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetEnvsWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("環境一覧取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "環境一覧取得")

		// レスポンス内容をアサート
		assert.AssertEnvsResponse(resp)

		t.Logf("環境一覧取得成功: 環境数=%d", len(resp.JSON200.Envs))
	})

	t.Run("環境作成", func(t *testing.T) {
		// 作成パラメータを準備
		envIDUint, err := strconv.ParseUint(testData.Envs.Create.Params.Id, 10, 64)
		if err != nil {
			t.Fatalf("環境ID変換エラー: %v", err)
		}
		createParam := authapi.CreateEnvParam{
			Id:          authapi.Id(envIDUint),
			Name:        testData.Envs.Create.Params.Name,
			DisplayName: testData.Envs.Create.Params.DisplayName,
		}

		// 環境を作成
		startTime := time.Now()
		createResp, err := client.Client.CreateEnvWithResponse(ctx, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("環境作成APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "環境作成")

		// ステータスコードをチェック
		assert.AssertStatusCode(201, createResp.StatusCode(), "環境作成")

		createdEnvID = testData.Envs.Create.Params.Id

		// リソース追跡に追加
		client.CreateTestResource(
			common.ResourceTypeEnv,
			createdEnvID,
			testData.Envs.Create.Params.Name,
			common.StoryBasicSetup,
			map[string]interface{}{
				"display_name": testData.Envs.Create.Params.DisplayName,
			},
		)

		t.Logf("環境作成成功: ID=%s", createdEnvID)
	})

	t.Run("環境詳細取得", func(t *testing.T) {
		if createdEnvID == "" {
			t.Skip("環境が作成されていないためスキップ")
		}

		startTime := time.Now()
		envIDUint, err := strconv.ParseUint(createdEnvID, 10, 64)
		if err != nil {
			t.Fatalf("環境ID変換エラー: %v", err)
		}
		resp, err := client.Client.GetEnvWithResponse(ctx, authapi.Id(envIDUint))
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("環境詳細取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "環境詳細取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "環境詳細取得")

		if resp.JSON200 != nil {
			// 環境情報をチェック
			assert.AssertEquals(createdEnvID, string(resp.JSON200.Id), "環境ID")
			assert.AssertEquals(testData.Envs.Create.Params.Name, resp.JSON200.Name, "環境名")
		}

		t.Logf("環境詳細取得成功: ID=%s", createdEnvID)
	})

	t.Run("環境更新", func(t *testing.T) {
		if createdEnvID == "" {
			t.Skip("環境が作成されていないためスキップ")
		}

		// 更新パラメータを準備
		updateParam := authapi.UpdateEnvParam{
			Name:        testData.Envs.Update.Params.Name,
			DisplayName: testData.Envs.Update.Params.DisplayName,
		}

		// 環境を更新
		startTime := time.Now()
		envIDUint, err := strconv.ParseUint(createdEnvID, 10, 64)
		if err != nil {
			t.Fatalf("環境ID変換エラー: %v", err)
		}
		updateResp, err := client.Client.UpdateEnvWithResponse(ctx, authapi.Id(envIDUint), updateParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("環境更新APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "環境更新")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, updateResp.StatusCode(), "環境更新")

		// 更新後の情報を確認
		envIDUint, err := strconv.ParseUint(createdEnvID, 10, 64)
		if err != nil {
			t.Fatalf("環境ID変換エラー: %v", err)
		}
		updatedResp, err := client.Client.GetEnvWithResponse(ctx, authapi.Id(envIDUint))
		if err != nil {
			t.Fatalf("更新後の環境詳細取得に失敗: %v", err)
		}

		if updatedResp.JSON200 != nil {
			// 更新が反映されているかチェック
			assert.AssertEquals(testData.Envs.Update.Params.Name, updatedResp.JSON200.Name, "更新後環境名")
		}

		t.Logf("環境更新成功: ID=%s", createdEnvID)
	})

	t.Run("環境削除", func(t *testing.T) {
		if createdEnvID == "" {
			t.Skip("環境が作成されていないためスキップ")
		}

		// 環境を削除
		startTime := time.Now()
		envIDUint, err := strconv.ParseUint(createdEnvID, 10, 64)
		if err != nil {
			t.Fatalf("環境ID変換エラー: %v", err)
		}
		deleteResp, err := client.Client.DeleteEnvWithResponse(ctx, authapi.Id(envIDUint))
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("環境削除APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "環境削除")

		// ステータスコードをチェック
		assert.AssertResourceDeleted(deleteResp.StatusCode(), "環境")

		// リソースを削除済みとしてマーク
		client.MarkResourceCleaned(createdEnvID, nil)

		// 削除確認
		envIDUint, err := strconv.ParseUint(createdEnvID, 10, 64)
		if err != nil {
			t.Fatalf("環境ID変換エラー: %v", err)
		}
		confirmResp, err := client.Client.GetEnvWithResponse(ctx, authapi.Id(envIDUint))
		if err == nil && confirmResp.StatusCode() != 404 {
			t.Errorf("環境が削除されていません: ステータスコード %d", confirmResp.StatusCode())
		}

		t.Logf("環境削除成功: ID=%s", createdEnvID)
	})
}

// testIdentityProviderManagement は外部IDプロバイダ管理のテストを実行します
func testIdentityProviderManagement(t *testing.T, client *common.ClientWrapper, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("外部IDプロバイダ取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetIdentityProvidersWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("外部IDプロバイダ取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "外部IDプロバイダ取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "外部IDプロバイダ取得")

		t.Log("外部IDプロバイダ取得成功")
	})
}

// testSignInSettingsManagement はサインイン設定管理のテストを実行します
func testSignInSettingsManagement(t *testing.T, client *common.ClientWrapper, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("サインイン設定取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetSignInSettingsWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("サインイン設定取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "サインイン設定取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "サインイン設定取得")

		t.Log("サインイン設定取得成功")
	})
}

// testCustomizeSettingsManagement はカスタマイズ設定管理のテストを実行します
func testCustomizeSettingsManagement(t *testing.T, client *common.ClientWrapper, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("カスタマイズページ設定取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetCustomizePageSettingsWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("カスタマイズページ設定取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "カスタマイズページ設定取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "カスタマイズページ設定取得")

		t.Log("カスタマイズページ設定取得成功")
	})

	t.Run("認証画面設定取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetCustomizePagesWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証画面設定取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "認証画面設定取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "認証画面設定取得")

		t.Log("認証画面設定取得成功")
	})

	t.Run("通知メールテンプレート取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.FindNotificationMessagesWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("通知メールテンプレート取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "通知メールテンプレート取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "通知メールテンプレート取得")

		t.Log("通知メールテンプレート取得成功")
	})
}