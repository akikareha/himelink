package render

import (
	"net/http"

	"github.com/akikareha/himelink/internal/config"
	"github.com/akikareha/himelink/internal/templates"
)

type RepoInfo struct {
	Description string
	Name        string
	ReadmeName  string
	ReadmePath  string
	URL         string
}

func RenderRepo(cfg *config.Config, w http.ResponseWriter, info RepoInfo) {
	tmpl := templates.New("repo.html")

	tmpl.Execute(w, struct {
		SiteName string
		Title    string
		Info     RepoInfo
	}{
		SiteName: cfg.Site.Name,
		Title:    info.Name,
		Info:     info,
	})
}
