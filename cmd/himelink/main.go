package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/akikareha/himelink/internal/config"
	"github.com/akikareha/himelink/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		print("Usage: " + os.Args[0] + " himelink.yaml\n")
		return
	}
	cfg := config.Load(os.Args[1])

	r := chi.NewRouter()
	server.RegisterRoutes(cfg, r)

	http.ListenAndServe(cfg.App.Addr, r)
}
