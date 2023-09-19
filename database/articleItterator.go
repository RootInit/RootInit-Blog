package database

import "blog/models"

type articleIttr struct {
	Article  models.Article
	Comments []models.Comment
	lastId   int
	db       DB
	Error    error
}

func (db DB) GetArticleIttr() *articleIttr {
	aI := articleIttr{
		db:     db,
		lastId: -1,
	}
	return &aI
}

func (aI *articleIttr) Next() bool {
	var err error
	aI.Article, err = aI.db.GetArticle(aI.lastId, true, true)
	if err != nil {
		aI.Error = err
		return false
	}
	aI.Comments, err = aI.db.GetArticleComments(aI.Article.Id)
	if err != nil {
		aI.Error = err
		return false
	}
	aI.lastId = aI.Article.Id
	return true
}
