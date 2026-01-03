package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"skill_hunter/internal/config"
)

var (
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "skill_hunter",
	Short: "Nix/Home-Manager パッケージ管理ツール",
	Long: `skill_hunterはNix/Home-ManagerのパッケージをHomebrewのような直感的なCLIで管理するツールです。

e.g.
	skill_hunter init		# 初期設定
	skill_hunter install ripgrep	# パッケージインストール
	skill_hunter list		# インストール済みパッケージ一覧
	skill_hunter uninstall ripgrep	# パッケージ削除
	skill_hunter search fzf		# パッケージ検索
	skill_hunter update ripgrep	# パッケージ更新`,
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

	if envPath := os.Getenv("SKILL_HUNTER_CONFIG"); envPath != "" {
		return envPath
	}

	candidates := []string{
		"./skill_hunter.toml",
		"~/.skill_hunter.toml",
		"~/.config/skill_hunter/config.toml",
	}

	for _, path := range candidates {
		if config.Exists(path) {
			return path
		}
	}

	return "./skill_hunter.toml"
}

func loadConfig() (*config.Config, error) {
	path := getConfigPath()

	if !config.Exists(path) {
		return nil, fmt.Errorf("設定ファイルが見つかりません: %s\n'skill_hunter initを実行して初期設定を行ってください", path)
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
