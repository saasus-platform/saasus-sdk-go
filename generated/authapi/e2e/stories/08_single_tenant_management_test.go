package stories

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

// TestStory08_SingleTenantManagement はシングルテナント管理ストーリーのE2Eテストです
func TestStory08_SingleTenantManagement(t *testing.T) {
	// テストデータローダーを初期化
	loader := common.NewTestDataLoader()
	testData, err := loader.LoadSingleTenantManagementData()
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
		if err := cleanup.CleanupByStory(ctx, common.StorySingleTenantManagement); err != nil {
			t.Logf("クリーンアップに失敗: %v", err)
		}
	}()

	// サブテストを順次実行
	t.Run("シングルテナント設定管理", func(t *testing.T) {
		testSingleTenantSettings(t, client, testData, assert)
	})

	t.Run("CloudFormationテンプレート管理", func(t *testing.T) {
		testCloudFormationTemplate(t, client, testData, assert)
	})

	t.Run("DDLテンプレート管理", func(t *testing.T) {
		testDDLTemplate(t, client, testData, assert)
	})

	t.Run("AWS連携設定", func(t *testing.T) {
		testAWSIntegration(t, client, testData, assert)
	})
}

// testSingleTenantSettings はシングルテナント設定のテストを実行します
func testSingleTenantSettings(t *testing.T, client *common.ClientWrapper, testData *common.SingleTenantManagementTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	t.Run("シングルテナント設定取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetSingleTenantSettingsWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("シングルテナント設定取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "シングルテナント設定取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "シングルテナント設定取得")

		if resp.JSON200 != nil {
			t.Logf("シングルテナント設定取得成功")
		}
	})

	t.Run("シングルテナント設定更新", func(t *testing.T) {
		// 更新パラメータを準備
		updateParam := authapi.UpdateSingleTenantSettingsParam{
			CustomizePageTitle: &testData.TestSettings.PageTitle,
			CustomizePageUrl:   &testData.TestSettings.PageUrl,
			CustomizeCssUrl:    &testData.TestSettings.CssUrl,
			CustomizeIconUrl:   &testData.TestSettings.IconUrl,
			CustomizeLogoUrl:   &testData.TestSettings.LogoUrl,
			PrivacyPolicyUrl:   &testData.TestSettings.PrivacyUrl,
			TermsOfServiceUrl:  &testData.TestSettings.TermsUrl,
		}

		// シングルテナント設定を更新
		startTime := time.Now()
		updateResp, err := client.Client.UpdateSingleTenantSettingsWithResponse(ctx, updateParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("シングルテナント設定更新APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "シングルテナント設定更新")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, updateResp.StatusCode(), "シングルテナント設定更新")

		t.Logf("シングルテナント設定更新成功")

		// 更新後の設定を確認
		getResp, err := client.Client.GetSingleTenantSettingsWithResponse(ctx)
		if err == nil && getResp.JSON200 != nil {
			// 更新された値を確認
			if getResp.JSON200.CustomizePageTitle != nil {
				assert.AssertEquals(testData.TestSettings.PageTitle, *getResp.JSON200.CustomizePageTitle, "ページタイトル")
			}
			if getResp.JSON200.CustomizePageUrl != nil {
				assert.AssertEquals(testData.TestSettings.PageUrl, *getResp.JSON200.CustomizePageUrl, "ページURL")
			}
		}
	})

	t.Run("部分更新テスト", func(t *testing.T) {
		// ページタイトルのみを更新
		partialUpdateParam := authapi.UpdateSingleTenantSettingsParam{
			CustomizePageTitle: common.StringPtr("部分更新テスト用タイトル"),
		}

		startTime := time.Now()
		updateResp, err := client.Client.UpdateSingleTenantSettingsWithResponse(ctx, partialUpdateParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("部分更新APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "部分更新")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, updateResp.StatusCode(), "部分更新")

		t.Log("部分更新成功")
	})

	t.Run("無効なURL形式エラー", func(t *testing.T) {
		// 無効なURL形式で更新を試行
		invalidUpdateParam := authapi.UpdateSingleTenantSettingsParam{
			CustomizePageUrl: common.StringPtr("invalid-url-format"),
		}

		startTime := time.Now()
		updateResp, err := client.Client.UpdateSingleTenantSettingsWithResponse(ctx, invalidUpdateParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("無効なURL更新APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "無効なURL形式")

		// エラーステータスコードをチェック
		if updateResp.StatusCode() == 200 {
			t.Error("無効なURL形式で成功レスポンスが返されました")
		}

		t.Log("無効なURL形式エラー確認成功")
	})
}

// testCloudFormationTemplate はCloudFormationテンプレートのテストを実行します
func testCloudFormationTemplate(t *testing.T, client *common.ClientWrapper, testData *common.SingleTenantManagementTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Run("CloudFormationテンプレート取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetSingleTenantSettingsCloudFormationTemplateWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("CloudFormationテンプレート取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "CloudFormationテンプレート取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "CloudFormationテンプレート取得")

		if resp.JSON200 != nil {
			t.Logf("CloudFormationテンプレート取得成功")
		}
	})

	t.Run("テンプレート形式確認", func(t *testing.T) {
		resp, err := client.Client.GetSingleTenantSettingsCloudFormationTemplateWithResponse(ctx)
		if err != nil {
			t.Fatalf("CloudFormationテンプレート取得APIの呼び出しに失敗: %v", err)
		}

		if resp.JSON200 != nil && resp.JSON200.Template != nil {
			// JSONとして解析可能かチェック
			var template interface{}
			if err := json.Unmarshal([]byte(*resp.JSON200.Template), &template); err != nil {
				t.Errorf("CloudFormationテンプレートがJSON形式ではありません: %v", err)
			} else {
				t.Log("CloudFormationテンプレートのJSON形式確認成功")
			}
		}
	})
}

// testDDLTemplate はDDLテンプレートのテストを実行します
func testDDLTemplate(t *testing.T, client *common.ClientWrapper, testData *common.SingleTenantManagementTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Run("DDLテンプレート取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetSingleTenantSettingsDdlTemplateWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("DDLテンプレート取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "DDLテンプレート取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "DDLテンプレート取得")

		if resp.JSON200 != nil {
			t.Logf("DDLテンプレート取得成功")
		}
	})

	t.Run("DDL形式確認", func(t *testing.T) {
		resp, err := client.Client.GetSingleTenantSettingsDdlTemplateWithResponse(ctx)
		if err != nil {
			t.Fatalf("DDLテンプレート取得APIの呼び出しに失敗: %v", err)
		}

		if resp.JSON200 != nil && resp.JSON200.Ddl != nil {
			ddl := *resp.JSON200.Ddl
			// DDLの基本的な形式をチェック（CREATE TABLEやALTER TABLEが含まれているか）
			if strings.Contains(strings.ToUpper(ddl), "CREATE TABLE") || 
			   strings.Contains(strings.ToUpper(ddl), "ALTER TABLE") {
				t.Log("DDLテンプレートの形式確認成功")
			} else {
				t.Error("DDLテンプレートに期待されるSQL文が含まれていません")
			}
		}
	})
}

// testAWSIntegration はAWS連携設定のテストを実行します
func testAWSIntegration(t *testing.T, client *common.ClientWrapper, testData *common.SingleTenantManagementTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	t.Run("ExternalID取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetSingleTenantSettingsExternalIdWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("ExternalID取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "ExternalID取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "ExternalID取得")

		if resp.JSON200 != nil && resp.JSON200.ExternalId != nil {
			externalID := *resp.JSON200.ExternalId
			// ExternalIDの形式をチェック（空でないこと）
			assert.AssertNotEmpty(externalID, "ExternalID")
			t.Logf("ExternalID取得成功: %s", externalID)
		}
	})

	t.Run("IAMロール設定", func(t *testing.T) {
		// IAMロール設定パラメータを準備
		roleParam := authapi.UpdateSingleTenantSettingsIamRoleParam{
			RoleArn:    testData.AWSIntegration.IAMRole.Params.RoleArn,
			ExternalId: testData.AWSIntegration.IAMRole.Params.ExternalID,
		}

		// IAMロールを設定
		startTime := time.Now()
		updateResp, err := client.Client.UpdateSingleTenantSettingsIamRoleWithResponse(ctx, roleParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("IAMロール設定APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "IAMロール設定")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, updateResp.StatusCode(), "IAMロール設定")

		t.Logf("IAMロール設定成功: %s", testData.AWSIntegration.IAMRole.Params.RoleArn)
	})

	t.Run("IAMロール取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetSingleTenantSettingsIamRoleWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("IAMロール取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "IAMロール取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "IAMロール取得")

		if resp.JSON200 != nil {
			// 設定されたIAMロールを確認
			if resp.JSON200.RoleArn != nil {
				assert.AssertEquals(testData.AWSIntegration.IAMRole.Params.RoleArn, *resp.JSON200.RoleArn, "IAMロールARN")
			}
			if resp.JSON200.ExternalId != nil {
				assert.AssertEquals(testData.AWSIntegration.IAMRole.Params.ExternalID, *resp.JSON200.ExternalId, "ExternalID")
			}
		}

		t.Log("IAMロール取得確認成功")
	})

	t.Run("無効なロールARNエラー", func(t *testing.T) {
		// 無効なロールARNで設定を試行
		invalidRoleParam := authapi.UpdateSingleTenantSettingsIamRoleParam{
			RoleArn:    "invalid-role-arn",
			ExternalId: testData.AWSIntegration.IAMRole.Params.ExternalID,
		}

		startTime := time.Now()
		updateResp, err := client.Client.UpdateSingleTenantSettingsIamRoleWithResponse(ctx, invalidRoleParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("無効なIAMロール設定APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "無効なロールARN")

		// エラーステータスコードをチェック
		if updateResp.StatusCode() == 200 {
			t.Error("無効なロールARNで成功レスポンスが返されました")
		}

		t.Log("無効なロールARNエラー確認成功")
	})
}
