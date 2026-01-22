package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"focus/internal/nix"
	"focus/internal/nixfile"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "インストール済みパッケージの一覧を表示する",
	Long: `focusでインストールしたパッケージの一覧を表示します。
各パッケージのバージョン情報も取得します。

例:
 focus list`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	manager := nixfile.NewManager(cfg.PackagesFilePath)

	packages, err := manager.ListPackages()
	if err != nil {
		return fmt.Errorf("パッケージ一覧の取得に失敗: %w", err)
	}

	if len(packages) == 0 {
		fmt.Println("インストール済みのパッケージはありません")
		return nil
	}

	nixClient := nix.NewClient()

	fmt.Printf("インストール済みパッケージ (%d個):\n", len(packages))

	for _, pkg := range packages {
		version, err := nixClient.GetPackageVersion(pkg)
		if err != nil {
			version = "unknown"
		}

		fmt.Printf("  - %s: %s\n", pkg, version)
	}

	return nil
}
