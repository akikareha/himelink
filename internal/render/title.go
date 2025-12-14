package render

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
)

type Heading struct {
	Level int
	Text  string
}

func extractHeadings(src []byte, doc ast.Node) (headings []Heading) {
	// scan AST
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		heading, ok := n.(*ast.Heading)
		if !ok || !entering {
			return ast.WalkContinue, nil
		}

		var buf bytes.Buffer
		for c := heading.FirstChild(); c != nil; c = c.NextSibling() {
			buf.Write(c.Text(src))
		}

		headings = append(headings, Heading{
			Level: heading.Level,
			Text:  buf.String(),
		})

		return ast.WalkContinue, nil
	})

	return headings
}

// return top heading
func extractTitle(filename string, src []byte, doc ast.Node) string {
	headings := extractHeadings(src, doc)
	if len(headings) == 0 {
		return filename
	}

	// find min heading level
	best := headings[0]
	for _, h := range headings {
		if h.Level < best.Level {
			best = h
		}
	}
	return best.Text
}
