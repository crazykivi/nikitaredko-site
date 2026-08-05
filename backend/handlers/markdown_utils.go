package handlers

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark/ast"
)

func renderNodeText(n ast.Node, source []byte) string {
	var buf bytes.Buffer
	ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch v := node.(type) {
		case *ast.Text:
			buf.Write(v.Segment.Value(source))
			if v.SoftLineBreak() || v.HardLineBreak() {
				buf.WriteByte(' ')
			}
		case *ast.String:
			buf.Write(v.Value)
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(buf.String())
}

func renderNodeTextSkipAutoLinks(n ast.Node, source []byte) string {
	var buf bytes.Buffer
	ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch v := node.(type) {
		case *ast.Text:
			buf.Write(v.Segment.Value(source))
			if v.SoftLineBreak() || v.HardLineBreak() {
				buf.WriteByte(' ')
			}
		case *ast.String:
			buf.Write(v.Value)
		case *ast.AutoLink:
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(buf.String())
}

func findFirstURL(n ast.Node, source []byte) string {
	var found string
	ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || found != "" {
			return ast.WalkContinue, nil
		}
		switch v := node.(type) {
		case *ast.Link:
			found = string(v.Destination)
			return ast.WalkSkipChildren, nil
		case *ast.AutoLink:
			found = string(v.URL(source))
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return found
}

func getFirstParagraph(n ast.Node) ast.Node {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if _, ok := c.(*ast.Paragraph); ok {
			return c
		}
		if _, ok := c.(*ast.TextBlock); ok {
			return c
		}
	}
	return nil
}
