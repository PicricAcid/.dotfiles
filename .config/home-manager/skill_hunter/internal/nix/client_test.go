package nix

import (
	"testing"
)

// TestNewClient tests creating a new client
func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("NewClient should not return nil")
	}
}

// TestSearchResult tests SearchResult struct
func TestSearchResult(t *testing.T) {
	result := SearchResult{
		Name:        "ripgrep",
		Description: "A search tool",
		Version:     "14.0.0",
	}

	if result.Name != "ripgrep" {
		t.Errorf("Name mismatch: got %s, want ripgrep", result.Name)
	}

	if result.Description != "A search tool" {
		t.Errorf("Description mismatch: got %s, want 'A search tool'", result.Description)
	}

	if result.Version != "14.0.0" {
		t.Errorf("Version mismatch: got %s, want 14.0.0", result.Version)
	}
}

// Note: 以下のテストは実際のNix環境が必要なため、統合テストとして扱う
// ユニットテストではモック化が必要だが、今回は基本的な構造テストのみ実装

// TestPackageExistsStructure tests the structure without actual execution
func TestPackageExistsStructure(t *testing.T) {
	client := NewClient()

	// この関数は実際にnixコマンドを実行するため、
	// Nix環境がない場合はスキップ
	t.Skip("Skipping test that requires Nix environment")

	_, err := client.PackageExists("ripgrep")
	// エラーの有無は環境依存なので、関数が実行できることのみ確認
	_ = err
}

// TestGetPackageVersionStructure tests the structure without actual execution
func TestGetPackageVersionStructure(t *testing.T) {
	client := NewClient()

	// この関数は実際にnixコマンドを実行するため、
	// Nix環境がない場合はスキップ
	t.Skip("Skipping test that requires Nix environment")

	_, err := client.GetPackageVersion("ripgrep")
	// エラーの有無は環境依存なので、関数が実行できることのみ確認
	_ = err
}

// TestSearchStructure tests the structure without actual execution
func TestSearchStructure(t *testing.T) {
	client := NewClient()

	// この関数は実際にnixコマンドを実行するため、
	// Nix環境がない場合はスキップ
	t.Skip("Skipping test that requires Nix environment")

	_, err := client.Search("ripgrep")
	// エラーの有無は環境依存なので、関数が実行できることのみ確認
	_ = err
}

// TestApplyHomeManagerStructure tests the structure without actual execution
func TestApplyHomeManagerStructure(t *testing.T) {
	client := NewClient()

	// この関数は実際にhome-manager switchを実行するため、
	// 環境がない場合はスキップ
	t.Skip("Skipping test that requires Home Manager environment")

	err := client.ApplyHomeManager("~/.config/home-manager/home.nix")
	// エラーの有無は環境依存なので、関数が実行できることのみ確認
	_ = err
}
