package static_gen

import (
	"blog/models"
	"bytes"
	"fmt"
	"log"
	"path"
)

func (sG StaticGen) buildCategoryIndex(category models.Category, cards []cardData) error {
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
			err := sG.writeCategoryIndexPage(&category, pageNum, pageCards, lastPage)
			if err != nil {
				return err
			}
			pageCards = make([]cardData, 0, sG.config.IndexSize)
			pageNum += 1
		}
	}
	return nil
}

func (sG StaticGen) buildCategoryIndexes() error {
	for _, category := range sG.resources.categoryCache.List {
		cards := sG.cards.getCardsFromCategoryId(category.Id)
		if err := sG.buildCategoryIndex(category, cards); err != nil {
			return err
		}
	}
	return nil
}

func (sG StaticGen) writeCategoryIndexPage(category *models.Category, pageNum int, cards []cardData, lastPage bool) error {
	var linkPrev, linkNext string
	if pageNum > 1 {
		linkPrev = sG.getAbsUrl(sG.getCategoryIndexPagePath(category, pageNum-1))
	}
	if !lastPage {
		linkNext = sG.getAbsUrl(sG.getCategoryIndexPagePath(category, pageNum+1))
	}
	urlPath := sG.getAbsUrl(sG.getCategoryIndexPagePath(category, pageNum))
	pageHead := headTmplData{
		Title:        sG.config.SiteName + "-Category Index",
		Description:  sG.config.IndexDesc,
		CanonicalUrl: urlPath,
		LinkPrev:     linkPrev,
		LinkNext:     linkNext,
	}
	indexPageData := indexPageData{
		Title: category.Name + " Articles:",
		Cards: cards,
	}
	pageData := sG.makePageData(pageHead, indexPageData)
	page := new(bytes.Buffer)
	if err := sG.tmplIndex.Execute(page, pageData); err != nil {
		return err
	}
	filePath := sG.getAbsPath(sG.getCategoryIndexPagePath(category, pageNum))
	return writeFile(filePath, page.Bytes())
}

func (sG StaticGen) getCategoryIndexPagePath(category *models.Category, pageNum int) string {
	var pagePath string
	if pageNum <= 0 {
		log.Fatalf(`getCategoryIndexRelPath got uninitialized paramater. pageNum: "%d"`, pageNum)
	} else if pageNum == 1 {
		// First page is root domain.tld/categorys/tag_name/index.html
		pagePath = path.Join(
			sG.config.OutputPaths.CatIndexDir,
			urlSafeName(category.Name),
			`index.html`,
		)
	} else {
		// Other pages follow scheme domain.tld/categorys/tag_name/page/2.html
		pagePath = path.Join(
			sG.config.OutputPaths.CatIndexDir,
			urlSafeName(category.Name),
			sG.config.OutputPaths.IndexPageDir,
			fmt.Sprintf(`%d.html`, pageNum),
		)
	}
	return pagePath
}
