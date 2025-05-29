package stories

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

// TestStory10_IntegrationTest は統合テストストーリーのE2Eテストです
func TestStory10_IntegrationTest(t *testing.T) {
	// テストデータローダーを初期化
	loader := common.NewTestDataLoader()
	testData, err := loader.LoadIntegrationTestData()
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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := cleanup.CleanupByStory(ctx, common.StoryIntegrationTest); err != nil {
			t.Logf("クリーンアップに失敗: %v", err)
		}
	}()

	// サブテストを順次実行
	t.Run("完全なSaaSセットアップフロー", func(t *testing.T) {
		testFullSaaSSetupFlow(t, client, testData, assert, cleanup)
	})

	t.Run("エンドツーエンドシナリオ", func(t *testing.T) {
		testEndToEndScenarios(t, client, testData, assert, cleanup)
	})

	t.Run("パフォーマンステスト", func(t *testing.T) {
		testPerformanceScenarios(t, client, testData, assert)
	})

	t.Run("障害復旧テスト", func(t *testing.T) {
		testFailureRecoveryScenarios(t, client, testData, assert)
	})
}

// testFullSaaSSetupFlow は完全なSaaSセットアップフローのテストを実行します
func testFullSaaSSetupFlow(t *testing.T, client *common.ClientWrapper, testData *common.IntegrationTestData, assert *common.AssertionHelper, cleanup *common.CleanupManager) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var (
		testEnvID    string
		testRoleName string
		testTenantID string
		testUserID   string
	)

	t.Run("1. 基本情報設定", func(t *testing.T) {
		// 基本情報を設定
		updateParam := authapi.UpdateBasicInfoParam{
			DomainName:        "integration-test.saasus.example.com",
			FromEmailAddress:  "noreply@integration-test.example.com",
			ReplyEmailAddress: common.StringPtr("support@integration-test.example.com"),
		}

		startTime := time.Now()
		resp, err := client.Client.UpdateBasicInfoWithResponse(ctx, updateParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("基本情報設定APIの呼び出しに失敗: %v", err)
		}

		assert.AssertResponseTime(duration, 15*time.Second, "基本情報設定")
		assert.AssertStatusCode(200, resp.StatusCode(), "基本情報設定")

		t.Log("基本情報設定完了")
	})

	t.Run("2. 環境作成", func(t *testing.T) {
		// テスト用環境を作成
		createParam := authapi.CreateEnvParam{
			Id:          "integration-test-env",
			Name:        "統合テスト環境",
			DisplayName: common.StringPtr("統合テスト用環境"),
		}

		startTime := time.Now()
		resp, err := client.Client.CreateEnvWithResponse(ctx, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("環境作成APIの呼び出しに失敗: %v", err)
		}

		assert.AssertResponseTime(duration, 15*time.Second, "環境作成")
		assert.AssertStatusCode(201, resp.StatusCode(), "環境作成")

		if resp.JSON201 != nil {
			testEnvID = resp.JSON201.Id
			client.CreateTestResource(
				common.ResourceTypeEnv,
				testEnvID,
				"統合テスト環境",
				common.StoryIntegrationTest,
				map[string]interface{}{"name": "統合テスト環境"},
			)
		}

		t.Logf("環境作成完了: ID=%s", testEnvID)
	})

	t.Run("3. ロール作成", func(t *testing.T) {
		// テスト用ロールを作成
		testRoleName = "integration-test-role"
		createParam := authapi.CreateRoleParam{
			RoleName:    testRoleName,
			DisplayName: "統合テスト用ロール",
		}

		startTime := time.Now()
		resp, err := client.Client.CreateRoleWithResponse(ctx, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("ロール作成APIの呼び出しに失敗: %v", err)
		}

		assert.AssertResponseTime(duration, 15*time.Second, "ロール作成")
		assert.AssertStatusCode(201, resp.StatusCode(), "ロール作成")

		client.CreateTestResource(
			common.ResourceTypeRole,
			testRoleName,
			"統合テスト用ロール",
			common.StoryIntegrationTest,
			map[string]interface{}{"display_name": "統合テスト用ロール"},
		)

		t.Logf("ロール作成完了: %s", testRoleName)
	})

	t.Run("4. テナント作成", func(t *testing.T) {
		// テスト用テナントを作成
		createParam := authapi.CreateTenantParam{
			Name:                 "統合テスト用テナント",
			BackOfficeStaffEmail: "admin@integration-test.example.com",
			Attributes:           &map[string]interface{}{},
		}

		startTime := time.Now()
		resp, err := client.Client.CreateTenantWithResponse(ctx, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("テナント作成APIの呼び出しに失敗: %v", err)
		}

		assert.AssertResponseTime(duration, 15*time.Second, "テナント作成")
		assert.AssertStatusCode(201, resp.StatusCode(), "テナント作成")

		if resp.JSON201 != nil {
			testTenantID = resp.JSON201.Id
			client.CreateTestResource(
				common.ResourceTypeTenant,
				testTenantID,
				"統合テスト用テナント",
				common.StoryIntegrationTest,
				map[string]interface{}{
					"name":                     "統合テスト用テナント",
					"back_office_staff_email": "admin@integration-test.example.com",
				},
			)
		}

		t.Logf("テナント作成完了: ID=%s", testTenantID)
	})

	t.Run("5. ユーザー作成", func(t *testing.T) {
		// テスト用ユーザーを作成
		createParam := authapi.CreateSaasUserParam{
			Email:      "integration-test-user@example.com",
			Attributes: &map[string]interface{}{},
		}

		startTime := time.Now()
		resp, err := client.Client.CreateSaasUserWithResponse(ctx, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("ユーザー作成APIの呼び出しに失敗: %v", err)
		}

		assert.AssertResponseTime(duration, 15*time.Second, "ユーザー作成")
		assert.AssertStatusCode(201, resp.StatusCode(), "ユーザー作成")

		if resp.JSON201 != nil {
			testUserID = resp.JSON201.Id
			client.CreateTestResource(
				common.ResourceTypeSaasUser,
				testUserID,
				"integration-test-user@example.com",
				common.StoryIntegrationTest,
				map[string]interface{}{"email": "integration-test-user@example.com"},
			)
		}

		t.Logf("ユーザー作成完了: ID=%s", testUserID)
	})

	t.Run("6. テナントユーザー追加", func(t *testing.T) {
		if testTenantID == "" || testUserID == "" {
			t.Skip("テナントIDまたはユーザーIDが取得できていないためスキップ")
		}

		// ユーザーをテナントに追加
		createParam := authapi.CreateTenantUserParam{
			UserId: testUserID,
			Roles:  &[]string{},
			Envs:   &[]string{testEnvID},
		}

		startTime := time.Now()
		resp, err := client.Client.CreateTenantUserWithResponse(ctx, testTenantID, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("テナントユーザー追加APIの呼び出しに失敗: %v", err)
		}

		assert.AssertResponseTime(duration, 15*time.Second, "テナントユーザー追加")
		assert.AssertStatusCode(201, resp.StatusCode(), "テナントユーザー追加")

		client.CreateTestResource(
			common.ResourceTypeTenantUser,
			testUserID,
			"統合テスト用テナントユーザー",
			common.StoryIntegrationTest,
			map[string]interface{}{
				"tenant_id": testTenantID,
				"user_id":   testUserID,
			},
		)

		t.Log("テナントユーザー追加完了")
	})

	t.Run("7. ロール割り当て", func(t *testing.T) {
		if testTenantID == "" || testUserID == "" || testEnvID == "" || testRoleName == "" {
			t.Skip("必要なIDが取得できていないためスキップ")
		}

		// ユーザーにロールを割り当て
		startTime := time.Now()
		resp, err := client.Client.AttachRoleWithResponse(ctx, testTenantID, testUserID, testEnvID, testRoleName)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("ロール割り当てAPIの呼び出しに失敗: %v", err)
		}

		assert.AssertResponseTime(duration, 15*time.Second, "ロール割り当て")
		assert.AssertStatusCode(200, resp.StatusCode(), "ロール割り当て")

		t.Log("ロール割り当て完了")
	})

	t.Run("8. 招待作成", func(t *testing.T) {
		if testTenantID == "" {
			t.Skip("テナントIDが取得できていないためスキップ")
		}

		// 新規ユーザーを招待
		createParam := authapi.CreateTenantInvitationParam{
			Email: "invited-user@integration-test.example.com",
			Roles: &[]string{testRoleName},
			Envs:  &[]string{testEnvID},
		}

		startTime := time.Now()
		resp, err := client.Client.CreateTenantInvitationWithResponse(ctx, testTenantID, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("招待作成APIの呼び出しに失敗: %v", err)
		}

		assert.AssertResponseTime(duration, 15*time.Second, "招待作成")
		assert.AssertStatusCode(201, resp.StatusCode(), "招待作成")

		if resp.JSON201 != nil {
			invitationID := resp.JSON201.Id
			client.CreateTestResource(
				common.ResourceTypeInvitation,
				invitationID,
				"invited-user@integration-test.example.com",
				common.StoryIntegrationTest,
				map[string]interface{}{
					"tenant_id": testTenantID,
					"email":     "invited-user@integration-test.example.com",
				},
			)
		}

		t.Log("招待作成完了")
	})

	t.Log("完全なSaaSセットアップフロー完了")
}

