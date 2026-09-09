package authapi

import (
	"path/filepath"

	"github.com/saasus-platform/saasus-sdk-go/tests/testlib/snapshot"
)

// GetAuthStorySnapshotConfig は Auth API のデフォルトスナップショット設定を返します。
//
// この関数は、スナップショットテストの基本設定を提供します。
// 各フェーズ（キャプチャ、比較、レポート、バリデーション）の有効/無効は
// コマンドラインフラグによって制御されます。
//
// 設定内容:
//   - 出力ディレクトリ: tests/e2e/snapshot
//   - モジュール名: auth
//   - キャプチャレベル: Full（完全なレスポンスをキャプチャ）
//   - バリデーションルール: Completion, Sequence, StateTransition, Timing
//
// 戻り値:
//   - *snapshot.StorySnapshotConfig: スナップショット設定
func GetAuthStorySnapshotConfig() *snapshot.StorySnapshotConfig {
	// Get absolute path to project root
	projectRoot := filepath.Join("..", "..", "..")
	outputDir := filepath.Join(projectRoot, "tests", "e2e", "snapshot")

	return &snapshot.StorySnapshotConfig{
		EnableCapture:    false, // Will be set by flag parser
		EnableComparison: false, // Will be set by flag parser
		EnableReporting:  false, // Will be set by flag parser
		EnableValidation: true,
		ComparisonMode:   "release",
		OutputDirectory:  outputDir,
		ModuleName:       "auth",
		SnapshotOnly:     false,
		FileNameFormat:   "story_snapshot_{tag}_{story_name}.json",
		CaptureLevel:     snapshot.CaptureLevelFull,
		ValidationRules:  GetAuthValidationRules(),
	}
}

// GetAuthValidationRules は Auth API 固有のバリデーションルールを返します。
//
// バリデーションルールは、スナップショットテストの検証基準を定義します。
//
// ルール一覧:
//   - Completion（Error レベル）: 全ステップが完了していることを検証
//   - Sequence（Error レベル）: ステップが正しい順序で実行されたことを検証
//   - StateTransition（Warning レベル）: 状態遷移が正しく行われたことを検証
//   - Timing（Info レベル、デフォルト無効）: 実行時間が許容範囲内であることを検証
//
// 戻り値:
//   - []snapshot.ValidationRule: バリデーションルールのスライス
func GetAuthValidationRules() []snapshot.ValidationRule {
	return []snapshot.ValidationRule{
		{
			Type:     snapshot.ValidationRuleCompletion,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityError,
			Parameters: map[string]interface{}{
				"description": "全ステップが完了していることを検証",
			},
		},
		{
			Type:     snapshot.ValidationRuleSequence,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityError,
			Parameters: map[string]interface{}{
				"description": "ステップが正しい順序で実行されたことを検証",
			},
		},
		{
			Type:     snapshot.ValidationRuleStateTransition,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityWarning,
			Parameters: map[string]interface{}{
				"description": "状態遷移が正しく行われたことを検証",
			},
		},
		{
			Type:     snapshot.ValidationRuleTiming,
			Enabled:  false, // デフォルトで無効（タイミングは変動する可能性があるため）
			Severity: snapshot.ValidationSeverityInfo,
			Parameters: map[string]interface{}{
				"description": "実行時間が許容範囲内であることを検証",
			},
		},
	}
}

// GetAuthSnapshotConfigForCapture はキャプチャ専用モードに最適化された設定を返します。
//
// このモードは、API レスポンスをキャプチャしてスナップショットファイルに保存します。
// 比較、レポート、バリデーションは実行されません。
//
// 使用例:
//   - 初回のスナップショット作成時
//   - API 仕様変更後のベースライン更新時
//
// 戻り値:
//   - *snapshot.StorySnapshotConfig: キャプチャ専用設定
func GetAuthSnapshotConfigForCapture() *snapshot.StorySnapshotConfig {
	config := GetAuthStorySnapshotConfig()
	config.EnableCapture = true
	config.EnableComparison = false
	config.EnableReporting = false
	config.EnableValidation = false
	config.SnapshotOnly = true
	return config
}

// GetAuthSnapshotConfigForComparison は比較専用モードに最適化された設定を返します。
//
// このモードは、現在の API レスポンスを既存のスナップショットと比較します。
// キャプチャ、レポート、バリデーションは実行されません。
//
// 使用例:
//   - 後方互換性の検証時
//   - リグレッションテスト時
//
// 戻り値:
//   - *snapshot.StorySnapshotConfig: 比較専用設定
func GetAuthSnapshotConfigForComparison() *snapshot.StorySnapshotConfig {
	config := GetAuthStorySnapshotConfig()
	config.EnableCapture = false
	config.EnableComparison = true
	config.EnableReporting = false
	config.EnableValidation = false
	return config
}

// GetAuthSnapshotConfigForReporting はレポート専用モードに最適化された設定を返します。
//
// このモードは、テスト結果のレポートを生成します。
// キャプチャ、比較、バリデーションは実行されません。
//
// 使用例:
//   - テスト結果の可視化時
//   - CI/CD パイプラインでのレポート生成時
//
// 戻り値:
//   - *snapshot.StorySnapshotConfig: レポート専用設定
func GetAuthSnapshotConfigForReporting() *snapshot.StorySnapshotConfig {
	config := GetAuthStorySnapshotConfig()
	config.EnableCapture = false
	config.EnableComparison = false
	config.EnableReporting = true
	config.EnableValidation = false
	return config
}

// GetAuthSnapshotConfigForIntegrated は統合実行（全フェーズ）用の設定を返します。
//
// このモードは、キャプチャ、比較、レポート、バリデーションの全フェーズを実行します。
//
// 使用例:
//   - 完全なスナップショットテストの実行時
//   - CI/CD パイプラインでの包括的なテスト時
//
// 戻り値:
//   - *snapshot.StorySnapshotConfig: 統合実行設定
func GetAuthSnapshotConfigForIntegrated() *snapshot.StorySnapshotConfig {
	config := GetAuthStorySnapshotConfig()
	config.EnableCapture = true
	config.EnableComparison = true
	config.EnableReporting = true
	config.EnableValidation = true
	config.SnapshotOnly = false
	return config
}
