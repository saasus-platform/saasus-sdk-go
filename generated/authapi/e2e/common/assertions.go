package common

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
)

// AssertionHelper はカスタムアサーションのヘルパー関数を提供します
type AssertionHelper struct {
	t       *testing.T
	verbose bool
}

// NewAssertionHelper は新しいアサーションヘルパーを作成します
func NewAssertionHelper(t *testing.T, verbose bool) *AssertionHelper {
	return &AssertionHelper{
		t:       t,
		verbose: verbose,
	}
}

// AssertStatusCode はHTTPステータスコードをアサートします
func (ah *AssertionHelper) AssertStatusCode(expected, actual int, message string) bool {
	if expected != actual {
		ah.t.Errorf("%s: 期待されるステータスコード %d, 実際 %d", message, expected, actual)
		return false
	}
	if ah.verbose {
		ah.t.Logf("✓ ステータスコード確認: %d", actual)
	}
	return true
}

// AssertNotNil はオブジェクトがnilでないことをアサートします
func (ah *AssertionHelper) AssertNotNil(obj interface{}, fieldName string) bool {
	if obj == nil {
		ah.t.Errorf("%s がnilです", fieldName)
		return false
	}
	if ah.verbose {
		ah.t.Logf("✓ %s は非nil", fieldName)
	}
	return true
}

// AssertNotEmpty は文字列が空でないことをアサートします
func (ah *AssertionHelper) AssertNotEmpty(value, fieldName string) bool {
	if value == "" {
		ah.t.Errorf("%s が空です", fieldName)
		return false
	}
	if ah.verbose {
		ah.t.Logf("✓ %s は非空: %s", fieldName, value)
	}
	return true
}

// AssertEquals は値が等しいことをアサートします
func (ah *AssertionHelper) AssertEquals(expected, actual interface{}, fieldName string) bool {
	if !reflect.DeepEqual(expected, actual) {
		ah.t.Errorf("%s: 期待値 %v, 実際 %v", fieldName, expected, actual)
		return false
	}
	if ah.verbose {
		ah.t.Logf("✓ %s 一致: %v", fieldName, actual)
	}
	return true
}

// AssertContains は文字列に部分文字列が含まれることをアサートします
func (ah *AssertionHelper) AssertContains(haystack, needle, fieldName string) bool {
	if !strings.Contains(haystack, needle) {
		ah.t.Errorf("%s: '%s' に '%s' が含まれていません", fieldName, haystack, needle)
		return false
	}
	if ah.verbose {
		ah.t.Logf("✓ %s に '%s' が含まれています", fieldName, needle)
	}
	return true
}

// AssertGreaterThan は値が指定値より大きいことをアサートします
func (ah *AssertionHelper) AssertGreaterThan(actual, threshold int, fieldName string) bool {
	if actual <= threshold {
		ah.t.Errorf("%s: %d は %d より大きくありません", fieldName, actual, threshold)
		return false
	}
	if ah.verbose {
		ah.t.Logf("✓ %s は閾値より大きい: %d > %d", fieldName, actual, threshold)
	}
	return true
}

// AssertBasicInfoResponse は基本情報レスポンスをアサートします
func (ah *AssertionHelper) AssertBasicInfoResponse(resp *authapi.GetBasicInfoResponse) bool {
	if !ah.AssertNotNil(resp, "基本情報レスポンス") {
		return false
	}

	if !ah.AssertStatusCode(200, resp.StatusCode(), "基本情報取得") {
		return false
	}

	if !ah.AssertNotNil(resp.JSON200, "基本情報データ") {
		return false
	}

	basicInfo := resp.JSON200
	success := true

	success = ah.AssertNotEmpty(basicInfo.DomainName, "ドメイン名") && success
	success = ah.AssertNotEmpty(basicInfo.FromEmailAddress, "送信元メールアドレス") && success
	success = ah.AssertNotNil(basicInfo.CertificateDnsRecord, "証明書DNSレコード") && success
	success = ah.AssertNotNil(basicInfo.CloudFrontDnsRecord, "CloudFront DNSレコード") && success

	// DnsRecordは構造体なので、フィールドが空でないかをチェック
	if basicInfo.CertificateDnsRecord.Name != "" || basicInfo.CertificateDnsRecord.Value != "" {
		success = ah.AssertNotEmpty(basicInfo.CertificateDnsRecord.Name, "証明書DNSレコード名") && success
		success = ah.AssertNotEmpty(basicInfo.CertificateDnsRecord.Value, "証明書DNSレコード値") && success
	}

	return success
}

