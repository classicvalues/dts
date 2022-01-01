package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v40/github"
)

type NPMRegistryPackage struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Version       string        `json:"version"`
	License       string        `json:"license"`
	Pika          bool          `json:"pika"`
	SideEffects   bool          `json:"sideEffects"`
	Keywords      []string      `json:"keywords"`
	Repository    Repository    `json:"repository"`
	Source        string        `json:"source"`
	Types         string        `json:"types"`
	Main          string        `json:"main"`
	Module        string        `json:"module"`
	GitHead       string        `json:"gitHead"`
	Homepage      string        `json:"homepage"`
	ID            string        `json:"_id"`
	NodeVersion   string        `json:"_nodeVersion"`
	NpmVersion    string        `json:"_npmVersion"`
	Dist          Dist          `json:"dist"`
	NpmUser       NpmUser       `json:"_npmUser"`
	Maintainers   []Maintainers `json:"maintainers"`
	HasShrinkwrap bool          `json:"_hasShrinkwrap"`
}
type Repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Dist struct {
	Integrity    string `json:"integrity"`
	Shasum       string `json:"shasum"`
	Tarball      string `json:"tarball"`
	FileCount    int    `json:"fileCount"`
	UnpackedSize int    `json:"unpackedSize"`
}
type NpmUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Maintainers struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetPackageGitHubRepo(packageName, version string) (string, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/%s", packageName, version)

	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	npmpkg := NPMRegistryPackage{}
	err = json.Unmarshal(body, &npmpkg)
	if err != nil {
		return "", err
	}

	if npmpkg.Repository.Type == "git" && strings.Contains(npmpkg.Repository.URL, "github.com") {
		/*PAT := os.Getenv("PAT")
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: PAT},
		)
		tc := oauth2.NewClient(ctx, ts)
		*/
		ctx := context.Background()
		client := github.NewClient(http.DefaultClient)
		repoParts := strings.Split(npmpkg.Repository.URL, "/")
		owner := repoParts[3]
		repo := repoParts[4]
		repo = strings.TrimRight(repo, ".git")
		repository, _, err := client.Repositories.Get(ctx, owner, repo)
		if err != nil {
			return "", err
		}

		return *repository.URL, nil
	}

	return "", nil
}
