package common

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TestDataLoader はテストデータの読み込みを管理します
type TestDataLoader struct {
	basePath string
}

// NewTestDataLoader は新しいテストデータローダーを作成します
func NewTestDataLoader() *TestDataLoader {
	// 実行時の作業ディレクトリからの相対パスで testdata ディレクトリを探す
	basePath := "testdata"
	
	// stories ディレクトリから実行された場合、一つ上のディレクトリの testdata を参照
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		basePath = filepath.Join("..", "testdata")
	}
	
	return &TestDataLoader{
		basePath: basePath,
	}
}

// LoadBasicSetupData は基本セットアップストーリーのテストデータを読み込みます
func (loader *TestDataLoader) LoadBasicSetupData() (*BasicSetupTestData, error) {
	var data BasicSetupTestData
	
	// 基本セットアップ用のYAMLファイルを読み込み
	filePath := filepath.Join(loader.basePath, "basic_setup", "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("基本セットアップデータの読み込みに失敗: %w", err)
	}

	return &data, nil
}

// LoadRoleManagementData はロール管理ストーリーのテストデータを読み込みます
func (loader *TestDataLoader) LoadRoleManagementData() (*RoleManagementTestData, error) {
	var data RoleManagementTestData
	
	// ロール管理用のYAMLファイルを読み込み
	filePath := filepath.Join(loader.basePath, "role_management", "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("ロール管理データの読み込みに失敗: %w", err)
	}

	return &data, nil
}

// LoadUserManagementData はユーザー管理ストーリーのテストデータを読み込みます
func (loader *TestDataLoader) LoadUserManagementData() (*UserManagementTestData, error) {
	var data UserManagementTestData
	
	// ユーザー管理用のYAMLファイルを読み込み
	filePath := filepath.Join(loader.basePath, "user_management", "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("ユーザー管理データの読み込みに失敗: %w", err)
	}

	return &data, nil
}

// LoadTenantManagementData はテナント管理ストーリーのテストデータを読み込みます
func (loader *TestDataLoader) LoadTenantManagementData() (*TenantManagementTestData, error) {
	var data TenantManagementTestData
	
	// テナント管理用のYAMLファイルを読み込み
	filePath := filepath.Join(loader.basePath, "tenant_management", "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("テナント管理データの読み込みに失敗: %w", err)
	}

	return &data, nil
}

// LoadTenantUserManagementData はテナントユーザー管理ストーリーのテストデータを読み込みます
func (loader *TestDataLoader) LoadTenantUserManagementData() (*TenantUserManagementTestData, error) {
	var data TenantUserManagementTestData
	
	// テナントユーザー管理用のYAMLファイルを読み込み
	filePath := filepath.Join(loader.basePath, "tenant_user_management", "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("テナントユーザー管理データの読み込みに失敗: %w", err)
	}

	return &data, nil
}

// LoadInvitationManagementData は招待管理ストーリーのテストデータを読み込みます
func (loader *TestDataLoader) LoadInvitationManagementData() (*InvitationManagementTestData, error) {
	var data InvitationManagementTestData
	
	// 招待管理用のYAMLファイルを読み込み
	filePath := filepath.Join(loader.basePath, "invitation_management", "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("招待管理データの読み込みに失敗: %w", err)
	}

	return &data, nil
}

// LoadAuthenticationFlowData は認証フローストーリーのテストデータを読み込みます
func (loader *TestDataLoader) LoadAuthenticationFlowData() (*AuthenticationFlowTestData, error) {
	var data AuthenticationFlowTestData
	
	// 認証フロー用のYAMLファイルを読み込み
	filePath := filepath.Join(loader.basePath, "authentication_flow", "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("認証フローデータの読み込みに失敗: %w", err)
	}

	return &data, nil
}

// LoadSingleTenantManagementData はシングルテナント管理ストーリーのテストデータを読み込みます
func (loader *TestDataLoader) LoadSingleTenantManagementData() (*SingleTenantManagementTestData, error) {
	var data SingleTenantManagementTestData
	
	// シングルテナント管理用のYAMLファイルを読み込み
	filePath := filepath.Join(loader.basePath, "single_tenant_management", "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("シングルテナント管理データの読み込みに失敗: %w", err)
	}

	return &data, nil
}

