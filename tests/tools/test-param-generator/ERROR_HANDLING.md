# エラーハンドリング実装ガイド

## 概要

Test Param Generatorは、ユーザーフレンドリーなエラーメッセージを提供するために、包括的なエラーハンドリングシステムを実装しています。

## エラー種別

### 1. ファイル未検出エラー (FileNotFoundError)

**発生条件:**
- レポートファイルが存在しない
- 型定義ファイルが見つからない

**エラーメッセージ例:**
```
❌ エラーが発生しました:
  ファイルが見つかりません
  設定の検証に失敗: ファイルが見つかりません: nonexistent.json

💡 ヒント: ファイルパスが正しいか確認してください
```

**使用方法:**
```go
return NewFileNotFoundError(filePath, err)
```

### 2. JSON解析エラー (InvalidJSONError)

**発生条件:**
- JSONファイルの構文が不正
- JSONのデコードに失敗

**エラーメッセージ例:**
```
❌ エラーが発生しました:
  JSONの形式が不正です
  レポートの解析に失敗: JSONの解析に失敗しました: report.json
詳細: unexpected end of JSON input

💡 ヒント: JSONファイルの構文を確認してください
```

**使用方法:**
```go
return NewInvalidJSONError(filePath, err)
```

### 3. レポート構造エラー (InvalidReportError)

**発生条件:**
- 必須フィールドが欠落
- レポートの構造が期待と異なる

**エラーメッセージ例:**
```
❌ エラーが発生しました:
  レポートの構造が不正です
  レポートの解析に失敗: レポートの構造が不正です: フィールド 'collection_name' - 必須フィールドが空です

💡 ヒント: postman-analyzerで生成された正しいレポートファイルを使用してください
```

**使用方法:**
```go
return NewInvalidReportError("collection_name", "必須フィールドが空です")
```

### 4. 出力ディレクトリ作成エラー (OutputDirCreateError)

**発生条件:**
- ディレクトリの作成に失敗
- パーミッションエラー

**使用方法:**
```go
return NewOutputDirCreateError(dirPath, err)
```

### 5. ファイル書き込みエラー (FileWriteError)

**発生条件:**
- ファイルの書き込みに失敗
- ディスク容量不足

**使用方法:**
```go
return NewFileWriteError(filePath, err)
```

### 6. テンプレート実行エラー (TemplateRenderError)

**発生条件:**
- テンプレートの実行に失敗
- テンプレートデータが不正

**使用方法:**
```go
return NewTemplateRenderError("helpers", err)
```

### 7. 設定エラー (InvalidConfigError)

**発生条件:**
- 必須パラメータが不足
- 設定値が不正

**使用方法:**
```go
return NewInvalidConfigError("report_path", "レポートファイルパスが指定されていません")
```

### 8. 型解析エラー (TypeAnalysisError)

**発生条件:**
- 型定義ファイルの解析に失敗
- 構造体定義が見つからない

**使用方法:**
```go
return NewTypeAnalysisError(typeName, err)
```

## エラーラッピング

エラーをラップしてコンテキストを追加する場合は、`WrapError`関数を使用します。

```go
if err := someOperation(); err != nil {
    return WrapError(err, "操作の実行に失敗")
}
```

これにより、エラーチェーンが保持され、`errors.Is`や`errors.As`で元のエラーを判定できます。

## エラー判定

特定のエラー種別かどうかを判定するヘルパー関数が提供されています。

```go
// ファイル未検出エラーの判定
if IsFileNotFoundError(err) {
    // ファイルが見つからない場合の処理
}

// JSON解析エラーの判定
if IsInvalidJSONError(err) {
    // JSON解析エラーの処理
}

// レポート構造エラーの判定
if IsInvalidReportError(err) {
    // レポート構造エラーの処理
}
```

## ユーザーフレンドリーなエラー表示

`main.go`の`printError`関数は、エラーの種類に応じて適切なメッセージとヒントを表示します。

```go
func printError(err error) {
    fmt.Fprintf(os.Stderr, "\n❌ エラーが発生しました:\n")

    switch {
    case IsFileNotFoundError(err):
        fmt.Fprintf(os.Stderr, "  ファイルが見つかりません\n")
        fmt.Fprintf(os.Stderr, "  %v\n", err)
        fmt.Fprintf(os.Stderr, "\n💡 ヒント: ファイルパスが正しいか確認してください\n")

    case IsInvalidJSONError(err):
        fmt.Fprintf(os.Stderr, "  JSONの形式が不正です\n")
        fmt.Fprintf(os.Stderr, "  %v\n", err)
        fmt.Fprintf(os.Stderr, "\n💡 ヒント: JSONファイルの構文を確認してください\n")

    // ... その他のエラー種別
    }
}
```

## ベストプラクティス

1. **適切なエラー型を使用する**: エラーの種類に応じて適切なエラー型を選択します。

2. **エラーをラップする**: `WrapError`を使用してエラーにコンテキストを追加します。

3. **エラーチェーンを保持する**: `%w`を使用してエラーをラップし、元のエラーを保持します。

4. **ユーザーフレンドリーなメッセージ**: エラーメッセージは明確で、問題の解決方法を示唆するものにします。

5. **詳細情報を含める**: ファイルパス、フィールド名など、デバッグに役立つ情報を含めます。

## 要件との対応

- **要件 8.1**: ファイル読み込みエラーにファイルパスを含める ✓
- **要件 8.2**: JSON解析エラーに詳細を表示 ✓
- **要件 8.3**: ファイル書き込みエラーに書き込み先パスを含める ✓

すべてのエラーハンドリング要件が実装されています。