// testEndToEndScenarios はエンドツーエンドシナリオのテストを実行します
func testEndToEndScenarios(t *testing.T, client *common.ClientWrapper, testData *common.IntegrationTestData, assert *common.AssertionHelper, cleanup *common.CleanupManager) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	t.Run("マルチテナント分離確認", func(t *testing.T) {
		// テナントAを作成
		tenantAParam := authapi.CreateTenantParam{
			Name:                 "テナントA",
			BackOfficeStaffEmail: "admin-a@example.com",
			Attributes:           &map[string]interface{}{},
		}

		tenantAResp, err := client.Client.CreateTenantWithResponse(ctx, tenantAParam)
		if err != nil || tenantAResp.JSON201 == nil {
			t.Fatalf("テナントA作成に失敗: %v", err)
		}
		tenantAID := tenantAResp.JSON201.Id

		// テナントBを作成
		tenantBParam := authapi.CreateTenantParam{
			Name:                 "テナントB",
			BackOfficeStaffEmail: "admin-b@example.com",
			Attributes:           &map[string]interface{}{},
		}

		tenantBResp, err := client.Client.CreateTenantWithResponse(ctx, tenantBParam)
		if err != nil || tenantBResp.JSON201 == nil {
			t.Fatalf("テナントB作成に失敗: %v", err)
		}
		tenantBID := tenantBResp.JSON201.Id

		// リソース追跡に追加
		client.CreateTestResource(common.ResourceTypeTenant, tenantAID, "テナントA", common.StoryIntegrationTest, nil)
		client.CreateTestResource(common.ResourceTypeTenant, tenantBID, "テナントB", common.StoryIntegrationTest, nil)

		// テナントAにユーザー1を作成
		userAParam := authapi.CreateSaasUserParam{
			Email:      "user-a@example.com",
			Attributes: &map[string]interface{}{},
		}

		userAResp, err := client.Client.CreateSaasUserWithResponse(ctx, userAParam)
		if err != nil || userAResp.JSON201 == nil {
			t.Fatalf("ユーザーA作成に失敗: %v", err)
		}
		userAID := userAResp.JSON201.Id

		// テナントBにユーザー2を作成
		userBParam := authapi.CreateSaasUserParam{
			Email:      "user-b@example.com",
			Attributes: &map[string]interface{}{},
		}

		userBResp, err := client.Client.CreateSaasUserWithResponse(ctx, userBParam)
		if err != nil || userBResp.JSON201 == nil {
			t.Fatalf("ユーザーB作成に失敗: %v", err)
		}
		userBID := userBResp.JSON201.Id

		// リソース追跡に追加
		client.CreateTestResource(common.ResourceTypeSaasUser, userAID, "user-a@example.com", common.StoryIntegrationTest, nil)
		client.CreateTestResource(common.ResourceTypeSaasUser, userBID, "user-b@example.com", common.StoryIntegrationTest, nil)

		// ユーザーAをテナントAに追加
		tenantUserAParam := authapi.CreateTenantUserParam{
			UserId: userAID,
			Roles:  &[]string{},
			Envs:   &[]string{},
		}

		_, err = client.Client.CreateTenantUserWithResponse(ctx, tenantAID, tenantUserAParam)
		if err != nil {
			t.Fatalf("テナントユーザーA追加に失敗: %v", err)
		}

		// ユーザーBをテナントBに追加
		tenantUserBParam := authapi.CreateTenantUserParam{
			UserId: userBID,
			Roles:  &[]string{},
			Envs:   &[]string{},
		}

		_, err = client.Client.CreateTenantUserWithResponse(ctx, tenantBID, tenantUserBParam)
		if err != nil {
			t.Fatalf("テナントユーザーB追加に失敗: %v", err)
		}

		// テナント分離確認: テナントAのユーザー一覧にユーザーBが含まれないことを確認
		tenantAUsersResp, err := client.Client.GetTenantUsersWithResponse(ctx, tenantAID)
		if err != nil {
			t.Fatalf("テナントAユーザー一覧取得に失敗: %v", err)
		}

		if tenantAUsersResp.JSON200 != nil {
			for _, user := range tenantAUsersResp.JSON200.Users {
				if user.Id == userBID {
					t.Error("テナント分離が正しく機能していません: テナントAにユーザーBが含まれています")
				}
			}
		}

		t.Log("マルチテナント分離確認完了")
	})

	t.Run("認証フロー統合テスト", func(t *testing.T) {
		// 認証情報を保存
		createParam := authapi.CreateAuthCredentialsJSONRequestBody{
			IdToken:     "integration-test-id-token",
			AccessToken: "integration-test-access-token",
		}

		createResp, err := client.Client.CreateAuthCredentialsWithResponse(ctx, createParam)
		if err != nil {
			t.Fatalf("認証情報保存に失敗: %v", err)
		}

		assert.AssertStatusCode(201, createResp.StatusCode(), "認証情報保存")

		if createResp.JSON201 != nil {
			tempCode := createResp.JSON201.Code

			// 一時コードで認証情報を取得
			params := &authapi.GetAuthCredentialsParams{
				Code: &tempCode,
			}

			getResp, err := client.Client.GetAuthCredentialsWithResponse(ctx, params)
			if err != nil {
				t.Fatalf("認証情報取得に失敗: %v", err)
			}

			assert.AssertStatusCode(200, getResp.StatusCode(), "認証情報取得")

			t.Log("認証フロー統合テスト完了")
		}
	})
}

