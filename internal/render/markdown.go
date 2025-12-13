package render

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM, // GitHub-like
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(), // allow embedding HTML
	),
)

func MarkdownToHTML(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := md.Convert(src, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
