package stories

import (
	"context"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi/e2e/common"
)

// TestStory07_AuthenticationFlow は認証フローストーリーのE2Eテストです
func TestStory07_AuthenticationFlow(t *testing.T) {
	// テストデータローダーを初期化
	loader := common.NewTestDataLoader()
	testData, err := loader.LoadAuthenticationFlowData()
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
		if err := cleanup.CleanupByStory(ctx, common.StoryAuthenticationFlow); err != nil {
			t.Logf("クリーンアップに失敗: %v", err)
		}
	}()

	// サブテストを順次実行
	t.Run("一時コード認証フロー", func(t *testing.T) {
		testTempCodeAuthFlow(t, client, testData, assert)
	})

	t.Run("リフレッシュトークン認証フロー", func(t *testing.T) {
		testRefreshTokenAuthFlow(t, client, testData, assert)
	})

	t.Run("一時コード有効期限テスト", func(t *testing.T) {
		testTempCodeExpiration(t, client, testData, assert)
	})

	t.Run("エラーハンドリング", func(t *testing.T) {
		testAuthenticationErrorHandling(t, client, testData, assert)
	})
}

// testTempCodeAuthFlow は一時コード認証フローのテストを実行します
func testTempCodeAuthFlow(t *testing.T, client *common.ClientWrapper, testData *common.AuthenticationFlowTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var tempCode string

	t.Run("認証情報保存", func(t *testing.T) {
		// 認証情報保存パラメータを準備
		createParam := authapi.CreateAuthCredentialsJSONRequestBody{
			IdToken:      testData.TestTokens.ValidIDToken,
			AccessToken:  testData.TestTokens.ValidAccessToken,
			RefreshToken: &testData.TestTokens.ValidRefreshToken,
		}

		// 認証情報を保存
		startTime := time.Now()
		createResp, err := client.Client.CreateAuthCredentialsWithResponse(ctx, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報保存APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "認証情報保存")

		// ステータスコードをチェック
		assert.AssertStatusCode(201, createResp.StatusCode(), "認証情報保存")

		if createResp.JSON201 != nil {
			tempCode = createResp.JSON201.Code
			t.Logf("認証情報保存成功: 一時コード=%s", tempCode)
		} else {
			t.Fatal("一時コードが取得できませんでした")
		}
	})

	t.Run("一時コードで認証情報取得", func(t *testing.T) {
		if tempCode == "" {
			t.Skip("一時コードが取得できていないためスキップ")
		}

		// 一時コードで認証情報を取得
		params := &authapi.GetAuthCredentialsParams{
			Code:     &tempCode,
			AuthFlow: (*authapi.GetAuthCredentialsParamsAuthFlow)(common.StringPtr("tempCodeAuth")),
		}

		startTime := time.Now()
		getResp, err := client.Client.GetAuthCredentialsWithResponse(ctx, params)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "認証情報取得")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, getResp.StatusCode(), "認証情報取得")

		if getResp.JSON200 != nil {
			// 取得した認証情報をチェック
			assert.AssertNotEmpty(getResp.JSON200.IdToken, "IDトークン")
			assert.AssertNotEmpty(getResp.JSON200.AccessToken, "アクセストークン")

			t.Logf("認証情報取得成功: IDトークン=%s...", getResp.JSON200.IdToken[:20])
		}
	})

	t.Run("同じ一時コードで再取得エラー", func(t *testing.T) {
		if tempCode == "" {
			t.Skip("一時コードが取得できていないためスキップ")
		}

		// 同じ一時コードで再度取得を試行
		params := &authapi.GetAuthCredentialsParams{
			Code:     &tempCode,
			AuthFlow: (*authapi.GetAuthCredentialsParamsAuthFlow)(common.StringPtr("tempCodeAuth")),
		}

		startTime := time.Now()
		getResp, err := client.Client.GetAuthCredentialsWithResponse(ctx, params)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "一時コード再使用")

		// エラーステータスコードをチェック（一時コードは一度使用すると無効）
		if getResp.StatusCode() == 200 {
			t.Error("一時コードが再使用できてしまいました")
		}

		t.Log("一時コード再使用エラー確認成功")
	})
}

