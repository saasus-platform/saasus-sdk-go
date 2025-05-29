package stories

import (
	"context"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

// TestStory04_TenantManagement はテナント管理ストーリーのE2Eテストです
func TestStory04_TenantManagement(t *testing.T) {
	// テストデータローダーを初期化
	loader := common.NewTestDataLoader()
	testData, err := loader.LoadTenantManagementData()
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
		if err := cleanup.CleanupByStory(ctx, common.StoryTenantManagement); err != nil {
			t.Logf("クリーンアップに失敗: %v", err)
		}
	}()

	// サブテストを順次実行
	t.Run("テナント属性管理", func(t *testing.T) {
		testTenantAttributeManagement(t, client, testData, assert, cleanup)
	})

	t.Run("テナント管理", func(t *testing.T) {
		testTenantManagement(t, client, testData, assert, cleanup)
	})

	t.Run("請求情報管理", func(t *testing.T) {
		testBillingInfoManagement(t, client, testData, assert)
	})
}

// testTenantAttributeManagement はテナント属性管理のテストを実行します
func testTenantAttributeManagement(t *testing.T, client *common.ClientWrapper, testData *common.TenantManagementTestData, assert *common.AssertionHelper, cleanup *common.CleanupManager) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var createdAttributes []string

	t.Run("テナント属性一覧取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetTenantAttributesWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("テナント属性一覧取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "テナント属性一覧取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "テナント属性一覧取得")

		t.Logf("テナント属性一覧取得成功: 属性数=%d", len(resp.JSON200.TenantAttributes))
	})

	t.Run("テナント属性作成", func(t *testing.T) {
		for i, param := range testData.TenantAttributes.Create.Params {
			t.Run(param.AttributeName, func(t *testing.T) {
				// 作成パラメータを準備
				createParam := authapi.CreateTenantAttributeParam{
					AttributeName: param.AttributeName,
					DisplayName:   param.DisplayName,
				}

				// テナント属性を作成
				startTime := time.Now()
				createResp, err := client.Client.CreateTenantAttributeWithResponse(ctx, createParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("テナント属性作成APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 15*time.Second, "テナント属性作成")

				// ステータスコードをチェック
				assert.AssertStatusCode(201, createResp.StatusCode(), "テナント属性作成")

				createdAttributes = append(createdAttributes, param.AttributeName)

				// リソース追跡に追加
				client.CreateTestResource(
					common.ResourceTypeTenantAttribute,
					param.AttributeName,
					param.DisplayName,
					common.StoryTenantManagement,
					map[string]interface{}{
						"display_name": param.DisplayName,
					},
				)

				t.Logf("テナント属性作成成功: %s", param.AttributeName)
			})
		}
	})

	t.Run("テナント属性削除", func(t *testing.T) {
		for _, attributeName := range createdAttributes {
			t.Run(attributeName, func(t *testing.T) {
				// テナント属性を削除
				startTime := time.Now()
				deleteResp, err := client.Client.DeleteTenantAttributeWithResponse(ctx, attributeName)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("テナント属性削除APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 15*time.Second, "テナント属性削除")

				// ステータスコードをチェック
				assert.AssertResourceDeleted(deleteResp.StatusCode(), "テナント属性")

				// リソースを削除済みとしてマーク
				client.MarkResourceCleaned(attributeName, nil)

				t.Logf("テナント属性削除成功: %s", attributeName)
			})
		}
	})
}

// testTenantManagement はテナント管理のテストを実行します
func testTenantManagement(t *testing.T, client *common.ClientWrapper, testData *common.TenantManagementTestData, assert *common.AssertionHelper, cleanup *common.CleanupManager) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	var createdTenantIDs []string

	t.Run("テナント一覧取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetTenantsWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("テナント一覧取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "テナント一覧取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "テナント一覧取得")

		t.Logf("テナント一覧取得成功: テナント数=%d", len(resp.JSON200.Tenants))
	})

	t.Run("テナント作成", func(t *testing.T) {
		for i, param := range testData.Tenants.Create.Params {
			t.Run(param.Name, func(t *testing.T) {
				// 作成パラメータを準備
				createParam := authapi.CreateTenantParam{
					Name:                 param.Name,
					BackOfficeStaffEmail: param.BackOfficeStaffEmail,
					Attributes:           &param.Attributes,
				}

				// テナントを作成
				startTime := time.Now()
				createResp, err := client.Client.CreateTenantWithResponse(ctx, createParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("テナント作成APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 30*time.Second, "テナント作成")

				// ステータスコードをチェック
				assert.AssertStatusCode(201, createResp.StatusCode(), "テナント作成")

				if createResp.JSON201 != nil {
					tenantID := createResp.JSON201.Id
					createdTenantIDs = append(createdTenantIDs, tenantID)

					// リソース追跡に追加
					client.CreateTestResource(
						common.ResourceTypeTenant,
						tenantID,
						param.Name,
						common.StoryTenantManagement,
						map[string]interface{}{
							"name":                     param.Name,
							"back_office_staff_email": param.BackOfficeStaffEmail,
							"attributes":              param.Attributes,
						},
					)

					t.Logf("テナント作成成功: ID=%s, Name=%s", tenantID, param.Name)
				}
			})
		}
	})

	t.Run("テナント詳細取得", func(t *testing.T) {
		for _, tenantID := range createdTenantIDs {
			t.Run(tenantID, func(t *testing.T) {
				startTime := time.Now()
				resp, err := client.Client.GetTenantWithResponse(ctx, tenantID)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("テナント詳細取得APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 10*time.Second, "テナント詳細取得")

				// ステータスコードをチェック
				assert.AssertStatusCode(200, resp.StatusCode(), "テナント詳細取得")

				if resp.JSON200 != nil {
					// テナント情報をチェック
					assert.AssertEquals(tenantID, resp.JSON200.Id, "テナントID")
				}

				t.Logf("テナント詳細取得成功: ID=%s", tenantID)
			})
		}
	})

	t.Run("テナント情報更新", func(t *testing.T) {
		if len(createdTenantIDs) > 0 {
			tenantID := createdTenantIDs[0]
			t.Run(tenantID, func(t *testing.T) {
				// 更新パラメータを準備
				updateParam := authapi.UpdateTenantParam{
					Name:                 testData.Tenants.Update.Params.Name,
					BackOfficeStaffEmail: testData.Tenants.Update.Params.BackOfficeStaffEmail,
					Attributes:           testData.Tenants.Update.Params.Attributes,
				}

				// テナント情報を更新
				startTime := time.Now()
				updateResp, err := client.Client.UpdateTenantWithResponse(ctx, tenantID, updateParam)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("テナント情報更新APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 20*time.Second, "テナント情報更新")

				// ステータスコードをチェック
				assert.AssertStatusCode(200, updateResp.StatusCode(), "テナント情報更新")

				// 更新後の情報を確認
				updatedResp, err := client.Client.GetTenantWithResponse(ctx, tenantID)
				if err != nil {
					t.Fatalf("更新後のテナント詳細取得に失敗: %v", err)
				}

				if updatedResp.JSON200 != nil {
					// 更新が反映されているかチェック
					assert.AssertEquals(testData.Tenants.Update.Params.Name, updatedResp.JSON200.Name, "更新後テナント名")
				}

				t.Logf("テナント情報更新成功: ID=%s", tenantID)
			})
		}
	})

	t.Run("テナント削除", func(t *testing.T) {
		for _, tenantID := range createdTenantIDs {
			t.Run(tenantID, func(t *testing.T) {
				// テナントを削除
				startTime := time.Now()
				deleteResp, err := client.Client.DeleteTenantWithResponse(ctx, tenantID)
				duration := time.Since(startTime)

				if err != nil {
					t.Fatalf("テナント削除APIの呼び出しに失敗: %v", err)
				}

				// レスポンス時間をチェック
				assert.AssertResponseTime(duration, 30*time.Second, "テナント削除")

				// ステータスコードをチェック
				assert.AssertResourceDeleted(deleteResp.StatusCode(), "テナント")

				// リソースを削除済みとしてマーク
				client.MarkResourceCleaned(tenantID, nil)

				// 削除確認
				confirmResp, err := client.Client.GetTenantWithResponse(ctx, tenantID)
				if err == nil && confirmResp.StatusCode() != 404 {
					t.Errorf("テナントが削除されていません: ステータスコード %d", confirmResp.StatusCode())
				}

				t.Logf("テナント削除成功: ID=%s", tenantID)
			})
		}
	})
}

// testBillingInfoManagement は請求情報管理のテストを実行します
func testBillingInfoManagement(t *testing.T, client *common.ClientWrapper, testData *common.TenantManagementTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// テスト用テナントを作成
	createParam := authapi.CreateTenantParam{
		Name:                 "請求情報テストテナント",
		BackOfficeStaffEmail: "billing-test@example.com",
		Attributes:           &map[string]interface{}{},
	}

	createResp, err := client.Client.CreateTenantWithResponse(ctx, createParam)
	if err != nil || createResp.JSON201 == nil {
		t.Skip("テスト用テナントの作成に失敗したためスキップ")
		return
	}

	testTenantID := createResp.JSON201.Id

	// テスト終了時にテナントを削除
	defer func() {
		client.Client.DeleteTenantWithResponse(ctx, testTenantID)
	}()

	t.Run("請求情報取得", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetTenantBillingInfoWithResponse(ctx, testTenantID)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("請求情報取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "請求情報取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, resp.StatusCode(), "請求情報取得")

		t.Log("請求情報取得成功")
	})

	t.Run("請求情報更新", func(t *testing.T) {
		// 更新パラメータを準備
		updateParam := authapi.UpdateTenantBillingInfoParam{
			Name: testData.BillingInfo.Update.Params.Name,
			Address: authapi.TenantBillingAddress{
				Street:                 testData.BillingInfo.Update.Params.Street,
				City:                   testData.BillingInfo.Update.Params.City,
				State:                  testData.BillingInfo.Update.Params.State,
				Country:                testData.BillingInfo.Update.Params.Country,
				PostalCode:             testData.BillingInfo.Update.Params.PostalCode,
				AdditionalAddressInfo:  nil,
			},
			InvoiceLanguage: authapi.InvoiceLanguage(testData.BillingInfo.Update.Params.InvoiceLanguage),
		}

		// 請求情報を更新
		startTime := time.Now()
		updateResp, err := client.Client.UpdateTenantBillingInfoWithResponse(ctx, testTenantID, updateParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("請求情報更新APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "請求情報更新")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, updateResp.StatusCode(), "請求情報更新")

		// 更新後の情報を確認
		updatedResp, err := client.Client.GetTenantBillingInfoWithResponse(ctx, testTenantID)
		if err != nil {
			t.Fatalf("更新後の請求情報取得に失敗: %v", err)
		}

		if updatedResp.JSON200 != nil {
			// 更新が反映されているかチェック
			assert.AssertEquals(testData.BillingInfo.Update.Params.Name, updatedResp.JSON200.Name, "更新後請求先名")
		}

		t.Log("請求情報更新成功")
	})
}