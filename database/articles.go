package database

import (
	"blog/models"
	"time"
)

func (db DB) createArticleTable() error {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "articles" (
			"id"			INTEGER,
			"title"			VARCHAR(150) NOT NULL,
			"url_title" 	VARCHAR(100) NOT NULL,
			"description"	VARCHAR(250) NOT NULL,
			"category_id" 	INTEGER		 NOT NULL,
			"timestamp"		TIMESTAMP	 NOT NULL,
			"body"			BLOB		 NOT NULL,
			"thumbnail"		VARCHAR(150) NOT NULL,
			PRIMARY KEY("id" AUTOINCREMENT)
			FOREIGN KEY("category_id") REFERENCES "categories"("id")
		)`,
	)
	return err
}

func (db DB) createArticleFileCache() (*models.ArticleFileCache, error) {
	stmt := `SELECT a.id, a.url_title FROM articles AS a`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var idList []int
	idMap := make(map[int]string)
	for rows.Next() {
		var id int
		var fileName string
		err = rows.Scan(
			&id,
			&fileName,
		)
		if err != nil {
			return nil, err
		}
		idList = append(idList, id)
		idMap[id] = fileName
	}
	articleFileCache := models.ArticleFileCache{
		IdList: idList,
		IdMap:  idMap,
	}
	return &articleFileCache, err
}

func (db DB) GetArticleCount() (int, error) {
	stmt := `SELECT COUNT(id) FROM articles`
	row := db.QueryRow(stmt)
	var id int
	err := row.Scan(&id)
	return id, err
}

// Method to get an article by Id, or of `afterId` is `true` the previous article.
// Call with an ID of `-1` to get the most recient article.
func (db DB) GetArticle(id int, beforeId, getComments bool) (models.Article, error) {
	var article models.Article
	var err error
	stmt :=
		`SELECT
		    a.id,
		    a.title,
			a.url_title,
		    a.description,
			a.category_id,
		    a.timestamp,
			a.body,
			a.thumbnail,
		    GROUP_CONCAT(DISTINCT at.tag_id) AS tag_list
		FROM articles AS a
		LEFT JOIN article_tags AS at ON at.article_id = a.id`
	if id == -1 {
		// Get first article
		stmt += ` WHERE a.id > ?`
	} else if beforeId {
		// Get article posted before id
		stmt += ` WHERE a.id < ?`
	} else {
		// Get matching article
		stmt += ` WHERE a.id = ?`
	}
	stmt +=
		` GROUP BY a.id
		ORDER BY a.id DESC
		LIMIT 1`
	row := db.QueryRow(stmt, id)
	var tagString string
	var categoryId int
	err = row.Scan(
		&article.Id,
		&article.Title,
		&article.UrlTitle,
		&article.Description,
		&categoryId,
		&article.Date,
		&article.Body,
		&article.Thumbnail,
		&tagString,
	)
	if err != nil {
		return article, err
	}
	// Add Tags
	tagIds := parseIntList(tagString)
	for _, id := range tagIds {
		tag := db.TagCache.IdMap[id]
		article.Tags = append(article.Tags, *tag)
	}
	//Add Category
	category := db.CategoryCache.IdMap[categoryId]
	article.Category = *category
	// Add Comments
	if getComments {
		article.Comments, err = db.GetArticleComments(article.Id)
		if err != nil {
			return article, err
		}
	}
	return article, err
}

func (db DB) NewArticle(
	title, urlTitle, desc string, categoryId int,
	date time.Time, body string, thumb string,
	tagIds []int) (models.Article, error) {
	var article models.Article
	// Create transaction
	var err error
	tx, err := db.Begin()
	if err != nil {
		return article, err
	}
	defer tx.Rollback()
	// Insert the article
	stmt :=
		`INSERT INTO articles
		(title, url_title, description, category_id, timestamp, body, thumbnail)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(stmt, title, urlTitle, desc, categoryId, date, body, thumb)
	if err != nil {
		return article, err
	}
	articleId, err := result.LastInsertId()
	if err != nil {
		return article, err
	}
	if len(tagIds) > 0 {
		// Insert the tag rows
		stmt =
			`INSERT INTO article_tags (article_id, tag_id)
			VALUES (?, ?)`
		var params []any
		params = append(params, articleId, tagIds[0])
		for i := 1; i < len(tagIds); i++ {
			stmt += `, (?, ?)`
			params = append(params, articleId, tagIds[i])
		}
		_, err = tx.Exec(stmt, params...)
		if err != nil {
			return article, err
		}
	}
	if err := tx.Commit(); err != nil {
		return article, err
	}
	// Return article object
	article, err = tx.GetArticle(int(articleId), false, false)
	return article, err
}

func (db DB) UpdateArticle(article *models.Article) error {
	// Create transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Run Update
	stmt :=
		`UPDATE articles SET
			id = id,
			title = ?,
			url_title = ?,
			description = ?,
			category_id = ?,
			timestamp = ?,
			body = ?,
			thumbnail = ?
		WHERE id = ?`
	_, err = db.Exec(
		stmt,
		article.Title,
		article.UrlTitle,
		article.Description,
		article.Category.Id,
		article.Date,
		article.Body,
		article.Thumbnail,
		article.Id,
	)
	if err == nil {
		return err
	}
	// Read back article from DB
	updatedArticle, err := tx.GetArticle(article.Id, false, false)
	if err != nil {
		return err
	}
	if updatedArticle.Id == 0 {
		err = &InvalidRowUpdate{
			Table:      "articles",
			PrimaryKey: article.Id,
			Message:    "Unable to read updated row back from database.",
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	// Update the article pointer
	article.Title = updatedArticle.Title
	article.UrlTitle = updatedArticle.UrlTitle
	article.Description = updatedArticle.Description
	article.Category = updatedArticle.Category
	article.Date = updatedArticle.Date
	article.Body = updatedArticle.Body
	article.Thumbnail = updatedArticle.Thumbnail
	return err
}
