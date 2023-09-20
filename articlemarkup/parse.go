package articlemarkup

import (
	"bytes"
	"io"
	"log"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/yaml.v3"

	"github.com/alecthomas/chroma/v2"
	chromaHtml "github.com/alecthomas/chroma/v2/formatters/html"
	chromaStyles "github.com/alecthomas/chroma/v2/styles"
)

var g struct {
	parser          *parser.Parser
	renderer        *html.Renderer
	chromaFormatter *chromaHtml.Formatter
	chromaStyle     *chroma.Style
}

func init() {
	g.parser = parser.New()
	g.parser.Opts.ParserHook = parserHook
	prev := g.parser.RegisterInline('!', nil)
	g.parser.RegisterInline('[', inlineImage(g.parser, prev))

	renderOpts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: renderHook,
	}
	g.renderer = html.NewRenderer(renderOpts)

	g.chromaFormatter = chromaHtml.New(chromaHtml.WithClasses(true), chromaHtml.TabWidth(4))
	g.chromaStyle = chromaStyles.Get("catppuccin-frappe")
	if g.chromaFormatter == nil || g.chromaStyle == nil {
		log.Println("Failed to initialize codeblock highlighting")
	}
}

/*
Unmarshals article data into passed pointer and
renders document content as html
*/
func ParseArticle(md []byte, articleData interface{}) []byte {
	metaOpener := []byte("```yaml")
	metaCloser := []byte("```")
	if !bytes.HasPrefix(md, metaOpener) {
		log.Println("Article missing metadata opener")
		return nil
	}
	md = md[len(metaOpener):]
	metaEnd := bytes.Index(md, metaCloser)
	if metaEnd < 0 {
		return nil
	}
	metadata := md[:metaEnd-1]
	metadata = bytes.Trim(metadata, "\n\t ")
	md = md[metaEnd+len(metaCloser):]
	md = bytes.Trim(md, "\n\t ")
	if err := yaml.Unmarshal(metadata, articleData); err != nil {
		return nil
	}
	html := markdown.ToHTML(md, g.parser, g.renderer)
	return html
}

func parserHook(data []byte) (ast.Node, []byte, int) {
	// Images
	if node, d, n := parseImageFigure(data); node != nil {
		return node, d, n
	}
	return nil, nil, 0
}

func renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	// Code blocks
	if code, ok := node.(*ast.CodeBlock); ok {
		renderCode(w, code, entering)
		return ast.GoToNext, true
	}
	// Images
	if imageFig, ok := node.(*imageFigure); ok {
		renderImageFigure(w, imageFig, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}
