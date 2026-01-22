package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"focus/internal/nix"
	"focus/internal/nixfile"
)

var installCmd = &cobra.Command{
	Use:   "install [package]",
	Short: "パッケージをインストールする",
	Long: `指定されたパッケージを focus-packages.nix に追加し、home-manager switch を実行してインストールします。

例:
 focus install ripgrep
 focus install fzf`,
	Args: cobra.ExactArgs(1),
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	packageName := args[0]

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	manager := nixfile.NewManager(cfg.PackagesFilePath)

	hasPackage, err := manager.HasPackage(packageName)
	if err != nil {
		return fmt.Errorf("パッケージチェックに失敗: %w", err)
	}

	if hasPackage {
		fmt.Printf("パッケージ '%s' は既にインストールされています\n", packageName)
		return nil
	}

	nixClient := nix.NewClient()

	fmt.Printf("パッケージ '%s' を検索しています...\n", packageName)
	exists, err := nixClient.PackageExists(packageName)
	if err != nil {
		return fmt.Errorf("パッケージの検索に失敗: %w", err)
	}

	if !exists {
		return fmt.Errorf("パッケージ '%s' が見つかりませんでした", packageName)
	}

	diff, err := manager.GetDiff(packageName, true)
	if err != nil {
		return fmt.Errorf("diff の生成に失敗: %w", err)
	}

	fmt.Println("\n変更内容:")
	fmt.Println(diff)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("続行しますか？ [y/N]: ")
	confirm, _ := reader.ReadString('\n')

	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(confirm)), "y") {
		fmt.Println("インストールをキャンセルしました")
		return nil
	}

	fmt.Printf("\nパッケージ '%s' を追加しています...\n", packageName)
	if err := manager.AddPackage(packageName); err != nil {
		return fmt.Errorf("パッケージの追加に失敗: %w", err)
	}

	fmt.Println("☑️ focus-packages.nix に追加しました")

	// Flake環境の場合、git addを実行
	if cfg.UseFlake {
		if err := gitAddFile(cfg, cfg.PackagesFilePath); err != nil {
			fmt.Fprintf(os.Stderr, "警告: git addに失敗しました: %v\n", err)
		}
	}

	fmt.Println("\nhome-manager switch を実行しています...")

	var switchErr error
	if cfg.UseFlake {
		switchErr = nixClient.(*nix.Client).ApplyHomeManagerWithFlake(cfg.FlakePath, cfg.FlakeConfig)
	} else {
		switchErr = nixClient.ApplyHomeManager(cfg.HomeNixPath)
	}

	if switchErr != nil {
		fmt.Fprintf(os.Stderr, "\nエラー: %v\n", switchErr)
		fmt.Println("ロールバックしています...")

		if rollbackErr := manager.Rollback(); rollbackErr != nil {
			return fmt.Errorf("ロールバックにも失敗しました: %w\n元のエラー: %v", rollbackErr, switchErr)
		}

		fmt.Println("✓ ロールバックが完了しました")
		return fmt.Errorf("home-manager switch に失敗しました")
	}

	fmt.Printf("\n☑️ パッケージ '%s' のインストールが完了しました\n", packageName)

	return nil
}
