package authapi

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/client"
)

// E2E テスト用の設定
const (
	SaaSusAPIEndpoint = "https://api.dev.saasus.io/v1/auth"
)

// TestParams はテストパラメータファイルの構造体
type TestParams struct {
	GetUserInfo struct {
		TestCases []struct {
			Name   string `json:"name"`
			Params struct {
				Token string `json:"token"`
			} `json:"params"`
		} `json:"testCases"`
		EdgeCases []struct {
			Name   string `json:"name"`
			Params struct {
				Token string `json:"token"`
			} `json:"params"`
			ExpectError bool `json:"expectError"`
		} `json:"edgeCases"`
	} `json:"getUserInfo"`
	GetUserInfoByEmail struct {
		TestCases []struct {
			Name   string `json:"name"`
			Params struct {
				Email string `json:"email"`
			} `json:"params"`
		} `json:"testCases"`
		EdgeCases []struct {
			Name   string `json:"name"`
			Params struct {
				Email string `json:"email"`
			} `json:"params"`
			ExpectError bool `json:"expectError"`
		} `json:"edgeCases"`
	} `json:"getUserInfoByEmail"`
	BasicInfo struct {
		TestCases []struct {
			Name   string   `json:"name"`
			Params struct{} `json:"params"`
		} `json:"testCases"`
		UpdateParams struct {
			DomainName        string  `json:"domainName"`
			FromEmailAddress  string  `json:"fromEmailAddress"`
			ReplyEmailAddress *string `json:"replyEmailAddress,omitempty"`
		} `json:"updateParams"`
	} `json:"basicInfo"`
	AuthInfo struct {
		TestCases []struct {
			Name   string   `json:"name"`
			Params struct{} `json:"params"`
		} `json:"testCases"`
		UpdateParams struct {
			CallbackUrl string `json:"callbackUrl"`
		} `json:"updateParams"`
	} `json:"authInfo"`
	SaasUsers struct {
		TestCases []struct {
			Name   string   `json:"name"`
			Params struct{} `json:"params"`
		} `json:"testCases"`
		CreateParams struct {
			Email    string  `json:"email"`
			Password *string `json:"password,omitempty"`
		} `json:"createParams"`
		TestUserIds    []string `json:"testUserIds"`
		InvalidUserIds []string `json:"invalidUserIds"`
	} `json:"saasUsers"`
	Tenants struct {
		TestCases []struct {
			Name   string   `json:"name"`
			Params struct{} `json:"params"`
		} `json:"testCases"`
		CreateParams struct {
			Name                 string                 `json:"name"`
			Attributes           map[string]interface{} `json:"attributes"`
			BackOfficeStaffEmail string                 `json:"backOfficeStaffEmail"`
		} `json:"createParams"`
		TestTenantIds    []string `json:"testTenantIds"`
		InvalidTenantIds []string `json:"invalidTenantIds"`
	} `json:"tenants"`
	TenantUsers struct {
		CreateParams struct {
			Email      string                 `json:"email"`
			Attributes map[string]interface{} `json:"attributes"`
		} `json:"createParams"`
		UpdateParams struct {
			Attributes map[string]interface{} `json:"attributes"`
		} `json:"updateParams"`
	} `json:"tenantUsers"`
	Roles struct {
		TestCases []struct {
			Name   string   `json:"name"`
			Params struct{} `json:"params"`
		} `json:"testCases"`
		CreateParams struct {
			RoleName    string `json:"roleName"`
			DisplayName string `json:"displayName"`
		} `json:"createParams"`
		TestRoleNames    []string `json:"testRoleNames"`
		InvalidRoleNames []string `json:"invalidRoleNames"`
	} `json:"roles"`
	UserAttributes struct {
		TestCases []struct {
			Name   string   `json:"name"`
			Params struct{} `json:"params"`
		} `json:"testCases"`
		CreateParams struct {
			AttributeName string        `json:"attributeName"`
			DisplayName   string        `json:"displayName"`
			AttributeType AttributeType `json:"attributeType"`
		} `json:"createParams"`
		TestAttributeNames    []string `json:"testAttributeNames"`
		InvalidAttributeNames []string `json:"invalidAttributeNames"`
	} `json:"userAttributes"`
	TenantAttributes struct {
		TestCases []struct {
			Name   string   `json:"name"`
			Params struct{} `json:"params"`
		} `json:"testCases"`
		CreateParams struct {
			AttributeName string        `json:"attributeName"`
			DisplayName   string        `json:"displayName"`
			AttributeType AttributeType `json:"attributeType"`
		} `json:"createParams"`
	} `json:"tenantAttributes"`
	Envs struct {
		TestCases []struct {
			Name   string   `json:"name"`
			Params struct{} `json:"params"`
		} `json:"testCases"`
		CreateParams struct {
			Id          Id      `json:"id"`
			Name        string  `json:"name"`
			DisplayName *string `json:"displayName,omitempty"`
		} `json:"createParams"`
		TestEnvIds    []Id `json:"testEnvIds"`
		InvalidEnvIds []Id `json:"invalidEnvIds"`
	} `json:"envs"`
	Credentials struct {
		TestCases []struct {
			Name   string `json:"name"`
			Params struct {
				Code         *string `json:"code,omitempty"`
				AuthFlow     *string `json:"authFlow,omitempty"`
				RefreshToken *string `json:"refreshToken,omitempty"`
			} `json:"params"`
		} `json:"testCases"`
		CreateParams struct {
			IdToken      string  `json:"idToken"`
			AccessToken  string  `json:"accessToken"`
			RefreshToken *string `json:"refreshToken,omitempty"`
		} `json:"createParams"`
		EdgeCases []struct {
			Name   string `json:"name"`
			Params struct {
				Code         *string `json:"code,omitempty"`
				AuthFlow     *string `json:"authFlow,omitempty"`
				RefreshToken *string `json:"refreshToken,omitempty"`
			} `json:"params"`
			ExpectError bool `json:"expectError"`
		} `json:"edgeCases"`
	} `json:"credentials"`
	SignUp struct {
		TestCases []struct {
			Name   string `json:"name"`
			Params struct {
				Email string `json:"email"`
			} `json:"params"`
		} `json:"testCases"`
		ResendParams struct {
			Email string `json:"email"`
		} `json:"resendParams"`
	} `json:"signUp"`
	IdentityProviders struct {
		UpdateParams struct {
			Provider              ProviderName           `json:"provider"`
			IdentityProviderProps *IdentityProviderProps `json:"identityProviderProps,omitempty"`
		} `json:"updateParams"`
	} `json:"identityProviders"`
	SignInSettings struct {
		UpdateParams struct {
			PasswordPolicy      *PasswordPolicy      `json:"passwordPolicy,omitempty"`
			DeviceConfiguration *DeviceConfiguration `json:"deviceConfiguration,omitempty"`
			MfaConfiguration    *MfaConfiguration    `json:"mfaConfiguration,omitempty"`
			SelfRegist          *SelfRegist          `json:"selfRegist,omitempty"`
		} `json:"updateParams"`
	} `json:"signInSettings"`
	CustomizePages struct {
		UpdateParams struct {
			SignUpPage        *CustomizePageProps `json:"signUpPage,omitempty"`
			SignInPage        *CustomizePageProps `json:"signInPage,omitempty"`
			PasswordResetPage *CustomizePageProps `json:"passwordResetPage,omitempty"`
		} `json:"updateParams"`
	} `json:"customizePages"`
	CustomizePageSettings struct {
		UpdateParams struct {
			Title                       string `json:"title"`
			TermsOfServiceUrl           string `json:"termsOfServiceUrl"`
			PrivacyPolicyUrl            string `json:"privacyPolicyUrl"`
			GoogleTagManagerContainerId string `json:"googleTagManagerContainerId"`
			Icon                        string `json:"icon"`
			Favicon                     string `json:"favicon"`
		} `json:"updateParams"`
	} `json:"customizePageSettings"`
	NotificationMessages struct {
		UpdateParams struct {
			SignUp     *MessageTemplate `json:"signUp,omitempty"`
			CreateUser *MessageTemplate `json:"createUser,omitempty"`
		} `json:"updateParams"`
	} `json:"notificationMessages"`
	SingleTenantSettings struct {
		UpdateParams struct {
			Enabled                *bool   `json:"enabled,omitempty"`
			RoleArn                *string `json:"roleArn,omitempty"`
			CloudformationTemplate *string `json:"cloudformationTemplate,omitempty"`
			DdlTemplate            *string `json:"ddlTemplate,omitempty"`
			RoleExternalId         *string `json:"roleExternalId,omitempty"`
		} `json:"updateParams"`
	} `json:"singleTenantSettings"`
	Performance struct {
		MaxResponseTime string `json:"maxResponseTime"`
		TestLimit       int64  `json:"testLimit"`
		LoadTestCases   []struct {
			Name            string `json:"name"`
			ExpectedMaxTime string `json:"expectedMaxTime"`
		} `json:"loadTestCases"`
	} `json:"performance"`
}

// stringPtr は文字列のポインタを返すヘルパー関数
func stringPtr(s string) *string {
	return &s
}

// loadTestParams はテストパラメータファイルを読み込みます
func loadTestParams(t *testing.T) *TestParams {
	paramFile := filepath.Join("test_params.json")
	data, err := os.ReadFile(paramFile)
	if err != nil {
		t.Fatalf("テストパラメータファイルの読み込みに失敗: %v", err)
	}

	var params TestParams
	if err := json.Unmarshal(data, &params); err != nil {
		t.Fatalf("テストパラメータファイルの解析に失敗: %v", err)
	}

	return &params
}

// setupE2EClient は認証付きのE2Eテスト用クライアントを作成します
func setupE2EClient(t *testing.T) *ClientWithResponses {
	// 必要な環境変数をチェック
	saasID := os.Getenv("SAASUS_SAAS_ID")
	apiKey := os.Getenv("SAASUS_API_KEY")
	secretKey := os.Getenv("SAASUS_SECRET_KEY")

	if saasID == "" || apiKey == "" || secretKey == "" {
		t.Skip("E2Eテストをスキップ: SAASUS_SAAS_ID, SAASUS_API_KEY, SAASUS_SECRET_KEY 環境変数が設定されていません")
	}

	// SaaSus認証を設定するリクエストエディター
	authEditor := func(ctx context.Context, req *http.Request) error {
		if err := client.SetSigV1(req); err != nil {
			return err
		}
		client.SetReferer(ctx, req)
		return nil
	}

	// 認証付きクライアントを作成
	clientWithAuth, err := NewClientWithResponses(
		SaaSusAPIEndpoint,
		WithRequestEditorFn(authEditor),
	)
	if err != nil {
		t.Fatalf("認証付きクライアントの作成に失敗: %v", err)
	}

	return clientWithAuth
}