// testRefreshTokenAuthFlow はリフレッシュトークン認証フローのテストを実行します
func testRefreshTokenAuthFlow(t *testing.T, client *common.ClientWrapper, testData *common.AuthenticationFlowTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var refreshToken string

	t.Run("認証情報保存でリフレッシュトークン取得", func(t *testing.T) {
		// 認証情報保存パラメータを準備
		createParam := authapi.CreateAuthCredentialsJSONRequestBody{
			IdToken:      testData.TestTokens.ValidIDToken,
			AccessToken:  testData.TestTokens.ValidAccessToken,
			RefreshToken: &testData.TestTokens.ValidRefreshToken,
		}

		// 認証情報を保存
		startTime := time.Now()
		createResp, err := client.Client.CreateAuthCredentialsWithResponse(ctx, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報保存APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 15*time.Second, "認証情報保存")

		// ステータスコードをチェック
		assert.AssertStatusCode(201, createResp.StatusCode(), "認証情報保存")

		refreshToken = testData.TestTokens.ValidRefreshToken
		t.Logf("リフレッシュトークン準備完了")
	})

	t.Run("リフレッシュトークンで認証情報取得", func(t *testing.T) {
		if refreshToken == "" {
			t.Skip("リフレッシュトークンが取得できていないためスキップ")
		}

		// リフレッシュトークンで認証情報を取得
		params := &authapi.GetAuthCredentialsParams{
			AuthFlow:     (*authapi.GetAuthCredentialsParamsAuthFlow)(common.StringPtr("refreshTokenAuth")),
			RefreshToken: &refreshToken,
		}

		startTime := time.Now()
		getResp, err := client.Client.GetAuthCredentialsWithResponse(ctx, params)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "リフレッシュトークン認証")

		// ステータスコードをチェック
		assert.AssertStatusCode(200, getResp.StatusCode(), "リフレッシュトークン認証")

		if getResp.JSON200 != nil {
			// 新しい認証情報をチェック
			assert.AssertNotEmpty(getResp.JSON200.IdToken, "新しいIDトークン")
			assert.AssertNotEmpty(getResp.JSON200.AccessToken, "新しいアクセストークン")

			t.Logf("リフレッシュトークン認証成功: 新しいアクセストークン=%s...", getResp.JSON200.AccessToken[:20])
		}
	})

	t.Run("無効なリフレッシュトークンエラー", func(t *testing.T) {
		invalidRefreshToken := "invalid-refresh-token"

		// 無効なリフレッシュトークンで認証情報取得を試行
		params := &authapi.GetAuthCredentialsParams{
			AuthFlow:     (*authapi.GetAuthCredentialsParamsAuthFlow)(common.StringPtr("refreshTokenAuth")),
			RefreshToken: &invalidRefreshToken,
		}

		startTime := time.Now()
		getResp, err := client.Client.GetAuthCredentialsWithResponse(ctx, params)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "無効なリフレッシュトークン")

		// エラーステータスコードをチェック
		if getResp.StatusCode() == 200 {
			t.Error("無効なリフレッシュトークンで成功レスポンスが返されました")
		}

		t.Log("無効なリフレッシュトークンエラー確認成功")
	})
}

