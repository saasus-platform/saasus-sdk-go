# OpenAPI Resolver テストデータ

このディレクトリには、OpenAPI Operation Resolverのテストに使用するテストデータが含まれています。

## ファイル一覧

### test_spec.yml
基本的なテストケースをカバーする最小限のOpenAPI仕様ファイル。

**カバーするケース:**
- 基本的なHTTPメソッド（GET, POST, PUT, DELETE, PATCH）
- リクエストボディの有無
- パスパラメータ（単一、複数、ネスト）
- オプショナルなリクエストボディ
- operationIdが定義されていないエンドポイント
- クエリパラメータ
- ハイフンとアンダースコアを含むパス

**主要なエンドポイント:**
- `/test` - 基本的なCRUD操作
- `/test/{id}` - パスパラメータを含むCRUD操作
- `/users/{user_id}/posts/{post_id}` - 複数のパスパラメータ
- `/optional-body` - オプショナルなリクエストボディ
- `/no-operation-id` - operationIdが欠落（エラーケース）
- `/search` - クエリパラメータ

### edge_cases_spec.yml
エッジケースと特殊なシナリオをテストするためのOpenAPI仕様ファイル。

**カバーするケース:**
- UUID形式のパラメータ
- 数値IDのパラメータ
- 特殊文字を含むパス
- 深くネストされたリソースパス
- すべてのHTTPメソッドを持つエンドポイント
- 複数のコンテンツタイプ
- レスポンスボディがないエンドポイント
- 複数のタグ
- 非常に長いoperationId
- 空の説明
- descriptionフィールドがないエンドポイント
- 多数のパラメータ

**主要なエンドポイント:**
- `/standard/{param_id}` - 標準的なパスパラメータ
- `/items/{item_uuid}` - UUID形式
- `/numeric/{numeric_id}` - 数値ID
- `/api/v1/organizations/{org_id}/departments/{dept_id}/...` - 深いネスト
- `/all-methods/{id}` - 全HTTPメソッド
- `/multi-content-type` - 複数のコンテンツタイプ
- `/many-params/{id}` - 多数のパラメータ

### invalid_spec.yml
不正なOpenAPI仕様のテストケース。

**カバーするケース:**
- operationIdが欠落しているエンドポイント
- 不正なHTTPメソッド
- パラメータ定義が欠落しているパスパラメータ

**用途:**
- エラーハンドリングのテスト
- バリデーション機能のテスト
- 警告メッセージの生成テスト

### malformed.yml
YAML構文エラーを含むファイル。

**カバーするケース:**
- 不正なインデント
- 閉じられていないブロック
- YAML構文エラー

**用途:**
- YAMLパーサーのエラーハンドリングテスト
- 適切なエラーメッセージの表示テスト

## 使用方法

### 基本的なテスト
```go
resolver, err := NewOperationResolver("testdata/test_spec.yml")
if err != nil {
    t.Fatal(err)
}

mapping, err := resolver.Resolve("GET", "/test")
// ...
```

### エッジケースのテスト
```go
resolver, err := NewOperationResolver("testdata/edge_cases_spec.yml")
if err != nil {
    t.Fatal(err)
}

mapping, err := resolver.Resolve("GET", "/items/{item_uuid}")
// ...
```

### エラーハンドリングのテスト
```go
// 不正な仕様ファイル
resolver, err := NewOperationResolver("testdata/invalid_spec.yml")
// エラーが発生することを確認

// YAML構文エラー
resolver, err := NewOperationResolver("testdata/malformed.yml")
// 適切なエラーメッセージが返されることを確認
```

## テストケースの追加

新しいテストケースを追加する場合は、以下のガイドラインに従ってください：

1. **test_spec.yml**: 基本的な機能のテストケース
2. **edge_cases_spec.yml**: 特殊なケースやエッジケース
3. **invalid_spec.yml**: 不正なOpenAPI仕様
4. **malformed.yml**: YAML構文エラー

各ファイルは独立してテストできるように設計されています。

## 要件との対応

- **要件1.4**: 複数のAPIファイルの処理 → 複数のテストファイルで検証
- **要件2.6**: パスパラメータの正規化 → edge_cases_spec.ymlで様々な形式をテスト
