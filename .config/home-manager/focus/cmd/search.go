package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"focus/internal/nix"
)

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "パッケージを検索する",
	Long: `nixpkgs から指定されたキーワードでパッケージを検索します。

例:
 focus search ripgrep
 focus search editor`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	keyword := args[0]

	nixClient := nix.NewClient()

	fmt.Printf("'%s' を検索しています...\n\n", keyword)

	results, err := nixClient.Search(keyword)
	if err != nil {
		return fmt.Errorf("検索に失敗: %w", err)
	}

	if len(results) == 0 {
		fmt.Printf("'%s' に一致するパッケージが見つかりませんでした\n", keyword)
		return nil
	}

	fmt.Printf("検索結果 (%d件):\n\n", len(results))

	for _, result := range results {
		fmt.Printf("  %s\n", result.Name)
		if result.Version != "" {
			fmt.Printf("	バージョン: %s\n", result.Version)
		}
		if result.Description != "" {
			fmt.Printf("	説明: %s\n", result.Description)
		}
		fmt.Println()
	}

	return nil
}
