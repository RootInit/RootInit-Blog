package static_gen

import (
	"blog/models"
	"bytes"
	"fmt"
	"log"
	"path"
)

func (sG StaticGen) buildTagIndex(tag models.Tag, cards []cardData) error {
	pageNum := 1
	lastIdx := len(cards) - 1
	var pageCards []cardData
	for i, card := range cards {
		pageCards = append(pageCards, card)
		// Check if out of cards or at max cards per page
		if i == lastIdx || len(pageCards) == sG.config.IndexSize {
			var lastPage bool
			if i == lastIdx {
				lastPage = true
			}
			err := sG.writeTagIndexPage(&tag, pageNum, pageCards, lastPage)
			if err != nil {
				return err
			}
			pageCards = make([]cardData, 0, sG.config.IndexSize)
			pageNum += 1
		}
	}
	return nil
}

func (sG StaticGen) buildTagIndexes() error {
	for _, tag := range sG.resources.tagCache.List {
		cards := sG.cards.getCardsFromTagId(tag.Id)
		if err := sG.buildTagIndex(tag, cards); err != nil {
			return err
		}
	}
	return nil
}
func (sG StaticGen) writeTagIndexPage(tag *models.Tag, pageNum int, cards []cardData, lastPage bool) error {
	var linkPrev, linkNext string
	if pageNum > 1 {
		linkPrev = sG.getAbsUrl(sG.getTagIndexPagePath(tag, pageNum-1))
	}
	if !lastPage {
		linkNext = sG.getAbsUrl(sG.getTagIndexPagePath(tag, pageNum+1))
	}
	urlPath := sG.getAbsUrl(sG.getTagIndexPagePath(tag, pageNum))
	pageHead := headTmplData{
		Title:        sG.config.SiteName + "-Tag Index",
		Description:  sG.config.IndexDesc,
		CanonicalUrl: urlPath,
		LinkPrev:     linkPrev,
		LinkNext:     linkNext,
	}
	indexPageData := indexPageData{
		Title: tag.Name + " Articles:",
		Cards: cards,
	}
	pageData := sG.makePageData(pageHead, indexPageData)
	page := new(bytes.Buffer)
	if err := sG.tmplIndex.Execute(page, pageData); err != nil {
		return err
	}
	filePath := sG.getAbsPath(sG.getTagIndexPagePath(tag, pageNum))
	return writeFile(filePath, page.Bytes())
}

func (sG StaticGen) getTagIndexPagePath(tag *models.Tag, pageNum int) string {
	var pagePath string
	if pageNum <= 0 {
		log.Fatalf(`getIndexRelPath got uninitialized paramater. pageNum: "%d"`, pageNum)
	} else if pageNum == 1 {
		// First page is root domain.tld/tags/tag_name/index.html
		pagePath = path.Join(
			sG.config.OutputPaths.TagIndexDir,
			urlSafeName(tag.Name),
			`index.html`,
		)
	} else {
		// Other pages follow scheme domain.tld/tags/tag_name/page/2.html
		pagePath = path.Join(
			sG.config.OutputPaths.TagIndexDir,
			urlSafeName(tag.Name),
			sG.config.OutputPaths.IndexPageDir,
			fmt.Sprintf(`%d.html`, pageNum),
		)
	}
	return pagePath
}
