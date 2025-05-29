package apilogapi

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/saasus-platform/saasus-sdk-go/client"
)

// E2E テスト用の設定
const (
	SaaSusAPIEndpoint = "https://api.dev.saasus.io/v1/apilog"
)

// TestParams はテストパラメータファイルの構造体
type TestParams struct {
	GetLogs struct {
		TestCases []struct {
			Name   string `json:"name"`
			Params struct {
				Limit       *int64  `json:"limit,omitempty"`
				CreatedDate *string `json:"created_date,omitempty"`
				CreatedAt   *string `json:"created_at,omitempty"`
				Cursor      *string `json:"cursor,omitempty"`
			} `json:"params"`
			ExpectError bool `json:"expectError,omitempty"`
		} `json:"testCases"`
		EdgeCases []struct {
			Name   string `json:"name"`
			Params struct {
				Limit       *int64  `json:"limit,omitempty"`
				CreatedDate *string `json:"created_date,omitempty"`
				CreatedAt   *string `json:"created_at,omitempty"`
				Cursor      *string `json:"cursor,omitempty"`
			} `json:"params"`
			ExpectError bool `json:"expectError"`
		} `json:"edgeCases"`
		Pagination struct {
			Limit int64 `json:"limit"`
		} `json:"pagination"`
	} `json:"getLogs"`
	GetLog struct {
		TestApiLogIds []string `json:"testApiLogIds"`
		InvalidIds    []string `json:"invalidIds"`
		EdgeCases     []struct {
			Name           string `json:"name"`
			Id             string `json:"id"`
			ExpectError    bool   `json:"expectError"`
			ExpectedStatus int    `json:"expectedStatus,omitempty"`
		} `json:"edgeCases"`
	} `json:"getLog"`
	Performance struct {
		MaxResponseTime string `json:"maxResponseTime"`
		TestLimit       int64  `json:"testLimit"`
		LoadTestCases   []struct {
			Name            string `json:"name"`
			Limit           int64  `json:"limit"`
			ExpectedMaxTime string `json:"expectedMaxTime"`
		} `json:"loadTestCases"`
	} `json:"performance"`
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

// convertToGetLogsParams は設定ファイルのパラメータをGetLogsParamsに変換します
func convertToGetLogsParams(paramConfig struct {
	Limit       *int64  `json:"limit,omitempty"`
	CreatedDate *string `json:"created_date,omitempty"`
	CreatedAt   *string `json:"created_at,omitempty"`
	Cursor      *string `json:"cursor,omitempty"`
}) *GetLogsParams {
	params := &GetLogsParams{}

	if paramConfig.Limit != nil {
		params.Limit = paramConfig.Limit
	}

	if paramConfig.CreatedDate != nil {
		if date, err := time.Parse("2006-01-02", *paramConfig.CreatedDate); err == nil {
			params.CreatedDate = &types.Date{Time: date}
		}
	}

	if paramConfig.CreatedAt != nil {
		if dateTime, err := time.Parse(time.RFC3339, *paramConfig.CreatedAt); err == nil {
			params.CreatedAt = &dateTime
		}
	}

	if paramConfig.Cursor != nil {
		params.Cursor = paramConfig.Cursor
	}

	return params
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

// TestE2E_GetLogsWithResponse は GetLogsWithResponse のE2Eテスト
func TestE2E_GetLogsWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)

	// 通常のテストケース
	for _, testCase := range testParams.GetLogs.TestCases {
		tc := testCase // ループ変数をキャプチャ
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			params := convertToGetLogsParams(tc.Params)
			resp, err := client.GetLogsWithResponse(ctx, params)

			if tc.ExpectError {
				if err == nil && (resp == nil || resp.StatusCode() < 400) {
					t.Error("エラーが期待されましたが、成功レスポンスが返されました")
				}
				return
			}

			if err != nil {
				t.Fatalf("リクエストエラー: %v", err)
			}

			if resp == nil {
				t.Fatal("レスポンスがnilです")
			}

			t.Logf("ステータスコード: %d", resp.StatusCode())
			t.Logf("ステータス: %s", resp.Status())

			switch resp.StatusCode() {
			case 200:
				// 成功レスポンスの検証
				if resp.JSON200 == nil {
					t.Error("200レスポンスが解析されませんでした")
					return
				}

				t.Logf("APIログ数: %d", len(resp.JSON200.ApiLogs))

				// 制限数パラメータの検証
				if params.Limit != nil {
					if int64(len(resp.JSON200.ApiLogs)) > *params.Limit {
						t.Errorf("取得件数が制限数を超えています: 制限=%d, 実際=%d", *params.Limit, len(resp.JSON200.ApiLogs))
					}
				}

				// カーソルの検証
				if resp.JSON200.Cursor != nil {
					t.Logf("カーソル: %s", *resp.JSON200.Cursor)
				}

				// 各ログの基本的な検証
				for i, log := range resp.JSON200.ApiLogs {
					if i >= 3 { // 最初の3件のみ詳細チェック
						break
					}

					t.Logf("ログ %d:", i+1)
					t.Logf("  - ID: %s", log.ApiLogId)
					t.Logf("  - トレースID: %s", log.TraceId)
					t.Logf("  - メソッド: %s", log.RequestMethod)
					t.Logf("  - URI: %s", log.RequestUri)
					t.Logf("  - ステータス: %s", log.ResponseStatus)
					t.Logf("  - 作成日時: %d", log.CreatedAt)
					t.Logf("  - 作成日: %s", log.CreatedDate)

					// 必須フィールドの検証
					if log.ApiLogId == "" {
						t.Errorf("ログ %d: APIログIDが空です", i+1)
					}
					if log.TraceId == "" {
						t.Errorf("ログ %d: トレースIDが空です", i+1)
					}
					if log.RequestMethod == "" {
						t.Errorf("ログ %d: リクエストメソッドが空です", i+1)
					}
					if log.RequestUri == "" {
						t.Errorf("ログ %d: リクエストURIが空です", i+1)
					}
					if log.ResponseStatus == "" {
						t.Errorf("ログ %d: レスポンスステータスが空です", i+1)
					}
					if log.CreatedAt == 0 {
						t.Errorf("ログ %d: 作成日時が0です", i+1)
					}
					if log.CreatedDate == "" {
						t.Errorf("ログ %d: 作成日が空です", i+1)
					}

					// 日付フィルタの検証
					if params.CreatedDate != nil {
						expectedDate := params.CreatedDate.Time.Format("2006-01-02")
						if log.CreatedDate != expectedDate {
							t.Logf("日付フィルタ: 期待=%s, 実際=%s", expectedDate, log.CreatedDate)
						}
					}
				}

			case 400:
				t.Logf("400エラー: 不正なリクエストパラメータ")
				// 境界値テストなどで期待される場合がある

			case 500:
				// サーバーエラーの検証
				if resp.JSON500 != nil {
					t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
				}
				t.Error("サーバーエラーが発生しました")

			default:
				t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
			}
		})
	}
}

