package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	HomeNixPath      string `toml:"home_nix_path"`
	PackagesFilePath string `toml:"packages_file_path"`
	UseFlake         bool   `toml:"use_flake"`
	FlakePath        string `toml:"flake_path"`
	FlakeConfig      string `toml:"flake_config"`
}

func DefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ホームディレクトリの取得に失敗: %w", err)
	}
	return filepath.Join(homeDir, ".config", "skill_hunter", "config.toml"), nil
}

func Load(configPath string) (*Config, error) {
	expandedPath, err := expandPath(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("設定ファイルの解析に失敗: %w", err)
	}

	config.HomeNixPath, err = expandPath(config.HomeNixPath)
	if err != nil {
		return nil, fmt.Errorf("home_nix_pathの展開に失敗: %w", err)
	}
	config.PackagesFilePath, err = expandPath(config.PackagesFilePath)
	if err != nil {
		return nil, fmt.Errorf("packages_file_pathの展開に失敗: %w", err)
	}

	if config.UseFlake && config.FlakePath != "" {
		config.FlakePath, err = expandPath(config.FlakePath)
		if err != nil {
			return nil, fmt.Errorf("flake_pathの展開に失敗: %w", err)
		}
	}

	return &config, nil
}

func Save(configPath string, config *Config) error {
	expandedPath, err := expandPath(configPath)
	if err != nil {
		return err
	}

	dir := filepath.Dir(expandedPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗: %w", err)
	}

	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("設定のシリアライズに失敗: %w", err)
	}

	if err := os.WriteFile(expandedPath, data, 0644); err != nil {
		return fmt.Errorf("設定ファイルの書き込みに失敗: %w", err)
	}

	return nil
}

func Exists(configPath string) bool {
	expandedPath, err := expandPath(configPath)
	if err != nil {
		return false
	}

	_, err = os.Stat(expandedPath)
	return err == nil
}

func expandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ホームディレクトリの取得に失敗: %w", err)
	}

	if len(path) == 1 {
		return homeDir, nil
	}

	if path[1] == '/' {
		return filepath.Join(homeDir, path[2:]), nil
	}

	return path, nil
}