// AssertAuthInfoResponse は認証情報レスポンスをアサートします
func (ah *AssertionHelper) AssertAuthInfoResponse(resp *authapi.GetAuthInfoResponse) bool {
	if !ah.AssertNotNil(resp, "認証情報レスポンス") {
		return false
	}

	if !ah.AssertStatusCode(200, resp.StatusCode(), "認証情報取得") {
		return false
	}

	if !ah.AssertNotNil(resp.JSON200, "認証情報データ") {
		return false
	}

	authInfo := resp.JSON200
	return ah.AssertNotEmpty(authInfo.CallbackUrl, "コールバックURL")
}

// AssertSaasUsersResponse はSaaSユーザー一覧レスポンスをアサートします
func (ah *AssertionHelper) AssertSaasUsersResponse(resp *authapi.GetSaasUsersResponse) bool {
	if !ah.AssertNotNil(resp, "SaaSユーザー一覧レスポンス") {
		return false
	}

	if !ah.AssertStatusCode(200, resp.StatusCode(), "SaaSユーザー一覧取得") {
		return false
	}

	if !ah.AssertNotNil(resp.JSON200, "SaaSユーザー一覧データ") {
		return false
	}

	users := resp.JSON200.Users
	success := true

	for i, user := range users {
		success = ah.AssertNotEmpty(user.Id, fmt.Sprintf("ユーザー[%d].ID", i)) && success
		success = ah.AssertNotEmpty(user.Email, fmt.Sprintf("ユーザー[%d].Email", i)) && success
	}

	if ah.verbose {
		ah.t.Logf("✓ SaaSユーザー数: %d", len(users))
	}

	return success
}

// AssertRolesResponse はロール一覧レスポンスをアサートします
func (ah *AssertionHelper) AssertRolesResponse(resp *authapi.GetRolesResponse) bool {
	if !ah.AssertNotNil(resp, "ロール一覧レスポンス") {
		return false
	}

	if !ah.AssertStatusCode(200, resp.StatusCode(), "ロール一覧取得") {
		return false
	}

	if !ah.AssertNotNil(resp.JSON200, "ロール一覧データ") {
		return false
	}

	roles := resp.JSON200.Roles
	success := true

	for i, role := range roles {
		success = ah.AssertNotEmpty(role.RoleName, fmt.Sprintf("ロール[%d].RoleName", i)) && success
		success = ah.AssertNotEmpty(role.DisplayName, fmt.Sprintf("ロール[%d].DisplayName", i)) && success
	}

	if ah.verbose {
		ah.t.Logf("✓ ロール数: %d", len(roles))
	}

	return success
}

// AssertTenantsResponse はテナント一覧レスポンスをアサートします
func (ah *AssertionHelper) AssertTenantsResponse(resp *authapi.GetTenantsResponse) bool {
	if !ah.AssertNotNil(resp, "テナント一覧レスポンス") {
		return false
	}

	if !ah.AssertStatusCode(200, resp.StatusCode(), "テナント一覧取得") {
		return false
	}

	if !ah.AssertNotNil(resp.JSON200, "テナント一覧データ") {
		return false
	}

	tenants := resp.JSON200.Tenants
	success := true

	for i, tenant := range tenants {
		success = ah.AssertNotEmpty(tenant.Id, fmt.Sprintf("テナント[%d].ID", i)) && success
		success = ah.AssertNotEmpty(tenant.Name, fmt.Sprintf("テナント[%d].Name", i)) && success
	}

	if ah.verbose {
		ah.t.Logf("✓ テナント数: %d", len(tenants))
	}

	return success
}

// AssertEnvsResponse は環境一覧レスポンスをアサートします
func (ah *AssertionHelper) AssertEnvsResponse(resp *authapi.GetEnvsResponse) bool {
	if !ah.AssertNotNil(resp, "環境一覧レスポンス") {
		return false
	}

	if !ah.AssertStatusCode(200, resp.StatusCode(), "環境一覧取得") {
		return false
	}

	if !ah.AssertNotNil(resp.JSON200, "環境一覧データ") {
		return false
	}

	envs := resp.JSON200.Envs
	success := true

	for i, env := range envs {
		success = ah.AssertNotEmpty(fmt.Sprintf("%d", env.Id), fmt.Sprintf("環境[%d].ID", i)) && success
		success = ah.AssertNotEmpty(env.Name, fmt.Sprintf("環境[%d].Name", i)) && success
	}

	if ah.verbose {
		ah.t.Logf("✓ 環境数: %d", len(envs))
	}

	return success
}

// AssertErrorResponse はエラーレスポンスをアサートします
func (ah *AssertionHelper) AssertErrorResponse(statusCode int, errorType, errorMessage string) bool {
	success := true

	if statusCode < 400 {
		ah.t.Errorf("エラーレスポンスのステータスコードが期待値未満: %d", statusCode)
		success = false
	}

	if errorType != "" {
		success = ah.AssertNotEmpty(errorType, "エラータイプ") && success
	}

	if errorMessage != "" {
		success = ah.AssertNotEmpty(errorMessage, "エラーメッセージ") && success
	}

	if ah.verbose && success {
		ah.t.Logf("✓ エラーレスポンス確認: %d %s %s", statusCode, errorType, errorMessage)
	}

	return success
}

