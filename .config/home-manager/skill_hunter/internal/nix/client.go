package nix

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// NixClient はNix操作のインターフェース
type NixClient interface {
	Search(keyword string) ([]SearchResult, error)
	PackageExists(packageName string) (bool, error)
	ApplyHomeManager(homeNixPath string) error
	GetPackageVersion(packageName string) (string, error)
}

// Client は実際のNixコマンドを実行するクライアント
type Client struct{}

// NewClient は新しいNixクライアントを作成する
func NewClient() NixClient {
	return &Client{}
}

func (c *Client) Search(keyword string) ([]SearchResult, error) {
	cmd := exec.Command("nix", "search", "nixpkgs", keyword, "--json")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("nix search の実行に失敗: %s\n%s", err, stderr.String())
	}

	output := stdout.String()

	results := []SearchResult{
		{
			Name:        keyword,
			Description: output,
		},
	}

	return results, nil
}

func (c *Client) PackageExists(packageName string) (bool, error) {
	cmd := exec.Command("nix", "search", "nixpkgs", packageName, "--json")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return false, nil
	}

	output := stdout.String()

	return len(strings.TrimSpace(output)) > 0 && output != "{}", nil
}

func (c *Client) ApplyHomeManager(homeNixPath string) error {
	// Note: homeNixPath is kept for backward compatibility but not used when flake is detected
	// For flake support, use ApplyHomeManagerWithConfig instead
	cmd := exec.Command("home-manager", "switch", "-f", homeNixPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("home-manager switch の実行に失敗: %s\n%s\n%s", err, stdout.String(), stderr.String())
	}

	fmt.Print(stdout.String())

	return nil
}

// ApplyHomeManagerWithFlake はFlake環境でhome-manager switchを実行
func (c *Client) ApplyHomeManagerWithFlake(flakePath, configName string) error {
	flakeRef := fmt.Sprintf("%s#%s", flakePath, configName)
	cmd := exec.Command("home-manager", "switch", "--flake", flakeRef)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("home-manager switch の実行に失敗: %s\n%s\n%s", err, stdout.String(), stderr.String())
	}

	fmt.Print(stdout.String())

	return nil
}

func (c *Client) GetPackageVersion(packageName string) (string, error) {
	cmd := exec.Command("nix", "eval", "nixpkgs#"+packageName+".version", "--raw")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "unknown", nil
	}

	version := strings.TrimSpace(stdout.String())
	if version == "" {
		return "unknown", nil
	}

	return version, nil
}

type SearchResult struct {
	Name        string
	Description string
	Version     string
}