// testTempCodeExpiration は一時コードの有効期限テストを実行します
func testTempCodeExpiration(t *testing.T, client *common.ClientWrapper, testData *common.AuthenticationFlowTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var tempCode string

	t.Run("認証情報保存", func(t *testing.T) {
		// 認証情報保存パラメータを準備
		createParam := authapi.CreateAuthCredentialsJSONRequestBody{
			IdToken:      testData.TestTokens.ValidIDToken,
			AccessToken:  testData.TestTokens.ValidAccessToken,
			RefreshToken: &testData.TestTokens.ValidRefreshToken,
		}

		// 認証情報を保存
		createResp, err := client.Client.CreateAuthCredentialsWithResponse(ctx, createParam)
		if err != nil || createResp.JSON201 == nil {
			t.Fatalf("認証情報保存に失敗: %v", err)
		}

		tempCode = createResp.JSON201.Code
		t.Logf("一時コード取得: %s", tempCode)
	})

	t.Run("10秒以内での取得成功", func(t *testing.T) {
		if tempCode == "" {
			t.Skip("一時コードが取得できていないためスキップ")
		}

		// 5秒待機（10秒以内）
		time.Sleep(5 * time.Second)

		// 一時コードで認証情報を取得
		params := &authapi.GetAuthCredentialsParams{
			Code:     &tempCode,
			AuthFlow: (*authapi.GetAuthCredentialsParamsAuthFlow)(common.StringPtr("tempCodeAuth")),
		}

		getResp, err := client.Client.GetAuthCredentialsWithResponse(ctx, params)
		if err != nil {
			t.Fatalf("認証情報取得APIの呼び出しに失敗: %v", err)
		}

		// 成功することを確認
		assert.AssertStatusCode(200, getResp.StatusCode(), "10秒以内での取得")

		t.Log("10秒以内での取得成功確認")
	})

	t.Run("期限切れ後の取得エラー", func(t *testing.T) {
		// 新しい一時コードを生成
		createParam := authapi.CreateAuthCredentialsJSONRequestBody{
			IdToken:      testData.TestTokens.ValidIDToken,
			AccessToken:  testData.TestTokens.ValidAccessToken,
			RefreshToken: &testData.TestTokens.ValidRefreshToken,
		}

		createResp, err := client.Client.CreateAuthCredentialsWithResponse(ctx, createParam)
		if err != nil || createResp.JSON201 == nil {
			t.Skip("新しい一時コードの生成に失敗したためスキップ")
		}

		newTempCode := createResp.JSON201.Code
		t.Logf("新しい一時コード取得: %s", newTempCode)

		// 15秒待機（10秒の有効期限を超過）
		t.Log("15秒待機中...")
		time.Sleep(15 * time.Second)

		// 期限切れ一時コードで認証情報取得を試行
		params := &authapi.GetAuthCredentialsParams{
			Code:     &newTempCode,
			AuthFlow: (*authapi.GetAuthCredentialsParamsAuthFlow)(common.StringPtr("tempCodeAuth")),
		}

		getResp, err := client.Client.GetAuthCredentialsWithResponse(ctx, params)
		if err != nil {
			t.Fatalf("認証情報取得APIの呼び出しに失敗: %v", err)
		}

		// エラーステータスコードをチェック
		if getResp.StatusCode() == 200 {
			t.Error("期限切れ一時コードで成功レスポンスが返されました")
		}

		t.Log("期限切れ一時コードエラー確認成功")
	})
}

// testAuthenticationErrorHandling は認証エラーハンドリングのテストを実行します
func testAuthenticationErrorHandling(t *testing.T, client *common.ClientWrapper, testData *common.AuthenticationFlowTestData, assert *common.AssertionHelper) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Run("無効なIDトークンエラー", func(t *testing.T) {
		// 無効なIDトークンで認証情報保存を試行
		createParam := authapi.CreateAuthCredentialsJSONRequestBody{
			IdToken:     testData.ErrorHandling.InvalidTokens.Params.InvalidIDToken,
			AccessToken: testData.TestTokens.ValidAccessToken,
		}

		startTime := time.Now()
		createResp, err := client.Client.CreateAuthCredentialsWithResponse(ctx, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報保存APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "無効なIDトークン")

		// エラーステータスコードをチェック
		if createResp.StatusCode() == 201 {
			t.Error("無効なIDトークンで成功レスポンスが返されました")
		}

		t.Log("無効なIDトークンエラー確認成功")
	})

	t.Run("必須パラメータ不足エラー", func(t *testing.T) {
		// IDトークンなしで認証情報保存を試行
		createParam := authapi.CreateAuthCredentialsJSONRequestBody{
			AccessToken: testData.TestTokens.ValidAccessToken,
		}

		startTime := time.Now()
		createResp, err := client.Client.CreateAuthCredentialsWithResponse(ctx, createParam)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報保存APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "必須パラメータ不足")

		// エラーステータスコードをチェック
		if createResp.StatusCode() == 201 {
			t.Error("必須パラメータ不足で成功レスポンスが返されました")
		}

		t.Log("必須パラメータ不足エラー確認成功")
	})

	t.Run("存在しない一時コードエラー", func(t *testing.T) {
		invalidTempCode := "non-existent-temp-code"

		// 存在しない一時コードで認証情報取得を試行
		params := &authapi.GetAuthCredentialsParams{
			Code:     &invalidTempCode,
			AuthFlow: (*authapi.GetAuthCredentialsParamsAuthFlow)(common.StringPtr("tempCodeAuth")),
		}

		startTime := time.Now()
		getResp, err := client.Client.GetAuthCredentialsWithResponse(ctx, params)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("認証情報取得APIの呼び出しに失敗: %v", err)
		}

		// レスポンス時間をチェック
		assert.AssertResponseTime(duration, 10*time.Second, "存在しない一時コード")

		// エラーステータスコードをチェック
		assert.AssertStatusCode(404, getResp.StatusCode(), "存在しない一時コード")

		t.Log("存在しない一時コードエラー確認成功")
	})
}
