package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"focus/internal/config"
	"focus/internal/nix"
	"focus/internal/nixfile"
)

const (
	testWorkspace = "./test_workspace"
	testUserHome  = "./test_workspace/testuser"
	testConfigDir = "./test_workspace/testuser/.config/focus"
)

// setupTestEnv はテスト環境をセットアップする
func setupTestEnv(t *testing.T) (*nix.MockClient, func()) {
	t.Helper()

	// 既存のテストワークスペースを削除（前のテストの影響を受けないように）
	os.RemoveAll(testWorkspace)

	// テストワークスペースを作成
	if err := os.MkdirAll(testConfigDir, 0755); err != nil {
		t.Fatalf("テストディレクトリの作成に失敗: %v", err)
	}

	// ホームディレクトリの作成
	if err := os.MkdirAll(testUserHome, 0755); err != nil {
		t.Fatalf("ホームディレクトリの作成に失敗: %v", err)
	}

	// テンプレートファイルのコピー
	homeNixTemplate := filepath.Join("testdata", "home.nix.template")
	destHomeNix := filepath.Join(testUserHome, "home.nix")

	templateContent, err := os.ReadFile(homeNixTemplate)
	if err != nil {
		t.Fatalf("テンプレートファイルの読み込みに失敗: %v", err)
	}

	if err := os.WriteFile(destHomeNix, templateContent, 0644); err != nil {
		t.Fatalf("home.nixのコピーに失敗: %v", err)
	}

	mockClient := nix.NewMockClient()

	// クリーンアップ関数（nixファイルを残すため、何もしない）
	cleanup := func() {
		t.Logf("テスト成果物は %s に残されています", testWorkspace)
	}

	return mockClient, cleanup
}

// TestInit はinitコマンドの統合テスト
func TestInit(t *testing.T) {
	mockClient, cleanup := setupTestEnv(t)
	defer cleanup()

	// 設定を作成
	cfg := &config.Config{
		HomeNixPath:      filepath.Join(testUserHome, "home.nix"),
		PackagesFilePath: filepath.Join(testUserHome, "focus-packages.nix"),
	}

	configPath := filepath.Join(testConfigDir, "config.toml")
	if err := config.Save(configPath, cfg); err != nil {
		t.Fatalf("設定の保存に失敗: %v", err)
	}

	// パッケージファイルを作成
	if err := createPackagesFile(cfg.PackagesFilePath); err != nil {
		t.Fatalf("パッケージファイルの作成に失敗: %v", err)
	}

	// home.nixにimportを追加
	if err := addImportToHomeNix(cfg.HomeNixPath, cfg.PackagesFilePath); err != nil {
		t.Fatalf("home.nixへのimport追加に失敗: %v", err)
	}

	// 検証1: config.tomlが存在する
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config.tomlが作成されていません")
	}

	// 検証2: focus-packages.nixが存在する
	if _, err := os.Stat(cfg.PackagesFilePath); os.IsNotExist(err) {
		t.Error("focus-packages.nixが作成されていません")
	}

	// 検証3: focus-packages.nixの内容が正しい
	content, err := os.ReadFile(cfg.PackagesFilePath)
	if err != nil {
		t.Fatalf("focus-packages.nixの読み込みに失敗: %v", err)
	}

	expectedContent := `# Focus managed packages
{ pkgs, ... }:
{
  home.packages = with pkgs; [
  ];
}
`
	if string(content) != expectedContent {
		t.Errorf("focus-packages.nixの内容が期待と異なります\n期待:\n%s\n実際:\n%s", expectedContent, string(content))
	}

	// 検証4: home.nixにimportが追加されている
	homeContent, err := os.ReadFile(cfg.HomeNixPath)
	if err != nil {
		t.Fatalf("home.nixの読み込みに失敗: %v", err)
	}

	if !strings.Contains(string(homeContent), "./focus-packages.nix") {
		t.Error("home.nixにfocus-packages.nixのimportが追加されていません")
	}

	_ = mockClient // モッククライアントは今回使用しない
}

