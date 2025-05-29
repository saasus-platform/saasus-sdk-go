package common

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/client"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
)

const (
	// SaaSusAPIEndpoint はSaaSus Auth APIのエンドポイント
	SaaSusAPIEndpoint = "https://api.dev.saasus.io/v1/auth"
	
	// DefaultTimeout はデフォルトのタイムアウト時間
	DefaultTimeout = 30 * time.Second
	
	// DefaultRetryCount はデフォルトのリトライ回数
	DefaultRetryCount = 3
)

var (
	// globalClient はグローバルクライアントインスタンス
	globalClient *ClientWrapper
	clientMutex  sync.RWMutex
)

// NewAuthenticatedClient は認証付きのクライアントを作成します
func NewAuthenticatedClient() (*ClientWrapper, error) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	// 既存のクライアントがある場合は再利用
	if globalClient != nil {
		return globalClient, nil
	}

	// 環境変数の確認
	saasID := os.Getenv("SAASUS_SAAS_ID")
	apiKey := os.Getenv("SAASUS_API_KEY")
	secretKey := os.Getenv("SAASUS_SECRET_KEY")

	if saasID == "" || apiKey == "" || secretKey == "" {
		return nil, fmt.Errorf("必要な環境変数が設定されていません: SAASUS_SAAS_ID, SAASUS_API_KEY, SAASUS_SECRET_KEY")
	}

	// SaaSus認証を設定するリクエストエディター
	authEditor := func(ctx context.Context, req *http.Request) error {
		if err := client.SetSigV1(req); err != nil {
			return fmt.Errorf("SigV1認証の設定に失敗: %w", err)
		}
		client.SetReferer(ctx, req)
		return nil
	}

	// 認証付きクライアントを作成
	clientWithAuth, err := authapi.NewClientWithResponses(
		SaaSusAPIEndpoint,
		authapi.WithRequestEditorFn(authEditor),
	)
	if err != nil {
		return nil, fmt.Errorf("認証付きクライアントの作成に失敗: %w", err)
	}

	// 設定を作成
	config := &TestConfig{
		APIEndpoint: SaaSusAPIEndpoint,
		Timeout:     DefaultTimeout,
		RetryCount:  DefaultRetryCount,
		Verbose:     false,
	}

	// クライアントラッパーを作成
	globalClient = &ClientWrapper{
		Client:    clientWithAuth,
		Config:    config,
		Resources: make([]*TestResource, 0),
	}

	return globalClient, nil
}

// GetClient は既存のクライアントを取得します
func GetClient() (*ClientWrapper, error) {
	clientMutex.RLock()
	defer clientMutex.RUnlock()

	if globalClient == nil {
		return NewAuthenticatedClient()
	}

	return globalClient, nil
}

// TestConnection はAPI接続をテストします
func TestConnection(ctx context.Context, wrapper *ClientWrapper) error {
	if wrapper == nil || wrapper.Client == nil {
		return fmt.Errorf("クライアントが初期化されていません")
	}

	// タイムアウト付きコンテキストを作成
	timeoutCtx, cancel := context.WithTimeout(ctx, wrapper.Config.Timeout)
	defer cancel()

	// 基本情報取得でAPI接続をテスト
	resp, err := wrapper.Client.GetBasicInfoWithResponse(timeoutCtx)
	if err != nil {
		return fmt.Errorf("API接続テストに失敗: %w", err)
	}

	if resp == nil {
		return fmt.Errorf("レスポンスがnilです")
	}

	// ステータスコードをチェック
	switch resp.StatusCode() {
	case 200:
		return nil
	case 401:
		return fmt.Errorf("認証に失敗しました。APIキーを確認してください")
	case 403:
		return fmt.Errorf("アクセスが拒否されました。権限を確認してください")
	case 500:
		return fmt.Errorf("サーバーエラーが発生しました")
	default:
		return fmt.Errorf("予期しないステータスコード: %d", resp.StatusCode())
	}
}

// AddResource はテスト中に作成されたリソースを追跡に追加します
func (w *ClientWrapper) AddResource(resource *TestResource) {
	if w.Resources == nil {
		w.Resources = make([]*TestResource, 0)
	}
	w.Resources = append(w.Resources, resource)
}

// GetResources は追跡中のリソース一覧を取得します
func (w *ClientWrapper) GetResources() []*TestResource {
	return w.Resources
}

// GetResourcesByType は指定されたタイプのリソース一覧を取得します
func (w *ClientWrapper) GetResourcesByType(resourceType string) []*TestResource {
	var resources []*TestResource
	for _, resource := range w.Resources {
		if resource.Type == resourceType {
			resources = append(resources, resource)
		}
	}
	return resources
}

