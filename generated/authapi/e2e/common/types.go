package common

import (
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
)

// TestConfig はE2Eテストの設定を表します
type TestConfig struct {
	APIEndpoint string        `yaml:"api_endpoint"`
	Timeout     time.Duration `yaml:"timeout"`
	RetryCount  int           `yaml:"retry_count"`
	Verbose     bool          `yaml:"verbose"`
}

// TestData は各ストーリーのテストデータを表します
type TestData struct {
	Story       string                 `yaml:"story"`
	Description string                 `yaml:"description"`
	TestCases   []TestCase             `yaml:"test_cases"`
	Setup       map[string]interface{} `yaml:"setup,omitempty"`
	Cleanup     map[string]interface{} `yaml:"cleanup,omitempty"`
}

// TestCase は個別のテストケースを表します
type TestCase struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Method      string                 `yaml:"method"`
	Endpoint    string                 `yaml:"endpoint"`
	Headers     map[string]string      `yaml:"headers,omitempty"`
	Params      map[string]interface{} `yaml:"params,omitempty"`
	Body        map[string]interface{} `yaml:"body,omitempty"`
	Expected    ExpectedResult         `yaml:"expected"`
	Retry       *RetryConfig           `yaml:"retry,omitempty"`
	Timeout     *time.Duration         `yaml:"timeout,omitempty"`
}

// ExpectedResult は期待される結果を表します
type ExpectedResult struct {
	StatusCode int                    `yaml:"status_code"`
	Headers    map[string]string      `yaml:"headers,omitempty"`
	Body       map[string]interface{} `yaml:"body,omitempty"`
	Schema     string                 `yaml:"schema,omitempty"`
	Error      *ExpectedError         `yaml:"error,omitempty"`
}

// ExpectedError は期待されるエラーを表します
type ExpectedError struct {
	Type    string `yaml:"type"`
	Message string `yaml:"message"`
}

// RetryConfig はリトライ設定を表します
type RetryConfig struct {
	Count    int           `yaml:"count"`
	Interval time.Duration `yaml:"interval"`
}

// BasicSetupTestData は基本セットアップストーリーのテストデータを表します
type BasicSetupTestData struct {
	BasicInfo struct {
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
		Update struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				DomainName        string  `yaml:"domain_name"`
				FromEmailAddress  string  `yaml:"from_email_address"`
				ReplyEmailAddress *string `yaml:"reply_email_address,omitempty"`
			} `yaml:"params"`
		} `yaml:"update"`
	} `yaml:"basic_info"`

	AuthInfo struct {
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
		Update struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				CallbackUrl string `yaml:"callback_url"`
			} `yaml:"params"`
		} `yaml:"update"`
	} `yaml:"auth_info"`

	Envs struct {
		List struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"list"`
		Create struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				Id          string  `yaml:"id"`
				Name        string  `yaml:"name"`
				DisplayName *string `yaml:"display_name,omitempty"`
			} `yaml:"params"`
		} `yaml:"create"`
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
		Update struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				Name        string  `yaml:"name"`
				DisplayName *string `yaml:"display_name,omitempty"`
			} `yaml:"params"`
		} `yaml:"update"`
		Delete struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"delete"`
	} `yaml:"envs"`

	TestEnvIds    []string `yaml:"test_env_ids"`
	InvalidEnvIds []string `yaml:"invalid_env_ids"`
}

// RoleManagementTestData はロール管理ストーリーのテストデータを表します
type RoleManagementTestData struct {
	Roles struct {
		List struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"list"`
		Create struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    []struct {
				RoleName    string `yaml:"role_name"`
				DisplayName string `yaml:"display_name"`
			} `yaml:"params"`
		} `yaml:"create"`
		Delete struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"delete"`
	} `yaml:"roles"`

	TestRoleNames    []string `yaml:"test_role_names"`
	InvalidRoleNames []string `yaml:"invalid_role_names"`
}