// TestE2E_GetLogsWithResponse_EdgeCases は GetLogsWithResponse のエッジケーステスト
func TestE2E_GetLogsWithResponse_EdgeCases(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)

	for _, edgeCase := range testParams.GetLogs.EdgeCases {
		tc := edgeCase // ループ変数をキャプチャ
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			params := convertToGetLogsParams(tc.Params)
			resp, err := client.GetLogsWithResponse(ctx, params)

			if tc.ExpectError {
				// エラーが期待される場合
				if err != nil {
					t.Logf("期待通りエラーが発生しました: %v", err)
					return
				}
				if resp != nil && resp.StatusCode() >= 400 {
					t.Logf("期待通りエラーレスポンスが返されました: %d", resp.StatusCode())
					return
				}
				t.Error("エラーが期待されましたが、成功レスポンスが返されました")
				return
			}

			// エラーが期待されない場合の通常処理
			if err != nil {
				t.Fatalf("予期しないリクエストエラー: %v", err)
			}

			if resp == nil {
				t.Fatal("レスポンスがnilです")
			}

			t.Logf("ステータスコード: %d", resp.StatusCode())
			t.Logf("ステータス: %s", resp.Status())

			// エッジケースでも200が返される場合の検証
			if resp.StatusCode() == 200 && resp.JSON200 != nil {
				t.Logf("エッジケースで成功レスポンス: APIログ数=%d", len(resp.JSON200.ApiLogs))
			}
		})
	}
}