// LoadErrorHandlingData はエラーハンドリングストーリーのテストデータを読み込みます
func (loader *TestDataLoader) LoadErrorHandlingData() (*ErrorHandlingTestData, error) {
	var data ErrorHandlingTestData
	
	// エラーハンドリング用のYAMLファイルを読み込み
	filePath := filepath.Join(loader.basePath, "error_handling", "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("エラーハンドリングデータの読み込みに失敗: %w", err)
	}

	return &data, nil
}

// LoadIntegrationTestData は統合テストストーリーのテストデータを読み込みます
func (loader *TestDataLoader) LoadIntegrationTestData() (*IntegrationTestData, error) {
	var data IntegrationTestData
	
	// 統合テスト用のYAMLファイルを読み込み
	filePath := filepath.Join(loader.basePath, "integration_test", "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("統合テストデータの読み込みに失敗: %w", err)
	}

	return &data, nil
}

// LoadTestData は汎用的なテストデータを読み込みます
func (loader *TestDataLoader) LoadTestData(storyName string) (*TestData, error) {
	var data TestData
	
	filePath := filepath.Join(loader.basePath, storyName, "test_data.yml")
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		return nil, fmt.Errorf("テストデータの読み込みに失敗 (%s): %w", storyName, err)
	}

	return &data, nil
}

// LoadTestConfig はテスト設定を読み込みます
func (loader *TestDataLoader) LoadTestConfig() (*TestConfig, error) {
	var config TestConfig
	
	filePath := filepath.Join(loader.basePath, "config.yml")
	if err := loader.loadYAMLFile(filePath, &config); err != nil {
		// 設定ファイルが存在しない場合はデフォルト設定を返す
		if os.IsNotExist(err) {
			return loader.getDefaultConfig(), nil
		}
		return nil, fmt.Errorf("テスト設定の読み込みに失敗: %w", err)
	}

	return &config, nil
}

// LoadTestCases は指定されたストーリーのテストケースを読み込みます
func (loader *TestDataLoader) LoadTestCases(storyName, fileName string) ([]TestCase, error) {
	var testCases []TestCase
	
	filePath := filepath.Join(loader.basePath, storyName, fileName)
	if err := loader.loadYAMLFile(filePath, &testCases); err != nil {
		return nil, fmt.Errorf("テストケースの読み込みに失敗 (%s/%s): %w", storyName, fileName, err)
	}

	return testCases, nil
}

// LoadStoryTestData は指定されたストーリーの全テストデータを読み込みます
func (loader *TestDataLoader) LoadStoryTestData(storyName string) (map[string]interface{}, error) {
	storyDir := filepath.Join(loader.basePath, storyName)
	
	// ディレクトリが存在するかチェック
	if _, err := os.Stat(storyDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("ストーリーディレクトリが存在しません: %s", storyDir)
	}

	// ディレクトリ内のYAMLファイルを読み込み
	files, err := filepath.Glob(filepath.Join(storyDir, "*.yml"))
	if err != nil {
		return nil, fmt.Errorf("YAMLファイルの検索に失敗: %w", err)
	}

	data := make(map[string]interface{})
	for _, file := range files {
		fileName := filepath.Base(file)
		fileKey := fileName[:len(fileName)-4] // .yml拡張子を除去

		var fileData interface{}
		if err := loader.loadYAMLFile(file, &fileData); err != nil {
			return nil, fmt.Errorf("ファイル %s の読み込みに失敗: %w", file, err)
		}

		data[fileKey] = fileData
	}

	return data, nil
}

// loadYAMLFile はYAMLファイルを読み込んで指定された構造体にデシリアライズします
func (loader *TestDataLoader) loadYAMLFile(filePath string, target interface{}) error {
	// ファイルの存在確認
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("ファイルが存在しません: %s", filePath)
	}

	// ファイル読み込み
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ファイルの読み込みに失敗: %w", err)
	}

	// YAML解析
	if err := yaml.Unmarshal(data, target); err != nil {
		return fmt.Errorf("YAML解析に失敗: %w", err)
	}

	return nil
}

