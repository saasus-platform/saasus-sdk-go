# Postman Collections Directory

このディレクトリにPostman Collectionファイルを配置します。

## 使い方

1. PostmanからCollectionをエクスポート（Collection v2.1形式）
2. このディレクトリに配置
3. アナライザーを実行

```bash
# Example
go run analyzer.go -collection "collections/SaaSus Analysis API Schema.postman_collection.json"
```

## セキュリティ注意事項

⚠️ **重要:** Postman CollectionファイルにはAPIキーやシークレットが含まれる可能性があるため、Gitには含めません。

- このディレクトリの`*.json`ファイルは`.gitignore`で除外されています
- 実際のCollectionファイルをコミットしないでください
- テスト用のキーのみを使用してください

## ファイル命名規則

```
SaaSus <Service Name> API Schema.postman_collection.json
```

例:
- `SaaSus Analysis API Schema.postman_collection.json`
- `SaaSus Billing API Schema.postman_collection.json`
- `SaaSus Auth API Schema.postman_collection.json`
