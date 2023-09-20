package articlemarkup

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	chromaStyles "github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown/ast"
)

func GetColormodeCss(lightStyle, darkStyle string) ([]byte, error) {
	var lightCss, darkCss, combinedCss []byte
	var lightDeclr, darkDeclr []byte
	var err error
	if lightCss, err = getCss(lightStyle); err != nil {
		return nil, err
	}
	lightDeclr, lightCss = extractCssColorVars(lightCss)
	lightDeclr = append([]byte("[colorMode=light] {\n"), lightDeclr...)
	lightDeclr = append(lightDeclr, []byte("}\n\n")...)
	if darkCss, err = getCss(darkStyle); err != nil {
		return nil, err
	}
	darkDeclr, _ = extractCssColorVars(darkCss)
	darkDeclr = append([]byte("[colorMode=dark] {\n"), darkDeclr...)
	darkDeclr = append(darkDeclr, []byte("}\n\n")...)
	if len(lightCss) != len(darkCss) {
		log.Println("Warning: Light and dark styles do not match.")
	}
	combinedCss = append(lightDeclr, darkDeclr...)
	combinedCss = append(combinedCss, lightCss...)
	return combinedCss, nil
}

func getCss(style string) ([]byte, error) {
	w := new(bytes.Buffer)
	err := g.chromaFormatter.WriteCSS(w, chromaStyles.Get(style))
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func extractCssColorVars(css []byte) ([]byte, []byte) {
	colorRegex := regexp.MustCompile(`(#[A-z0-9]{6})(?:;)?`)
	matches := colorRegex.FindAllSubmatch(css, -1)
	// Find all unique color hex codes
	var uniqueColors [][]byte
	var uniqueMap = make(map[string]bool)
	for _, match := range matches {
		color := match[1]
		strCol := string(color)
		if uniqueMap[strCol] {
			continue
		}
		uniqueColors = append(uniqueColors, color)
		uniqueMap[strCol] = true
	}
	// Make variables
	var declarations []byte
	for i, col := range uniqueColors {
		cVar := []byte(fmt.Sprintf("var(--codeblock-color-%d);", i))
		css = bytes.ReplaceAll(css, col, cVar)
		cDeclr := []byte(fmt.Sprintf("--codeblock-color-%d: %s;\n", i, col))
		declarations = append(declarations, cDeclr...)
	}
	return declarations, css
}
func highlightCodeBlock(w io.Writer, source, lang string) error {
	var l chroma.Lexer
	if l = lexers.Get(lang); l == nil {
		if l = lexers.Analyse(source); l == nil {
			l = lexers.Fallback
		}
	}
	l = chroma.Coalesce(l)
	ittr, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return g.chromaFormatter.Format(w, g.chromaStyle, ittr)
}

func renderCode(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	lang := string(codeBlock.Info)
	err := highlightCodeBlock(w, string(codeBlock.Literal), lang)
	if err != nil {
		log.Println("Error syntax highlighting codeblock: ", err)
	}
}
