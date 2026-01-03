package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"skill_hunter/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "skill_hunterの初期設定を行う",
	Long: `対話的にskill_hunterの初期設定を行います。
以下の情報を設定します:
- 設定ファイルの保存先
- home.nixのパス
- skill-hunter-packages.nixのパス`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("skill_hunterの初期設定を開始します")
	fmt.Println()

	fmt.Printf("設定ファイルの保存先 [./skill_hunter.toml]: ")
	configPathInput, _ := reader.ReadString('\n')
	configPathInput = strings.TrimSpace(configPathInput)

	savePath := "./skill_hunter.toml"
	if configPathInput != "" {
		savePath = configPathInput
	}

	if config.Exists(savePath) {
		fmt.Printf("設定ファイル '%s' は既に存在します。上書きしますか？ [y/N]: ", savePath)
		confirm, _ := reader.ReadString('\n')
		if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(confirm)), "y") {
			fmt.Println("初期設定をキャンセルしました")
			return nil
		}
	}

	fmt.Printf("home.nixのパス [~/.config/home-manager/home.nix]: ")
	homeNixInput, _ := reader.ReadString('\n')
	homeNixInput = strings.TrimSpace(homeNixInput)

	homeNixPath := "~/.config/home-manager/home.nix"
	if homeNixInput != "" {
		homeNixPath = homeNixInput
	}

	expandedHomeNix, err := expandPathForInit(homeNixPath)
	if err != nil {
		return fmt.Errorf("home.nixのパス展開に失敗: %w", err)
	}

	if _, err := os.Stat(expandedHomeNix); os.IsNotExist(err) {
		fmt.Printf("警告: '%s'が見つかりません\n", expandedHomeNix)
		fmt.Print("続行しますか？ [y/N]: ")
		confirm, _ := reader.ReadString('\n')
		if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(confirm)), "y") {
			return nil
		}
	}

	homeNixDir := filepath.Dir(expandedHomeNix)
	packagesFilePath := filepath.Join(homeNixDir, "skill-hunter-packages.nix")

	fmt.Printf("skill-hunter-packages.nixのパス [%s]: ", packagesFilePath)
	packagesInput, _ := reader.ReadString('\n')
	packagesInput = strings.TrimSpace(packagesInput)

	if packagesInput != "" {
		packagesFilePath = packagesInput
		expandedPackages, err := expandPathForInit(packagesFilePath)
		if err != nil {
			return fmt.Errorf("packagesファイルのパス展開に失敗: %w", err)
		}
		packagesFilePath = expandedPackages
	}

	cfg := &config.Config{
		HomeNixPath:      homeNixPath,
		PackagesFilePath: packagesFilePath,
	}

	if err := config.Save(savePath, cfg); err != nil {
		return fmt.Errorf("設定ファイルの保存に失敗: %w", err)
	}

	fmt.Printf("\n☑️ 設定ファイルを保存しました: %s\n", savePath)

	if err := createPackagesFile(packagesFilePath); err != nil {
		return fmt.Errorf("packages ファイルの作成に失敗: %w", err)
	}

	fmt.Printf("☑️ %sを作成しました\n", packagesFilePath)

	if err := addImportToHomeNix(expandedHomeNix, packagesFilePath); err != nil {
		return fmt.Errorf("home.nixへのimport追加に失敗: %w", err)
	}

	fmt.Printf("☑️ %sにimport文を追加しました\n", expandedHomeNix)

	// 設定ファイルが既に存在する場合（再初期化）、Flake設定があればgit addを試みる
	if loadedCfg, err := config.Load(savePath); err == nil && loadedCfg.UseFlake {
		if gitErr := gitAddFile(loadedCfg, packagesFilePath); gitErr != nil {
			fmt.Fprintf(os.Stderr, "警告: git addに失敗しました: %v\n", gitErr)
		}
	}

	fmt.Println()
	fmt.Println("初期設定が完了しました！")
	fmt.Println("次のコマンドでパッケージをインストールできます:")
	fmt.Println("	skill_hunter install <package>")

	return nil
}

func createPackagesFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	content := `{ pkgs, ... }: {
		home.packages = with pkgs; [
			# skill_hunterでインストールしたパッケージ
		];
	}
	`

	return os.WriteFile(path, []byte(content), 0644)
}

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

	backupPath := homeNixPath + ".bak"
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("バックアップの作成に失敗: %w", err)
	}

	if err := os.WriteFile(homeNixPath, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func expandPathForInit(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if len(path) == 1 {
		return homeDir, nil
	}

	if path[1] == '/' {
		return filepath.Join(homeDir, path[2:]), nil
	}

	return path, nil
}
