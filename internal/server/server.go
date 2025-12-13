package server

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"

	"github.com/akikareha/himelink/internal/config"
	"github.com/akikareha/himelink/internal/fetch"
	"github.com/akikareha/himelink/internal/render"
)

func renderMarkdown(w http.ResponseWriter, raw []byte) {
	htmlBytes, err := render.MarkdownToHTML(raw)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
		<!doctype html>
		<html>
		<head>
		  <meta charset="utf-8" />
		</head>
		<body>
	`))

	w.Write(htmlBytes)

	w.Write([]byte(`
		</body>
		</html>
	`))
}

var validName = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func isValid(name string) bool {
    return validName.MatchString(name)
}

func handleGitea(
	cfg *config.Config,
	w http.ResponseWriter,
	r *http.Request,
) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if !isValid(owner) || !isValid(repo) {
		http.Error(w, "invalid repo name", 400)
		return
	}

	info, err := fetch.GiteaFetchRepoInfo(cfg.Gitea.ApiBase, owner, repo)
	if err != nil {
		http.Error(w, "cannot get repo info: "+err.Error(), 500)
		return
	}

	branch := info.DefaultBranch
	if branch == "" {
		branch = "main"
	}

	raw, err := fetch.GiteaFetchRaw(
		cfg.Gitea.ApiBase,
		owner,
		repo,
		branch,
		"README.md",
	)
	if err != nil {
		http.Error(w, "cannot fetch README.md: "+err.Error(), 500)
		return
	}

	renderMarkdown(w, raw)
}

func giteaHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleGitea(cfg, w, r)
	}
}

func handleGitHub(
	cfg *config.Config,
	w http.ResponseWriter,
	r *http.Request,
) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if !isValid(owner) || !isValid(repo) {
		http.Error(w, "invalid repo name", 400)
		return
	}

	info, err := fetch.GitHubFetchRepoInfo(cfg.GitHub.ApiBase, owner, repo)
	if err != nil {
		http.Error(w, "cannot get repo info: "+err.Error(), 500)
		return
	}

	branch := info.DefaultBranch
	if branch == "" {
		branch = "main"
	}

	raw, err := fetch.GitHubFetchRaw(
		cfg.GitHub.RawBase,
		owner,
		repo,
		branch,
		"README.md",
	)
	if err != nil {
		http.Error(w, "cannot fetch README.md: "+err.Error(), 500)
		return
	}

	renderMarkdown(w, raw)
}

func gitHubHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleGitHub(cfg, w, r)
	}
}

func RegisterRoutes(cfg *config.Config, r chi.Router) {
	r.Get("/gitea/{owner}/{repo}", giteaHandler(cfg))
	r.Get("/github/{owner}/{repo}", gitHubHandler(cfg))

	fileServer := http.FileServer(http.Dir(cfg.Site.Static))
	r.Handle("/*", fileServer)
}
