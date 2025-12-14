package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/akikareha/himelink/internal/config"
	"github.com/akikareha/himelink/internal/render"
)

func isValid(name string) bool {
	return validName.MatchString(name)
}

type repoInfo struct {
	DefaultBranch string `json:"default_branch"`
	Description string `json:"description"`
	HtmlUrl string `json:"html_url"`
	Name string `json:"name"`
}

func fetchRepoInfo(baseURL, owner, repo string) (repoInfo, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", baseURL, owner, repo)

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

type readmeInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func fetchReadmeInfo(baseURL, owner, repo string) (readmeInfo, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/readme", baseURL, owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		return readmeInfo{}, err
	}
	defer resp.Body.Close()

	var info readmeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return readmeInfo{}, err
	}

	return info, nil
}

func fetchRaw(baseURL, owner, repo, branch, path string) ([]byte, error) {
	url := fmt.Sprintf(
		"%s/%s/%s/%s/%s",
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

func handleRepo(
	cfg *config.Config,
	route config.Route,
	w http.ResponseWriter,
	r *http.Request,
	slash bool,
) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if !isValid(owner) || !isValid(repo) {
		http.Error(w, "invalid repo name", 400)
		return
	}

	info, err := fetchRepoInfo(route.API, owner, repo)
	if err != nil {
		http.Error(w, "cannot get repo info", 500)
		return
	}

	readme, err := fetchReadmeInfo(route.API, owner, repo)
	if err != nil {
		http.Error(w, "cannot get readme info", 500)
		return
	}

	var readmePath string
	if slash {
		readmePath = "blob/" + readme.Path
	} else {
		readmePath = repo + "/blob/" + readme.Path
	}

	repoInfo := render.RepoInfo{
		Description: info.Description,
		Name: info.Name,
		ReadmeName: readme.Name,
		ReadmePath: readmePath,
		URL: info.HtmlUrl,
	}
	render.RenderRepo(cfg, w, repoInfo)
}

func RepoHandler(cfg *config.Config, route config.Route, slash bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleRepo(cfg, route, w, r, slash)
	}
}

func handlePath(
	cfg *config.Config,
	route config.Route,
	w http.ResponseWriter,
	r *http.Request,
) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	mode := chi.URLParam(r, "mode")
	path := chi.URLParam(r, "*")

	if !isValid(owner) || !isValid(repo) {
		http.Error(w, "invalid repo name", 400)
		return
	}
	if mode != "blob" {
		http.Error(w, "invalid mode", 400)
		return
	}

	if strings.Contains(path, "..") {
		http.Error(w, "invalid path", 500)
		return
	}
	if strings.HasPrefix(path, "/") {
		http.Error(w, "invalid path", 500)
		return
	}

	info, err := fetchRepoInfo(route.API, owner, repo)
	if err != nil {
		http.Error(w, "cannot get repo info", 500)
		return
	}
	branch := info.DefaultBranch
	if branch == "" {
		http.Error(w, "invalid branch", 500)
		return
	}

	raw, err := fetchRaw(
		route.Raw,
		owner,
		repo,
		branch,
		path,
	)
	if err != nil {
		http.Error(w, "cannot fetch", 500)
		return
	}

	ext := filepath.Ext(path)
	if ext == ".md" {
		render.RenderMarkdown(cfg, w, raw)
	} else {
		http.Error(w, "unsupported extension " + ext, 500)
		return
	}
}

func PathHandler(cfg *config.Config, route config.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlePath(cfg, route, w, r)
	}
}
