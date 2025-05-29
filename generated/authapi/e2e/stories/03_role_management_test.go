package stories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

// TestStory03_RoleManagement はロール管理ストーリーのE2Eテストです
func TestStory03_RoleManagement(t *testing.T) {
	// テストデータローダーを初期化
	loader := common.NewTestDataLoader()
	testData, err := loader.LoadRoleManagementData()
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
		if err := cleanup.CleanupByStory(ctx, common.StoryRoleManagement); err != nil {
			t.Logf("クリーンアップに失敗: %v", err)
		}
	}()

	// サブテストを順次実行
	t.Run("ロール一覧取得", func(t *testing.T) {
		testGetRoles(t, client, assert)
	})

	t.Run("ロール作成", func(t *testing.T) {
		testCreateRoles(t, client, testData, assert)
	})

	t.Run("ロール作成エラーケース", func(t *testing.T) {
		testCreateRoleErrorCases(t, client, testData, assert)
	})

	t.Run("作成されたロールの確認", func(t *testing.T) {
		testVerifyCreatedRoles(t, client, testData, assert)
	})

	t.Run("ロール削除", func(t *testing.T) {
		testDeleteRoles(t, client, testData, assert)
	})

	t.Run("ロール削除エラーケース", func(t *testing.T) {
		testDeleteRoleErrorCases(t, client, assert)
	})

	t.Run("パフォーマンステスト", func(t *testing.T) {
		testRolePerformance(t, client, assert)
	})
}

