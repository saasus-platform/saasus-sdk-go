package common

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
)

// CleanupManager はテストリソースのクリーンアップを管理します
type CleanupManager struct {
	client  *ClientWrapper
	verbose bool
}

// NewCleanupManager は新しいクリーンアップマネージャーを作成します
func NewCleanupManager(client *ClientWrapper, verbose bool) *CleanupManager {
	return &CleanupManager{
		client:  client,
		verbose: verbose,
	}
}

// CleanupTestResources はテスト中に作成されたリソースをクリーンアップします
func CleanupTestResources(ctx context.Context, client *ClientWrapper) error {
	manager := NewCleanupManager(client, true)
	return manager.CleanupAll(ctx)
}

// CleanupAll は全てのテストリソースをクリーンアップします
func (cm *CleanupManager) CleanupAll(ctx context.Context) error {
	if cm.client == nil {
		return fmt.Errorf("クライアントが初期化されていません")
	}

	resources := cm.client.GetUncleanedResources()
	if len(resources) == 0 {
		if cm.verbose {
			log.Println("クリーンアップ対象のリソースはありません")
		}
		return nil
	}

	if cm.verbose {
		log.Printf("クリーンアップ対象リソース数: %d", len(resources))
	}

	var errors []error

	// リソースタイプ別にクリーンアップ順序を制御
	cleanupOrder := []string{
		ResourceTypeInvitation,     // 招待
		ResourceTypeTenantUser,     // テナントユーザー
		ResourceTypeTenant,         // テナント
		ResourceTypeSaasUser,       // SaaSユーザー
		ResourceTypeRole,           // ロール
		ResourceTypeUserAttribute,  // ユーザー属性
		ResourceTypeTenantAttribute, // テナント属性
		ResourceTypeEnv,            // 環境
	}

	for _, resourceType := range cleanupOrder {
		typeResources := cm.client.GetResourcesByType(resourceType)
		for _, resource := range typeResources {
			if resource.CleanedUp {
				continue
			}

			if err := cm.cleanupResource(ctx, resource); err != nil {
				errors = append(errors, fmt.Errorf("リソース %s (%s) のクリーンアップに失敗: %w", 
					resource.ID, resource.Type, err))
				cm.client.MarkResourceCleaned(resource.ID, err)
			} else {
				cm.client.MarkResourceCleaned(resource.ID, nil)
				if cm.verbose {
					log.Printf("リソース %s (%s) をクリーンアップしました", resource.ID, resource.Type)
				}
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("一部のリソースのクリーンアップに失敗しました: %v", errors)
	}

	return nil
}

// CleanupByStory は指定されたストーリーのリソースをクリーンアップします
func (cm *CleanupManager) CleanupByStory(ctx context.Context, story string) error {
	resources := cm.client.GetResourcesByStory(story)
	if len(resources) == 0 {
		if cm.verbose {
			log.Printf("ストーリー %s のクリーンアップ対象リソースはありません", story)
		}
		return nil
	}

	if cm.verbose {
		log.Printf("ストーリー %s のクリーンアップ対象リソース数: %d", story, len(resources))
	}

	var errors []error
	for _, resource := range resources {
		if resource.CleanedUp {
			continue
		}

		if err := cm.cleanupResource(ctx, resource); err != nil {
			errors = append(errors, fmt.Errorf("リソース %s のクリーンアップに失敗: %w", resource.ID, err))
			cm.client.MarkResourceCleaned(resource.ID, err)
		} else {
			cm.client.MarkResourceCleaned(resource.ID, nil)
			if cm.verbose {
				log.Printf("リソース %s をクリーンアップしました", resource.ID)
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("一部のリソースのクリーンアップに失敗しました: %v", errors)
	}

	return nil
}

// cleanupResource は個別のリソースをクリーンアップします
func (cm *CleanupManager) cleanupResource(ctx context.Context, resource *TestResource) error {
	switch resource.Type {
	case ResourceTypeEnv:
		return cm.cleanupEnv(ctx, resource.ID)
	case ResourceTypeRole:
		return cm.cleanupRole(ctx, resource.ID)
	case ResourceTypeSaasUser:
		return cm.cleanupSaasUser(ctx, resource.ID)
	case ResourceTypeTenant:
		return cm.cleanupTenant(ctx, resource.ID)
	case ResourceTypeTenantUser:
		return cm.cleanupTenantUser(ctx, resource)
	case ResourceTypeUserAttribute:
		return cm.cleanupUserAttribute(ctx, resource.ID)
	case ResourceTypeTenantAttribute:
		return cm.cleanupTenantAttribute(ctx, resource.ID)
	case ResourceTypeInvitation:
		return cm.cleanupInvitation(ctx, resource)
	default:
		return fmt.Errorf("未対応のリソースタイプ: %s", resource.Type)
	}
}

// cleanupEnv は環境をクリーンアップします
func (cm *CleanupManager) cleanupEnv(ctx context.Context, envID string) error {
	// stringからuint64への変換
	envIDUint, err := strconv.ParseUint(envID, 10, 64)
	if err != nil {
		return fmt.Errorf("環境ID変換エラー: %w", err)
	}
	
	resp, err := cm.client.Client.DeleteEnvWithResponse(ctx, authapi.Id(envIDUint))
	if err != nil {
		return fmt.Errorf("環境削除APIの呼び出しに失敗: %w", err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 404 {
		return fmt.Errorf("環境削除に失敗: ステータスコード %d", resp.StatusCode())
	}

	return nil
}

// cleanupRole はロールをクリーンアップします
func (cm *CleanupManager) cleanupRole(ctx context.Context, roleName string) error {
	resp, err := cm.client.Client.DeleteRoleWithResponse(ctx, roleName)
	if err != nil {
		return fmt.Errorf("ロール削除APIの呼び出しに失敗: %w", err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 404 {
		return fmt.Errorf("ロール削除に失敗: ステータスコード %d", resp.StatusCode())
	}

	return nil
}

// cleanupSaasUser はSaaSユーザーをクリーンアップします
func (cm *CleanupManager) cleanupSaasUser(ctx context.Context, userID string) error {
	resp, err := cm.client.Client.DeleteSaasUserWithResponse(ctx, userID)
	if err != nil {
		return fmt.Errorf("SaaSユーザー削除APIの呼び出しに失敗: %w", err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 404 {
		return fmt.Errorf("SaaSユーザー削除に失敗: ステータスコード %d", resp.StatusCode())
	}

	return nil
}

// cleanupTenant はテナントをクリーンアップします
func (cm *CleanupManager) cleanupTenant(ctx context.Context, tenantID string) error {
	resp, err := cm.client.Client.DeleteTenantWithResponse(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("テナント削除APIの呼び出しに失敗: %w", err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 404 {
		return fmt.Errorf("テナント削除に失敗: ステータスコード %d", resp.StatusCode())
	}

	return nil
}

// cleanupTenantUser はテナントユーザーをクリーンアップします
func (cm *CleanupManager) cleanupTenantUser(ctx context.Context, resource *TestResource) error {
	// メタデータからテナントIDを取得
	tenantID, ok := resource.Metadata["tenant_id"].(string)
	if !ok {
		return fmt.Errorf("テナントIDがメタデータに見つかりません")
	}

	resp, err := cm.client.Client.DeleteTenantUserWithResponse(ctx, tenantID, resource.ID)
	if err != nil {
		return fmt.Errorf("テナントユーザー削除APIの呼び出しに失敗: %w", err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 404 {
		return fmt.Errorf("テナントユーザー削除に失敗: ステータスコード %d", resp.StatusCode())
	}

	return nil
}

// cleanupUserAttribute はユーザー属性をクリーンアップします
func (cm *CleanupManager) cleanupUserAttribute(ctx context.Context, attributeName string) error {
	resp, err := cm.client.Client.DeleteUserAttributeWithResponse(ctx, attributeName)
	if err != nil {
		return fmt.Errorf("ユーザー属性削除APIの呼び出しに失敗: %w", err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 404 {
		return fmt.Errorf("ユーザー属性削除に失敗: ステータスコード %d", resp.StatusCode())
	}

	return nil
}

// cleanupTenantAttribute はテナント属性をクリーンアップします
func (cm *CleanupManager) cleanupTenantAttribute(ctx context.Context, attributeName string) error {
	resp, err := cm.client.Client.DeleteTenantAttributeWithResponse(ctx, attributeName)
	if err != nil {
		return fmt.Errorf("テナント属性削除APIの呼び出しに失敗: %w", err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 404 {
		return fmt.Errorf("テナント属性削除に失敗: ステータスコード %d", resp.StatusCode())
	}

	return nil
}

// cleanupInvitation は招待をクリーンアップします
func (cm *CleanupManager) cleanupInvitation(ctx context.Context, resource *TestResource) error {
	// メタデータからテナントIDを取得
	tenantID, ok := resource.Metadata["tenant_id"].(string)
	if !ok {
		return fmt.Errorf("テナントIDがメタデータに見つかりません")
	}

	resp, err := cm.client.Client.DeleteTenantInvitationWithResponse(ctx, tenantID, resource.ID)
	if err != nil {
		return fmt.Errorf("招待削除APIの呼び出しに失敗: %w", err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 404 {
		return fmt.Errorf("招待削除に失敗: ステータスコード %d", resp.StatusCode())
	}

	return nil
}

// CleanupWithRetry はリトライ機能付きでクリーンアップを実行します
func (cm *CleanupManager) CleanupWithRetry(ctx context.Context, resource *TestResource, maxRetries int) error {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			// リトライ前に待機
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(i) * time.Second):
			}

			if cm.verbose {
				log.Printf("リソース %s のクリーンアップをリトライします (%d/%d)", resource.ID, i, maxRetries)
			}
		}

		lastErr = cm.cleanupResource(ctx, resource)
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

// GetCleanupSummary はクリーンアップの要約を取得します
func (cm *CleanupManager) GetCleanupSummary() *CleanupSummary {
	resources := cm.client.GetResources()
	summary := &CleanupSummary{
		TotalResources:   len(resources),
		CleanedResources: 0,
		FailedResources:  0,
		Results:          make([]CleanupResult, 0),
	}

	for _, resource := range resources {
		if resource.CleanedUp {
			summary.CleanedResources++
			if resource.CleanupErr != "" {
				summary.FailedResources++
			}
		}

		result := CleanupResult{
			ResourceType: resource.Type,
			ResourceID:   resource.ID,
			Success:      resource.CleanedUp && resource.CleanupErr == "",
		}
		if resource.CleanupErr != "" {
			result.Error = resource.CleanupErr
		}

		summary.Results = append(summary.Results, result)
	}

	return summary
}

// CleanupSummary はクリーンアップの要約を表します
type CleanupSummary struct {
	TotalResources   int             `json:"total_resources"`
	CleanedResources int             `json:"cleaned_resources"`
	FailedResources  int             `json:"failed_resources"`
	Results          []CleanupResult `json:"results"`
}

// ForceCleanupAll は強制的に全てのテストリソースをクリーンアップします
// 通常のクリーンアップで失敗したリソースも含めて削除を試行します
func (cm *CleanupManager) ForceCleanupAll(ctx context.Context) error {
	if cm.verbose {
		log.Println("強制クリーンアップを開始します")
	}

	// 既知のテストリソースパターンを使用してクリーンアップ
	if err := cm.forceCleanupTestRoles(ctx); err != nil {
		log.Printf("テストロールの強制クリーンアップに失敗: %v", err)
	}

	if err := cm.forceCleanupTestEnvs(ctx); err != nil {
		log.Printf("テスト環境の強制クリーンアップに失敗: %v", err)
	}

	if err := cm.forceCleanupTestUsers(ctx); err != nil {
		log.Printf("テストユーザーの強制クリーンアップに失敗: %v", err)
	}

	return nil
}

// forceCleanupTestRoles はテスト用ロールを強制クリーンアップします
func (cm *CleanupManager) forceCleanupTestRoles(ctx context.Context) error {
	// テスト用ロール名のパターン
	testRolePatterns := []string{"test_", "e2e_", "temp_"}

	resp, err := cm.client.Client.GetRolesWithResponse(ctx)
	if err != nil {
		return err
	}

	if resp.JSON200 == nil {
		return nil
	}

	for _, role := range resp.JSON200.Roles {
		for _, pattern := range testRolePatterns {
			if len(role.RoleName) > len(pattern) && role.RoleName[:len(pattern)] == pattern {
				if err := cm.cleanupRole(ctx, role.RoleName); err != nil {
					log.Printf("テストロール %s の削除に失敗: %v", role.RoleName, err)
				} else if cm.verbose {
					log.Printf("テストロール %s を削除しました", role.RoleName)
				}
				break
			}
		}
	}

	return nil
}

// forceCleanupTestEnvs はテスト用環境を強制クリーンアップします
func (cm *CleanupManager) forceCleanupTestEnvs(ctx context.Context) error {
	testEnvPatterns := []string{"test", "e2e", "temp"}

	resp, err := cm.client.Client.GetEnvsWithResponse(ctx)
	if err != nil {
		return err
	}

	if resp.JSON200 == nil {
		return nil
	}

	for _, env := range resp.JSON200.Envs {
		for _, pattern := range testEnvPatterns {
			envIdStr := fmt.Sprintf("%d", env.Id)
			if len(envIdStr) > len(pattern) && envIdStr[:len(pattern)] == pattern {
				if err := cm.cleanupEnv(ctx, envIdStr); err != nil {
					log.Printf("テスト環境 %d の削除に失敗: %v", env.Id, err)
				} else if cm.verbose {
					log.Printf("テスト環境 %d を削除しました", env.Id)
				}
				break
			}
		}
	}

	return nil
}

// forceCleanupTestUsers はテスト用ユーザーを強制クリーンアップします
func (cm *CleanupManager) forceCleanupTestUsers(ctx context.Context) error {
	testEmailPatterns := []string{"test@", "e2e@", "temp@"}

	resp, err := cm.client.Client.GetSaasUsersWithResponse(ctx)
	if err != nil {
		return err
	}

	if resp.JSON200 == nil {
		return nil
	}

	for _, user := range resp.JSON200.Users {
		for _, pattern := range testEmailPatterns {
			if len(user.Email) > len(pattern) && user.Email[:len(pattern)] == pattern {
				if err := cm.cleanupSaasUser(ctx, user.Id); err != nil {
					log.Printf("テストユーザー %s の削除に失敗: %v", user.Id, err)
				} else if cm.verbose {
					log.Printf("テストユーザー %s を削除しました", user.Id)
				}
				break
			}
		}
	}

	return nil
}