// UserManagementTestData はユーザー管理ストーリーのテストデータを表します
type UserManagementTestData struct {
	UserAttributes struct {
		List struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"list"`
		Create struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    []struct {
				AttributeName string `yaml:"attribute_name"`
				DisplayName   string `yaml:"display_name"`
			} `yaml:"params"`
		} `yaml:"create"`
		Delete struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"delete"`
	} `yaml:"user_attributes"`

	SaasUsers struct {
		List struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"list"`
		Create struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    []struct {
				Email      string                 `yaml:"email"`
				Attributes map[string]interface{} `yaml:"attributes,omitempty"`
			} `yaml:"params"`
		} `yaml:"create"`
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
		Update struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				Attributes map[string]interface{} `yaml:"attributes"`
			} `yaml:"params"`
		} `yaml:"update"`
		Delete struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"delete"`
	} `yaml:"saas_users"`

	UserInfo struct {
		UpdateEmail struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				Email string `yaml:"email"`
			} `yaml:"params"`
		} `yaml:"update_email"`
		UpdatePassword struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				Password string `yaml:"password"`
			} `yaml:"params"`
		} `yaml:"update_password"`
	} `yaml:"user_info"`

	TestUserEmails       []string `yaml:"test_user_emails"`
	TestAttributeNames   []string `yaml:"test_attribute_names"`
	InvalidUserIds       []string `yaml:"invalid_user_ids"`
	InvalidAttributeNames []string `yaml:"invalid_attribute_names"`
}

// TenantManagementTestData はテナント管理ストーリーのテストデータを表します
type TenantManagementTestData struct {
	TenantAttributes struct {
		List struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"list"`
		Create struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    []struct {
				AttributeName string `yaml:"attribute_name"`
				DisplayName   string `yaml:"display_name"`
			} `yaml:"params"`
		} `yaml:"create"`
		Delete struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"delete"`
	} `yaml:"tenant_attributes"`

	Tenants struct {
		List struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"list"`
		Create struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    []struct {
				Name                   string                 `yaml:"name"`
				BackOfficeStaffEmail   string                 `yaml:"back_office_staff_email"`
				Attributes             map[string]interface{} `yaml:"attributes,omitempty"`
			} `yaml:"params"`
		} `yaml:"create"`
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
		Update struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				Name                   string                 `yaml:"name"`
				BackOfficeStaffEmail   string                 `yaml:"back_office_staff_email"`
				Attributes             map[string]interface{} `yaml:"attributes"`
			} `yaml:"params"`
		} `yaml:"update"`
		Delete struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"delete"`
	} `yaml:"tenants"`

	BillingInfo struct {
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
		Update struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				Name            string `yaml:"name"`
				Street          string `yaml:"street"`
				City            string `yaml:"city"`
				State           string `yaml:"state"`
				Country         string `yaml:"country"`
				PostalCode      string `yaml:"postal_code"`
				InvoiceLanguage string `yaml:"invoice_language"`
			} `yaml:"params"`
		} `yaml:"update"`
	} `yaml:"billing_info"`

	TestTenantNames      []string `yaml:"test_tenant_names"`
	TestAttributeNames   []string `yaml:"test_attribute_names"`
	InvalidTenantIds     []string `yaml:"invalid_tenant_ids"`
	InvalidAttributeNames []string `yaml:"invalid_attribute_names"`
}

