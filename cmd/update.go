package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/blang/semver"
	"github.com/google/go-github/github"
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
)

const (
	platform = runtime.GOOS + "-" + runtime.GOARCH
)

var (
	ErrorNoBinary = errors.New("no binary for the update found")
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for update",
	Run:   startUpdate,
}

func init() {
	RootCmd.AddCommand(updateCmd)
}

func startUpdate(cmd *cobra.Command, args []string) {
	var version string
	fmt.Sscanf(Version, "v%s", &version)
	u := &Updater{
		CurrentVersion: version,
		GithubOwner:    "jckimble",
		GithubRepo:     "lighttower",
	}
	available, err := u.CheckUpdateAvailable()
	if err != nil {
		log.Printf("Unable to check Update: %s\n", err)
	}

	if available != "" {
		log.Printf("Version %s available\n", available)
		err := u.Update()
		if err != nil {
			log.Printf("Unable to Update: %s\n", err)
		}
	}
}

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
