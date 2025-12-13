package gitea

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"

	"github.com/akikareha/himelink/internal/config"
	"github.com/akikareha/himelink/internal/render"
)

type repoInfo struct {
	DefaultBranch string `json:"default_branch"`
}

func fetchRepoInfo(baseURL, owner, repo string) (repoInfo, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s", baseURL, owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		return repoInfo{}, err
	}
	defer resp.Body.Close()

	var info repoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return repoInfo{}, err
	}

	return info, nil
}

func fetchRaw(baseURL, owner, repo, branch, path string) ([]byte, error) {
	url := fmt.Sprintf(
		"%s/api/v1/repos/%s/%s/raw/%s/%s",
		baseURL, owner, repo, branch, path,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

var validName = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func isValid(name string) bool {
	return validName.MatchString(name)
}

func handle(
	cfg *config.Config,
	route config.Route,
	w http.ResponseWriter,
	r *http.Request,
) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if !isValid(owner) || !isValid(repo) {
		http.Error(w, "invalid repo name", 400)
		return
	}

	info, err := fetchRepoInfo(route.API, owner, repo)
	if err != nil {
		http.Error(w, "cannot get repo info: "+err.Error(), 500)
		return
	}

	branch := info.DefaultBranch
	if branch == "" {
		branch = "main"
	}

	raw, err := fetchRaw(
		route.API,
		owner,
		repo,
		branch,
		"README.md",
	)
	if err != nil {
		http.Error(w, "cannot fetch README.md: "+err.Error(), 500)
		return
	}

	render.RenderMarkdown(cfg, w, raw)
}

func Handler(cfg *config.Config, route config.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(cfg, route, w, r)
	}
}
