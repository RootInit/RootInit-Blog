package static_gen

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

func (sG StaticGen) getMinifier() *minify.M {
	m := minify.New()
	m.AddFunc("css", css.Minify)
	m.AddFunc("html", html.Minify)
	m.AddFunc("js", js.Minify)
	m.AddFunc("svg", svg.Minify)
	return m
}

/* Funtion to read, process, and output a CSS file */
func (sG StaticGen) processCss(m *minify.M, src, dest string) error {
	// Process and write css file
	fileBytes, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	// Left Merge CSS files imported with @import
	importRegex := regexp.MustCompile(`(?m)^@import\s+"(_[\w\-\\\/]+\.css)";$`)
	cssImports := importRegex.FindAllSubmatch(fileBytes, 64)
	for _, imprt := range cssImports {
		// imprt[0] == import statement imprt[1] == imported file
		impPath := path.Join(filepath.Dir(src), string(imprt[1]))
		importBytes, err := os.ReadFile(impPath)
		if err != nil {
			log.Println("Warning: Unable to resolve CSS import:", string(imprt[1]))
			continue
		}
		fileBytes = bytes.Replace(fileBytes, imprt[0], importBytes, 1)
	}
	// Minify CSS
	fileBytes, err = m.Bytes("css", fileBytes)
	if err != nil {
		return err
	}
	err = writeFile(dest, fileBytes)
	return err
}

/* Funtion to read, process, and output a HTML file */
func (sG StaticGen) processHtml(m *minify.M, src, dest string) error {
	fileBytes, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	// Inline SVG file <object> embeds
	assetDir := path.Join(
		sG.config.SourcePaths.Root,
		`assets`,
	)
	svgObgRegex := regexp.MustCompile(`(?m)<object\s+data="([\w\-\\\/]+\.svg)"><\/object>`)
	svgObjects := svgObgRegex.FindAllSubmatch(fileBytes, 64)
	for _, imprt := range svgObjects {
		// imprt[0] == import statement
		// imprt[1] == imported file
		impPath := path.Join(assetDir, string(imprt[1]))
		importBytes, err := os.ReadFile(impPath)
		if err != nil {
			log.Println("Warning: Unable to resolve SVG object embed:", string(imprt[1]))
			continue
		}
		importBytes, err = m.Bytes("svg", importBytes)
		if err != nil {
			return err
		}
		fileBytes = bytes.Replace(
			fileBytes, imprt[0], importBytes, 1,
		)
	}

	// Minify HTML
	fileBytes, err = m.Bytes("html", fileBytes)
	if err != nil {
		return err
	}
	err = writeFile(dest, fileBytes)
	return err
}

/* Funtion to read, process, and output a JS file */
func (sG StaticGen) processJs(m *minify.M, src, dest string) error {
	fileBytes, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	fileBytes, err = m.Bytes("js", fileBytes)
	if err != nil {
		return err
	}
	err = writeFile(dest, fileBytes)
	return err
}

/* Funtion to read, process, and output a SVG file */
func (sG StaticGen) processSvg(m *minify.M, src, dest string) error {
	fileBytes, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	fileBytes, err = m.Bytes("svg", fileBytes)
	if err != nil {
		return err
	}
	err = writeFile(dest, fileBytes)
	return err
}
