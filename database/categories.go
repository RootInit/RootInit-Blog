package database

import (
	"blog/models"
)

func createCategoryCache(categories []models.Category) *models.CategoryCache {
	categoryIdMap := make(map[int]*models.Category, len(categories))
	for i := 0; i < len(categories); i++ {
		categoryIdMap[categories[i].Id] = &categories[i]
	}
	categoryCache := models.CategoryCache{
		List:  categories,
		IdMap: categoryIdMap,
	}
	return &categoryCache
}

func (db DB) createCategoriesTable() error {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "categories" (
			"id"		INTEGER,
			"name"		VARCHAR(100) NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		)
	`)
	return err
}

func (db DB) GetCategory(categoryId int) (models.Category, error) {
	var category models.Category
	var err error
	stmt :=
		`SELECT 
			c.id,
			c.name
		FROM categories AS c
		WHERE c.id = ?`
	row := db.QueryRow(stmt, categoryId)
	err = row.Scan(
		&category.Id,
		&category.Name,
	)
	return category, err
}

func (db DB) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	var err error
	stmt :=
		`SELECT 
			c.id,
			c.name
		FROM categories AS c`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.Category
		err = rows.Scan(
			&c.Id,
			&c.Name,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, err
}

func (db DB) NewCategory(name string) (models.Category, error) {
	var category models.Category
	tx, err := db.Begin()
	if err != nil {
		return category, err
	}
	defer tx.Rollback()
	stmt := `INSERT INTO categories (name) VALUES (?)`
	result, err := tx.Exec(stmt, name)
	if err != nil {
		return category, err
	}
	resultId, err := result.LastInsertId()
	if err != nil {
		return category, err
	}
	category, err = tx.GetCategory(int(resultId))
	if err != nil {
		return category, err
	}
	err = tx.Commit()
	return category, err
}

func (db DB) UpdateCategory(category *models.Category) error {
	// Create transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Run Update
	stmt :=
		`UPDATE tags SET
			name = ?
		WHERE id = ?`
	_, err = db.Exec(stmt, category.Name)
	if err == nil {
		return err
	}
	// Read back cache from DB
	updatedCategory, err := tx.GetCategory(category.Id)
	if err != nil {
		return err
	}
	if updatedCategory.Id == 0 {
		err = &InvalidRowUpdate{
			Table:      "articles",
			PrimaryKey: category.Id,
			Message:    "Unable to read updated row back from database.",
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	// Update the cache pointer
	category.Name = updatedCategory.Name
	return err
}