// TestE2E_GetLogWithResponse は GetLogWithResponse のE2Eテスト
func TestE2E_GetLogWithResponse(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)

	// まず、ログ一覧を取得して実際のAPIログIDを取得
	ctx := context.Background()
	logsResp, err := client.GetLogsWithResponse(ctx, &GetLogsParams{
		Limit: func() *int64 { v := int64(1); return &v }(), // 1件のみ取得
	})

	if err != nil {
		t.Fatalf("ログ一覧取得エラー: %v", err)
	}

	if logsResp == nil || logsResp.StatusCode() != 200 || logsResp.JSON200 == nil {
		t.Skip("ログ一覧の取得に失敗したため、個別ログ取得テストをスキップします")
	}

	if len(logsResp.JSON200.ApiLogs) == 0 {
		t.Skip("利用可能なAPIログがないため、個別ログ取得テストをスキップします")
	}

	// 実際のAPIログIDを使用してテスト
	actualLogId := logsResp.JSON200.ApiLogs[0].ApiLogId
	t.Logf("テスト対象のAPIログID: %s", actualLogId)

	// テストケースを構築（実際のIDと設定ファイルからの無効なIDを組み合わせ）
	testCases := []struct {
		name     string
		apiLogId ApiLogId
		expectOK bool
	}{
		{
			name:     "実際のAPIログID",
			apiLogId: actualLogId,
			expectOK: true,
		},
	}

	// 設定ファイルからの無効なIDを追加
	for _, invalidId := range testParams.GetLog.InvalidIds {
		testCases = append(testCases, struct {
			name     string
			apiLogId ApiLogId
			expectOK bool
		}{
			name:     "無効なAPIログID: " + invalidId,
			apiLogId: ApiLogId(invalidId),
			expectOK: false,
		})
	}

	// エッジケースを追加
	for _, edgeCase := range testParams.GetLog.EdgeCases {
		testCases = append(testCases, struct {
			name     string
			apiLogId ApiLogId
			expectOK bool
		}{
			name:     "エッジケース: " + edgeCase.Name,
			apiLogId: ApiLogId(edgeCase.Id),
			expectOK: !edgeCase.ExpectError,
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.GetLogWithResponse(ctx, tc.apiLogId)

			if err != nil {
				t.Fatalf("リクエストエラー: %v", err)
			}

			if resp == nil {
				t.Fatal("レスポンスがnilです")
			}

			t.Logf("ステータスコード: %d", resp.StatusCode())
			t.Logf("ステータス: %s", resp.Status())

			switch resp.StatusCode() {
			case 200:
				if !tc.expectOK {
					t.Error("存在しないIDで200レスポンスが返されました")
					return
				}

				// 成功レスポンスの検証
				if resp.JSON200 == nil {
					t.Error("200レスポンスが解析されませんでした")
					return
				}

				log := resp.JSON200
				t.Logf("APIログ詳細:")
				t.Logf("  - ID: %s", log.ApiLogId)
				t.Logf("  - トレースID: %s", log.TraceId)
				t.Logf("  - SaaSID: %s", log.SaasId)
				t.Logf("  - APIキー: %s", log.ApiKey)
				t.Logf("  - リクエストメソッド: %s", log.RequestMethod)
				t.Logf("  - リクエストURI: %s", log.RequestUri)
				t.Logf("  - レスポンスステータス: %s", log.ResponseStatus)
				t.Logf("  - リモートアドレス: %s", log.RemoteAddress)
				t.Logf("  - リファラー: %s", log.Referer)
				t.Logf("  - 作成日時: %d", log.CreatedAt)
				t.Logf("  - 作成日: %s", log.CreatedDate)
				t.Logf("  - TTL: %d", log.Ttl)

				// 必須フィールドの検証
				if log.ApiLogId != string(tc.apiLogId) {
					t.Errorf("APIログIDが一致しません: 期待=%s, 実際=%s", tc.apiLogId, log.ApiLogId)
				}
				if log.TraceId == "" {
					t.Error("トレースIDが空です")
				}
				if log.SaasId == "" {
					t.Error("SaaSIDが空です")
				}
				if log.ApiKey == "" {
					t.Error("APIキーが空です")
				}
				if log.RequestMethod == "" {
					t.Error("リクエストメソッドが空です")
				}
				if log.RequestUri == "" {
					t.Error("リクエストURIが空です")
				}
				if log.ResponseStatus == "" {
					t.Error("レスポンスステータスが空です")
				}
				if log.RemoteAddress == "" {
					t.Error("リモートアドレスが空です")
				}
				if log.CreatedAt == 0 {
					t.Error("作成日時が0です")
				}
				if log.CreatedDate == "" {
					t.Error("作成日が空です")
				}
				if log.Ttl == 0 {
					t.Error("TTLが0です")
				}

			case 400:
				if tc.expectOK {
					t.Error("存在するIDで400レスポンスが返されました")
				} else {
					t.Log("期待通り400レスポンスが返されました")
				}

			case 404:
				if tc.expectOK {
					t.Error("存在するIDで404レスポンスが返されました")
				} else {
					t.Log("期待通り404レスポンスが返されました")
				}

			case 500:
				// サーバーエラーの検証
				if resp.JSON500 != nil {
					t.Logf("サーバーエラー: Type=%s, Message=%s", resp.JSON500.Type, resp.JSON500.Message)
				}
				t.Error("サーバーエラーが発生しました")

			default:
				t.Errorf("予期しないステータスコード: %d", resp.StatusCode())
			}
		})
	}
}

