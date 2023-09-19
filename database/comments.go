package database

import (
	"blog/models"
	"database/sql"
	"time"
)

func (db DB) createCommentTable() error {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "comments" (
			"id"			INTEGER,
			"article_id"	INTEGER NOT NULL,
			"parent_id"		INTEGER NOT NULL DEFAULT 0,
			"user_id"		INTEGER NOT NULL,
			"timestamp"		TIMESTAMP NOT NULL,
			"body"			varchar(2000) NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("parent_id") REFERENCES "comments"("id"),
		FOREIGN KEY("user_id") REFERENCES "users"("id"),
		FOREIGN KEY("article_id") REFERENCES "articles"("id")
		)`,
	)
	return err
}

/*
Method to get a comment by id.
Comment replies will NOT be retreived.
*/
func (db DB) GetComment(commentId int) (models.Comment, error) {
	var comment models.Comment
	var err error
	stmt :=
		`SELECT 
			c.id,
			c.article_id,
			c.user_id,
			c.timestamp,
			c.body
		FROM comments AS c
		WHERE c.id = ?`
	row := db.QueryRow(stmt, commentId)
	var userId int
	err = row.Scan(
		&comment.Id,
		&comment.ArticleId,
		&userId,
		&comment.Date,
		&comment.Body,
	)
	if err != nil {
		return comment, err
	}
	// Add the comment author User
	author, exists := db.UserCache.IdMap[userId]
	if !exists {
		return comment, err
	}
	comment.Author = *author
	return comment, err
}

type preComment struct {
	comment  models.Comment
	parentId int
}

func makeCommentTree(preComments []preComment) []models.Comment {
	var rootComments []models.Comment
	var replies []preComment
	for _, pc := range preComments {
		if pc.parentId == 0 {
			// Top level comment
			rootComments = append(rootComments, pc.comment)
		} else {
			// Reply to another comment
			replies = append(replies, pc)
		}
	}
	/* Recursive function to add replies */
	var addReplies func(com *models.Comment, replyPool *[]preComment)
	addReplies = func(com *models.Comment, replyPool *[]preComment) {
		for i := 0; i < len(*replyPool); {
			reply := (*replyPool)[i]
			if reply.parentId == com.Id {
				replyComment := reply.comment
				// Remove the reply from the pool
				*replyPool = append((*replyPool)[:i], (*replyPool)[i+1:]...)
				addReplies(&replyComment, replyPool)
				com.Replies = append(com.Replies, replyComment)
			} else {
				i += 1
			}
		}
	}
	// Add replies to root comments
	for i := 0; i < len(replies); i++ {
		addReplies(&rootComments[i], &replies)
	}
	return rootComments
}

func (db DB) GetArticleComments(articleId int) ([]models.Comment, error) {
	var err error
	stmt :=
		`SELECT 
			c.id,
			c.article_id,
			c.parent_id,
			c.user_id,
			c.timestamp,
			c.body
		FROM comments AS c
		WHERE c.article_id = ?`
	rows, err := db.Query(stmt, articleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create map to hold comments with parentId
	var preComments []preComment
	// Itterate through rows
	for rows.Next() {
		var c models.Comment
		var parentIdNull sql.NullInt32
		var userId int
		err = rows.Scan(
			&c.Id,
			&c.ArticleId,
			&parentIdNull,
			&userId,
			&c.Date,
			&c.Body,
		)
		if err != nil {
			return nil, err
		}
		// Add the comment author User
		author, exists := db.UserCache.IdMap[userId]
		if !exists {
			return nil, err
		}
		c.Author = *author

		if !parentIdNull.Valid || parentIdNull.Int32 == 0 {
			// Top level comment
			preComments = append(preComments, preComment{
				comment:  c,
				parentId: 0,
			})
		} else {
			// Reply
			preComments = append(preComments, preComment{
				comment:  c,
				parentId: int(parentIdNull.Int32),
			})
		}
	}

	commentTree := makeCommentTree(preComments)
	return commentTree, err
}

func (db DB) NewComment(articleId, parentId, userId int, date time.Time, body string) (models.Comment, error) {
	var comment models.Comment
	// Create transaction
	tx, err := db.Begin()
	if err != nil {
		return comment, err
	}
	defer tx.Rollback()
	// Insert the comment
	stmt := `INSERT INTO comments (article_id, parent_id, user_id, timestamp, body) VALUES (?, ?, ?, ?, ?)`
	result, err := tx.Exec(stmt, articleId, parentId, userId, time.Now(), body)
	if err != nil {
		return comment, err
	}
	resultId, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}
	// Read back the inserted comment
	comment, err = tx.GetComment(int(resultId))
	if err != nil {
		return comment, err
	}
	err = tx.Commit()
	return comment, err
}

func (db DB) UpdateComment(comment *models.Comment) error {
	// Create transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Run Update
	stmt :=
		`UPDATE comments SET
			timestamp = ?,
			body = ?
		WHERE id = ?`
	_, err = db.Exec(stmt, comment.Date, comment.Body, comment.Id)
	if err == nil {
		return err
	}
	// Read back comment from DB
	updatedComment, err := tx.GetComment(comment.Id)
	if err != nil {
		return err
	}
	if updatedComment.Id == 0 {
		err = &InvalidRowUpdate{
			Table:      "comments",
			PrimaryKey: comment.Id,
			Message:    "Unable to read updated row back from database.",
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	// Update the comment pointer
	comment.Date = updatedComment.Date
	comment.Body = updatedComment.Body
	return err
}