// testPerformanceScenarios はパフォーマンステストを実行します
func testPerformanceScenarios(t *testing.T, client *common.ClientWrapper, testData *common.IntegrationTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	t.Run("同時ユーザー作成負荷テスト", func(t *testing.T) {
		concurrentUsers := testData.PerformanceTests.LoadTest.Params.ConcurrentUsers
		var wg sync.WaitGroup
		var successCount int32
		var errorCount int32

		startTime := time.Now()

		for i := 0; i < concurrentUsers; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				createParam := authapi.CreateSaasUserParam{
					Email:      fmt.Sprintf("load-test-user-%d@example.com", index),
					Attributes: &map[string]interface{}{},
				}

				resp, err := client.Client.CreateSaasUserWithResponse(ctx, createParam)
				if err != nil || resp.StatusCode() != 201 {
					errorCount++
				} else {
					successCount++
					if resp.JSON201 != nil {
						// テスト用リソースとして追跡
						client.CreateTestResource(
							common.ResourceTypeSaasUser,
							resp.JSON201.Id,
							fmt.Sprintf("load-test-user-%d@example.com", index),
							common.StoryIntegrationTest,
							nil,
						)
					}
				}
			}(i)
		}

		wg.Wait()
		duration := time.Since(startTime)

		// パフォーマンス結果を評価
		assert.AssertResponseTime(duration, 60*time.Second, "同時ユーザー作成負荷テスト")

		successRate := float64(successCount) / float64(concurrentUsers) * 100
		if successRate < 95.0 {
			t.Errorf("成功率が低すぎます: %.2f%% (期待値: 95%%以上)", successRate)
		}

		t.Logf("負荷テスト完了: 成功=%d, エラー=%d, 成功率=%.2f%%, 実行時間=%v",
			successCount, errorCount, successRate, duration)
	})

	t.Run("認証情報取得負荷テスト", func(t *testing.T) {
		// 事前に認証情報を保存
		createParam := authapi.CreateAuthCredentialsJSONRequestBody{
			IdToken:     "load-test-id-token",
			AccessToken: "load-test-access-token",
		}

		createResp, err := client.Client.CreateAuthCredentialsWithResponse(ctx, createParam)
		if err != nil || createResp.JSON201 == nil {
			t.Skip("認証情報の事前保存に失敗")
			return
		}

		tempCode := createResp.JSON201.Code
		concurrentRequests := 20
		var wg sync.WaitGroup
		var successCount int32
		var errorCount int32

		startTime := time.Now()

		for i := 0; i < concurrentRequests; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				params := &authapi.GetAuthCredentialsParams{
					Code: &tempCode,
				}

				resp, err := client.Client.GetAuthCredentialsWithResponse(ctx, params)
				if err != nil || (resp.StatusCode() != 200 && resp.StatusCode() != 404) {
					errorCount++
				} else {
					successCount++
				}
			}()
		}

		wg.Wait()
		duration := time.Since(startTime)

		assert.AssertResponseTime(duration, 30*time.Second, "認証情報取得負荷テスト")

		t.Logf("認証負荷テスト完了: 成功=%d, エラー=%d, 実行時間=%v",
			successCount, errorCount, duration)
	})
}

