package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"skill_hunter/internal/nix"
	"skill_hunter/internal/nixfile"
)

var updateCmd = &cobra.Command{
	Use:   "update [package]",
	Short: "パッケージを更新する",
	Long: `パッケージを最新版に更新します。
パッケージ名を指定しない場合は、全てのパッケージを更新します。

注意: Nixではパッケージのバージョンはnixpkgsのバージョンに依存します。
このコマンドは home-manager switch を実行して更新を適用します。

例:
 skill_hunter update		# 全パッケージ更新
 skill_hunter update ripgrep	# ripgrepが含まれることを確認してから更新`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if len(args) > 0 {
		packageName := args[0]

		manager := nixfile.NewManager(cfg.PackagesFilePath)

		hasPackage, err := manager.HasPackage(packageName)
		if err != nil {
			return fmt.Errorf("パッケージチェックに失敗: %w", err)
		}

		if !hasPackage {
			return fmt.Errorf("パッケージ '%s' はインストールされていません", packageName)
		}

		fmt.Printf("パッケージ '%s' を含む全パッケージを更新します...\n\n", packageName)
	} else {
		fmt.Println("全パッケージを更新します...")
		fmt.Println()
	}

	nixClient := nix.NewClient()

	fmt.Println("home-manager switch を実行しています...")

	var switchErr error
	if cfg.UseFlake {
		switchErr = nixClient.(*nix.Client).ApplyHomeManagerWithFlake(cfg.FlakePath, cfg.FlakeConfig)
	} else {
		switchErr = nixClient.ApplyHomeManager(cfg.HomeNixPath)
	}

	if switchErr != nil {
		return fmt.Errorf("home-manager switch に失敗: %w", switchErr)
	}

	fmt.Println("\n✓ 更新が完了しました")

	if len(args) > 0 {
		packageName := args[0]
		version, err := nixClient.GetPackageVersion(packageName)
		if err == nil && version != "unknown" {
			fmt.Printf("	%s: %s\n", packageName, version)
		}
	}

	return nil
}
