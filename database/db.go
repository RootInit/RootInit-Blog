package database

import (
	"blog/models"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

/*
Wrapper arround *sql.DB with support for transactions
Includes a "cache" of some tables.
*/
type DB struct {
	db               *sql.DB
	tx               *sql.Tx
	UserCache        *models.UserCache
	TagCache         *models.TagCache
	CategoryCache    *models.CategoryCache
	ArticleFileCache *models.ArticleFileCache
}

func (m *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.tx != nil {
		return m.tx.Query(query, args...)
	}
	return m.db.Query(query, args...)
}

func (m *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.tx != nil {
		return m.tx.QueryRow(query, args...)
	}
	return m.db.QueryRow(query, args...)
}

func (m *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.tx != nil {
		return m.tx.Exec(query, args...)
	}
	return m.db.Exec(query, args...)
}

func (db *DB) Begin() (DB, error) {
	var err error
	db.tx, err = db.db.Begin()
	return *db, err
}

func (db *DB) Rollback() error {
	var err error
	if db.tx != nil {
		err = db.tx.Rollback()
		db.tx = nil
	}
	return err
}

func (db *DB) Commit() error {
	var err error
	if db.tx != nil {
		err = db.tx.Commit()
		db.tx = nil
	}
	return err
}

func InitializeDB(path string) (DB, error) {
	// Open database
	dbConn, err := sql.Open("sqlite3", path)
	if err != nil {
		return DB{}, err
	}
	db := DB{db: dbConn}

	// Ensure tables are created
	err = db.createTables()
	if err != nil {
		return db, err
	}
	// Load cache
	if err := db.LoadCaches(); err != nil {
		return db, err
	}
	return db, err
}

func (db DB) createTables() error {
	err := db.createArticleTable()
	if err != nil {
		log.Println(`Failed to load table: "articles"`)
		return err
	}
	err = db.createTagsTable()
	if err != nil {
		log.Println(`Failed to load table: "tags"`)
		return err
	}
	err = db.createCategoriesTable()
	if err != nil {
		log.Println(`Failed to load table: "categories"`)
		return err
	}
	err = db.createArticleTagsTable()
	if err != nil {
		log.Println(`Failed to load table: "article_tags"`)
		return err
	}
	err = db.createUserTable()
	if err != nil {
		log.Println(`Failed to load table: "users"`)
		return err
	}
	err = db.createCommentTable()
	if err != nil {
		log.Println(`Failed to load table: "comments"`)
		return err
	}
	return err
}

func (db *DB) LoadCaches() error {
	categories, err := db.GetCategories()
	if err != nil {
		return err
	}
	db.CategoryCache = createCategoryCache(categories)

	tags, err := db.GetTags()
	if err != nil {
		return err
	}
	db.TagCache = createTagCache(tags)

	users, err := db.GetUsers()
	if err != nil {
		return err
	}
	db.UserCache = createUserCache(users)

	db.ArticleFileCache, err = db.createArticleFileCache()
	return err
}