// TenantUserManagementTestData はテナントユーザー管理ストーリーのテストデータを表します
type TenantUserManagementTestData struct {
	TenantUsers struct {
		List struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"list"`
		Create struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    []struct {
				UserID string                 `yaml:"user_id"`
				Roles  []string               `yaml:"roles,omitempty"`
				Envs   []string               `yaml:"envs,omitempty"`
			} `yaml:"params"`
		} `yaml:"create"`
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
		Update struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				Roles []string `yaml:"roles"`
				Envs  []string `yaml:"envs"`
			} `yaml:"params"`
		} `yaml:"update"`
		Delete struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"delete"`
	} `yaml:"tenant_users"`

	AllTenantUsers struct {
		List struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"list"`
		Create struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    []struct {
				Email      string                 `yaml:"email"`
				Attributes map[string]interface{} `yaml:"attributes,omitempty"`
				Roles      []string               `yaml:"roles,omitempty"`
				Envs       []string               `yaml:"envs,omitempty"`
			} `yaml:"params"`
		} `yaml:"create"`
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
		Update struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				Attributes map[string]interface{} `yaml:"attributes"`
			} `yaml:"params"`
		} `yaml:"update"`
		Delete struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"delete"`
	} `yaml:"all_tenant_users"`

	RoleManagement struct {
		AttachRole struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				EnvID    string `yaml:"env_id"`
				RoleName string `yaml:"role_name"`
			} `yaml:"params"`
		} `yaml:"attach_role"`
		DetachRole struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				EnvID    string `yaml:"env_id"`
				RoleName string `yaml:"role_name"`
			} `yaml:"params"`
		} `yaml:"detach_role"`
	} `yaml:"role_management"`

	TestUserEmails       []string `yaml:"test_user_emails"`
	TestTenantIds        []string `yaml:"test_tenant_ids"`
	TestEnvIds           []string `yaml:"test_env_ids"`
	TestRoleNames        []string `yaml:"test_role_names"`
	InvalidUserIds       []string `yaml:"invalid_user_ids"`
	InvalidTenantIds     []string `yaml:"invalid_tenant_ids"`
}

// InvitationManagementTestData は招待管理ストーリーのテストデータを表します
type InvitationManagementTestData struct {
	Invitations struct {
		List struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"list"`
		Create struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    []struct {
				Email string   `yaml:"email"`
				Roles []string `yaml:"roles,omitempty"`
				Envs  []string `yaml:"envs,omitempty"`
			} `yaml:"params"`
		} `yaml:"create"`
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
		Delete struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"delete"`
	} `yaml:"invitations"`

	InvitationValidation struct {
		ValidateCode struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				InvitationID string `yaml:"invitation_id"`
				Code         string `yaml:"code"`
			} `yaml:"params"`
		} `yaml:"validate_code"`
		ValidateInvitation struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				InvitationID string `yaml:"invitation_id"`
			} `yaml:"params"`
		} `yaml:"validate_invitation"`
	} `yaml:"invitation_validation"`

	TestInvitationEmails []string `yaml:"test_invitation_emails"`
	TestTenantIds        []string `yaml:"test_tenant_ids"`
	TestEnvIds           []string `yaml:"test_env_ids"`
	TestRoleNames        []string `yaml:"test_role_names"`
	InvalidInvitationIds []string `yaml:"invalid_invitation_ids"`
	InvalidCodes         []string `yaml:"invalid_codes"`
}

