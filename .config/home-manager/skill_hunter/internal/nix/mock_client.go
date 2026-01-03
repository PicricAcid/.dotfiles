package nix

import "fmt"

// MockClient はテスト用のNixクライアント
type MockClient struct {
	// PackageExistsの戻り値を制御
	ShouldPackageExist bool
	// ApplyHomeManagerが失敗するかを制御
	ShouldApplyFail bool
	// GetPackageVersionの戻り値
	PackageVersions map[string]string
}

// NewMockClient は新しいモッククライアントを作成する
func NewMockClient() *MockClient {
	return &MockClient{
		ShouldPackageExist: true,
		ShouldApplyFail:    false,
		PackageVersions:    make(map[string]string),
	}
}

// Search はダミーの検索結果を返す
func (m *MockClient) Search(keyword string) ([]SearchResult, error) {
	return []SearchResult{
		{
			Name:        keyword,
			Description: fmt.Sprintf("Mock package: %s", keyword),
			Version:     "1.0.0",
		},
	}, nil
}

// PackageExists は設定に応じてパッケージの存在を返す
func (m *MockClient) PackageExists(packageName string) (bool, error) {
	return m.ShouldPackageExist, nil
}

// ApplyHomeManager は設定に応じて成功/失敗を返す（実際には何もしない）
func (m *MockClient) ApplyHomeManager(homeNixPath string) error {
	if m.ShouldApplyFail {
		return fmt.Errorf("mock: home-manager switch failed")
	}
	return nil
}

// GetPackageVersion は設定されたバージョンを返す
func (m *MockClient) GetPackageVersion(packageName string) (string, error) {
	if version, ok := m.PackageVersions[packageName]; ok {
		return version, nil
	}
	return "1.0.0", nil
}
