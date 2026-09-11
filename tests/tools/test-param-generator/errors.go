package main

import (
	"errors"
	"fmt"
	"os"
)

// エラー種別の定義
var (
	// ErrFileNotFound はファイルが存在しない場合のエラー
	ErrFileNotFound = errors.New("ファイルが見つかりません")

	// ErrInvalidJSON はJSON解析に失敗した場合のエラー
	ErrInvalidJSON = errors.New("JSONの形式が不正です")

	// ErrInvalidReport はレポート構造が不正な場合のエラー
	ErrInvalidReport = errors.New("レポートの構造が不正です")

	// ErrOutputDirCreate は出力ディレクトリの作成に失敗した場合のエラー
	ErrOutputDirCreate = errors.New("出力ディレクトリの作成に失敗しました")

	// ErrFileWrite はファイル書き込みに失敗した場合のエラー
	ErrFileWrite = errors.New("ファイルの書き込みに失敗しました")

	// ErrTemplateRender はテンプレート実行に失敗した場合のエラー
	ErrTemplateRender = errors.New("テンプレートの実行に失敗しました")

	// ErrInvalidConfig は設定が不正な場合のエラー
	ErrInvalidConfig = errors.New("設定が不正です")

	// ErrTypeAnalysis は型解析に失敗した場合のエラー
	ErrTypeAnalysis = errors.New("型解析に失敗しました")
)

// FileNotFoundError はファイルが見つからない場合のエラー詳細を提供します
type FileNotFoundError struct {
	Path string
	Err  error
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("ファイルが見つかりません: %s", e.Path)
}

func (e *FileNotFoundError) Unwrap() error {
	return e.Err
}

// NewFileNotFoundError はFileNotFoundErrorを作成します
func NewFileNotFoundError(path string, err error) error {
	return &FileNotFoundError{
		Path: path,
		Err:  err,
	}
}

// InvalidJSONError はJSON解析エラーの詳細を提供します
type InvalidJSONError struct {
	Path string
	Err  error
}

func (e *InvalidJSONError) Error() string {
	return fmt.Sprintf("JSONの解析に失敗しました: %s\n詳細: %v", e.Path, e.Err)
}

func (e *InvalidJSONError) Unwrap() error {
	return e.Err
}

// NewInvalidJSONError はInvalidJSONErrorを作成します
func NewInvalidJSONError(path string, err error) error {
	return &InvalidJSONError{
		Path: path,
		Err:  err,
	}
}

// InvalidReportError はレポート構造エラーの詳細を提供します
type InvalidReportError struct {
	Field   string
	Message string
}

func (e *InvalidReportError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("レポートの構造が不正です: フィールド '%s' - %s", e.Field, e.Message)
	}
	return fmt.Sprintf("レポートの構造が不正です: %s", e.Message)
}

// NewInvalidReportError はInvalidReportErrorを作成します
func NewInvalidReportError(field, message string) error {
	return &InvalidReportError{
		Field:   field,
		Message: message,
	}
}

// OutputDirCreateError は出力ディレクトリ作成エラーの詳細を提供します
type OutputDirCreateError struct {
	Path string
	Err  error
}

func (e *OutputDirCreateError) Error() string {
	return fmt.Sprintf("出力ディレクトリの作成に失敗しました: %s\n詳細: %v", e.Path, e.Err)
}

func (e *OutputDirCreateError) Unwrap() error {
	return e.Err
}

// NewOutputDirCreateError はOutputDirCreateErrorを作成します
func NewOutputDirCreateError(path string, err error) error {
	return &OutputDirCreateError{
		Path: path,
		Err:  err,
	}
}

// FileWriteError はファイル書き込みエラーの詳細を提供します
type FileWriteError struct {
	Path string
	Err  error
}

func (e *FileWriteError) Error() string {
	return fmt.Sprintf("ファイルの書き込みに失敗しました: %s\n詳細: %v", e.Path, e.Err)
}

func (e *FileWriteError) Unwrap() error {
	return e.Err
}

// NewFileWriteError はFileWriteErrorを作成します
func NewFileWriteError(path string, err error) error {
	return &FileWriteError{
		Path: path,
		Err:  err,
	}
}

// TemplateRenderError はテンプレート実行エラーの詳細を提供します
type TemplateRenderError struct {
	Template string
	Err      error
}

func (e *TemplateRenderError) Error() string {
	return fmt.Sprintf("テンプレートの実行に失敗しました: %s\n詳細: %v", e.Template, e.Err)
}

func (e *TemplateRenderError) Unwrap() error {
	return e.Err
}

// NewTemplateRenderError はTemplateRenderErrorを作成します
func NewTemplateRenderError(template string, err error) error {
	return &TemplateRenderError{
		Template: template,
		Err:      err,
	}
}

// InvalidConfigError は設定エラーの詳細を提供します
type InvalidConfigError struct {
	Field   string
	Message string
}

func (e *InvalidConfigError) Error() string {
	return fmt.Sprintf("設定が不正です: %s - %s", e.Field, e.Message)
}

// NewInvalidConfigError はInvalidConfigErrorを作成します
func NewInvalidConfigError(field, message string) error {
	return &InvalidConfigError{
		Field:   field,
		Message: message,
	}
}

// TypeAnalysisError は型解析エラーの詳細を提供します
type TypeAnalysisError struct {
	TypeName string
	Err      error
}

func (e *TypeAnalysisError) Error() string {
	return fmt.Sprintf("型解析に失敗しました: %s\n詳細: %v", e.TypeName, e.Err)
}

func (e *TypeAnalysisError) Unwrap() error {
	return e.Err
}

// NewTypeAnalysisError はTypeAnalysisErrorを作成します
func NewTypeAnalysisError(typeName string, err error) error {
	return &TypeAnalysisError{
		TypeName: typeName,
		Err:      err,
	}
}

// IsFileNotFoundError はエラーがファイル未検出エラーかどうかを判定します
func IsFileNotFoundError(err error) bool {
	var fnfErr *FileNotFoundError
	if errors.As(err, &fnfErr) {
		return true
	}
	return errors.Is(err, os.ErrNotExist)
}

// IsInvalidJSONError はエラーがJSON解析エラーかどうかを判定します
func IsInvalidJSONError(err error) bool {
	var ijErr *InvalidJSONError
	return errors.As(err, &ijErr)
}

// IsInvalidReportError はエラーがレポート構造エラーかどうかを判定します
func IsInvalidReportError(err error) bool {
	var irErr *InvalidReportError
	return errors.As(err, &irErr)
}

// WrapError はエラーをラップしてコンテキストを追加します
func WrapError(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", message, err)
}
