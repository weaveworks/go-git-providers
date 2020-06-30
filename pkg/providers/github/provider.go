package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v31/github"
	giturls "github.com/whilp/git-urls"
	"golang.org/x/oauth2"

	"github.com/weaveworks/go-git-providers/pkg/key"
)

const (
	EnvVarGitHubToken = "GITHUB_TOKEN"
)

var (
	sshFull = regexp.MustCompile(`ssh://git@github.com/([^/]+)/([^.]+).git`)
	sshShort = regexp.MustCompile(`git@github.com:([^/]+)/([^.]+).git`)

	patterns = []*regexp.Regexp{
		sshFull,
		sshShort,
	}
)

// GitHubProvider accesses the Github API
type GitHubProvider struct {
	owner, repo string
	githubToken string
}

func NewGitHubProvider(repoURL string) (*GitHubProvider, error) {
	githubToken := os.Getenv(EnvVarGitHubToken)
	if githubToken == "" {
		return nil,fmt.Errorf("%s is not set. Cannot authenticate to github.com", EnvVarGitHubToken)
	}

	repo, err := repoName(repoURL)
	if err != nil {
		return nil, err
	}

	owner, err := repoOwner(repoURL)
	if err != nil {
		return nil, err
	}
	return &GitHubProvider{
			githubToken: githubToken,
			owner:       owner,
			repo:        repo,
	}, nil
}
func (p *GitHubProvider) ListKeys(ctx context.Context) ([]key.SSHKey, error) {
	gh := p.getGitHubAPIClient(ctx)

	keys, resp, err := gh.Repositories.ListKeys(ctx, p.owner, p.repo, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to list deploy keys from %s/%s", p.owner, p.repo)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to list deploy keys from %s/%s. Got response %s", p.owner, p.repo, resp.Status)
	}

	allKeys := make([]key.SSHKey, len(keys))
	for _, k := range keys {
		allKeys = append(allKeys, key.SSHKey{
			Title:    *k.Title,
			Key:      *k.Key,
			ReadOnly: *k.ReadOnly,
		})
	}
	return allKeys, nil
}

func (p *GitHubProvider) AuthorizeSSHKey(ctx context.Context, key key.SSHKey) error {
	gh := p.getGitHubAPIClient(ctx)

	_, resp, err := gh.Repositories.CreateKey(ctx, p.owner, p.repo, &github.Key{
		Key:      &key.Key,
		Title:    &key.Title,
		ReadOnly: &key.ReadOnly,
	})

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unable to authorize SSH Key %q. Got StatusCode %s", key.Title, resp.Status)
	}

	return nil
}

func (p *GitHubProvider) DeleteSSHKey(ctx context.Context, title string) error {
	gh := p.getGitHubAPIClient(ctx)

	keys, _, err := gh.Repositories.ListKeys(ctx, p.owner, p.repo, &github.ListOptions{})
	if err != nil {
		return err
	}

	var keyID int64

	for _, key := range keys {
		if key.GetTitle() == title {
			keyID = key.GetID()

			break
		}
	}

	if keyID == 0 {
		return nil
	}

	if _, err := gh.Repositories.DeleteKey(ctx, p.owner, p.repo, keyID); err != nil {
		return err
	}

	return nil
}

func (p *GitHubProvider) getGitHubAPIClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	gh := github.NewClient(tc)

	return gh
}

func repoName(repoURL string) (string, error) {
		u, err := giturls.Parse(repoURL)
		if err != nil {
			return "", fmt.Errorf("unable to parse git URL '%s': %s", repoURL, err.Error())
		}
		parts := strings.Split(u.Path, "/")
		if len(parts) == 0 {
			return "", fmt.Errorf("could not find name of repository %s", repoURL)
		}

		lastPathSegment := parts[len(parts)-1]
		return strings.TrimRight(lastPathSegment, ".git"), nil
}

func repoOwner(repoURL string) (string, error) {
		u, err := giturls.Parse(repoURL)
		if err != nil {
			return "", fmt.Errorf("unable to parse git URL '%s': %s", repoURL, err.Error())
		}
		parts := strings.Split(u.Path, "/")
		if len(parts) == 0 {
			return "", fmt.Errorf("could not find name of repository %s", repoURL)
		}

		firstPathSegment := parts[0]
		return firstPathSegment, nil

}