package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// NewConfig はCLIConfigからConfigを作成します
func NewConfig(cliConfig *CLIConfig) *Config {
	return &Config{
		ReportPath: cliConfig.ReportPath,
		OutputDir:  cliConfig.OutputDir,
		ModuleName: cliConfig.ModuleName,
		Overwrite:  cliConfig.Overwrite,
		Verbose:    cliConfig.Verbose,
	}
}

// Validate は設定の妥当性を検証します
func (c *Config) Validate() error {
	// 必須パラメータのチェック
	if c.ReportPath == "" {
		return NewInvalidConfigError("report_path", "レポートファイルパスが指定されていません")
	}

	if c.ModuleName == "" {
		return NewInvalidConfigError("module_name", "モジュール名が指定されていません")
	}

	if c.OutputDir == "" {
		return NewInvalidConfigError("output_dir", "出力ディレクトリが指定されていません")
	}

	// レポートファイルの存在確認
	if _, err := os.Stat(c.ReportPath); os.IsNotExist(err) {
		return NewFileNotFoundError(c.ReportPath, err)
	}

	// レポートファイルが読み取り可能か確認
	file, err := os.Open(c.ReportPath)
	if err != nil {
		return WrapError(err, "レポートファイルを開けません: %s", c.ReportPath)
	}
	file.Close()

	return nil
}

// EnsureOutputDir は出力ディレクトリを作成します
func (c *Config) EnsureOutputDir() error {
	// ディレクトリが既に存在するか確認
	if info, err := os.Stat(c.OutputDir); err == nil {
		// 存在する場合、ディレクトリかどうか確認
		if !info.IsDir() {
			return NewInvalidConfigError("output_dir", fmt.Sprintf("出力パスはディレクトリではありません: %s", c.OutputDir))
		}
		// 既に存在する場合は何もしない
		return nil
	}

	// ディレクトリを作成（親ディレクトリも含めて）
	if err := os.MkdirAll(c.OutputDir, 0755); err != nil {
		return NewOutputDirCreateError(c.OutputDir, err)
	}

	if c.Verbose {
		fmt.Printf("出力ディレクトリを作成しました: %s\n", c.OutputDir)
	}

	return nil
}

// GetOutputPaths は出力ファイルのパスを返します
func (c *Config) GetOutputPaths() (paramsPath, helpersPath string) {
	paramsPath = filepath.Join(c.OutputDir, "test_params.json")
	helpersPath = filepath.Join(c.OutputDir, fmt.Sprintf("%s_test_helpers.go", c.ModuleName))
	return
}

// ShouldOverwriteFile はファイルを上書きすべきかどうかを判定します
func (c *Config) ShouldOverwriteFile(filePath string) bool {
	// ファイルが存在しない場合は上書き確認不要
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return true
	}

	// -overwrite フラグが指定されている場合は上書き
	if c.Overwrite {
		if c.Verbose {
			fmt.Printf("既存ファイルを上書きします: %s\n", filePath)
		}
		return true
	}

	// -overwrite フラグが指定されていない場合は上書きしない
	return false
}

// LogVerbose は詳細ログモードが有効な場合にメッセージを出力します
func (c *Config) LogVerbose(format string, args ...any) {
	if c.Verbose {
		fmt.Printf(format+"\n", args...)
	}
}

// LogInfo は情報メッセージを出力します
func (c *Config) LogInfo(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

// LogError はエラーメッセージを出力します
func (c *Config) LogError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "エラー: "+format+"\n", args...)
}

// defaultOutputDir はモジュール名に応じた出力先を返します
func defaultOutputDir(moduleName string) string {
	base := filepath.Join("tests", "tools", "test-param-generator", "tests", "e2e")
	if moduleName != "" {
		base = filepath.Join(base, moduleName)
	}
	if root := projectRoot(); root != "" {
		return filepath.Join(root, base)
	}
	return base
}

func projectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}
	return ""
}