// GetResourcesByStory は指定されたストーリーのリソース一覧を取得します
func (w *ClientWrapper) GetResourcesByStory(story string) []*TestResource {
	var resources []*TestResource
	for _, resource := range w.Resources {
		if resource.Story == story {
			resources = append(resources, resource)
		}
	}
	return resources
}

// MarkResourceCleaned はリソースをクリーンアップ済みとしてマークします
func (w *ClientWrapper) MarkResourceCleaned(resourceID string, err error) {
	for _, resource := range w.Resources {
		if resource.ID == resourceID {
			resource.CleanedUp = true
			if err != nil {
				resource.CleanupErr = err.Error()
			}
			break
		}
	}
}

// GetUncleanedResources はクリーンアップされていないリソース一覧を取得します
func (w *ClientWrapper) GetUncleanedResources() []*TestResource {
	var resources []*TestResource
	for _, resource := range w.Resources {
		if !resource.CleanedUp {
			resources = append(resources, resource)
		}
	}
	return resources
}

// ExecuteWithRetry はリトライ機能付きでAPIリクエストを実行します
func (w *ClientWrapper) ExecuteWithRetry(ctx context.Context, operation func() error) error {
	var lastErr error
	
	for i := 0; i <= w.Config.RetryCount; i++ {
		if i > 0 {
			// リトライ前に少し待機
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(i) * time.Second):
			}
		}

		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		// リトライ可能なエラーかチェック
		if !isRetryableError(lastErr) {
			return lastErr
		}
	}

	return fmt.Errorf("リトライ回数を超過しました。最後のエラー: %w", lastErr)
}

// isRetryableError はエラーがリトライ可能かどうかを判定します
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// ネットワークエラーやタイムアウトはリトライ可能
	errStr := err.Error()
	retryableErrors := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"temporary failure",
		"service unavailable",
	}

	for _, retryableErr := range retryableErrors {
		if contains(errStr, retryableErr) {
			return true
		}
	}

	return false
}

// contains は文字列に部分文字列が含まれているかチェックします
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
			(s[:len(substr)] == substr || 
			 s[len(s)-len(substr):] == substr ||
			 containsSubstring(s, substr))))
}

// containsSubstring は文字列内の部分文字列を検索します
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// CreateTestResource はテストリソースを作成し、追跡に追加します
func (w *ClientWrapper) CreateTestResource(resourceType, resourceID, resourceName, story string, metadata map[string]interface{}) *TestResource {
	resource := &TestResource{
		Type:      resourceType,
		ID:        resourceID,
		Name:      resourceName,
		CreatedAt: time.Now(),
		Story:     story,
		Metadata:  metadata,
		CleanedUp: false,
	}

	w.AddResource(resource)
	return resource
}

// ValidateEnvironment は実行環境が適切に設定されているかを検証します
func ValidateEnvironment() error {
	required := map[string]string{
		"SAASUS_SAAS_ID":    os.Getenv("SAASUS_SAAS_ID"),
		"SAASUS_API_KEY":    os.Getenv("SAASUS_API_KEY"),
		"SAASUS_SECRET_KEY": os.Getenv("SAASUS_SECRET_KEY"),
	}

	var missing []string
	for key, value := range required {
		if value == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("以下の環境変数が設定されていません: %v", missing)
	}

	return nil
}

// ResetClient はグローバルクライアントをリセットします（テスト用）
func ResetClient() {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	globalClient = nil
}

// NewClientWithCustomAuth はカスタム認証トークンでクライアントを作成します（テスト用）
func NewClientWithCustomAuth(customAuth string) (*ClientWrapper, error) {
	// カスタム認証を設定するリクエストエディター
	authEditor := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+customAuth)
		client.SetReferer(ctx, req)
		return nil
	}

	// カスタム認証付きクライアントを作成
	clientWithAuth, err := authapi.NewClientWithResponses(
		SaaSusAPIEndpoint,
		authapi.WithRequestEditorFn(authEditor),
	)
	if err != nil {
		return nil, fmt.Errorf("カスタム認証付きクライアントの作成に失敗: %w", err)
	}

	// 設定を作成
	config := &TestConfig{
		APIEndpoint: SaaSusAPIEndpoint,
		Timeout:     DefaultTimeout,
		RetryCount:  DefaultRetryCount,
		Verbose:     false,
	}

	// クライアントラッパーを作成
	wrapper := &ClientWrapper{
		Client:    clientWithAuth,
		Config:    config,
		Resources: make([]*TestResource, 0),
	}

	return wrapper, nil
}