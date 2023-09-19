package static_gen

import (
	"bytes"
	"fmt"
	"log"
	"path"
)

func (sG StaticGen) buildMainIndex() error {
	cards := sG.cards.getAllCards()
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
			err := sG.writeMainIndexPage(pageNum, pageCards, lastPage)
			if err != nil {
				return err
			}
			pageCards = make([]cardData, 0, sG.config.IndexSize)
			pageNum += 1
		}
	}
	return nil
}

func (sG StaticGen) writeMainIndexPage(pageNum int, cards []cardData, lastPage bool) error {
	var linkPrev, linkNext string
	if pageNum > 1 {
		linkPrev = sG.getAbsUrl(sG.getIndexRelPath(pageNum - 1))
	}
	if !lastPage {
		linkNext = sG.getAbsUrl(sG.getIndexRelPath(pageNum + 1))
	}
	urlPath := sG.getAbsUrl(sG.getIndexRelPath(pageNum))
	pageHead := headTmplData{
		Title:        sG.config.SiteName + "-Article Index",
		Description:  sG.config.IndexDesc,
		CanonicalUrl: urlPath,
		LinkPrev:     linkPrev,
		LinkNext:     linkNext,
	}
	indexPageData := indexPageData{
		Title: "Articles:",
		Cards: cards,
	}
	pageData := sG.makePageData(pageHead, indexPageData)
	page := new(bytes.Buffer)
	if err := sG.tmplIndex.Execute(page, pageData); err != nil {
		return err
	}
	filePath := sG.getAbsPath(sG.getIndexRelPath(pageNum))
	return writeFile(filePath, page.Bytes())
}

func (sG StaticGen) getIndexRelPath(pageNum int) string {
	var pagePath string
	if pageNum <= 0 {
		log.Fatalf(`getIndexRelPath got uninitialized paramater. pageNum: "%d"`, pageNum)
	} else if pageNum == 1 {
		// First page is root domain.tld/index.html
		pagePath = `index.html`
	} else {
		// Other pages follow scheme domain.tld/page/2.html
		pagePath = path.Join(
			sG.config.OutputPaths.IndexPageDir,
			fmt.Sprintf(`%d.html`, pageNum),
		)
	}
	return pagePath
}