// TestInstall はinstallコマンドの統合テスト
func TestInstall(t *testing.T) {
	mockClient, cleanup := setupTestEnv(t)
	defer cleanup()

	// 設定を作成
	cfg := &config.Config{
		HomeNixPath:      filepath.Join(testUserHome, "home.nix"),
		PackagesFilePath: filepath.Join(testUserHome, "focus-packages.nix"),
	}

	configPath := filepath.Join(testConfigDir, "config.toml")
	if err := config.Save(configPath, cfg); err != nil {
		t.Fatalf("設定の保存に失敗: %v", err)
	}

	// 初期化
	if err := createPackagesFile(cfg.PackagesFilePath); err != nil {
		t.Fatalf("パッケージファイルの作成に失敗: %v", err)
	}

	manager := nixfile.NewManager(cfg.PackagesFilePath)

	// モッククライアントの設定
	mockClient.ShouldPackageExist = true
	mockClient.PackageVersions["curl"] = "8.0.0"

	// パッケージのインストール
	packageName := "curl"

	exists, err := mockClient.PackageExists(packageName)
	if err != nil {
		t.Fatalf("パッケージの存在確認に失敗: %v", err)
	}
	if !exists {
		t.Fatalf("パッケージ %s が存在しません", packageName)
	}

	if err := manager.AddPackage(packageName); err != nil {
		t.Fatalf("パッケージの追加に失敗: %v", err)
	}

	// 検証1: focus-packages.nixにパッケージが追加されている
	content, err := os.ReadFile(cfg.PackagesFilePath)
	if err != nil {
		t.Fatalf("focus-packages.nixの読み込みに失敗: %v", err)
	}

	if !strings.Contains(string(content), "curl") {
		t.Error("curlがfocus-packages.nixに追加されていません")
	}

	// 検証2: パッケージリストに含まれている
	packages, err := manager.ListPackages()
	if err != nil {
		t.Fatalf("パッケージ一覧の取得に失敗: %v", err)
	}

	found := false
	for _, pkg := range packages {
		if pkg == "curl" {
			found = true
			break
		}
	}
	if !found {
		t.Error("インストールしたパッケージがリストに含まれていません")
	}
}

// TestUninstall はuninstallコマンドの統合テスト
func TestUninstall(t *testing.T) {
	mockClient, cleanup := setupTestEnv(t)
	defer cleanup()

	// 設定を作成
	cfg := &config.Config{
		HomeNixPath:      filepath.Join(testUserHome, "home.nix"),
		PackagesFilePath: filepath.Join(testUserHome, "focus-packages.nix"),
	}

	configPath := filepath.Join(testConfigDir, "config.toml")
	if err := config.Save(configPath, cfg); err != nil {
		t.Fatalf("設定の保存に失敗: %v", err)
	}

	// 初期化とパッケージ追加
	if err := createPackagesFile(cfg.PackagesFilePath); err != nil {
		t.Fatalf("パッケージファイルの作成に失敗: %v", err)
	}

	manager := nixfile.NewManager(cfg.PackagesFilePath)

	if err := manager.AddPackage("curl"); err != nil {
		t.Fatalf("パッケージの追加に失敗: %v", err)
	}

	// パッケージのアンインストール
	if err := manager.RemovePackage("curl"); err != nil {
		t.Fatalf("パッケージの削除に失敗: %v", err)
	}

	// 検証1: focus-packages.nixからパッケージが削除されている
	content, err := os.ReadFile(cfg.PackagesFilePath)
	if err != nil {
		t.Fatalf("focus-packages.nixの読み込みに失敗: %v", err)
	}

	if strings.Contains(string(content), "curl") {
		t.Error("curlがfocus-packages.nixから削除されていません")
	}

	// 検証2: パッケージリストに含まれていない
	packages, err := manager.ListPackages()
	if err != nil {
		t.Fatalf("パッケージ一覧の取得に失敗: %v", err)
	}

	for _, pkg := range packages {
		if pkg == "curl" {
			t.Error("削除したパッケージがリストに含まれています")
		}
	}

	_ = mockClient
}