// getDefaultConfig はデフォルトのテスト設定を返します
func (loader *TestDataLoader) getDefaultConfig() *TestConfig {
	return &TestConfig{
		APIEndpoint: SaaSusAPIEndpoint,
		Timeout:     DefaultTimeout,
		RetryCount:  DefaultRetryCount,
		Verbose:     false,
	}
}

// ValidateTestData はテストデータの妥当性を検証します
func (loader *TestDataLoader) ValidateTestData(data *TestData) error {
	if data == nil {
		return fmt.Errorf("テストデータがnilです")
	}

	if data.Story == "" {
		return fmt.Errorf("ストーリー名が設定されていません")
	}

	if len(data.TestCases) == 0 {
		return fmt.Errorf("テストケースが設定されていません")
	}

	// 各テストケースの妥当性をチェック
	for i, testCase := range data.TestCases {
		if err := loader.validateTestCase(&testCase, i); err != nil {
			return fmt.Errorf("テストケース %d の検証に失敗: %w", i, err)
		}
	}

	return nil
}

// validateTestCase は個別のテストケースの妥当性を検証します
func (loader *TestDataLoader) validateTestCase(testCase *TestCase, index int) error {
	if testCase.Name == "" {
		return fmt.Errorf("テストケース名が設定されていません")
	}

	if testCase.Method == "" {
		return fmt.Errorf("HTTPメソッドが設定されていません")
	}

	if testCase.Endpoint == "" {
		return fmt.Errorf("エンドポイントが設定されていません")
	}

	if testCase.Expected.StatusCode == 0 {
		return fmt.Errorf("期待されるステータスコードが設定されていません")
	}

	return nil
}

// GetTestDataPath はテストデータファイルのパスを取得します
func (loader *TestDataLoader) GetTestDataPath(storyName, fileName string) string {
	return filepath.Join(loader.basePath, storyName, fileName)
}

// FileExists はファイルが存在するかチェックします
func (loader *TestDataLoader) FileExists(storyName, fileName string) bool {
	filePath := loader.GetTestDataPath(storyName, fileName)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// ListStoryDirectories は利用可能なストーリーディレクトリ一覧を取得します
func (loader *TestDataLoader) ListStoryDirectories() ([]string, error) {
	entries, err := os.ReadDir(loader.basePath)
	if err != nil {
		return nil, fmt.Errorf("testdataディレクトリの読み込みに失敗: %w", err)
	}

	var stories []string
	for _, entry := range entries {
		if entry.IsDir() {
			stories = append(stories, entry.Name())
		}
	}

	return stories, nil
}

// LoadEnvironmentSpecificData は環境固有のテストデータを読み込みます
func (loader *TestDataLoader) LoadEnvironmentSpecificData(storyName, environment string) (map[string]interface{}, error) {
	fileName := fmt.Sprintf("test_data_%s.yml", environment)
	filePath := filepath.Join(loader.basePath, storyName, fileName)

	var data map[string]interface{}
	if err := loader.loadYAMLFile(filePath, &data); err != nil {
		// 環境固有ファイルが存在しない場合はデフォルトを読み込み
		if os.IsNotExist(err) {
			return loader.LoadStoryTestData(storyName)
		}
		return nil, fmt.Errorf("環境固有データの読み込みに失敗: %w", err)
	}

	return data, nil
}

// MergeTestData は複数のテストデータをマージします
func (loader *TestDataLoader) MergeTestData(base, override map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// ベースデータをコピー
	for key, value := range base {
		result[key] = value
	}

	// オーバーライドデータで上書き
	for key, value := range override {
		result[key] = value
	}

	return result
}

// SaveTestResult はテスト結果をファイルに保存します
func (loader *TestDataLoader) SaveTestResult(result *TestResult, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("出力ディレクトリの作成に失敗: %w", err)
	}

	fileName := fmt.Sprintf("result_%s_%s.yml", result.Story, result.TestName)
	filePath := filepath.Join(outputDir, fileName)

	data, err := yaml.Marshal(result)
	if err != nil {
		return fmt.Errorf("テスト結果のシリアライズに失敗: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("テスト結果の保存に失敗: %w", err)
	}

	return nil
}