// AuthenticationFlowTestData は認証フローストーリーのテストデータを表します
type AuthenticationFlowTestData struct {
	Credentials struct {
		Create struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				IDToken      string `yaml:"id_token"`
				AccessToken  string `yaml:"access_token"`
				RefreshToken string `yaml:"refresh_token"`
			} `yaml:"params"`
		} `yaml:"create"`
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
			TempCodeAuth struct {
				TestCases []TestCase `yaml:"test_cases"`
			} `yaml:"temp_code_auth"`
			RefreshTokenAuth struct {
				TestCases []TestCase `yaml:"test_cases"`
			} `yaml:"refresh_token_auth"`
		} `yaml:"get"`
	} `yaml:"credentials"`

	ErrorHandling struct {
		InvalidTokens struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				InvalidIDToken      string `yaml:"invalid_id_token"`
				InvalidAccessToken  string `yaml:"invalid_access_token"`
				InvalidRefreshToken string `yaml:"invalid_refresh_token"`
			} `yaml:"params"`
		} `yaml:"invalid_tokens"`
		ExpiredCodes struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"expired_codes"`
	} `yaml:"error_handling"`

	TestTokens struct {
		ValidIDToken      string `yaml:"valid_id_token"`
		ValidAccessToken  string `yaml:"valid_access_token"`
		ValidRefreshToken string `yaml:"valid_refresh_token"`
	} `yaml:"test_tokens"`
}

// SingleTenantManagementTestData はシングルテナント管理ストーリーのテストデータを表します
type SingleTenantManagementTestData struct {
	SingleTenantSettings struct {
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
		Update struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				CustomizePageTitle       string `yaml:"customize_page_title"`
				CustomizePageUrl         string `yaml:"customize_page_url"`
				CustomizeCssUrl          string `yaml:"customize_css_url"`
				CustomizeIconUrl         string `yaml:"customize_icon_url"`
				CustomizeLogoUrl         string `yaml:"customize_logo_url"`
				PrivacyPolicyUrl         string `yaml:"privacy_policy_url"`
				TermsOfServiceUrl        string `yaml:"terms_of_service_url"`
			} `yaml:"params"`
		} `yaml:"update"`
	} `yaml:"single_tenant_settings"`

	CloudFormationTemplate struct {
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
	} `yaml:"cloudformation_template"`

	DDLTemplate struct {
		Get struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"get"`
	} `yaml:"ddl_template"`

	AWSIntegration struct {
		IAMRole struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				RoleArn    string `yaml:"role_arn"`
				ExternalID string `yaml:"external_id"`
			} `yaml:"params"`
		} `yaml:"iam_role"`
		ExternalID struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"external_id"`
	} `yaml:"aws_integration"`

	TestSettings struct {
		PageTitle    string `yaml:"page_title"`
		PageUrl      string `yaml:"page_url"`
		CssUrl       string `yaml:"css_url"`
		IconUrl      string `yaml:"icon_url"`
		LogoUrl      string `yaml:"logo_url"`
		PrivacyUrl   string `yaml:"privacy_url"`
		TermsUrl     string `yaml:"terms_url"`
	} `yaml:"test_settings"`
}

// ErrorHandlingTestData はエラーハンドリングストーリーのテストデータを表します
type ErrorHandlingTestData struct {
	ErrorEndpoints struct {
		InternalServerError struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"internal_server_error"`
		BadRequest struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"bad_request"`
		Unauthorized struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"unauthorized"`
		Forbidden struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"forbidden"`
		NotFound struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"not_found"`
	} `yaml:"error_endpoints"`

	ErrorResponseValidation struct {
		ResponseFormat struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"response_format"`
		ErrorMessages struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"error_messages"`
	} `yaml:"error_response_validation"`

	NetworkErrors struct {
		Timeout struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"timeout"`
		ConnectionError struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"connection_error"`
	} `yaml:"network_errors"`

	ExpectedErrors []struct {
		StatusCode int    `yaml:"status_code"`
		Type       string `yaml:"type"`
		Message    string `yaml:"message"`
	} `yaml:"expected_errors"`
}

// IntegrationTestData は統合テストストーリーのテストデータを表します
type IntegrationTestData struct {
	FullSetupFlow struct {
		BasicSetup struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"basic_setup"`
		UserCreation struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"user_creation"`
		TenantCreation struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"tenant_creation"`
		RoleAssignment struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"role_assignment"`
		InvitationFlow struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"invitation_flow"`
	} `yaml:"full_setup_flow"`

	EndToEndScenarios struct {
		UserJourney struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"user_journey"`
		AdminWorkflow struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"admin_workflow"`
		MultiTenantScenario struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"multi_tenant_scenario"`
	} `yaml:"end_to_end_scenarios"`

	PerformanceTests struct {
		LoadTest struct {
			TestCases []TestCase `yaml:"test_cases"`
			Params    struct {
				ConcurrentUsers int `yaml:"concurrent_users"`
				Duration        int `yaml:"duration_seconds"`
			} `yaml:"params"`
		} `yaml:"load_test"`
		StressTest struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"stress_test"`
	} `yaml:"performance_tests"`

	FailureRecovery struct {
		ServiceRecovery struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"service_recovery"`
		DataConsistency struct {
			TestCases []TestCase `yaml:"test_cases"`
		} `yaml:"data_consistency"`
	} `yaml:"failure_recovery"`

	TestScenarios []struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Steps       []string `yaml:"steps"`
		Expected    string   `yaml:"expected"`
	} `yaml:"test_scenarios"`
}