// TestList はlistコマンドの統合テスト
func TestList(t *testing.T) {
	mockClient, cleanup := setupTestEnv(t)
	defer cleanup()

	// 設定を作成
	cfg := &config.Config{
		HomeNixPath:      filepath.Join(testUserHome, "home.nix"),
		PackagesFilePath: filepath.Join(testUserHome, "focus-packages.nix"),
	}

	configPath := filepath.Join(testConfigDir, "config.toml")
	if err := config.Save(configPath, cfg); err != nil {
		t.Fatalf("設定の保存に失敗: %v", err)
	}

	// 初期化と複数パッケージ追加
	if err := createPackagesFile(cfg.PackagesFilePath); err != nil {
		t.Fatalf("パッケージファイルの作成に失敗: %v", err)
	}

	manager := nixfile.NewManager(cfg.PackagesFilePath)

	testPackages := []string{"curl", "git", "vim"}
	for _, pkg := range testPackages {
		if err := manager.AddPackage(pkg); err != nil {
			t.Fatalf("パッケージ %s の追加に失敗: %v", pkg, err)
		}
	}

	// パッケージ一覧の取得
	packages, err := manager.ListPackages()
	if err != nil {
		t.Fatalf("パッケージ一覧の取得に失敗: %v", err)
	}

	// 検証: すべてのパッケージがリストに含まれている
	if len(packages) != len(testPackages) {
		t.Errorf("パッケージ数が一致しません: 期待=%d, 実際=%d", len(testPackages), len(packages))
	}

	for _, expected := range testPackages {
		found := false
		for _, actual := range packages {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("パッケージ %s がリストに含まれていません", expected)
		}
	}

	_ = mockClient
}

// TestErrorHandling はエラーハンドリングの統合テスト
func TestErrorHandling(t *testing.T) {
	mockClient, cleanup := setupTestEnv(t)
	defer cleanup()

	// ApplyHomeManagerの失敗をシミュレート
	mockClient.ShouldApplyFail = true

	err := mockClient.ApplyHomeManager(filepath.Join(testUserHome, "home.nix"))
	if err == nil {
		t.Error("ApplyHomeManagerが失敗すべきところで成功しています")
	}

	expectedErrMsg := "mock: home-manager switch failed"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("エラーメッセージが期待と異なります: 期待=%s, 実際=%s", expectedErrMsg, err.Error())
	}

	// パッケージが存在しない場合のテスト
	mockClient.ShouldPackageExist = false

	exists, err := mockClient.PackageExists("nonexistent")
	if err != nil {
		t.Fatalf("PackageExistsでエラー: %v", err)
	}
	if exists {
		t.Error("存在しないパッケージがexistsになっています")
	}
}

// createPackagesFile はパッケージファイルを作成する
func createPackagesFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	content := `# Focus managed packages
{ pkgs, ... }:
{
  home.packages = with pkgs; [
  ];
}
`

	return os.WriteFile(path, []byte(content), 0644)
}

// addImportToHomeNix はhome.nixにimportを追加する
func addImportToHomeNix(homeNixPath, packagesFilePath string) error {
	data, err := os.ReadFile(homeNixPath)
	if err != nil {
		return err
	}

	content := string(data)

	homeNixDir := filepath.Dir(homeNixPath)
	relPath, err := filepath.Rel(homeNixDir, packagesFilePath)
	if err != nil {
		relPath = packagesFilePath
	}

	importLine := fmt.Sprintf("	./%s", relPath)

	if strings.Contains(content, importLine) || strings.Contains(content, packagesFilePath) {
		return nil
	}

	importsStart := strings.Index(content, "imports = [")
	if importsStart == -1 {
		// Nix関数の }: { または in { の後にimportsセクションを挿入
		var insertPos int
		moduleStart := strings.Index(content, "}: {")
		if moduleStart != -1 {
			insertPos = moduleStart + 4 // "}: {" の後
		} else {
			// let式を使っている場合は "in {" を探す
			inStart := strings.Index(content, "in {")
			if inStart == -1 {
				return fmt.Errorf("home.nixの構造を解析できません")
			}
			insertPos = inStart + 4 // "in {" の後
		}

		newImports := fmt.Sprintf("\n	imports = [\n%s\n	];\n", importLine)
		content = content[:insertPos] + newImports + content[insertPos:]
	} else {
		closeBracket := strings.Index(content[importsStart:], "];")
		if closeBracket == -1 {
			return fmt.Errorf("importsの終了位置が見つかりません")
		}

		insertPos := importsStart + closeBracket
		content = content[:insertPos] + fmt.Sprintf("\n%s\n	", importLine) + content[insertPos:]
	}

	return os.WriteFile(homeNixPath, []byte(content), 0644)
}
