package static_gen

import (
	"blog/models"
)

type headTmplData struct {
	Title        string
	Description  string
	CanonicalUrl string
	LinkPrev     string
	LinkNext     string
}

type sidebarTmplData struct {
	ArticleTags []tagData
}

func (sG StaticGen) makeSidebarData() sidebarTmplData {
	sbData := sidebarTmplData{
		ArticleTags: sG.convertTags(&sG.resources.tagCache.List),
	}
	// TODO Future additions...
	return sbData
}

type indexPageData struct {
	Title string
	Cards []cardData
}

type pageTmplData struct {
	Head    headTmplData
	Sidebar sidebarTmplData
	Primary interface{}
	Strings map[string]string
}

/*
Makes a pageTmplData from a `pageHeadData` and a `bodyData` struct
`bodyData` can be any of: articleTmplData, indexTmplData, contactTmplData
*/
func (sG StaticGen) makePageData(pageHeadData headTmplData, bodyData interface{}) pageTmplData {
	pageData := pageTmplData{
		Head:    pageHeadData,
		Sidebar: sG.makeSidebarData(),
		Primary: bodyData,
		Strings: sG.config.SubStrings,
	}
	return pageData
}

type commentTmplData struct {
	Username  string
	Date      string
	Body      string
	ReplyLink string
	Replies   []commentTmplData
}

/* Converts a models.Comment into a commentTmplData recursively */
func (sG StaticGen) convertComment(c *models.Comment) commentTmplData {
	comment := commentTmplData{
		Username: c.Author.Username,
		Date:     c.Date.Format(sG.config.OutputOpts.DateFormat),
		Body:     c.Body,
		// ReplyLink: "TODO",
	}
	for _, reply := range c.Replies {
		comment.Replies = append(comment.Replies, sG.convertComment(&reply))
	}
	return comment
}

type tagData struct {
	Name string
	Icon string
	Link string
}

/* Converts a models.Tag into a tagTmplData */
func (sG StaticGen) convertTag(t *models.Tag) tagData {
	comment := tagData{
		Name: t.Name,
		Icon: t.Icon,
		Link: sG.getTagIndexPagePath(t, 1),
	}
	return comment
}

/* Converts []models.Tag into a []tagTmplData */
func (sG StaticGen) convertTags(tags *[]models.Tag) []tagData {
	tagTmpls := make([]tagData, len(*tags))
	for i, tag := range *tags {
		tagTmpls[i] = sG.convertTag(&tag)
	}
	return tagTmpls
}

type cardData struct {
	Title       string
	Date        string
	Description string
	Tags        []tagData
	Thumbnail   string
	Link        string
}

/* Converts []models.Article into a []cardTmplData */
func (sG StaticGen) convertArticleCard(a *models.Article) cardData {
	cardData := cardData{
		Title: a.Title,
		Date: a.Date.UTC().Format(
			sG.config.OutputOpts.DateFormat,
		),
		Description: a.Description,
		Tags:        sG.convertTags(&a.Tags),
		Thumbnail:   a.Thumbnail,
		Link:        sG.getAbsUrl(sG.getArticleRelPath(a.Id, a.UrlTitle)),
	}
	return cardData
}