// testFailureRecoveryScenarios は障害復旧テストを実行します
func testFailureRecoveryScenarios(t *testing.T, client *common.ClientWrapper, testData *common.IntegrationTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	t.Run("データ整合性確認", func(t *testing.T) {
		// ユーザーを作成
		createParam := authapi.CreateSaasUserParam{
			Email:      "consistency-test@example.com",
			Attributes: &map[string]interface{}{},
		}

		createResp, err := client.Client.CreateSaasUserWithResponse(ctx, createParam)
		if err != nil || createResp.JSON201 == nil {
			t.Fatalf("ユーザー作成に失敗: %v", err)
		}

		userID := createResp.JSON201.Id
		client.CreateTestResource(common.ResourceTypeSaasUser, userID, "consistency-test@example.com", common.StoryIntegrationTest, nil)

		// 作成直後にユーザーを取得して整合性を確認
		getResp, err := client.Client.GetSaasUserWithResponse(ctx, userID)
		if err != nil {
			t.Fatalf("ユーザー取得に失敗: %v", err)
		}

		assert.AssertStatusCode(200, getResp.StatusCode(), "ユーザー取得")

		if getResp.JSON200 != nil {
			assert.AssertEquals(userID, getResp.JSON200.Id, "ユーザーID")
			assert.AssertEquals("consistency-test@example.com", getResp.JSON200.Email, "メールアドレス")
		}

		t.Log("データ整合性確認完了")
	})

	t.Run("トランザクション整合性確認", func(t *testing.T) {
		// テナントを作成
		tenantParam := authapi.CreateTenantParam{
			Name:                 "トランザクションテスト用テナント",
			BackOfficeStaffEmail: "transaction-test@example.com",
			Attributes:           &map[string]interface{}{},
		}

		tenantResp, err := client.Client.CreateTenantWithResponse(ctx, tenantParam)
		if err != nil || tenantResp.JSON201 == nil {
			t.Fatalf("テナント作成に失敗: %v", err)
		}

		tenantID := tenantResp.JSON201.Id
		client.CreateTestResource(common.ResourceTypeTenant, tenantID, "トランザクションテスト用テナント", common.StoryIntegrationTest, nil)

		// ユーザーを作成
		userParam := authapi.CreateSaasUserParam{
			Email:      "transaction-user@example.com",
			Attributes: &map[string]interface{}{},
		}

		userResp, err := client.Client.CreateSaasUserWithResponse(ctx, userParam)
		if err != nil || userResp.JSON201 == nil {
			t.Fatalf("ユーザー作成に失敗: %v", err)
		}

		userID := userResp.JSON201.Id
		client.CreateTestResource(common.ResourceTypeSaasUser, userID, "transaction-user@example.com", common.StoryIntegrationTest, nil)

		// ユーザーをテナントに追加
		tenantUserParam := authapi.CreateTenantUserParam{
			UserId: userID,
			Roles:  &[]string{},
			Envs:   &[]string{},
		}

		tenantUserResp, err := client.Client.CreateTenantUserWithResponse(ctx, tenantID, tenantUserParam)
		if err != nil {
			t.Fatalf("テナントユーザー追加に失敗: %v", err)
		}

		assert.AssertStatusCode(201, tenantUserResp.StatusCode(), "テナントユーザー追加")

		// テナントユーザー一覧でユーザーが存在することを確認
		tenantUsersResp, err := client.Client.GetTenantUsersWithResponse(ctx, tenantID)
		if err != nil {
			t.Fatalf("テナントユーザー一覧取得に失敗: %v", err)
		}

		found := false
		if tenantUsersResp.JSON200 != nil {
			for _, user := range tenantUsersResp.JSON200.Users {
				if user.Id == userID {
					found = true
					break
				}
			}
		}

		if !found {
			t.Error("トランザクション整合性エラー: 追加したユーザーがテナントユーザー一覧に存在しません")
		}

		t.Log("トランザクション整合性確認完了")
	})
}
