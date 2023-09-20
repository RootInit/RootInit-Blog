```yaml
Article_Title: Making a Custom Golang Markup Language For Blog Articles
URL_Title: Golang_Custom_Markup_Language_Articles
Description: Creating a custom human-readable markup language for typesetting of blog articles in Golang
Publish_Date: 2023-09-19
Thumbnail: image.png
Category: Programming
Tags:
    - Golang
    - Parsing
    - Regex
```

While working on a recent "scratch" built static site generator for this blog, I have experienced a great deal of indecision on the method by which blog articles would be typeset. I initially planned on simply writing the HTML for each article manually, however with desired additions such as alignable images and figures with captions, tables, quotes, and codeblocks, all having their own CSS classes this was clearly not going to work. Writing the articles in Markdown seemed like an obvious solution which would remove any need for such an interface.

example |  table
--------|--------
row 1   |  data
row 2   |  data

```go
func init() {
    htmlFormatter = html.New(html.WithClasses(true), html.TabWidth(2))
    if htmlFormatter == nil {
        panic("couldn't create html formatter")
    }
    styleName := "monokailight"
    highlightStyle = styles.Get(styleName)
    if highlightStyle == nil {
        panic(fmt.Sprintf("didn't find style '%s'", styleName))
    }
}

// based on https://github.com/alecthomas/chroma/blob/master/quick/quick.go
func htmlHighlight(w io.Writer, source, lang, defaultLang string) error {
    if lang == "" {
        lang = defaultLang
    }
    l := lexers.Get(lang)
    if l == nil {
        l = lexers.Analyse(source)
    }
    if l == nil {
        l = lexers.Fallback
    }
    l = chroma.Coalesce(l)

    it, err := l.Tokenise(nil, source)
    if err != nil {
        return err
    }
    return htmlFormatter.Format(w, highlightStyle, it)
}
```

This is an exampe of a simple block level image embed which follows standard markdown rules.

![Alt text](../source/assets/img/seacable_feat.webp)


And this is an example of an inline image right aligned and set to 50% page size. A caption can also be set ![Alt text](http://lenna.org/len_std.jpg)[left][50]{Figure 2. Lenna}

```md ![Alt text](http://lenna.org/len_std.jpg)[left][50]{caption text here}```