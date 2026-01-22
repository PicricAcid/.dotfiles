package nixfile

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Manager struct {
	filePath string
}

func NewManager(filePath string) *Manager {
	return &Manager{
		filePath: filePath,
	}
}

func (m *Manager) ListPackages() ([]string, error) {
	content, err := os.ReadFile(m.filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗: %w", err)
	}

	packages := m.parsePackages(string(content))
	return packages, nil
}

func (m *Manager) AddPackage(packageName string) error {
	if err := m.backup(); err != nil {
		return fmt.Errorf("バックアップの作成に失敗: %w", err)
	}

	content, err := os.ReadFile(m.filePath)
	if err != nil {
		return fmt.Errorf("ファイルの読み込みに失敗: %w", err)
	}

	contentStr := string(content)
	packages := m.parsePackages(contentStr)

	for _, pkg := range packages {
		if pkg == packageName {
			return fmt.Errorf("パッケージ '%s' は既にインストールされています", packageName)
		}
	}

	packages = append(packages, packageName)
	sort.Strings(packages)

	newContent := m.generateContent(packages)

	if err := os.WriteFile(m.filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗: %w", err)
	}

	return nil
}

func (m *Manager) RemovePackage(packageName string) error {
	if err := m.backup(); err != nil {
		return fmt.Errorf("バックアップの作成に失敗: %w", err)
	}

	content, err := os.ReadFile(m.filePath)
	if err != nil {
		return fmt.Errorf("ファイルの読み込みに失敗: %w", err)
	}

	contentStr := string(content)
	packages := m.parsePackages(contentStr)

	found := false
	newPackages := make([]string, 0, len(packages))
	for _, pkg := range packages {
		if pkg == packageName {
			found = true
			continue
		}
		newPackages = append(newPackages, pkg)
	}

	if !found {
		return fmt.Errorf("パッケージ '%s' は見つかりませんでした", packageName)
	}

	newContent := m.generateContent(newPackages)

	if err := os.WriteFile(m.filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗: %w", err)
	}

	return nil
}

func (m *Manager) HasPackage(packageName string) (bool, error) {
	packages, err := m.ListPackages()
	if err != nil {
		return false, err
	}

	for _, pkg := range packages {
		if pkg == packageName {
			return true, nil
		}
	}

	return false, nil
}

func (m *Manager) GetDiff(packageName string, isAdd bool) (string, error) {
	packages, err := m.ListPackages()
	if err != nil {
		return "", err
	}

	var before, after []string

	if isAdd {
		before = packages
		after = append(packages, packageName)
		sort.Strings(after)
	} else {
		before = packages
		after = make([]string, 0, len(packages))
		for _, pkg := range packages {
			if pkg != packageName {
				after = append(after, pkg)
			}
		}
	}

	diff := " home.packages = with pkgs; [\n"

	for _, pkg := range before {
		found := false
		for _, p := range after {
			if p == pkg {
				found = true
				break
			}
		}
		if !found {
			diff += fmt.Sprintf("-	%s\n", pkg)
		}
	}

	for _, pkg := range after {
		found := false
		for _, p := range before {
			if p == pkg {
				found = true
				break
			}
		}
		if found {
			diff += fmt.Sprintf("	%s\n", pkg)
		} else {
			diff += fmt.Sprintf("+	%s\n", pkg)
		}
	}

	diff += " ];"

	return diff, nil
}

func (m *Manager) Rollback() error {
	backupPath := m.filePath + ".bak"

	content, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("バックアップファイルの読み込みに失敗: %w", err)
	}

	if err := os.WriteFile(m.filePath, content, 0644); err != nil {
		return fmt.Errorf("ロールバックに失敗: %w", err)
	}

	return nil
}

func (m *Manager) parsePackages(content string) []string {
	re := regexp.MustCompile(`home\.packages\s*=\s*with\s+pkgs;\s*\[\s*([\s\S]*?)\s*\];`)
	matches := re.FindStringSubmatch(content)

	if len(matches) < 2 {
		return []string{}
	}

	packagesBlock := matches[1]

	lines := strings.Split(packagesBlock, "\n")
	packages := make([]string, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		if line != "" {
			packages = append(packages, line)
		}
	}

	return packages
}

func (m *Manager) generateContent(packages []string) string {
	var builder strings.Builder

	builder.WriteString("{ pkgs, ... }: {\n")
	builder.WriteString("	home.packages = with pkgs; [\n")

	if len(packages) == 0 {
		builder.WriteString("	# focus でインストールしたパッケージ\n")
	} else {
		for _, pkg := range packages {
			builder.WriteString(fmt.Sprintf("	%s\n", pkg))
		}
	}

	builder.WriteString("	];\n")
	builder.WriteString("}\n")

	return builder.String()
}

func (m *Manager) backup() error {
	content, err := os.ReadFile(m.filePath)
	if err != nil {
		return err
	}

	backupPath := m.filePath + ".bak"
	return os.WriteFile(backupPath, content, 0644)
}