// testGetRoles はロール一覧取得のテストを実行します
func testGetRoles(t *testing.T, client *common.ClientWrapper, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("ロール一覧取得_正常系", func(t *testing.T) {
		startTime := time.Now()
		resp, err := client.Client.GetRolesWithResponse(ctx)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("ロール一覧取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "ロール一覧取得")

		// レスポンス内容をアサート
		assert.AssertRolesResponse(resp)

		t.Logf("ロール一覧取得成功: ロール数=%d", len(resp.JSON200.Roles))
	})
}

// testCreateRoles はロール作成のテストを実行します
func testCreateRoles(t *testing.T, client *common.ClientWrapper, testData *common.RoleManagementTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// テスト用ロールを順次作成
	for _, roleParam := range testData.Roles.Create.Params {
		t.Run(fmt.Sprintf("ロール作成_%s", roleParam.RoleName), func(t *testing.T) {
			// 作成パラメータを準備
			createParam := authapi.CreateRoleParam{
				RoleName:    roleParam.RoleName,
				DisplayName: roleParam.DisplayName,
			}

			// ロールを作成
			startTime := time.Now()
			createResp, err := client.Client.CreateRoleWithResponse(ctx, createParam)
			duration := time.Since(startTime)

			if err != nil {
				t.Fatalf("ロール作成APIの呼び出しに失敗: %v", err)
			}

			// レスポンス時間をチェック
			assert.AssertResponseTime(duration, 15*time.Second, "ロール作成")

			// ステータスコードをチェック
			assert.AssertStatusCode(201, createResp.StatusCode(), "ロール作成")

			// リソース追跡に追加
			client.CreateTestResource(
				common.ResourceTypeRole,
				roleParam.RoleName,
				roleParam.DisplayName,
				common.StoryRoleManagement,
				map[string]interface{}{
					"display_name": roleParam.DisplayName,
				},
			)

			t.Logf("ロール作成成功: %s (%s)", roleParam.RoleName, roleParam.DisplayName)

			// 作成間隔を設ける（API負荷軽減）
			time.Sleep(500 * time.Millisecond)
		})
	}
}

// testCreateRoleErrorCases はロール作成のエラーケースをテストします
func testCreateRoleErrorCases(t *testing.T, client *common.ClientWrapper, testData *common.RoleManagementTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	t.Run("重複ロール名エラー", func(t *testing.T) {
		// 既存のロール名で作成を試行
		if len(testData.Roles.Create.Params) > 0 {
			existingRole := testData.Roles.Create.Params[0]
			createParam := authapi.CreateRoleParam{
				RoleName:    existingRole.RoleName,
				DisplayName: "重複テスト",
			}

			resp, err := client.Client.CreateRoleWithResponse(ctx, createParam)
			if err != nil {
				t.Fatalf("ロール作成APIの呼び出しに失敗: %v", err)
			}

			// 409 Conflictエラーを期待
			if resp.StatusCode() != 409 {
				t.Errorf("重複ロール作成で期待されるステータスコード 409, 実際 %d", resp.StatusCode())
			} else {
				t.Log("重複ロール名エラーケース確認成功")
			}
		}
	})

	t.Run("必須項目不足エラー", func(t *testing.T) {
		// ロール名が空の場合
		createParam := authapi.CreateRoleParam{
			RoleName:    "",
			DisplayName: "表示名のみ",
		}

		resp, err := client.Client.CreateRoleWithResponse(ctx, createParam)
		if err != nil {
			t.Fatalf("ロール作成APIの呼び出しに失敗: %v", err)
		}

		// 400 Bad Requestエラーを期待
		if resp.StatusCode() != 400 {
			t.Errorf("必須項目不足で期待されるステータスコード 400, 実際 %d", resp.StatusCode())
		} else {
			t.Log("必須項目不足エラーケース確認成功")
		}
	})

	t.Run("不正文字エラー", func(t *testing.T) {
		// 不正な文字を含むロール名
		createParam := authapi.CreateRoleParam{
			RoleName:    "invalid@role#name",
			DisplayName: "不正文字テスト",
		}

		resp, err := client.Client.CreateRoleWithResponse(ctx, createParam)
		if err != nil {
			t.Fatalf("ロール作成APIの呼び出しに失敗: %v", err)
		}

		// 400 Bad Requestエラーを期待
		if resp.StatusCode() != 400 {
			t.Errorf("不正文字で期待されるステータスコード 400, 実際 %d", resp.StatusCode())
		} else {
			t.Log("不正文字エラーケース確認成功")
		}
	})
}

// testVerifyCreatedRoles は作成されたロールの確認テストを実行します
func testVerifyCreatedRoles(t *testing.T, client *common.ClientWrapper, testData *common.RoleManagementTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	t.Run("作成ロール一覧確認", func(t *testing.T) {
		// ロール一覧を取得
		resp, err := client.Client.GetRolesWithResponse(ctx)
		if err != nil {
			t.Fatalf("ロール一覧取得APIの呼び出しに失敗: %v", err)
		}

		assert.AssertRolesResponse(resp)

		// 作成したロールが含まれているかチェック
		createdRoles := make(map[string]bool)
		for _, role := range resp.JSON200.Roles {
			createdRoles[role.RoleName] = true
		}

		for _, roleParam := range testData.Roles.Create.Params {
			if !createdRoles[roleParam.RoleName] {
				t.Errorf("作成したロール %s が一覧に含まれていません", roleParam.RoleName)
			} else {
				t.Logf("作成ロール確認成功: %s", roleParam.RoleName)
			}
		}
	})

	// 個別のロール詳細確認
	for _, roleParam := range testData.Roles.Create.Params {
		t.Run(fmt.Sprintf("ロール詳細確認_%s", roleParam.RoleName), func(t *testing.T) {
			// ロール詳細を取得（GetRoleは存在しない場合があるため、一覧から検索）
			resp, err := client.Client.GetRolesWithResponse(ctx)
			if err != nil {
				t.Fatalf("ロール一覧取得APIの呼び出しに失敗: %v", err)
			}

			var foundRole *authapi.Role
			for _, role := range resp.JSON200.Roles {
				if role.RoleName == roleParam.RoleName {
					foundRole = &role
					break
				}
			}

			if foundRole == nil {
				t.Errorf("ロール %s が見つかりません", roleParam.RoleName)
				return
			}

			// ロール情報をチェック
			assert.AssertEquals(roleParam.RoleName, foundRole.RoleName, "ロール名")
			assert.AssertEquals(roleParam.DisplayName, foundRole.DisplayName, "表示名")

			t.Logf("ロール詳細確認成功: %s", roleParam.RoleName)
		})
	}
}

// testDeleteRoles はロール削除のテストを実行します
func testDeleteRoles(t *testing.T, client *common.ClientWrapper, testData *common.RoleManagementTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 作成したロールを削除
	for _, roleParam := range testData.Roles.Create.Params {
		t.Run(fmt.Sprintf("ロール削除_%s", roleParam.RoleName), func(t *testing.T) {
			// ロールを削除
			startTime := time.Now()
			deleteResp, err := client.Client.DeleteRoleWithResponse(ctx, roleParam.RoleName)
			duration := time.Since(startTime)

			if err != nil {
				t.Fatalf("ロール削除APIの呼び出しに失敗: %v", err)
			}

			// レスポンス時間をチェック
			assert.AssertResponseTime(duration, 15*time.Second, "ロール削除")

			// ステータスコードをチェック
			assert.AssertResourceDeleted(deleteResp.StatusCode(), "ロール")

			// リソースを削除済みとしてマーク
			client.MarkResourceCleaned(roleParam.RoleName, nil)

			t.Logf("ロール削除成功: %s", roleParam.RoleName)

			// 削除間隔を設ける（API負荷軽減）
			time.Sleep(500 * time.Millisecond)
		})
	}

	t.Run("削除確認", func(t *testing.T) {
		// ロール一覧を取得して削除を確認
		resp, err := client.Client.GetRolesWithResponse(ctx)
		if err != nil {
			t.Fatalf("ロール一覧取得APIの呼び出しに失敗: %v", err)
		}

		// 削除したロールが含まれていないかチェック
		remainingRoles := make(map[string]bool)
		for _, role := range resp.JSON200.Roles {
			remainingRoles[role.RoleName] = true
		}

		for _, roleParam := range testData.Roles.Create.Params {
			if remainingRoles[roleParam.RoleName] {
				t.Errorf("削除したロール %s がまだ一覧に含まれています", roleParam.RoleName)
			} else {
				t.Logf("ロール削除確認成功: %s", roleParam.RoleName)
			}
		}
	})
}

// testDeleteRoleErrorCases はロール削除のエラーケースをテストします
func testDeleteRoleErrorCases(t *testing.T, client *common.ClientWrapper, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("存在しないロール削除エラー", func(t *testing.T) {
		// 存在しないロール名で削除を試行
		resp, err := client.Client.DeleteRoleWithResponse(ctx, "non_existent_role")
		if err != nil {
			t.Fatalf("ロール削除APIの呼び出しに失敗: %v", err)
		}

		// 404 Not Foundエラーを期待
		if resp.StatusCode() != 404 {
			t.Errorf("存在しないロール削除で期待されるステータスコード 404, 実際 %d", resp.StatusCode())
		} else {
			t.Log("存在しないロール削除エラーケース確認成功")
		}
	})
}

// testRolePerformance はロール管理のパフォーマンステストを実行します
func testRolePerformance(t *testing.T, client *common.ClientWrapper, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	t.Run("ロール一覧取得パフォーマンス", func(t *testing.T) {
		// 複数回実行して平均レスポンス時間を測定
		const iterations = 5
		var totalDuration time.Duration

		for i := 0; i < iterations; i++ {
			startTime := time.Now()
			resp, err := client.Client.GetRolesWithResponse(ctx)
			duration := time.Since(startTime)

			if err != nil {
				t.Fatalf("ロール一覧取得APIの呼び出しに失敗 (iteration %d): %v", i+1, err)
			}

			assert.AssertStatusCode(200, resp.StatusCode(), fmt.Sprintf("ロール一覧取得_iteration_%d", i+1))
			totalDuration += duration

			// 間隔を設ける
			time.Sleep(200 * time.Millisecond)
		}

		averageDuration := totalDuration / iterations
		assert.AssertResponseTime(averageDuration, 5*time.Second, "ロール一覧取得平均")

		t.Logf("ロール一覧取得パフォーマンステスト完了: 平均レスポンス時間=%v", averageDuration)
	})

	t.Run("ロール作成削除パフォーマンス", func(t *testing.T) {
		const testRoleCount = 3
		var createdRoles []string

		// テスト用ロールを作成
		for i := 0; i < testRoleCount; i++ {
			roleName := fmt.Sprintf("perf_test_role_%d_%d", time.Now().Unix(), i)
			createParam := authapi.CreateRoleParam{
				RoleName:    roleName,
				DisplayName: fmt.Sprintf("パフォーマンステストロール %d", i+1),
			}

			startTime := time.Now()
			createResp, err := client.Client.CreateRoleWithResponse(ctx, createParam)
			duration := time.Since(startTime)

			if err != nil {
				t.Fatalf("パフォーマンステスト用ロール作成に失敗: %v", err)
			}

			assert.AssertStatusCode(201, createResp.StatusCode(), fmt.Sprintf("パフォーマンステストロール作成_%d", i+1))
			assert.AssertResponseTime(duration, 10*time.Second, fmt.Sprintf("ロール作成_%d", i+1))

			createdRoles = append(createdRoles, roleName)

			// リソース追跡に追加
			client.CreateTestResource(
				common.ResourceTypeRole,
				roleName,
				createParam.DisplayName,
				common.StoryRoleManagement,
				map[string]interface{}{
					"performance_test": true,
				},
			)

			// 作成間隔を設ける
			time.Sleep(300 * time.Millisecond)
		}

		// 作成したロールを削除
		for i, roleName := range createdRoles {
			startTime := time.Now()
			deleteResp, err := client.Client.DeleteRoleWithResponse(ctx, roleName)
			duration := time.Since(startTime)

			if err != nil {
				t.Fatalf("パフォーマンステスト用ロール削除に失敗: %v", err)
			}

			assert.AssertResourceDeleted(deleteResp.StatusCode(), "パフォーマンステストロール")
			assert.AssertResponseTime(duration, 10*time.Second, fmt.Sprintf("ロール削除_%d", i+1))

			// リソースを削除済みとしてマーク
			client.MarkResourceCleaned(roleName, nil)

			// 削除間隔を設ける
			time.Sleep(300 * time.Millisecond)
		}

		t.Logf("ロール作成削除パフォーマンステスト完了: %d個のロールを処理", testRoleCount)
	})
}