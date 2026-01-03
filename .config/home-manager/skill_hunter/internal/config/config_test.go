package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSaveAndLoad tests saving and loading config
func TestSaveAndLoad(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.toml")

	// テスト用の設定を作成
	testConfig := &Config{
		HomeNixPath:      "/test/home.nix",
		PackagesFilePath: "/test/packages.nix",
	}

	// 保存
	err := Save(configPath, testConfig)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// ファイルが存在することを確認
	if !Exists(configPath) {
		t.Fatal("Config file does not exist after Save")
	}

	// 読み込み
	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// 内容を確認
	if loadedConfig.HomeNixPath != testConfig.HomeNixPath {
		t.Errorf("HomeNixPath mismatch: got %s, want %s", loadedConfig.HomeNixPath, testConfig.HomeNixPath)
	}

	if loadedConfig.PackagesFilePath != testConfig.PackagesFilePath {
		t.Errorf("PackagesFilePath mismatch: got %s, want %s", loadedConfig.PackagesFilePath, testConfig.PackagesFilePath)
	}
}

// TestExists tests the Exists function
func TestExists(t *testing.T) {
	tmpDir := t.TempDir()

	// 存在しないファイル
	nonExistentPath := filepath.Join(tmpDir, "nonexistent.toml")
	if Exists(nonExistentPath) {
		t.Error("Exists returned true for non-existent file")
	}

	// 存在するファイルを作成
	existingPath := filepath.Join(tmpDir, "existing.toml")
	err := os.WriteFile(existingPath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if !Exists(existingPath) {
		t.Error("Exists returned false for existing file")
	}
}

// TestLoadNonExistent tests loading a non-existent file
func TestLoadNonExistent(t *testing.T) {
	_, err := Load("/nonexistent/path/config.toml")
	if err == nil {
		t.Error("Load should fail for non-existent file")
	}
}

// TestSaveCreateDirectory tests that Save creates directory if needed
func TestSaveCreateDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.toml")

	testConfig := &Config{
		HomeNixPath:      "/test/home.nix",
		PackagesFilePath: "/test/packages.nix",
	}

	err := Save(configPath, testConfig)
	if err != nil {
		t.Fatalf("Save failed to create directory: %v", err)
	}

	if !Exists(configPath) {
		t.Error("Config file was not created in subdirectory")
	}
}

// TestExpandPathWithTilde tests expandPath function indirectly through Load
func TestExpandPathWithTilde(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	// チルダを含むパスを設定
	testConfig := &Config{
		HomeNixPath:      "~/test/home.nix",
		PackagesFilePath: "~/test/packages.nix",
	}

	err := Save(configPath, testConfig)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// チルダが展開されていることを確認
	homeDir, _ := os.UserHomeDir()
	expectedHomePath := filepath.Join(homeDir, "test/home.nix")
	expectedPackagesPath := filepath.Join(homeDir, "test/packages.nix")

	if loadedConfig.HomeNixPath != expectedHomePath {
		t.Errorf("Tilde not expanded in HomeNixPath: got %s, want %s", loadedConfig.HomeNixPath, expectedHomePath)
	}

	if loadedConfig.PackagesFilePath != expectedPackagesPath {
		t.Errorf("Tilde not expanded in PackagesFilePath: got %s, want %s", loadedConfig.PackagesFilePath, expectedPackagesPath)
	}
}
