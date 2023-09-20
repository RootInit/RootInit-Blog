package articlemarkup

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// ![Alt](Src)[Align][Width]{Caption}
var imageRegex = regexp.MustCompile(
	`!\[(.*?)]\((.+?)\)(?:\[\s?(l|left|r|right|c|center|)\s?\])?(?:\[\s?([0-9]+)?\s?\])?(?:\{\s?(.*?)\s?\})?`,
)

type imageFigure struct {
	ast.Leaf
	Alt     string
	Src     string
	Align   string
	Width   int
	Caption string
}

/*
Parses an imageFigure embed.
*/
func parseImageFigure(data []byte) (ast.Node, []byte, int) {
	if !bytes.HasPrefix(data, []byte("![")) {
		return nil, nil, 0
	}
	matches := imageRegex.FindSubmatch(data)
	if matches == nil {
		return nil, nil, 0
	}
	end := len(matches[0])
	data = data[:end]
	width, _ := strconv.Atoi(string(matches[4]))
	res := &imageFigure{
		Leaf:    ast.Leaf{Literal: data},
		Alt:     string(matches[1]),
		Src:     string(matches[2]),
		Align:   string(matches[3]),
		Width:   width,
		Caption: string(matches[5]),
	}
	return res, nil, end
}

/*
Parses an inline imageFigure.

Wrapper function is used to allow calling the previous definition if
there is not an imageFigure at the current token
*/
func inlineImage(p *parser.Parser, fn func(p *parser.Parser, data []byte, offset int) (int, ast.Node)) func(p *parser.Parser, data []byte, offset int) (int, ast.Node) {
	return func(p *parser.Parser, original []byte, offset int) (int, ast.Node) {
		if offset < 1 {
			return fn(p, original, offset)
		}
		data := original[offset-1:]
		var end int
		node, _, end := parseImageFigure(data)
		if node == nil {
			return fn(p, original, offset)
		}
		return end, node
	}
}

/*
Renders an imageFigure in HTML
*/
func renderImageFigure(w io.Writer, imgFig *imageFigure, entering bool) {
	var style string
	switch imgFig.Align {
	case "r", "right":
		style += "float:right; margin-right:0;"
	case "l", "left":
		style += "float:left; margin-left:0;"
	case "c", "center":
		style += "margin-inline:auto;"
	}
	if imgFig.Width > 0 && imgFig.Width < 100 {
		style += "width: " + fmt.Sprint(imgFig.Width) + "%;"
	}
	figureHtml := fmt.Sprintf(`<figure><img alt="%s" src="%s" style="%s"></img><figcaption>%s</figcaption></figure>`, imgFig.Alt, imgFig.Src, style, imgFig.Caption)
	io.WriteString(w, figureHtml)
}