// TestE2E_GetLogsWithPagination はページネーション機能のE2Eテスト
func TestE2E_GetLogsWithPagination(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)

	ctx := context.Background()

	// 最初のページを取得
	firstPageResp, err := client.GetLogsWithResponse(ctx, &GetLogsParams{
		Limit: &testParams.GetLogs.Pagination.Limit, // 設定ファイルから取得数を取得
	})

	if err != nil {
		t.Fatalf("最初のページ取得エラー: %v", err)
	}

	if firstPageResp == nil || firstPageResp.StatusCode() != 200 || firstPageResp.JSON200 == nil {
		t.Skip("最初のページの取得に失敗したため、ページネーションテストをスキップします")
	}

	t.Logf("最初のページ: %d件のログ", len(firstPageResp.JSON200.ApiLogs))

	// カーソルが存在する場合、次のページを取得
	if firstPageResp.JSON200.Cursor != nil && *firstPageResp.JSON200.Cursor != "" {
		t.Logf("カーソル: %s", *firstPageResp.JSON200.Cursor)

		secondPageResp, err := client.GetLogsWithResponse(ctx, &GetLogsParams{
			Limit:  &testParams.GetLogs.Pagination.Limit,
			Cursor: firstPageResp.JSON200.Cursor,
		})

		if err != nil {
			t.Fatalf("2番目のページ取得エラー: %v", err)
		}

		if secondPageResp != nil && secondPageResp.StatusCode() == 200 && secondPageResp.JSON200 != nil {
			t.Logf("2番目のページ: %d件のログ", len(secondPageResp.JSON200.ApiLogs))

			// 異なるログが取得されることを確認
			if len(firstPageResp.JSON200.ApiLogs) > 0 && len(secondPageResp.JSON200.ApiLogs) > 0 {
				firstLogId := firstPageResp.JSON200.ApiLogs[0].ApiLogId
				secondLogId := secondPageResp.JSON200.ApiLogs[0].ApiLogId

				if firstLogId == secondLogId {
					t.Error("ページネーションで同じログが返されました")
				} else {
					t.Log("ページネーションが正常に動作しています")
				}
			}
		}
	} else {
		t.Log("カーソルが存在しないため、ページネーションテストを終了します")
	}
}

