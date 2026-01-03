package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"skill_hunter/internal/nix"
	"skill_hunter/internal/nixfile"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [package]",
	Short: "パッケージをアンインストールする",
	Long: `指定されたパッケージを skill-hunter-packages.nix から削除し、
home-manager switch を実行してアンインストールします。

例:
 skill_hunter uninstall ripgrep
 skill_hunter uninstall fzf`,
	Args: cobra.ExactArgs(1),
	RunE: runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) error {
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

	if !hasPackage {
		fmt.Printf("パッケージ '%s' はインストールされていません\n", packageName)
		return nil
	}

	diff, err := manager.GetDiff(packageName, false)
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
		fmt.Println("アンインストールをキャンセルしました")
		return nil
	}

	fmt.Printf("\nパッケージ '%s' を削除しています...\n", packageName)
	if err := manager.RemovePackage(packageName); err != nil {
		return fmt.Errorf("パッケージの削除に失敗: %w", err)
	}

	fmt.Println("☑️ skill-hunter-packages.nix から削除しました")

	// Flake環境の場合、git addを実行
	if cfg.UseFlake {
		if err := gitAddFile(cfg, cfg.PackagesFilePath); err != nil {
			fmt.Fprintf(os.Stderr, "警告: git addに失敗しました: %v\n", err)
		}
	}

	nixClient := nix.NewClient()

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

		fmt.Println("☑️ ロールバックが完了しました")
		return fmt.Errorf("home-manager switchに失敗しました")
	}

	fmt.Printf("\n☑️ パッケージ '%s' のアンインストールが完了しました\n", packageName)

	return nil
}
