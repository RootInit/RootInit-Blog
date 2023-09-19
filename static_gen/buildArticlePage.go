package static_gen

import (
	"blog/models"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path"
)

type articlePageData struct {
	Title    string
	Date     string
	Tags     []tagData
	Body     template.HTML //HTML page body
	Comments []commentTmplData
}

/*  */
func (sG StaticGen) getArticlePageData(a *models.Article) articlePageData {
	article := articlePageData{
		Title: a.Title,
		Date: a.Date.UTC().Format(
			sG.config.OutputOpts.DateFormat,
		),
		Tags: sG.convertTags(&a.Tags),
		Body: template.HTML(a.Body),
	}
	for _, c := range a.Comments {
		article.Comments = append(article.Comments, sG.convertComment(&c))
	}
	return article
}

func (sG StaticGen) writeArticlePage(article *models.Article) error {
	urlPath := sG.getAbsUrl(sG.getArticleRelPath(article.Id, article.UrlTitle))
	linkNext, linkPrev := sG.getArticleNextPrevLink(article)
	pageHead := headTmplData{
		Title:        sG.config.SiteName + "-" + article.Title,
		Description:  article.Description,
		CanonicalUrl: urlPath,
		LinkPrev:     linkPrev,
		LinkNext:     linkNext,
	}
	articlePageData := sG.getArticlePageData(article)
	pageData := sG.makePageData(pageHead, articlePageData)
	page := new(bytes.Buffer)
	if err := sG.tmplArticle.Execute(page, pageData); err != nil {
		return err
	}
	filePath := sG.getAbsPath(sG.getArticleRelPath(article.Id, article.UrlTitle))
	return writeFile(filePath, page.Bytes())
}

func (sG StaticGen) getArticleNextPrevLink(article *models.Article) (string, string) {
	var prevArticleId, nextArticleId int
	var prevArticleLink, nextArticleLink string
	idList := &sG.resources.articleFileCache.IdList
	idMap := &sG.resources.articleFileCache.IdMap
	lastIdx := len(*idList) - 1
	for i, a := range *idList {
		if a != article.Id {
			continue
		}
		if i > 0 {
			prevArticleId = (*idList)[i-1]
		}
		if i < lastIdx {
			nextArticleId = (*idList)[i+1]
		}
		break
	}
	if nextArticleId != 0 {
		urlTitle := (*idMap)[nextArticleId]
		path := sG.getArticleRelPath(nextArticleId, urlTitle)
		nextArticleLink = sG.getAbsUrl(path)
	}
	if prevArticleId != 0 {
		urlTitle := (*idMap)[prevArticleId]
		path := sG.getArticleRelPath(prevArticleId, urlTitle)
		prevArticleLink = sG.getAbsUrl(path)

	}
	return prevArticleLink, nextArticleLink
}
func (sG StaticGen) getArticleRelPath(articleId int, articleUrlTitle string) string {
	if articleId == 0 || articleUrlTitle == "" {
		log.Printf(`Ignored Error: getArticleRelPath got uninitialized paramater. articleId: "%d" articleUrlTitle "%s"`, articleId, articleUrlTitle)
		return "."
	}
	articlePath := path.Join(
		sG.config.OutputPaths.ArticleDir,
		fmt.Sprint(articleId),
		articleUrlTitle+`.html`,
	)
	return articlePath
}