// TestE2E_GetAuthInfoWithResponse は GetAuthInfoWithResponse のE2Eテスト
func TestE2E_GetAuthInfoWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetAuthInfoWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("認証情報: CallbackUrl=%s", resp.JSON200.CallbackUrl)
		if resp.JSON200.CallbackUrl == "" {
			t.Error("CallbackUrlが空です")
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateAuthInfoWithResponse は UpdateAuthInfoWithResponse のE2Eテスト
func TestE2E_UpdateAuthInfoWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	updateParam := UpdateAuthInfoParam{
		CallbackUrl: testParams.AuthInfo.UpdateParams.CallbackUrl,
	}

	resp, err := client.UpdateAuthInfoWithResponse(ctx, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("認証情報の更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetBasicInfoWithResponse は GetBasicInfoWithResponse のE2Eテスト
func TestE2E_GetBasicInfoWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetBasicInfoWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("基本情報: DomainName=%s", resp.JSON200.DomainName)
		t.Logf("DNS検証済み: %t", resp.JSON200.IsDnsValidated)
		t.Logf("SESサンドボックス解除済み: %t", resp.JSON200.IsSesSandboxGranted)

		if resp.JSON200.DomainName == "" {
			t.Error("ドメイン名が空です")
		}
		if resp.JSON200.FromEmailAddress == "" {
			t.Error("送信元メールアドレスが空です")
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateBasicInfoWithResponse は UpdateBasicInfoWithResponse のE2Eテスト
func TestE2E_UpdateBasicInfoWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	updateParam := UpdateBasicInfoParam{
		DomainName:        testParams.BasicInfo.UpdateParams.DomainName,
		FromEmailAddress:  testParams.BasicInfo.UpdateParams.FromEmailAddress,
		ReplyEmailAddress: testParams.BasicInfo.UpdateParams.ReplyEmailAddress,
	}

	resp, err := client.UpdateBasicInfoWithResponse(ctx, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("基本情報の更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetSaasUsersWithResponse は GetSaasUsersWithResponse のE2Eテスト
func TestE2E_GetSaasUsersWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetSaasUsersWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("SaaSユーザー数: %d", len(resp.JSON200.Users))

		for i, user := range resp.JSON200.Users {
			if i >= 3 { // 最初の3件のみ詳細チェック
				break
			}
			t.Logf("ユーザー %d: ID=%s, Email=%s", i+1, user.Id, user.Email)

			if user.Id == "" {
				t.Errorf("ユーザー %d: IDが空です", i+1)
			}
			if user.Email == "" {
				t.Errorf("ユーザー %d: メールアドレスが空です", i+1)
			}
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateSaasUserWithResponse は CreateSaasUserWithResponse のE2Eテスト
func TestE2E_CreateSaasUserWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用のユニークなメールアドレスを生成
	timestamp := time.Now().Unix()
	testEmail := "test-" + string(rune(timestamp)) + "@example.com"

	createParam := CreateSaasUserParam{
		Email:    testEmail,
		Password: testParams.SaasUsers.CreateParams.Password,
	}

	resp, err := client.CreateSaasUserWithResponse(ctx, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("作成されたユーザー: ID=%s, Email=%s", resp.JSON201.Id, resp.JSON201.Email)

		if resp.JSON201.Id == "" {
			t.Error("作成されたユーザーのIDが空です")
		}
		if resp.JSON201.Email != testEmail {
			t.Errorf("メールアドレスが一致しません: 期待=%s, 実際=%s", testEmail, resp.JSON201.Email)
		}
	case 400:
		if resp.JSON400 != nil {
			t.Logf("バリデーションエラー: Type=%s, Message=%s", resp.JSON400.Type, resp.JSON400.Message)
		}
		t.Log("バリデーションエラーが発生しました（既存ユーザーの可能性）")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetSaasUserWithResponse は GetSaasUserWithResponse のE2Eテスト
func TestE2E_GetSaasUserWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、個別ユーザー取得テストをスキップします")
	}

	// 実際のユーザーIDを使用してテスト
	actualUserId := UserId(usersResp.JSON200.Users[0].Id)
	t.Logf("テスト対象のユーザーID: %s", actualUserId)

	resp, err := client.GetSaasUserWithResponse(ctx, actualUserId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("ユーザー詳細: ID=%s, Email=%s", resp.JSON200.Id, resp.JSON200.Email)

		if resp.JSON200.Id != string(actualUserId) {
			t.Errorf("ユーザーIDが一致しません: 期待=%s, 実際=%s", actualUserId, resp.JSON200.Id)
		}
		if resp.JSON200.Email == "" {
			t.Error("メールアドレスが空です")
		}
	case 404:
		t.Log("ユーザーが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}

	// 無効なユーザーIDでのテスト
	for _, invalidId := range testParams.SaasUsers.InvalidUserIds {
		t.Run("無効なユーザーID_"+invalidId, func(t *testing.T) {
			resp, err := client.GetSaasUserWithResponse(ctx, UserId(invalidId))

			if err != nil {
				t.Logf("リクエストエラー: %v", err)
				return
			}

			if resp != nil && (resp.StatusCode() == 400 || resp.StatusCode() == 404) {
				t.Logf("期待通りエラーレスポンス: %d", resp.StatusCode())
			} else {
				t.Logf("予期しないレスポンス: %d", resp.StatusCode())
			}
		})
	}
}

// TestE2E_DeleteSaasUserWithResponse は DeleteSaasUserWithResponse のE2Eテスト
func TestE2E_DeleteSaasUserWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// テスト用ユーザーを作成
	timestamp := time.Now().Unix()
	testEmail := "delete-test-" + string(rune(timestamp)) + "@example.com"

	createParam := CreateSaasUserParam{
		Email: testEmail,
	}

	createResp, err := client.CreateSaasUserWithResponse(ctx, createParam)
	if err != nil || createResp == nil || createResp.StatusCode() != 201 || createResp.JSON201 == nil {
		t.Skip("テスト用ユーザーの作成に失敗したため、削除テストをスキップします")
	}

	createdUserId := UserId(createResp.JSON201.Id)
	t.Logf("削除対象のユーザーID: %s", createdUserId)

	// ユーザーを削除
	resp, err := client.DeleteSaasUserWithResponse(ctx, createdUserId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("ユーザーの削除が成功しました")
	case 404:
		t.Log("ユーザーが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetTenantsWithResponse は GetTenantsWithResponse のE2Eテスト
func TestE2E_GetTenantsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetTenantsWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("テナント数: %d", len(resp.JSON200.Tenants))

		for i, tenant := range resp.JSON200.Tenants {
			if i >= 3 { // 最初の3件のみ詳細チェック
				break
			}
			t.Logf("テナント %d: ID=%s, Name=%s", i+1, tenant.Id, tenant.Name)

			if tenant.Id == "" {
				t.Errorf("テナント %d: IDが空です", i+1)
			}
			if tenant.Name == "" {
				t.Errorf("テナント %d: 名前が空です", i+1)
			}
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateTenantWithResponse は CreateTenantWithResponse のE2Eテスト
func TestE2E_CreateTenantWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用のユニークなテナント名を生成
	timestamp := time.Now().Unix()
	testTenantName := "Test Tenant " + string(rune(timestamp))

	createParam := CreateTenantParam{
		Name:                 testTenantName,
		Attributes:           testParams.Tenants.CreateParams.Attributes,
		BackOfficeStaffEmail: testParams.Tenants.CreateParams.BackOfficeStaffEmail,
	}

	resp, err := client.CreateTenantWithResponse(ctx, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("作成されたテナント: ID=%s, Name=%s", resp.JSON201.Id, resp.JSON201.Name)

		if resp.JSON201.Id == "" {
			t.Error("作成されたテナントのIDが空です")
		}
		if resp.JSON201.Name != testTenantName {
			t.Errorf("テナント名が一致しません: 期待=%s, 実際=%s", testTenantName, resp.JSON201.Name)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetRolesWithResponse は GetRolesWithResponse のE2Eテスト
func TestE2E_GetRolesWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetRolesWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("役割数: %d", len(resp.JSON200.Roles))

		for i, role := range resp.JSON200.Roles {
			if i >= 3 { // 最初の3件のみ詳細チェック
				break
			}
			t.Logf("役割 %d: RoleName=%s, DisplayName=%s", i+1, role.RoleName, role.DisplayName)

			if role.RoleName == "" {
				t.Errorf("役割 %d: 役割名が空です", i+1)
			}
			if role.DisplayName == "" {
				t.Errorf("役割 %d: 表示名が空です", i+1)
			}
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateRoleWithResponse は CreateRoleWithResponse のE2Eテスト
func TestE2E_CreateRoleWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用のユニークな役割名を生成
	timestamp := time.Now().Unix()
	testRoleName := "test_role_" + string(rune(timestamp))

	createParam := CreateRoleParam{
		RoleName:    testRoleName,
		DisplayName: testParams.Roles.CreateParams.DisplayName,
	}

	resp, err := client.CreateRoleWithResponse(ctx, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("作成された役割: RoleName=%s, DisplayName=%s", resp.JSON201.RoleName, resp.JSON201.DisplayName)

		if resp.JSON201.RoleName != testRoleName {
			t.Errorf("役割名が一致しません: 期待=%s, 実際=%s", testRoleName, resp.JSON201.RoleName)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetEnvsWithResponse は GetEnvsWithResponse のE2Eテスト
func TestE2E_GetEnvsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetEnvsWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("環境数: %d", len(resp.JSON200.Envs))

		for i, env := range resp.JSON200.Envs {
			if i >= 3 { // 最初の3件のみ詳細チェック
				break
			}
			t.Logf("環境 %d: ID=%d, Name=%s", i+1, env.Id, env.Name)

			if env.Id == 0 {
				t.Errorf("環境 %d: IDが0です", i+1)
			}
			if env.Name == "" {
				t.Errorf("環境 %d: 名前が空です", i+1)
			}
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetUserAttributesWithResponse は GetUserAttributesWithResponse のE2Eテスト
func TestE2E_GetUserAttributesWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetUserAttributesWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("ユーザー属性数: %d", len(resp.JSON200.UserAttributes))

		for i, attr := range resp.JSON200.UserAttributes {
			if i >= 3 { // 最初の3件のみ詳細チェック
				break
			}
			t.Logf("属性 %d: Name=%s, DisplayName=%s, Type=%s", i+1, attr.AttributeName, attr.DisplayName, attr.AttributeType)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetTenantAttributesWithResponse は GetTenantAttributesWithResponse のE2Eテスト
func TestE2E_GetTenantAttributesWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetTenantAttributesWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("テナント属性数: %d", len(resp.JSON200.TenantAttributes))

		for i, attr := range resp.JSON200.TenantAttributes {
			if i >= 3 { // 最初の3件のみ詳細チェック
				break
			}
			t.Logf("属性 %d: Name=%s, DisplayName=%s, Type=%s", i+1, attr.AttributeName, attr.DisplayName, attr.AttributeType)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetSignInSettingsWithResponse は GetSignInSettingsWithResponse のE2Eテスト
func TestE2E_GetSignInSettingsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetSignInSettingsWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("サインイン設定: パスワード最小長=%d", resp.JSON200.PasswordPolicy.MinimumLength)
		t.Logf("デバイス記憶設定: %s", resp.JSON200.DeviceConfiguration.DeviceRemembering)
		t.Logf("MFA設定: %s", resp.JSON200.MfaConfiguration.MfaConfiguration)
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetIdentityProvidersWithResponse は GetIdentityProvidersWithResponse のE2Eテスト
func TestE2E_GetIdentityProvidersWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetIdentityProvidersWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("Google IDプロバイダー: ApplicationId=%s", resp.JSON200.Google.ApplicationId)
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetCustomizePagesWithResponse は GetCustomizePagesWithResponse のE2Eテスト
func TestE2E_GetCustomizePagesWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetCustomizePagesWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("カスタマイズページ設定を取得しました")
		t.Logf("サインアップページ: 利用規約=%t, プライバシーポリシー=%t",
			resp.JSON200.SignUpPage.IsTermsOfService, resp.JSON200.SignUpPage.IsPrivacyPolicy)
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetCustomizePageSettingsWithResponse は GetCustomizePageSettingsWithResponse のE2Eテスト
func TestE2E_GetCustomizePageSettingsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetCustomizePageSettingsWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("ページ設定: Title=%s", resp.JSON200.Title)
		t.Logf("利用規約URL: %s", resp.JSON200.TermsOfServiceUrl)
		t.Logf("プライバシーポリシーURL: %s", resp.JSON200.PrivacyPolicyUrl)
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_FindNotificationMessagesWithResponse は FindNotificationMessagesWithResponse のE2Eテスト
func TestE2E_FindNotificationMessagesWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.FindNotificationMessagesWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("通知メッセージ設定を取得しました")
		t.Logf("サインアップメッセージ: Subject=%s", resp.JSON200.SignUp.Subject)
		t.Logf("ユーザー作成メッセージ: Subject=%s", resp.JSON200.CreateUser.Subject)
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetSingleTenantSettingsWithResponse は GetSingleTenantSettingsWithResponse のE2Eテスト
func TestE2E_GetSingleTenantSettingsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetSingleTenantSettingsWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("シングルテナント設定: Enabled=%t", resp.JSON200.Enabled)
		if resp.JSON200.Enabled {
			t.Logf("RoleArn: %s", resp.JSON200.RoleArn)
			t.Logf("RoleExternalId: %s", resp.JSON200.RoleExternalId)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetCloudFormationLaunchStackLinkForSingleTenantWithResponse は GetCloudFormationLaunchStackLinkForSingleTenantWithResponse のE2Eテスト
func TestE2E_GetCloudFormationLaunchStackLinkForSingleTenantWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetCloudFormationLaunchStackLinkForSingleTenantWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("CloudFormationスタック起動リンク: %s", resp.JSON200.Link)
		if resp.JSON200.Link == "" {
			t.Error("リンクが空です")
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_ResetPlanWithResponse は ResetPlanWithResponse のE2Eテスト
func TestE2E_ResetPlanWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.ResetPlanWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("プランリセットが成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateTenantAndPricingWithResponse は CreateTenantAndPricingWithResponse のE2Eテスト
func TestE2E_CreateTenantAndPricingWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.CreateTenantAndPricingWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("Stripe初期設定が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_DeleteStripeTenantAndPricingWithResponse は DeleteStripeTenantAndPricingWithResponse のE2Eテスト
func TestE2E_DeleteStripeTenantAndPricingWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.DeleteStripeTenantAndPricingWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("Stripe顧客・商品情報の削除が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_SignUpWithResponse は SignUpWithResponse のE2Eテスト
func TestE2E_SignUpWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用のユニークなメールアドレスを生成
	timestamp := time.Now().Unix()
	testEmail := "signup-test-" + string(rune(timestamp)) + "@example.com"

	signUpParam := SignUpParam{
		Email: testEmail,
	}

	resp, err := client.SignUpWithResponse(ctx, signUpParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("サインアップ成功: ID=%s, Email=%s", resp.JSON201.Id, resp.JSON201.Email)

		if resp.JSON201.Email != testEmail {
			t.Errorf("メールアドレスが一致しません: 期待=%s, 実際=%s", testEmail, resp.JSON201.Email)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}

	// 確認メール再送信のテスト
	resendParam := ResendSignUpConfirmationEmailParam{
		Email: testParams.SignUp.ResendParams.Email,
	}

	resendResp, err := client.ResendSignUpConfirmationEmailWithResponse(ctx, resendParam)
	if err != nil {
		t.Logf("確認メール再送信エラー: %v", err)
	} else if resendResp != nil {
		t.Logf("確認メール再送信ステータス: %d", resendResp.StatusCode())
	}
}

// TestE2E_PerformanceBasic は基本的なパフォーマンステスト
func TestE2E_PerformanceBasic(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// レスポンス時間の測定
	start := time.Now()
	resp, err := client.GetBasicInfoWithResponse(ctx)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("パフォーマンステストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("レスポンス時間: %v", duration)
	t.Logf("ステータスコード: %d", resp.StatusCode())

	// 設定ファイルから最大レスポンス時間を取得して検証
	maxDuration, err := time.ParseDuration(testParams.Performance.MaxResponseTime)
	if err != nil {
		t.Fatalf("最大レスポンス時間の解析に失敗: %v", err)
	}

	if duration > maxDuration {
		t.Errorf("レスポンス時間が長すぎます: %v (%v以内が期待される)", duration, maxDuration)
	}

	// 成功レスポンスの場合、データの整合性を確認
	if resp.StatusCode() == 200 && resp.JSON200 != nil {
		t.Logf("基本情報取得成功")
	}
}

// TestE2E_PerformanceLoad は負荷別パフォーマンステスト
func TestE2E_PerformanceLoad(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)

	for _, loadTest := range testParams.Performance.LoadTestCases {
		tc := loadTest // ループ変数をキャプチャ
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			// レスポンス時間の測定
			start := time.Now()
			resp, err := client.GetSaasUsersWithResponse(ctx)
			duration := time.Since(start)

			if err != nil {
				t.Fatalf("パフォーマンステストエラー: %v", err)
			}

			if resp == nil {
				t.Fatal("レスポンスがnilです")
			}

			t.Logf("レスポンス時間: %v", duration)
			t.Logf("ステータスコード: %d", resp.StatusCode())

			// 期待される最大時間を取得して検証
			expectedMaxDuration, err := time.ParseDuration(tc.ExpectedMaxTime)
			if err != nil {
				t.Fatalf("期待最大レスポンス時間の解析に失敗: %v", err)
			}

			if duration > expectedMaxDuration {
				t.Errorf("レスポンス時間が期待値を超えています: %v (%v以内が期待される)", duration, expectedMaxDuration)
			}

			// 成功レスポンスの場合、データの整合性を確認
			if resp.StatusCode() == 200 && resp.JSON200 != nil {
				userCount := len(resp.JSON200.Users)
				t.Logf("取得したユーザー数: %d", userCount)

				// スループットの計算
				if duration > 0 {
					throughput := float64(userCount) / duration.Seconds()
					t.Logf("スループット: %.2f件/秒", throughput)
				}
			}
		})
	}
}

// TestE2E_ConcurrentRequests は同時リクエストのパフォーマンステスト
func TestE2E_ConcurrentRequests(t *testing.T) {
	client := setupE2EClient(t)
	concurrency := 3 // 同時実行数
	ctx := context.Background()

	// 同時実行用のチャネル
	results := make(chan struct {
		duration time.Duration
		err      error
		status   int
	}, concurrency)

	// 同時リクエストを開始
	start := time.Now()
	for i := 0; i < concurrency; i++ {
		go func() {
			reqStart := time.Now()
			resp, err := client.GetBasicInfoWithResponse(ctx)
			reqDuration := time.Since(reqStart)

			result := struct {
				duration time.Duration
				err      error
				status   int
			}{
				duration: reqDuration,
				err:      err,
				status:   0,
			}

			if resp != nil {
				result.status = resp.StatusCode()
			}

			results <- result
		}()
	}

	// 結果を収集
	var totalDuration time.Duration
	var successCount int
	var errorCount int

	for i := 0; i < concurrency; i++ {
		result := <-results
		totalDuration += result.duration

		if result.err != nil {
			errorCount++
			t.Logf("リクエスト %d エラー: %v", i+1, result.err)
		} else if result.status == 200 {
			successCount++
			t.Logf("リクエスト %d 成功: %v", i+1, result.duration)
		} else {
			t.Logf("リクエスト %d ステータス: %d, 時間: %v", i+1, result.status, result.duration)
		}
	}

	totalTime := time.Since(start)
	avgDuration := totalDuration / time.Duration(concurrency)

	t.Logf("同時実行結果:")
	t.Logf("  - 同時実行数: %d", concurrency)
	t.Logf("  - 総実行時間: %v", totalTime)
	t.Logf("  - 平均レスポンス時間: %v", avgDuration)
	t.Logf("  - 成功数: %d", successCount)
	t.Logf("  - エラー数: %d", errorCount)

	// 成功率の検証
	successRate := float64(successCount) / float64(concurrency) * 100
	t.Logf("  - 成功率: %.1f%%", successRate)

	if successRate < 80.0 {
		t.Errorf("成功率が低すぎます: %.1f%% (80%%以上が期待される)", successRate)
	}
}

// TestE2E_GetTenantWithResponse は GetTenantWithResponse のE2Eテスト
func TestE2E_GetTenantWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// まず、テナント一覧を取得して実際のテナントIDを取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、個別テナント取得テストをスキップします")
	}

	// 実際のテナントIDを使用してテスト
	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)
	t.Logf("テスト対象のテナントID: %s", actualTenantId)

	resp, err := client.GetTenantWithResponse(ctx, actualTenantId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("テナント詳細: ID=%s, Name=%s", resp.JSON200.Id, resp.JSON200.Name)

		if resp.JSON200.Id != string(actualTenantId) {
			t.Errorf("テナントIDが一致しません: 期待=%s, 実際=%s", actualTenantId, resp.JSON200.Id)
		}
		if resp.JSON200.Name == "" {
			t.Error("テナント名が空です")
		}
	case 404:
		t.Log("テナントが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}

	// 無効なテナントIDでのテスト
	for _, invalidId := range testParams.Tenants.InvalidTenantIds {
		t.Run("無効なテナントID_"+invalidId, func(t *testing.T) {
			resp, err := client.GetTenantWithResponse(ctx, TenantId(invalidId))

			if err != nil {
				t.Logf("リクエストエラー: %v", err)
				return
			}

			if resp != nil && (resp.StatusCode() == 400 || resp.StatusCode() == 404) {
				t.Logf("期待通りエラーレスポンス: %d", resp.StatusCode())
			} else {
				t.Logf("予期しないレスポンス: %d", resp.StatusCode())
			}
		})
	}
}

// TestE2E_DeleteTenantWithResponse は DeleteTenantWithResponse のE2Eテスト
func TestE2E_DeleteTenantWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用テナントを作成
	timestamp := time.Now().Unix()
	testTenantName := "Delete Test Tenant " + string(rune(timestamp))

	createParam := CreateTenantParam{
		Name:                 testTenantName,
		Attributes:           testParams.Tenants.CreateParams.Attributes,
		BackOfficeStaffEmail: testParams.Tenants.CreateParams.BackOfficeStaffEmail,
	}

	createResp, err := client.CreateTenantWithResponse(ctx, createParam)
	if err != nil || createResp == nil || createResp.StatusCode() != 201 || createResp.JSON201 == nil {
		t.Skip("テスト用テナントの作成に失敗したため、削除テストをスキップします")
	}

	createdTenantId := TenantId(createResp.JSON201.Id)
	t.Logf("削除対象のテナントID: %s", createdTenantId)

	// テナントを削除
	resp, err := client.DeleteTenantWithResponse(ctx, createdTenantId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナントの削除が成功しました")
	case 404:
		t.Log("テナントが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetTenantUsersWithResponse は GetTenantUsersWithResponse のE2Eテスト
func TestE2E_GetTenantUsersWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得して実際のテナントIDを取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナントユーザー取得テストをスキップします")
	}

	// 実際のテナントIDを使用してテスト
	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)
	t.Logf("テスト対象のテナントID: %s", actualTenantId)

	resp, err := client.GetTenantUsersWithResponse(ctx, actualTenantId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("テナントユーザー数: %d", len(resp.JSON200.Users))

		for i, user := range resp.JSON200.Users {
			if i >= 3 { // 最初の3件のみ詳細チェック
				break
			}
			t.Logf("ユーザー %d: ID=%s, Email=%s", i+1, user.Id, user.Email)

			if user.Id == "" {
				t.Errorf("ユーザー %d: IDが空です", i+1)
			}
			if user.Email == "" {
				t.Errorf("ユーザー %d: メールアドレスが空です", i+1)
			}
		}
	case 404:
		t.Log("テナントが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateTenantUserWithResponse は CreateTenantUserWithResponse のE2Eテスト
func TestE2E_CreateTenantUserWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// まず、テナント一覧を取得して実際のテナントIDを取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナントユーザー作成テストをスキップします")
	}

	// 実際のテナントIDを使用してテスト
	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)
	t.Logf("テスト対象のテナントID: %s", actualTenantId)

	// テスト用のユニークなメールアドレスを生成
	timestamp := time.Now().Unix()
	testEmail := "tenant-user-" + string(rune(timestamp)) + "@example.com"

	createParam := CreateTenantUserParam{
		Email:      testEmail,
		Attributes: testParams.TenantUsers.CreateParams.Attributes,
	}

	resp, err := client.CreateTenantUserWithResponse(ctx, actualTenantId, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("作成されたテナントユーザー: ID=%s, Email=%s", resp.JSON201.Id, resp.JSON201.Email)

		if resp.JSON201.Id == "" {
			t.Error("作成されたユーザーのIDが空です")
		}
		if resp.JSON201.Email != testEmail {
			t.Errorf("メールアドレスが一致しません: 期待=%s, 実際=%s", testEmail, resp.JSON201.Email)
		}
	case 400:
		t.Log("バリデーションエラーが発生しました（既存ユーザーの可能性）")
		t.Log("バリデーションエラーが発生しました（既存ユーザーの可能性）")
	case 404:
		t.Log("テナントが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetRoleWithResponse は GetRoleWithResponse のE2Eテスト（存在しない場合は個別役割取得をスキップ）
func TestE2E_GetRoleWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// まず、役割一覧を取得して実際の役割名を取得
	rolesResp, err := client.GetRolesWithResponse(ctx)
	if err != nil {
		t.Fatalf("役割一覧取得エラー: %v", err)
	}

	if rolesResp == nil || rolesResp.StatusCode() != 200 || rolesResp.JSON200 == nil || len(rolesResp.JSON200.Roles) == 0 {
		t.Skip("利用可能な役割がないため、個別役割取得テストをスキップします")
	}

	// 実際の役割名を使用してテスト（個別取得APIが存在しない場合はスキップ）
	t.Log("個別役割取得APIは存在しないため、役割一覧取得で代替します")

	// 無効な役割名でのテスト（削除APIを使用）
	for _, invalidRoleName := range testParams.Roles.InvalidRoleNames {
		t.Run("無効な役割名削除_"+invalidRoleName, func(t *testing.T) {
			resp, err := client.DeleteRoleWithResponse(ctx, RoleName(invalidRoleName))

			if err != nil {
				t.Logf("リクエストエラー: %v", err)
				return
			}

			if resp != nil && (resp.StatusCode() == 400 || resp.StatusCode() == 404) {
				t.Logf("期待通りエラーレスポンス: %d", resp.StatusCode())
			} else {
				t.Logf("予期しないレスポンス: %d", resp.StatusCode())
			}
		})
	}
}

// TestE2E_DeleteRoleWithResponse は DeleteRoleWithResponse のE2Eテスト
func TestE2E_DeleteRoleWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用役割を作成
	timestamp := time.Now().Unix()
	testRoleName := "delete_test_role_" + string(rune(timestamp))

	createParam := CreateRoleParam{
		RoleName:    testRoleName,
		DisplayName: testParams.Roles.CreateParams.DisplayName,
	}

	createResp, err := client.CreateRoleWithResponse(ctx, createParam)
	if err != nil || createResp == nil || createResp.StatusCode() != 201 || createResp.JSON201 == nil {
		t.Skip("テスト用役割の作成に失敗したため、削除テストをスキップします")
	}

	createdRoleName := RoleName(createResp.JSON201.RoleName)
	t.Logf("削除対象の役割名: %s", createdRoleName)

	// 役割を削除
	resp, err := client.DeleteRoleWithResponse(ctx, createdRoleName)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("役割の削除が成功しました")
	case 404:
		t.Log("役割が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateUserAttributeWithResponse は CreateUserAttributeWithResponse のE2Eテスト
func TestE2E_CreateUserAttributeWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用のユニークな属性名を生成
	timestamp := time.Now().Unix()
	testAttributeName := "test_user_attr_" + string(rune(timestamp))

	createParam := CreateUserAttributeParam{
		AttributeName: testAttributeName,
		DisplayName:   testParams.UserAttributes.CreateParams.DisplayName,
		AttributeType: testParams.UserAttributes.CreateParams.AttributeType,
	}

	resp, err := client.CreateUserAttributeWithResponse(ctx, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("作成されたユーザー属性: Name=%s, DisplayName=%s, Type=%s",
			resp.JSON201.AttributeName, resp.JSON201.DisplayName, resp.JSON201.AttributeType)

		if resp.JSON201.AttributeName != testAttributeName {
			t.Errorf("属性名が一致しません: 期待=%s, 実際=%s", testAttributeName, resp.JSON201.AttributeName)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_DeleteUserAttributeWithResponse は DeleteUserAttributeWithResponse のE2Eテスト
func TestE2E_DeleteUserAttributeWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用ユーザー属性を作成
	timestamp := time.Now().Unix()
	testAttributeName := "delete_test_user_attr_" + string(rune(timestamp))

	createParam := CreateUserAttributeParam{
		AttributeName: testAttributeName,
		DisplayName:   testParams.UserAttributes.CreateParams.DisplayName,
		AttributeType: testParams.UserAttributes.CreateParams.AttributeType,
	}

	createResp, err := client.CreateUserAttributeWithResponse(ctx, createParam)
	if err != nil || createResp == nil || createResp.StatusCode() != 201 || createResp.JSON201 == nil {
		t.Skip("テスト用ユーザー属性の作成に失敗したため、削除テストをスキップします")
	}

	createdAttributeName := createResp.JSON201.AttributeName
	t.Logf("削除対象の属性名: %s", createdAttributeName)

	// ユーザー属性を削除
	resp, err := client.DeleteUserAttributeWithResponse(ctx, createdAttributeName)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("ユーザー属性の削除が成功しました")
	case 404:
		t.Log("ユーザー属性が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateTenantAttributeWithResponse は CreateTenantAttributeWithResponse のE2Eテスト
func TestE2E_CreateTenantAttributeWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用のユニークな属性名を生成
	timestamp := time.Now().Unix()
	testAttributeName := "test_tenant_attr_" + string(rune(timestamp))

	createParam := CreateTenantAttributeParam{
		AttributeName: testAttributeName,
		DisplayName:   testParams.TenantAttributes.CreateParams.DisplayName,
		AttributeType: testParams.TenantAttributes.CreateParams.AttributeType,
	}

	resp, err := client.CreateTenantAttributeWithResponse(ctx, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("作成されたテナント属性: Name=%s, DisplayName=%s, Type=%s",
			resp.JSON201.AttributeName, resp.JSON201.DisplayName, resp.JSON201.AttributeType)

		if resp.JSON201.AttributeName != testAttributeName {
			t.Errorf("属性名が一致しません: 期待=%s, 実際=%s", testAttributeName, resp.JSON201.AttributeName)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_DeleteTenantAttributeWithResponse は DeleteTenantAttributeWithResponse のE2Eテスト
func TestE2E_DeleteTenantAttributeWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用テナント属性を作成
	timestamp := time.Now().Unix()
	testAttributeName := "delete_test_tenant_attr_" + string(rune(timestamp))

	createParam := CreateTenantAttributeParam{
		AttributeName: testAttributeName,
		DisplayName:   testParams.TenantAttributes.CreateParams.DisplayName,
		AttributeType: testParams.TenantAttributes.CreateParams.AttributeType,
	}

	createResp, err := client.CreateTenantAttributeWithResponse(ctx, createParam)
	if err != nil || createResp == nil || createResp.StatusCode() != 201 || createResp.JSON201 == nil {
		t.Skip("テスト用テナント属性の作成に失敗したため、削除テストをスキップします")
	}

	createdAttributeName := createResp.JSON201.AttributeName
	t.Logf("削除対象の属性名: %s", createdAttributeName)

	// テナント属性を削除
	resp, err := client.DeleteTenantAttributeWithResponse(ctx, createdAttributeName)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナント属性の削除が成功しました")
	case 404:
		t.Log("テナント属性が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateEnvWithResponse は CreateEnvWithResponse のE2Eテスト
func TestE2E_CreateEnvWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用のユニークな環境名を生成
	timestamp := time.Now().Unix()
	testEnvName := "test_env_" + string(rune(timestamp))

	createParam := CreateEnvParam{
		Id:          testParams.Envs.CreateParams.Id,
		Name:        testEnvName,
		DisplayName: testParams.Envs.CreateParams.DisplayName,
	}

	resp, err := client.CreateEnvWithResponse(ctx, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("作成された環境: ID=%d, Name=%s", resp.JSON201.Id, resp.JSON201.Name)

		if resp.JSON201.Name != testEnvName {
			t.Errorf("環境名が一致しません: 期待=%s, 実際=%s", testEnvName, resp.JSON201.Name)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetEnvWithResponse は GetEnvWithResponse のE2Eテスト
func TestE2E_GetEnvWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// まず、環境一覧を取得して実際の環境IDを取得
	envsResp, err := client.GetEnvsWithResponse(ctx)
	if err != nil {
		t.Fatalf("環境一覧取得エラー: %v", err)
	}

	if envsResp == nil || envsResp.StatusCode() != 200 || envsResp.JSON200 == nil || len(envsResp.JSON200.Envs) == 0 {
		t.Skip("利用可能な環境がないため、個別環境取得テストをスキップします")
	}

	// 実際の環境IDを使用してテスト
	actualEnvId := EnvId(envsResp.JSON200.Envs[0].Id)
	t.Logf("テスト対象の環境ID: %d", actualEnvId)

	resp, err := client.GetEnvWithResponse(ctx, actualEnvId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("環境詳細: ID=%d, Name=%s", resp.JSON200.Id, resp.JSON200.Name)

		if resp.JSON200.Id != uint64(actualEnvId) {
			t.Errorf("環境IDが一致しません: 期待=%d, 実際=%d", actualEnvId, resp.JSON200.Id)
		}
		if resp.JSON200.Name == "" {
			t.Error("環境名が空です")
		}
	case 404:
		t.Log("環境が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}

	// 無効な環境IDでのテスト
	for _, invalidId := range testParams.Envs.InvalidEnvIds {
		t.Run("無効な環境ID_"+string(rune(invalidId)), func(t *testing.T) {
			resp, err := client.GetEnvWithResponse(ctx, EnvId(invalidId))

			if err != nil {
				t.Logf("リクエストエラー: %v", err)
				return
			}

			if resp != nil && (resp.StatusCode() == 400 || resp.StatusCode() == 404) {
				t.Logf("期待通りエラーレスポンス: %d", resp.StatusCode())
			} else {
				t.Logf("予期しないレスポンス: %d", resp.StatusCode())
			}
		})
	}
}

// TestE2E_DeleteEnvWithResponse は DeleteEnvWithResponse のE2Eテスト
func TestE2E_DeleteEnvWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テスト用環境を作成
	timestamp := time.Now().Unix()
	testEnvName := "delete_test_env_" + string(rune(timestamp))

	createParam := CreateEnvParam{
		Id:          testParams.Envs.CreateParams.Id,
		Name:        testEnvName,
		DisplayName: testParams.Envs.CreateParams.DisplayName,
	}

	createResp, err := client.CreateEnvWithResponse(ctx, createParam)
	if err != nil || createResp == nil || createResp.StatusCode() != 201 || createResp.JSON201 == nil {
		t.Skip("テスト用環境の作成に失敗したため、削除テストをスキップします")
	}

	createdEnvId := EnvId(createResp.JSON201.Id)
	t.Logf("削除対象の環境ID: %d", createdEnvId)

	// 環境を削除
	resp, err := client.DeleteEnvWithResponse(ctx, createdEnvId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("環境の削除が成功しました")
	case 404:
		t.Log("環境が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateSignInSettingsWithResponse は UpdateSignInSettingsWithResponse のE2Eテスト
func TestE2E_UpdateSignInSettingsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	updateParam := UpdateSignInSettingsParam{
		PasswordPolicy:      testParams.SignInSettings.UpdateParams.PasswordPolicy,
		DeviceConfiguration: testParams.SignInSettings.UpdateParams.DeviceConfiguration,
		MfaConfiguration:    testParams.SignInSettings.UpdateParams.MfaConfiguration,
		SelfRegist:          testParams.SignInSettings.UpdateParams.SelfRegist,
	}

	resp, err := client.UpdateSignInSettingsWithResponse(ctx, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("サインイン設定の更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateCustomizePagesWithResponse は UpdateCustomizePagesWithResponse のE2Eテスト
func TestE2E_UpdateCustomizePagesWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	updateParam := UpdateCustomizePagesParam{
		SignUpPage:        testParams.CustomizePages.UpdateParams.SignUpPage,
		SignInPage:        testParams.CustomizePages.UpdateParams.SignInPage,
		PasswordResetPage: testParams.CustomizePages.UpdateParams.PasswordResetPage,
	}

	resp, err := client.UpdateCustomizePagesWithResponse(ctx, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("カスタマイズページの更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateCustomizePageSettingsWithResponse は UpdateCustomizePageSettingsWithResponse のE2Eテスト
func TestE2E_UpdateCustomizePageSettingsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	updateParam := UpdateCustomizePageSettingsParam{
		Title:                       testParams.CustomizePageSettings.UpdateParams.Title,
		TermsOfServiceUrl:           testParams.CustomizePageSettings.UpdateParams.TermsOfServiceUrl,
		PrivacyPolicyUrl:            testParams.CustomizePageSettings.UpdateParams.PrivacyPolicyUrl,
		GoogleTagManagerContainerId: testParams.CustomizePageSettings.UpdateParams.GoogleTagManagerContainerId,
		Icon:                        testParams.CustomizePageSettings.UpdateParams.Icon,
		Favicon:                     testParams.CustomizePageSettings.UpdateParams.Favicon,
	}

	resp, err := client.UpdateCustomizePageSettingsWithResponse(ctx, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("ページ設定の更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateNotificationMessagesWithResponse は UpdateNotificationMessagesWithResponse のE2Eテスト
func TestE2E_UpdateNotificationMessagesWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	updateParam := UpdateNotificationMessagesParam{
		SignUp:     testParams.NotificationMessages.UpdateParams.SignUp,
		CreateUser: testParams.NotificationMessages.UpdateParams.CreateUser,
	}

	resp, err := client.UpdateNotificationMessagesWithResponse(ctx, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("通知メッセージの更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateSingleTenantSettingsWithResponse は UpdateSingleTenantSettingsWithResponse のE2Eテスト
func TestE2E_UpdateSingleTenantSettingsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	updateParam := UpdateSingleTenantSettingsParam{
		Enabled:                testParams.SingleTenantSettings.UpdateParams.Enabled,
		RoleArn:                testParams.SingleTenantSettings.UpdateParams.RoleArn,
		CloudformationTemplate: testParams.SingleTenantSettings.UpdateParams.CloudformationTemplate,
		DdlTemplate:            testParams.SingleTenantSettings.UpdateParams.DdlTemplate,
		RoleExternalId:         testParams.SingleTenantSettings.UpdateParams.RoleExternalId,
	}

	resp, err := client.UpdateSingleTenantSettingsWithResponse(ctx, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("シングルテナント設定の更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetAuthCredentialsWithResponse は GetAuthCredentialsWithResponse のE2Eテスト
func TestE2E_GetAuthCredentialsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テストケースを実行
	for _, testCase := range testParams.Credentials.TestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			params := &GetAuthCredentialsParams{
				Code:         testCase.Params.Code,
				AuthFlow:     (*GetAuthCredentialsParamsAuthFlow)(testCase.Params.AuthFlow),
				RefreshToken: testCase.Params.RefreshToken,
			}

			resp, err := client.GetAuthCredentialsWithResponse(ctx, params)

			if err != nil {
				t.Logf("リクエストエラー: %v", err)
				return
			}

			if resp == nil {
				t.Fatal("レスポンスがnilです")
			}

			t.Logf("ステータスコード: %d", resp.StatusCode())

			switch resp.StatusCode() {
			case 200:
				if resp.JSON200 == nil {
					t.Error("200レスポンスが解析されませんでした")
					return
				}
				t.Logf("認証情報取得成功: AccessToken=%s", resp.JSON200.AccessToken)
			case 400, 401:
				t.Log("認証エラーが発生しました（期待される動作）")
			case 500:
				if resp.JSON500 != nil {
					t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
				}
				t.Error("サーバーエラーが発生しました")
			default:
				t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
			}
		})
	}

	// エッジケースのテスト
	for _, edgeCase := range testParams.Credentials.EdgeCases {
		t.Run("EdgeCase_"+edgeCase.Name, func(t *testing.T) {
			params := &GetAuthCredentialsParams{
				Code:         edgeCase.Params.Code,
				AuthFlow:     (*GetAuthCredentialsParamsAuthFlow)(edgeCase.Params.AuthFlow),
				RefreshToken: edgeCase.Params.RefreshToken,
			}

			resp, err := client.GetAuthCredentialsWithResponse(ctx, params)

			if err != nil {
				if edgeCase.ExpectError {
					t.Logf("期待通りエラーが発生: %v", err)
				} else {
					t.Errorf("予期しないエラー: %v", err)
				}
				return
			}

			if resp != nil && edgeCase.ExpectError && (resp.StatusCode() == 400 || resp.StatusCode() == 401) {
				t.Logf("期待通りエラーレスポンス: %d", resp.StatusCode())
			}
		})
	}
}

// TestE2E_CreateAuthCredentialsWithResponse は CreateAuthCredentialsWithResponse のE2Eテスト
func TestE2E_CreateAuthCredentialsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	createParam := CreateAuthCredentialsParam{
		IdToken:      testParams.Credentials.CreateParams.IdToken,
		AccessToken:  testParams.Credentials.CreateParams.AccessToken,
		RefreshToken: testParams.Credentials.CreateParams.RefreshToken,
	}

	resp, err := client.CreateAuthCredentialsWithResponse(ctx, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		t.Log("認証情報の作成が成功しました")
	case 400:
		t.Log("バリデーションエラーが発生しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateEnvWithResponse は UpdateEnvWithResponse のE2Eテスト
func TestE2E_UpdateEnvWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// まず、環境一覧を取得して実際の環境IDを取得
	envsResp, err := client.GetEnvsWithResponse(ctx)
	if err != nil {
		t.Fatalf("環境一覧取得エラー: %v", err)
	}

	if envsResp == nil || envsResp.StatusCode() != 200 || envsResp.JSON200 == nil || len(envsResp.JSON200.Envs) == 0 {
		t.Skip("利用可能な環境がないため、環境更新テストをスキップします")
	}

	// 実際の環境IDを使用してテスト
	actualEnvId := EnvId(envsResp.JSON200.Envs[0].Id)
	t.Logf("テスト対象の環境ID: %d", actualEnvId)

	updateParam := UpdateEnvParam{
		Name:        testParams.Envs.CreateParams.Name,
		DisplayName: testParams.Envs.CreateParams.DisplayName,
	}

	resp, err := client.UpdateEnvWithResponse(ctx, actualEnvId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("環境の更新が成功しました")
	case 404:
		t.Log("環境が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_LinkAwsMarketplaceWithResponse は LinkAwsMarketplaceWithResponse のE2Eテスト
func TestE2E_LinkAwsMarketplaceWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// AWS Marketplace連携のテストパラメータ
	linkParam := LinkAwsMarketplaceParam{
		AccessToken:       "test_access_token",
		RegistrationToken: "test_aws_marketplace_token",
		TenantId:          "test_tenant_id",
	}

	resp, err := client.LinkAwsMarketplaceWithResponse(ctx, linkParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("AWS Marketplace連携が成功しました")
	case 400:
		t.Log("バリデーションエラーが発生しました（無効なトークンの可能性）")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_SignUpWithAwsMarketplaceWithResponse は SignUpWithAwsMarketplaceWithResponse のE2Eテスト
func TestE2E_SignUpWithAwsMarketplaceWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// テスト用のユニークなメールアドレスを生成
	timestamp := time.Now().Unix()
	testEmail := "aws-signup-" + string(rune(timestamp)) + "@example.com"

	signUpParam := SignUpWithAwsMarketplaceParam{
		Email:             testEmail,
		RegistrationToken: "test_aws_marketplace_token",
	}

	resp, err := client.SignUpWithAwsMarketplaceWithResponse(ctx, signUpParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		t.Log("AWS Marketplaceサインアップが成功しました")
	case 400:
		t.Log("バリデーションエラーが発生しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_ConfirmSignUpWithAwsMarketplaceWithResponse は ConfirmSignUpWithAwsMarketplaceWithResponse のE2Eテスト
func TestE2E_ConfirmSignUpWithAwsMarketplaceWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	confirmParam := ConfirmSignUpWithAwsMarketplaceParam{
		AccessToken:       "test_access_token",
		RegistrationToken: "test_registration_token",
		TenantName:        stringPtr("Test Tenant"),
	}

	resp, err := client.ConfirmSignUpWithAwsMarketplaceWithResponse(ctx, confirmParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("AWS Marketplaceサインアップ確認が成功しました")
	case 400:
		t.Log("バリデーションエラーが発生しました（無効な確認コードの可能性）")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_RequestExternalUserLinkWithResponse は RequestExternalUserLinkWithResponse のE2Eテスト
func TestE2E_RequestExternalUserLinkWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	requestParam := RequestExternalUserLinkParam{
		AccessToken: "test_access_token",
	}

	resp, err := client.RequestExternalUserLinkWithResponse(ctx, requestParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("外部ユーザー連携要求が成功しました")
	case 400:
		t.Log("バリデーションエラーが発生しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_ConfirmExternalUserLinkWithResponse は ConfirmExternalUserLinkWithResponse のE2Eテスト
func TestE2E_ConfirmExternalUserLinkWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	confirmParam := ConfirmExternalUserLinkParam{
		AccessToken: "test_access_token",
	}

	resp, err := client.ConfirmExternalUserLinkWithResponse(ctx, confirmParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("外部ユーザー連携確認が成功しました")
	case 400:
		t.Log("バリデーションエラーが発生しました（無効なアクセストークンの可能性）")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateIdentityProviderWithResponse は UpdateIdentityProviderWithResponse のE2Eテスト
func TestE2E_UpdateIdentityProviderWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	updateParam := UpdateIdentityProviderParam{
		Provider:              testParams.IdentityProviders.UpdateParams.Provider,
		IdentityProviderProps: testParams.IdentityProviders.UpdateParams.IdentityProviderProps,
	}

	resp, err := client.UpdateIdentityProviderWithResponse(ctx, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("IDプロバイダーの更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetAllTenantUsersWithResponse は GetAllTenantUsersWithResponse のE2Eテスト
func TestE2E_GetAllTenantUsersWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	resp, err := client.GetAllTenantUsersWithResponse(ctx)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("全テナントユーザー数: %d", len(resp.JSON200.Users))

		for i, user := range resp.JSON200.Users {
			if i >= 3 { // 最初の3件のみ詳細チェック
				break
			}
			t.Logf("ユーザー %d: ID=%s, Email=%s", i+1, user.Id, user.Email)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetAllTenantUserWithResponse は GetAllTenantUserWithResponse のE2Eテスト
func TestE2E_GetAllTenantUserWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、全テナントユーザーを取得して実際のユーザーIDを取得
	usersResp, err := client.GetAllTenantUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("全テナントユーザー取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なテナントユーザーがないため、個別ユーザー取得テストをスキップします")
	}

	// 実際のユーザーIDを使用してテスト
	actualUserId := UserId(usersResp.JSON200.Users[0].Id)
	t.Logf("テスト対象のユーザーID: %s", actualUserId)

	resp, err := client.GetAllTenantUserWithResponse(ctx, actualUserId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		// GetAllTenantUserResponseはUsersタイプを返すため、最初のユーザーを取得
		if len(resp.JSON200.Users) == 0 {
			t.Error("ユーザーが見つかりませんでした")
			return
		}
		user := resp.JSON200.Users[0]
		t.Logf("全テナントユーザー詳細: ID=%s, Email=%s", user.Id, user.Email)

		if user.Id != string(actualUserId) {
			t.Errorf("ユーザーIDが一致しません: 期待=%s, 実際=%s", actualUserId, user.Id)
		}
	case 404:
		t.Log("ユーザーが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateTenantWithResponse は UpdateTenantWithResponse のE2Eテスト
func TestE2E_UpdateTenantWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// まず、テナント一覧を取得して実際のテナントIDを取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナント更新テストをスキップします")
	}

	// 実際のテナントIDを使用してテスト
	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)
	t.Logf("テスト対象のテナントID: %s", actualTenantId)

	updateParam := UpdateTenantParam{
		Name:       testParams.Tenants.CreateParams.Name,
		Attributes: testParams.Tenants.CreateParams.Attributes,
	}

	resp, err := client.UpdateTenantWithResponse(ctx, actualTenantId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナントの更新が成功しました")
	case 404:
		t.Log("テナントが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetTenantUserWithResponse は GetTenantUserWithResponse のE2Eテスト
func TestE2E_GetTenantUserWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナントユーザー詳細取得テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	// テナントユーザー一覧を取得
	usersResp, err := client.GetTenantUsersWithResponse(ctx, actualTenantId)
	if err != nil {
		t.Fatalf("テナントユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なテナントユーザーがないため、個別ユーザー取得テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)
	t.Logf("テスト対象のテナントID: %s, ユーザーID: %s", actualTenantId, actualUserId)

	resp, err := client.GetTenantUserWithResponse(ctx, actualTenantId, actualUserId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("テナントユーザー詳細: ID=%s, Email=%s", resp.JSON200.Id, resp.JSON200.Email)
	case 404:
		t.Log("テナントユーザーが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateTenantUserWithResponse は UpdateTenantUserWithResponse のE2Eテスト
func TestE2E_UpdateTenantUserWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナントユーザー更新テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	// テナントユーザー一覧を取得
	usersResp, err := client.GetTenantUsersWithResponse(ctx, actualTenantId)
	if err != nil {
		t.Fatalf("テナントユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なテナントユーザーがないため、ユーザー更新テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)
	t.Logf("テスト対象のテナントID: %s, ユーザーID: %s", actualTenantId, actualUserId)

	updateParam := UpdateTenantUserParam{
		Attributes: testParams.TenantUsers.UpdateParams.Attributes,
	}

	resp, err := client.UpdateTenantUserWithResponse(ctx, actualTenantId, actualUserId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナントユーザーの更新が成功しました")
	case 404:
		t.Log("テナントユーザーが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_DeleteTenantUserWithResponse は DeleteTenantUserWithResponse のE2Eテスト
func TestE2E_DeleteTenantUserWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナントユーザー削除テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	// テスト用テナントユーザーを作成
	timestamp := time.Now().Unix()
	testEmail := "delete-tenant-user-" + string(rune(timestamp)) + "@example.com"

	createParam := CreateTenantUserParam{
		Email:      testEmail,
		Attributes: testParams.TenantUsers.CreateParams.Attributes,
	}

	createResp, err := client.CreateTenantUserWithResponse(ctx, actualTenantId, createParam)
	if err != nil || createResp == nil || createResp.StatusCode() != 201 || createResp.JSON201 == nil {
		t.Skip("テスト用テナントユーザーの作成に失敗したため、削除テストをスキップします")
	}

	createdUserId := UserId(createResp.JSON201.Id)
	t.Logf("削除対象のテナントID: %s, ユーザーID: %s", actualTenantId, createdUserId)

	// テナントユーザーを削除
	resp, err := client.DeleteTenantUserWithResponse(ctx, actualTenantId, createdUserId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナントユーザーの削除が成功しました")
	case 404:
		t.Log("テナントユーザーが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_ValidateInvitationWithResponse は ValidateInvitationWithResponse のE2Eテスト
func TestE2E_ValidateInvitationWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// テスト用の招待IDを使用
	testInvitationId := InvitationId("test-invitation-id")

	validateParam := ValidateInvitationParam{
		AccessToken: stringPtr("test_access_token"),
		Email:       stringPtr("test@example.com"),
		Password:    stringPtr("TestPassword123!"),
	}

	resp, err := client.ValidateInvitationWithResponse(ctx, testInvitationId, validateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("招待の検証が成功しました")
	case 400:
		t.Log("バリデーションエラーが発生しました")
	case 404:
		t.Log("招待が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetInvitationValidityWithResponse は GetInvitationValidityWithResponse のE2Eテスト
func TestE2E_GetInvitationValidityWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// テスト用の招待IDを使用
	testInvitationId := InvitationId("test-invitation-id")

	resp, err := client.GetInvitationValidityWithResponse(ctx, testInvitationId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("招待の有効性: IsValid=%t", resp.JSON200.IsValid)
	case 404:
		t.Log("招待が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateTenantBillingInfoWithResponse は UpdateTenantBillingInfoWithResponse のE2Eテスト
func TestE2E_UpdateTenantBillingInfoWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナント請求情報更新テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	updateParam := UpdateTenantBillingInfoParam{
		Name: "Test Billing Company",
		Address: BillingAddress{
			Street:     "123 Test Street",
			City:       "Test City",
			State:      "Test State",
			PostalCode: "12345",
			Country:    "US",
		},
		InvoiceLanguage: EnUS,
	}

	resp, err := client.UpdateTenantBillingInfoWithResponse(ctx, actualTenantId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナント請求情報の更新が成功しました")
	case 400:
		t.Log("バリデーションエラーが発生しました")
	case 404:
		t.Log("テナントが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetTenantIdentityProvidersWithResponse は GetTenantIdentityProvidersWithResponse のE2Eテスト
func TestE2E_GetTenantIdentityProvidersWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナントIDプロバイダー取得テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	resp, err := client.GetTenantIdentityProvidersWithResponse(ctx, actualTenantId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Log("テナントIDプロバイダー情報を取得しました")
		if resp.JSON200.Saml != nil {
			t.Logf("SAML設定: EmailAttribute=%s", resp.JSON200.Saml.EmailAttribute)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateTenantIdentityProviderWithResponse は UpdateTenantIdentityProviderWithResponse のE2Eテスト
func TestE2E_UpdateTenantIdentityProviderWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナントIDプロバイダー更新テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	updateParam := UpdateTenantIdentityProviderParam{
		ProviderType:          SAML,
		IdentityProviderProps: nil, // 無効化する場合
	}

	resp, err := client.UpdateTenantIdentityProviderWithResponse(ctx, actualTenantId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナントIDプロバイダーの更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetTenantInvitationsWithResponse は GetTenantInvitationsWithResponse のE2Eテスト
func TestE2E_GetTenantInvitationsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナント招待一覧取得テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	resp, err := client.GetTenantInvitationsWithResponse(ctx, actualTenantId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("テナント招待数: %d", len(resp.JSON200.Invitations))

		for i, invitation := range resp.JSON200.Invitations {
			if i >= 3 { // 最初の3件のみ詳細チェック
				break
			}
			t.Logf("招待 %d: ID=%s, Email=%s, Status=%s", i+1, invitation.Id, invitation.Email, invitation.Status)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateTenantInvitationWithResponse は CreateTenantInvitationWithResponse のE2Eテスト
func TestE2E_CreateTenantInvitationWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナント招待作成テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	// 環境一覧を取得
	envsResp, err := client.GetEnvsWithResponse(ctx)
	if err != nil || envsResp == nil || envsResp.StatusCode() != 200 || envsResp.JSON200 == nil || len(envsResp.JSON200.Envs) == 0 {
		t.Skip("利用可能な環境がないため、テナント招待作成テストをスキップします")
	}

	// テスト用のユニークなメールアドレスを生成
	timestamp := time.Now().Unix()
	testEmail := "invitation-" + string(rune(timestamp)) + "@example.com"

	createParam := CreateTenantInvitationParam{
		AccessToken: "test_access_token",
		Email:       testEmail,
		Envs: []struct {
			Id        Id       `json:"id"`
			RoleNames []string `json:"role_names"`
		}{
			{
				Id:        envsResp.JSON200.Envs[0].Id,
				RoleNames: []string{"admin"},
			},
		},
	}

	resp, err := client.CreateTenantInvitationWithResponse(ctx, actualTenantId, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("作成されたテナント招待: ID=%s, Email=%s", resp.JSON201.Id, resp.JSON201.Email)
	case 400:
		t.Log("バリデーションエラーが発生しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetTenantInvitationWithResponse は GetTenantInvitationWithResponse のE2Eテスト
func TestE2E_GetTenantInvitationWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナント招待詳細取得テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	// テナント招待一覧を取得
	invitationsResp, err := client.GetTenantInvitationsWithResponse(ctx, actualTenantId)
	if err != nil {
		t.Fatalf("テナント招待一覧取得エラー: %v", err)
	}

	if invitationsResp == nil || invitationsResp.StatusCode() != 200 || invitationsResp.JSON200 == nil || len(invitationsResp.JSON200.Invitations) == 0 {
		t.Skip("利用可能なテナント招待がないため、個別招待取得テストをスキップします")
	}

	actualInvitationId := InvitationId(invitationsResp.JSON200.Invitations[0].Id)
	t.Logf("テスト対象のテナントID: %s, 招待ID: %s", actualTenantId, actualInvitationId)

	resp, err := client.GetTenantInvitationWithResponse(ctx, actualTenantId, actualInvitationId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("テナント招待詳細: ID=%s, Email=%s, Status=%s", resp.JSON200.Id, resp.JSON200.Email, resp.JSON200.Status)
	case 404:
		t.Log("テナント招待が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_DeleteTenantInvitationWithResponse は DeleteTenantInvitationWithResponse のE2Eテスト
func TestE2E_DeleteTenantInvitationWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナント招待削除テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	// テスト用の招待IDを使用（実際の招待を作成してから削除するのが理想的）
	testInvitationId := InvitationId("test-invitation-id")

	resp, err := client.DeleteTenantInvitationWithResponse(ctx, actualTenantId, testInvitationId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナント招待の削除が成功しました")
	case 404:
		t.Log("テナント招待が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateTenantPlanWithResponse は UpdateTenantPlanWithResponse のE2Eテスト
func TestE2E_UpdateTenantPlanWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナントプラン更新テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	prorationBehavior := CreateProrations
	deleteUsage := false
	updateParam := UpdateTenantPlanParam{
		NextPlanId:        stringPtr("test-plan-id"),
		ProrationBehavior: &prorationBehavior,
		DeleteUsage:       &deleteUsage,
	}

	resp, err := client.UpdateTenantPlanWithResponse(ctx, actualTenantId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナントプランの更新が成功しました")
	case 400:
		t.Log("バリデーションエラーが発生しました")
	case 404:
		t.Log("テナントが見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetStripeCustomerWithResponse は GetStripeCustomerWithResponse のE2Eテスト
func TestE2E_GetStripeCustomerWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、Stripe顧客情報取得テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	resp, err := client.GetStripeCustomerWithResponse(ctx, actualTenantId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("Stripe顧客情報: CustomerID=%s", resp.JSON200.CustomerId)
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateTenantUserRolesWithResponse は CreateTenantUserRolesWithResponse のE2Eテスト
func TestE2E_CreateTenantUserRolesWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナントユーザー役割作成テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	// テナントユーザー一覧を取得
	usersResp, err := client.GetTenantUsersWithResponse(ctx, actualTenantId)
	if err != nil {
		t.Fatalf("テナントユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なテナントユーザーがないため、ユーザー役割作成テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	// 環境一覧を取得
	envsResp, err := client.GetEnvsWithResponse(ctx)
	if err != nil || envsResp == nil || envsResp.StatusCode() != 200 || envsResp.JSON200 == nil || len(envsResp.JSON200.Envs) == 0 {
		t.Skip("利用可能な環境がないため、ユーザー役割作成テストをスキップします")
	}

	actualEnvId := EnvId(envsResp.JSON200.Envs[0].Id)

	createParam := CreateTenantUserRolesParam{
		RoleNames: []string{"admin", "user"},
	}

	resp, err := client.CreateTenantUserRolesWithResponse(ctx, actualTenantId, actualUserId, actualEnvId, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナントユーザー役割の作成が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_DeleteTenantUserRoleWithResponse は DeleteTenantUserRoleWithResponse のE2Eテスト
func TestE2E_DeleteTenantUserRoleWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、テナント一覧を取得
	tenantsResp, err := client.GetTenantsWithResponse(ctx)
	if err != nil {
		t.Fatalf("テナント一覧取得エラー: %v", err)
	}

	if tenantsResp == nil || tenantsResp.StatusCode() != 200 || tenantsResp.JSON200 == nil || len(tenantsResp.JSON200.Tenants) == 0 {
		t.Skip("利用可能なテナントがないため、テナントユーザー役割削除テストをスキップします")
	}

	actualTenantId := TenantId(tenantsResp.JSON200.Tenants[0].Id)

	// テナントユーザー一覧を取得
	usersResp, err := client.GetTenantUsersWithResponse(ctx, actualTenantId)
	if err != nil {
		t.Fatalf("テナントユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なテナントユーザーがないため、ユーザー役割削除テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	// 環境一覧を取得
	envsResp, err := client.GetEnvsWithResponse(ctx)
	if err != nil || envsResp == nil || envsResp.StatusCode() != 200 || envsResp.JSON200 == nil || len(envsResp.JSON200.Envs) == 0 {
		t.Skip("利用可能な環境がないため、ユーザー役割削除テストをスキップします")
	}

	actualEnvId := EnvId(envsResp.JSON200.Envs[0].Id)
	testRoleName := RoleName("admin")

	resp, err := client.DeleteTenantUserRoleWithResponse(ctx, actualTenantId, actualUserId, actualEnvId, testRoleName)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("テナントユーザー役割の削除が成功しました")
	case 404:
		t.Log("テナントユーザー役割が見つかりませんでした")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetUserInfoWithResponse は GetUserInfoWithResponse のE2Eテスト
func TestE2E_GetUserInfoWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)
	ctx := context.Background()

	// テストケースを実行
	for _, testCase := range testParams.GetUserInfo.TestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			params := &GetUserInfoParams{
				Token: testCase.Params.Token,
			}

			resp, err := client.GetUserInfoWithResponse(ctx, params)

			if err != nil {
				t.Logf("リクエストエラー: %v", err)
				return
			}

			if resp == nil {
				t.Fatal("レスポンスがnilです")
			}

			t.Logf("ステータスコード: %d", resp.StatusCode())

			switch resp.StatusCode() {
			case 200:
				if resp.JSON200 == nil {
					t.Error("200レスポンスが解析されませんでした")
					return
				}
				t.Logf("ユーザー情報取得成功: ID=%s, Email=%s", resp.JSON200.Id, resp.JSON200.Email)
			case 401:
				t.Log("認証エラーが発生しました（期待される動作）")
			case 500:
				if resp.JSON500 != nil {
					t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
				}
				t.Error("サーバーエラーが発生しました")
			default:
				t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
			}
		})
	}

	// エッジケースのテスト
	for _, edgeCase := range testParams.GetUserInfo.EdgeCases {
		t.Run("EdgeCase_"+edgeCase.Name, func(t *testing.T) {
			params := &GetUserInfoParams{
				Token: edgeCase.Params.Token,
			}

			resp, err := client.GetUserInfoWithResponse(ctx, params)

			if err != nil {
				if edgeCase.ExpectError {
					t.Logf("期待通りエラーが発生: %v", err)
				} else {
					t.Errorf("予期しないエラー: %v", err)
				}
				return
			}

			if resp != nil && edgeCase.ExpectError && (resp.StatusCode() == 400 || resp.StatusCode() == 401) {
				t.Logf("期待通りエラーレスポンス: %d", resp.StatusCode())
			}
		})
	}
}

// TestE2E_UpdateSaasUserAttributesWithResponse は UpdateSaasUserAttributesWithResponse のE2Eテスト
func TestE2E_UpdateSaasUserAttributesWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、ユーザー属性更新テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	updateParam := UpdateSaasUserAttributesParam{
		Attributes: map[string]interface{}{
			"department": "Engineering",
			"role":       "Developer",
		},
	}

	resp, err := client.UpdateSaasUserAttributesWithResponse(ctx, actualUserId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("SaaSユーザー属性の更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateSaasUserEmailWithResponse は UpdateSaasUserEmailWithResponse のE2Eテスト
func TestE2E_UpdateSaasUserEmailWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、ユーザーメール更新テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	// テスト用のユニークなメールアドレスを生成
	timestamp := time.Now().Unix()
	newEmail := "updated-" + string(rune(timestamp)) + "@example.com"

	updateParam := UpdateSaasUserEmailParam{
		Email: newEmail,
	}

	resp, err := client.UpdateSaasUserEmailWithResponse(ctx, actualUserId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("SaaSユーザーメールアドレスの更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_ConfirmEmailUpdateWithResponse は ConfirmEmailUpdateWithResponse のE2Eテスト
func TestE2E_ConfirmEmailUpdateWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、メール更新確認テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	confirmParam := ConfirmEmailUpdateParam{
		AccessToken: "test_access_token",
		Code:        "123456",
	}

	resp, err := client.ConfirmEmailUpdateWithResponse(ctx, actualUserId, confirmParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("メール更新確認が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_RequestEmailUpdateWithResponse は RequestEmailUpdateWithResponse のE2Eテスト
func TestE2E_RequestEmailUpdateWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、メール更新要求テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	// テスト用のユニークなメールアドレスを生成
	timestamp := time.Now().Unix()
	newEmail := "request-update-" + string(rune(timestamp)) + "@example.com"

	requestParam := RequestEmailUpdateParam{
		AccessToken: "test_access_token",
		Email:       newEmail,
	}

	resp, err := client.RequestEmailUpdateWithResponse(ctx, actualUserId, requestParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("メール更新要求が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_GetUserMfaPreferenceWithResponse は GetUserMfaPreferenceWithResponse のE2Eテスト
func TestE2E_GetUserMfaPreferenceWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、ユーザーMFA設定取得テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	resp, err := client.GetUserMfaPreferenceWithResponse(ctx, actualUserId)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		if resp.JSON200 == nil {
			t.Error("200レスポンスが解析されませんでした")
			return
		}
		t.Logf("ユーザーMFA設定: Enabled=%t", resp.JSON200.Enabled)
		if resp.JSON200.Method != nil {
			t.Logf("MFA方法: %s", *resp.JSON200.Method)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateUserMfaPreferenceWithResponse は UpdateUserMfaPreferenceWithResponse のE2Eテスト
func TestE2E_UpdateUserMfaPreferenceWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、ユーザーMFA設定更新テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	mfaMethod := SoftwareToken
	updateParam := UpdateUserMfaPreferenceParam{
		Enabled: true,
		Method:  &mfaMethod,
	}

	resp, err := client.UpdateUserMfaPreferenceWithResponse(ctx, actualUserId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("ユーザーMFA設定の更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateSoftwareTokenWithResponse は UpdateSoftwareTokenWithResponse のE2Eテスト
func TestE2E_UpdateSoftwareTokenWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、ソフトウェアトークン更新テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	updateParam := UpdateSoftwareTokenParam{
		AccessToken:      "test_access_token",
		VerificationCode: "123456",
	}

	resp, err := client.UpdateSoftwareTokenWithResponse(ctx, actualUserId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("ソフトウェアトークンの更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateSecretCodeWithResponse は CreateSecretCodeWithResponse のE2Eテスト
func TestE2E_CreateSecretCodeWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、シークレットコード作成テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	createParam := CreateSecretCodeParam{
		AccessToken: "test_access_token",
	}

	resp, err := client.CreateSecretCodeWithResponse(ctx, actualUserId, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("シークレットコード作成成功: SecretCode=%s", resp.JSON201.SecretCode)
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UpdateSaasUserPasswordWithResponse は UpdateSaasUserPasswordWithResponse のE2Eテスト
func TestE2E_UpdateSaasUserPasswordWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、ユーザーパスワード更新テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)

	updateParam := UpdateSaasUserPasswordParam{
		Password: "NewPassword123!",
	}

	resp, err := client.UpdateSaasUserPasswordWithResponse(ctx, actualUserId, updateParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("SaaSユーザーパスワードの更新が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_UnlinkProviderWithResponse は UnlinkProviderWithResponse のE2Eテスト
func TestE2E_UnlinkProviderWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// まず、ユーザー一覧を取得して実際のユーザーIDを取得
	usersResp, err := client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}

	if usersResp == nil || usersResp.StatusCode() != 200 || usersResp.JSON200 == nil || len(usersResp.JSON200.Users) == 0 {
		t.Skip("利用可能なSaaSユーザーがないため、プロバイダー連携解除テストをスキップします")
	}

	actualUserId := UserId(usersResp.JSON200.Users[0].Id)
	providerName := "Google"

	resp, err := client.UnlinkProviderWithResponse(ctx, actualUserId, providerName)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 200:
		t.Log("プロバイダー連携解除が成功しました")
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// TestE2E_CreateSaasUserAttributeWithResponse は CreateSaasUserAttributeWithResponse のE2Eテスト
func TestE2E_CreateSaasUserAttributeWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	ctx := context.Background()

	// テスト用のユニークな属性名を生成
	timestamp := time.Now().Unix()
	testAttributeName := "test_saas_user_attr_" + string(rune(timestamp))

	createParam := CreateSaasUserAttributeParam{
		AttributeName: testAttributeName,
		DisplayName:   "Test SaaS User Attribute",
		AttributeType: String,
	}

	resp, err := client.CreateSaasUserAttributeWithResponse(ctx, createParam)

	if err != nil {
		t.Fatalf("リクエストエラー: %v", err)
	}

	if resp == nil {
		t.Fatal("レスポンスがnilです")
	}

	t.Logf("ステータスコード: %d", resp.StatusCode())

	switch resp.StatusCode() {
	case 201:
		if resp.JSON201 == nil {
			t.Error("201レスポンスが解析されませんでした")
			return
		}
		t.Logf("作成されたSaaSユーザー属性: Name=%s, DisplayName=%s, Type=%s",
			resp.JSON201.AttributeName, resp.JSON201.DisplayName, resp.JSON201.AttributeType)

		if resp.JSON201.AttributeName != testAttributeName {
			t.Errorf("属性名が一致しません: 期待=%s, 実際=%s", testAttributeName, resp.JSON201.AttributeName)
		}
	case 500:
		if resp.JSON500 != nil {
			t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
		}
		t.Error("サーバーエラーが発生しました")
	default:
		t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}