// TestE2E_ErrorHandling はエラーハンドリングのE2Eテスト
func TestE2E_ErrorHandling(t *testing.T) {
	client := setupE2EClient(t)

	ctx := context.Background()

	t.Run("無効なAPIログID形式", func(t *testing.T) {
		invalidId := ApiLogId("invalid-uuid-format")
		resp, err := client.GetLogWithResponse(ctx, invalidId)

		if err != nil {
			t.Logf("リクエストエラー: %v", err)
			return
		}

		if resp == nil {
			t.Fatal("レスポンスがnilです")
		}

		t.Logf("無効なID形式でのレスポンス - ステータス: %d", resp.StatusCode())

		// 400 Bad Request または 404 Not Found が期待される
		if resp.StatusCode() != 400 && resp.StatusCode() != 404 {
			t.Logf("予期しないステータスコード: %d (400または404が期待される)", resp.StatusCode())
		}
	})

	t.Run("空のAPIログID", func(t *testing.T) {
		emptyId := ApiLogId("")
		resp, err := client.GetLogWithResponse(ctx, emptyId)

		if err != nil {
			t.Logf("リクエストエラー: %v", err)
			return
		}

		if resp == nil {
			t.Fatal("レスポンスがnilです")
		}

		t.Logf("空のID形式でのレスポンス - ステータス: %d", resp.StatusCode())

		// 400 Bad Request または 404 Not Found が期待される
		if resp.StatusCode() != 400 && resp.StatusCode() != 404 {
			t.Logf("予期しないステータスコード: %d (400または404が期待される)", resp.StatusCode())
		}
	})
}

// TestE2E_PerformanceBasic は基本的なパフォーマンステスト
func TestE2E_PerformanceBasic(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)

	ctx := context.Background()

	// レスポンス時間の測定
	start := time.Now()
	resp, err := client.GetLogsWithResponse(ctx, &GetLogsParams{
		Limit: &testParams.Performance.TestLimit, // 設定ファイルから制限数を取得
	})
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
		t.Logf("取得したログ数: %d", len(resp.JSON200.ApiLogs))
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
			resp, err := client.GetLogsWithResponse(ctx, &GetLogsParams{
				Limit: &tc.Limit,
			})
			duration := time.Since(start)

			if err != nil {
				t.Fatalf("パフォーマンステストエラー: %v", err)
			}

			if resp == nil {
				t.Fatal("レスポンスがnilです")
			}

			t.Logf("制限数: %d", tc.Limit)
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
				actualCount := int64(len(resp.JSON200.ApiLogs))
				t.Logf("取得したログ数: %d", actualCount)

				// 制限数の検証
				if actualCount > tc.Limit {
					t.Errorf("取得件数が制限数を超えています: 制限=%d, 実際=%d", tc.Limit, actualCount)
				}

				// スループットの計算
				if duration > 0 {
					throughput := float64(actualCount) / duration.Seconds()
					t.Logf("スループット: %.2f件/秒", throughput)
				}
			}
		})
	}
}

// TestE2E_ConcurrentRequests は同時リクエストのパフォーマンステスト
func TestE2E_ConcurrentRequests(t *testing.T) {
	client := setupE2EClient(t)
	testParams := loadTestParams(t)

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
			resp, err := client.GetLogsWithResponse(ctx, &GetLogsParams{
				Limit: &testParams.Performance.TestLimit,
			})
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
