package util

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/blang/semver"
	"github.com/google/go-github/github"
	"github.com/inconshreveable/go-update"
)

const (
	platform = runtime.GOOS + "-" + runtime.GOARCH
)

var (
	ErrorNoBinary = errors.New("no binary for the update found")
)

type Updater struct {
	CurrentVersion     string // Currently running version.
	GithubOwner        string // The owner of the repo like "pcdummy"
	GithubRepo         string // The repository like "go-githubupdate"
	latestReleasesResp *github.RepositoryRelease
}

func (u *Updater) CheckUpdateAvailable() (string, error) {
	client := github.NewClient(nil)

	ctx := context.Background()
	release, _, err := client.Repositories.GetLatestRelease(ctx, u.GithubOwner, u.GithubRepo)
	if err != nil {
		return "", err
	}

	u.latestReleasesResp = release

	var updateVersion string
	fmt.Sscanf(*u.latestReleasesResp.TagName, "v%s", &updateVersion)
	current, err := semver.Make(u.CurrentVersion)
	update, err := semver.Make(updateVersion)
	if current.LT(update) {
		return *u.latestReleasesResp.TagName, nil
	}

	return "", nil
}

func (u *Updater) Update() error {
	reqFilename := u.GithubRepo + "-" + platform
	var foundAsset github.ReleaseAsset
	for _, asset := range u.latestReleasesResp.Assets {
		if *asset.Name == reqFilename {
			foundAsset = asset
			break
		}
	}

	if foundAsset.Name == nil {
		return ErrorNoBinary
	}

	dlURL := *foundAsset.BrowserDownloadURL

	resp, err := http.Get(dlURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		return err
	}
	return nil
}
