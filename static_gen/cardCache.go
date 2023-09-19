package static_gen

import (
	"blog/models"
	"sort"
)

/* Holds article cards for use in indices. */
type cardCache struct {
	articleIds    []int
	all           map[int]cardData // articleId:card
	tagArticleIds map[int][]int    // tagId:[]articleId
	catArticleIds map[int][]int    // categoryId:[]articleId
}

func newCardCache() *cardCache {
	cardCache := cardCache{
		all:           make(map[int]cardData),
		tagArticleIds: make(map[int][]int),
		catArticleIds: make(map[int][]int),
	}
	return &cardCache
}

func (sG StaticGen) addCardToCache(article *models.Article) {
	sG.cards.articleIds = append(sG.cards.articleIds, article.Id)
	cardData := sG.convertArticleCard(article)
	sG.cards.all[article.Id] = cardData
	for _, tag := range article.Tags {
		sG.cards.tagArticleIds[tag.Id] = append(
			sG.cards.tagArticleIds[tag.Id],
			article.Id,
		)
	}
	sG.cards.catArticleIds[article.Category.Id] = append(
		sG.cards.catArticleIds[article.Category.Id],
		article.Id,
	)
}

func (cC cardCache) getAllCards() []cardData {
	cards := make([]cardData, len(cC.articleIds))
	for i, id := range cC.articleIds {
		cards[i] = cC.all[id]
	}
	return cards
}
func (cC cardCache) getCardsByIds(cardIds []int) []cardData {
	cards := make([]cardData, len(cardIds))
	for i, id := range cardIds {
		cards[i] = cC.all[id]
	}
	return cards
}
func (cC cardCache) getCardsFromTagId(tagId int) []cardData {
	articleIds := cC.tagArticleIds[tagId]
	sort.Slice(articleIds, func(a, b int) bool {
		return articleIds[b] < articleIds[a]
	})
	cards := cC.getCardsByIds(articleIds)
	return cards
}
func (cC cardCache) getCardsFromCategoryId(categoryId int) []cardData {
	articleIds, exists := cC.catArticleIds[categoryId]
	sort.Slice(articleIds, func(a, b int) bool {
		return articleIds[b] < articleIds[a]
	})
	if !exists {
		return nil
	}
	cards := cC.getCardsByIds(articleIds)
	return cards
}