// TestResource はテスト中に作成されるリソースを追跡するための構造体です
type TestResource struct {
	Type       string                 `json:"type"`
	ID         string                 `json:"id"`
	Name       string                 `json:"name,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	Story      string                 `json:"story"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CleanedUp  bool                   `json:"cleaned_up"`
	CleanupErr string                 `json:"cleanup_error,omitempty"`
}

// TestResult はテスト結果を表します
type TestResult struct {
	TestName     string        `json:"test_name"`
	Story        string        `json:"story"`
	Success      bool          `json:"success"`
	Duration     time.Duration `json:"duration"`
	Error        string        `json:"error,omitempty"`
	Response     *APIResponse  `json:"response,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	ResourcesUsed []string     `json:"resources_used,omitempty"`
}

// APIResponse はAPIレスポンスの詳細を表します
type APIResponse struct {
	StatusCode int                    `json:"status_code"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body,omitempty"`
	Duration   time.Duration          `json:"duration"`
}

// ClientWrapper は認証付きクライアントのラッパーです
type ClientWrapper struct {
	Client    *authapi.ClientWithResponses
	Config    *TestConfig
	Resources []*TestResource
}

// AssertionResult はアサーション結果を表します
type AssertionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Expected interface{} `json:"expected,omitempty"`
	Actual   interface{} `json:"actual,omitempty"`
}

// ValidationError はバリデーションエラーを表します
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// CleanupResult はクリーンアップ結果を表します
type CleanupResult struct {
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
	Duration     time.Duration `json:"duration"`
}

// TestSummary はテスト実行の要約を表します
type TestSummary struct {
	TotalTests    int           `json:"total_tests"`
	PassedTests   int           `json:"passed_tests"`
	FailedTests   int           `json:"failed_tests"`
	SkippedTests  int           `json:"skipped_tests"`
	TotalDuration time.Duration `json:"total_duration"`
	Stories       []StorySummary `json:"stories"`
	CreatedAt     time.Time     `json:"created_at"`
}

// StorySummary はストーリーテストの要約を表します
type StorySummary struct {
	Name          string        `json:"name"`
	TestCount     int           `json:"test_count"`
	PassedCount   int           `json:"passed_count"`
	FailedCount   int           `json:"failed_count"`
	Duration      time.Duration `json:"duration"`
	ResourcesUsed int           `json:"resources_used"`
}

// Constants for resource types
const (
	ResourceTypeEnv            = "env"
	ResourceTypeRole           = "role"
	ResourceTypeSaasUser       = "saas_user"
	ResourceTypeTenant         = "tenant"
	ResourceTypeTenantUser     = "tenant_user"
	ResourceTypeUserAttribute  = "user_attribute"
	ResourceTypeTenantAttribute = "tenant_attribute"
	ResourceTypeInvitation     = "invitation"
)

// Constants for test stories
const (
	StoryBasicSetup         = "01_basic_setup"
	StoryUserManagement     = "02_user_management"
	StoryRoleManagement     = "03_role_management"
	StoryTenantManagement   = "04_tenant_management"
	StoryTenantUserManagement = "05_tenant_user_management"
	StoryInvitationManagement = "06_invitation_management"
	StoryAuthenticationFlow = "07_authentication_flow"
	StorySingleTenantManagement = "08_single_tenant_management"
	StoryErrorHandling      = "09_error_handling"
	StoryIntegrationTest    = "10_integration_test"
)
// StringPtr は文字列のポインタを返すヘルパー関数です
func StringPtr(s string) *string {
	return &s
}