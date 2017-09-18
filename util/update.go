package util

import (
	"context"
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/google/go-github/github"
	"github.com/inconshreveable/go-update"
)

const (
	platform = runtime.GOOS + "-" + runtime.GOARCH
)

var (
	ErrorNoBinary    = errors.New("no binary for the update found")
	ErrorNoCheckSum  = errors.New("no checksum for the update found")
	ErrorNoSignature = errors.New("no signature for the update found")
)

type Updater struct {
	CurrentVersion     string // Currently running version.
	GithubOwner        string // The owner of the repo like "pcdummy"
	GithubRepo         string // The repository like "go-githubupdate"
	PublicKey          string
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
	var binaryAsset, checksumAsset, signatureAsset github.ReleaseAsset
	for _, asset := range u.latestReleasesResp.Assets {
		if *asset.Name == reqFilename {
			binaryAsset = asset
		} else if *asset.Name == reqFilename+".sha256" {
			checksumAsset = asset
		} else if *asset.Name == reqFilename+".sig" {
			signatureAsset = asset
		}
	}

	if binaryAsset.Name == nil {
		return ErrorNoBinary
	} else if checksumAsset.Name == nil {
		return ErrorNoCheckSum
	} else if signatureAsset.Name == nil {
		return ErrorNoSignature
	}

	dlURL := *binaryAsset.BrowserDownloadURL

	resp, err := http.Get(dlURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	checksum, err := parseChecksum(*checksumAsset.BrowserDownloadURL)
	if err != nil {
		return err
	}
	opts := update.Options{
		Checksum: checksum,
		Hash:     crypto.SHA256,
	}
	if u.PublicKey != "" {
		err = opts.SetPublicKeyPEM([]byte(u.PublicKey))
		if err != nil {
			return err
		}
		opts.Signature, err = getSignature(*signatureAsset.BrowserDownloadURL)
		if err != nil {
			return err
		}
		opts.Verifier = update.NewRSAVerifier()
	}
	err = update.Apply(resp.Body, opts)
	if err != nil {
		return err
	}
	return nil
}

func getSignature(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func parseChecksum(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	sha256, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.Trim(string(sha256), " \n"))
}
