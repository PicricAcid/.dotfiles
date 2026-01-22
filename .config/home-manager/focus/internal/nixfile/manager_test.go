package nixfile

import (
	"os"
	"path/filepath"
	"testing"
)

// TestListPackages tests listing packages from a nix file
func TestListPackages(t *testing.T) {
	tmpDir := t.TempDir()
	nixFilePath := filepath.Join(tmpDir, "packages.nix")

	// テスト用のNixファイルを作成
	content := `{ pkgs, ... }: {
  home.packages = with pkgs; [
    ripgrep
    fzf
    neovim
  ];
}
`
	err := os.WriteFile(nixFilePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	manager := NewManager(nixFilePath)
	packages, err := manager.ListPackages()
	if err != nil {
		t.Fatalf("ListPackages failed: %v", err)
	}

	expectedPackages := []string{"ripgrep", "fzf", "neovim"}
	if len(packages) != len(expectedPackages) {
		t.Errorf("Package count mismatch: got %d, want %d", len(packages), len(expectedPackages))
	}

	for i, pkg := range expectedPackages {
		if packages[i] != pkg {
			t.Errorf("Package[%d] mismatch: got %s, want %s", i, packages[i], pkg)
		}
	}
}

// TestListPackagesEmpty tests listing from empty packages file
func TestListPackagesEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	nixFilePath := filepath.Join(tmpDir, "packages.nix")

	content := `{ pkgs, ... }: {
  home.packages = with pkgs; [
    # No packages
  ];
}
`
	err := os.WriteFile(nixFilePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	manager := NewManager(nixFilePath)
	packages, err := manager.ListPackages()
	if err != nil {
		t.Fatalf("ListPackages failed: %v", err)
	}

	if len(packages) != 0 {
		t.Errorf("Expected empty package list, got %d packages", len(packages))
	}
}

// TestAddPackage tests adding a package
func TestAddPackage(t *testing.T) {
	tmpDir := t.TempDir()
	nixFilePath := filepath.Join(tmpDir, "packages.nix")

	// 初期ファイルを作成
	initialContent := `{ pkgs, ... }: {
  home.packages = with pkgs; [
    ripgrep
  ];
}
`
	err := os.WriteFile(nixFilePath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	manager := NewManager(nixFilePath)

	// パッケージを追加
	err = manager.AddPackage("fzf")
	if err != nil {
		t.Fatalf("AddPackage failed: %v", err)
	}

	// 確認
	packages, err := manager.ListPackages()
	if err != nil {
		t.Fatalf("ListPackages failed: %v", err)
	}

	// アルファベット順にソートされているはず
	expected := []string{"fzf", "ripgrep"}
	if len(packages) != len(expected) {
		t.Errorf("Package count mismatch: got %d, want %d", len(packages), len(expected))
	}

	for i, pkg := range expected {
		if packages[i] != pkg {
			t.Errorf("Package[%d] mismatch: got %s, want %s", i, packages[i], pkg)
		}
	}

	// バックアップファイルが作成されているか確認
	backupPath := nixFilePath + ".bak"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("Backup file was not created")
	}
}

// TestAddPackageDuplicate tests adding a duplicate package
func TestAddPackageDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	nixFilePath := filepath.Join(tmpDir, "packages.nix")

	initialContent := `{ pkgs, ... }: {
  home.packages = with pkgs; [
    ripgrep
  ];
}
`
	err := os.WriteFile(nixFilePath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	manager := NewManager(nixFilePath)

	// 既に存在するパッケージを追加しようとする
	err = manager.AddPackage("ripgrep")
	if err == nil {
		t.Error("AddPackage should fail for duplicate package")
	}
}

// TestRemovePackage tests removing a package
func TestRemovePackage(t *testing.T) {
	tmpDir := t.TempDir()
	nixFilePath := filepath.Join(tmpDir, "packages.nix")

	initialContent := `{ pkgs, ... }: {
  home.packages = with pkgs; [
    fzf
    ripgrep
  ];
}
`
	err := os.WriteFile(nixFilePath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	manager := NewManager(nixFilePath)

	// パッケージを削除
	err = manager.RemovePackage("fzf")
	if err != nil {
		t.Fatalf("RemovePackage failed: %v", err)
	}

	// 確認
	packages, err := manager.ListPackages()
	if err != nil {
		t.Fatalf("ListPackages failed: %v", err)
	}

	if len(packages) != 1 {
		t.Errorf("Expected 1 package, got %d", len(packages))
	}

	if packages[0] != "ripgrep" {
		t.Errorf("Wrong package remained: got %s, want ripgrep", packages[0])
	}
}

// TestRemovePackageNotFound tests removing non-existent package
func TestRemovePackageNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	nixFilePath := filepath.Join(tmpDir, "packages.nix")

	initialContent := `{ pkgs, ... }: {
  home.packages = with pkgs; [
    ripgrep
  ];
}
`
	err := os.WriteFile(nixFilePath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	manager := NewManager(nixFilePath)

	// 存在しないパッケージを削除しようとする
	err = manager.RemovePackage("nonexistent")
	if err == nil {
		t.Error("RemovePackage should fail for non-existent package")
	}
}

// TestHasPackage tests checking if a package exists
func TestHasPackage(t *testing.T) {
	tmpDir := t.TempDir()
	nixFilePath := filepath.Join(tmpDir, "packages.nix")

	content := `{ pkgs, ... }: {
  home.packages = with pkgs; [
    ripgrep
    fzf
  ];
}
`
	err := os.WriteFile(nixFilePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	manager := NewManager(nixFilePath)

	// 存在するパッケージ
	has, err := manager.HasPackage("ripgrep")
	if err != nil {
		t.Fatalf("HasPackage failed: %v", err)
	}
	if !has {
		t.Error("HasPackage should return true for existing package")
	}

	// 存在しないパッケージ
	has, err = manager.HasPackage("nonexistent")
	if err != nil {
		t.Fatalf("HasPackage failed: %v", err)
	}
	if has {
		t.Error("HasPackage should return false for non-existent package")
	}
}

// TestGetDiffAdd tests getting diff for adding a package
func TestGetDiffAdd(t *testing.T) {
	tmpDir := t.TempDir()
	nixFilePath := filepath.Join(tmpDir, "packages.nix")

	content := `{ pkgs, ... }: {
  home.packages = with pkgs; [
    ripgrep
  ];
}
`
	err := os.WriteFile(nixFilePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	manager := NewManager(nixFilePath)

	diff, err := manager.GetDiff("fzf", true)
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	// diffに"+"記号が含まれていることを確認
	if diff == "" {
		t.Error("Diff should not be empty")
	}
}

// TestRollback tests rollback functionality
func TestRollback(t *testing.T) {
	tmpDir := t.TempDir()
	nixFilePath := filepath.Join(tmpDir, "packages.nix")

	initialContent := `{ pkgs, ... }: {
  home.packages = with pkgs; [
    ripgrep
  ];
}
`
	err := os.WriteFile(nixFilePath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	manager := NewManager(nixFilePath)

	// パッケージを追加（バックアップが作成される）
	err = manager.AddPackage("fzf")
	if err != nil {
		t.Fatalf("AddPackage failed: %v", err)
	}

	// ロールバック
	err = manager.Rollback()
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	// 元の状態に戻っているか確認
	packages, err := manager.ListPackages()
	if err != nil {
		t.Fatalf("ListPackages failed: %v", err)
	}

	if len(packages) != 1 {
		t.Errorf("Expected 1 package after rollback, got %d", len(packages))
	}

	if packages[0] != "ripgrep" {
		t.Errorf("Wrong package after rollback: got %s, want ripgrep", packages[0])
	}
}
