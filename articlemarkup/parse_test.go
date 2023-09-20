package articlemarkup

import (
	"fmt"
	"os"
	"testing"
)

type articleData struct {
	Title string `yaml:"Title"`
}

func TestParser(t *testing.T) {
	markdown, err := os.ReadFile("testArticle.md")
	if err != nil {
		t.Fatal(err)
	}
	var a articleData
	html := ParseArticle(markdown, &a)
	fmt.Println(string(html))
}

func TestGetCSS(t *testing.T) {
	// Dark Mode  catppuccin-frappe
	// Light Mode catppuccin-latte
	css, err := GetColormodeCss("catppuccin-latte", "catppuccin-frappe")
	cssStr := string(css)
	fmt.Println(cssStr, err)
}
