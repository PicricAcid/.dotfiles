package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"focus/internal/config"
)

var (
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "focus",
	Short: "Nix/Home-Manager パッケージ管理ツール",
	Long: `focusはNix/Home-ManagerのパッケージをHomebrewのような直感的なCLIで管理するツールです。

e.g.
	focus init		# 初期設定
	focus install ripgrep	# パッケージインストール
	focus list		# インストール済みパッケージ一覧
	focus uninstall ripgrep	# パッケージ削除
	focus search fzf	# パッケージ検索
	focus update ripgrep	# パッケージ更新`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "設定ファイルのパス")
}

func getConfigPath() string {
	if configPath != "" {
		return configPath
	}

	if envPath := os.Getenv("FOCUS_CONFIG"); envPath != "" {
		return envPath
	}

	candidates := []string{
		"./focus.toml",
		"~/.focus.toml",
		"~/.config/focus/config.toml",
	}

	for _, path := range candidates {
		if config.Exists(path) {
			return path
		}
	}

	return "./focus.toml"
}

func loadConfig() (*config.Config, error) {
	path := getConfigPath()

	if !config.Exists(path) {
		return nil, fmt.Errorf("設定ファイルが見つかりません: %s\n'focus init'を実行して初期設定を行ってください", path)
	}

	return config.Load(path)
}

// gitAddFile はFlake環境の場合にファイルをgit addする
func gitAddFile(cfg *config.Config, filePath string) error {
	// Flakeを使っていない場合は何もしない
	if !cfg.UseFlake {
		return nil
	}

	// git repositoryかチェック
	checkCmd := exec.Command("git", "-C", cfg.FlakePath, "rev-parse", "--git-dir")
	if err := checkCmd.Run(); err != nil {
		// git repositoryでない場合はスキップ（エラーにしない）
		return nil
	}

	// git add を実行
	addCmd := exec.Command("git", "-C", cfg.FlakePath, "add", filePath)
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("git addに失敗: %w", err)
	}

	return nil
}