// AssertResponseTime はレスポンス時間をアサートします
func (ah *AssertionHelper) AssertResponseTime(duration time.Duration, maxDuration time.Duration, operation string) bool {
	if duration > maxDuration {
		ah.t.Errorf("%s のレスポンス時間が上限を超過: %v > %v", operation, duration, maxDuration)
		return false
	}

	if ah.verbose {
		ah.t.Logf("✓ %s レスポンス時間: %v", operation, duration)
	}

	return true
}

// AssertValidEmail はメールアドレスの形式をアサートします
func (ah *AssertionHelper) AssertValidEmail(email, fieldName string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		ah.t.Errorf("%s の形式が不正: %s", fieldName, email)
		return false
	}

	if ah.verbose {
		ah.t.Logf("✓ %s の形式確認: %s", fieldName, email)
	}

	return true
}

// AssertValidURL はURLの形式をアサートします
func (ah *AssertionHelper) AssertValidURL(url, fieldName string) bool {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		ah.t.Errorf("%s の形式が不正: %s", fieldName, url)
		return false
	}

	if ah.verbose {
		ah.t.Logf("✓ %s の形式確認: %s", fieldName, url)
	}

	return true
}

// AssertArrayLength は配列の長さをアサートします
func (ah *AssertionHelper) AssertArrayLength(array interface{}, expectedLength int, fieldName string) bool {
	v := reflect.ValueOf(array)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		ah.t.Errorf("%s は配列またはスライスではありません", fieldName)
		return false
	}

	actualLength := v.Len()
	if actualLength != expectedLength {
		ah.t.Errorf("%s の長さ: 期待値 %d, 実際 %d", fieldName, expectedLength, actualLength)
		return false
	}

	if ah.verbose {
		ah.t.Logf("✓ %s の長さ確認: %d", fieldName, actualLength)
	}

	return true
}

// AssertArrayNotEmpty は配列が空でないことをアサートします
func (ah *AssertionHelper) AssertArrayNotEmpty(array interface{}, fieldName string) bool {
	v := reflect.ValueOf(array)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		ah.t.Errorf("%s は配列またはスライスではありません", fieldName)
		return false
	}

	if v.Len() == 0 {
		ah.t.Errorf("%s が空です", fieldName)
		return false
	}

	if ah.verbose {
		ah.t.Logf("✓ %s は非空: 長さ %d", fieldName, v.Len())
	}

	return true
}

// AssertBooleanValue はブール値をアサートします
func (ah *AssertionHelper) AssertBooleanValue(expected, actual bool, fieldName string) bool {
	if expected != actual {
		ah.t.Errorf("%s: 期待値 %t, 実際 %t", fieldName, expected, actual)
		return false
	}

	if ah.verbose {
		ah.t.Logf("✓ %s 確認: %t", fieldName, actual)
	}

	return true
}

// AssertResourceCreated はリソースが正常に作成されたことをアサートします
func (ah *AssertionHelper) AssertResourceCreated(resourceID, resourceType string) bool {
	success := true
	success = ah.AssertNotEmpty(resourceID, fmt.Sprintf("%s ID", resourceType)) && success

	if ah.verbose && success {
		ah.t.Logf("✓ %s が作成されました: %s", resourceType, resourceID)
	}

	return success
}

// AssertResourceDeleted はリソースが正常に削除されたことをアサートします
func (ah *AssertionHelper) AssertResourceDeleted(statusCode int, resourceType string) bool {
	// 200 (削除成功) または 404 (既に存在しない) を許可
	if statusCode != 200 && statusCode != 404 {
		ah.t.Errorf("%s の削除に失敗: ステータスコード %d", resourceType, statusCode)
		return false
	}

	if ah.verbose {
		ah.t.Logf("✓ %s が削除されました: ステータスコード %d", resourceType, statusCode)
	}

	return true
}

// AssertTestResult はテスト結果全体をアサートします
func (ah *AssertionHelper) AssertTestResult(result *TestResult, expectedSuccess bool) bool {
	if !ah.AssertNotNil(result, "テスト結果") {
		return false
	}

	success := true
	success = ah.AssertEquals(expectedSuccess, result.Success, "テスト成功フラグ") && success
	success = ah.AssertNotEmpty(result.TestName, "テスト名") && success
	success = ah.AssertNotEmpty(result.Story, "ストーリー名") && success

	if result.Duration > 0 {
		success = ah.AssertResponseTime(result.Duration, 30*time.Second, result.TestName) && success
	}

	return success
}