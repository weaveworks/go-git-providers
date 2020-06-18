package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/weaveworks/go-git-providers/pkg/key"
	"github.com/weaveworks/go-git-providers/pkg/providers/github"
)

type Provider interface {
	AuthorizeSSHKey(ctx context.Context, key key.SSHKey) error
	DeleteSSHKey(ctx context.Context, keyId string) error
	ListKeys(ctx context.Context) ([]key.SSHKey, error)
}

// GetProvider returns the appropriate provider for the URL
func GetProvider(repoURL string) (Provider, error) {
	if strings.Contains(repoURL, "github.com") {
		return github.NewGitHubProvider(repoURL)
	}
	return nil, fmt.Errorf("provider not found for URL %q", repoURL)
}
