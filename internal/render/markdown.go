package render

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/akikareha/himelink/internal/config"
	"github.com/akikareha/himelink/internal/templates"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM, // GitHub-like
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(), // allow embedding HTML
	),
)

func markdownToHTML(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := md.Convert(src, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func RenderMarkdown(cfg *config.Config, w http.ResponseWriter, raw []byte) {
	title := ExtractTitle(raw)

	htmlBytes, err := markdownToHTML(raw)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl := templates.New("markdown.html")

	tmpl.Execute(w, struct {
		SiteName string
		Title    string
		Rendered template.HTML
	}{
		SiteName: cfg.Site.Name,
		Title:    title,
		Rendered: template.HTML(string(htmlBytes)),
	})
}
