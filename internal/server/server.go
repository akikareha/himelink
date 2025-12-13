package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/akikareha/himelink/internal/config"
	"github.com/akikareha/himelink/internal/gitea"
	"github.com/akikareha/himelink/internal/github"
)

func RegisterRoutes(cfg *config.Config, r chi.Router) {
	for _, route := range cfg.Routes {
		if route.Protocol == "gitea" {
			path := "/" + route.Path + "/{owner}/{repo}"
			r.Get(path, gitea.Handler(cfg, route))
			r.Get(path+"/", gitea.Handler(cfg, route))
			r.Get(path+"/{rest:*}", gitea.Handler(cfg, route))
		} else if route.Protocol == "github" {
			path := "/" + route.Path + "/{owner}/{repo}"
			r.Get(path, github.Handler(cfg, route))
			r.Get(path+"/", github.Handler(cfg, route))
			r.Get(path+"/{rest:*}", github.Handler(cfg, route))
		}
	}

	fileServer := http.FileServer(http.Dir(cfg.Site.Static))
	r.Handle("/*", fileServer)
}
