package database

import "blog/models"

func createTagCache(tags []models.Tag) *models.TagCache {
	tagIdMap := make(map[int]*models.Tag, len(tags))
	for i := 0; i < len(tags); i++ {
		tagIdMap[tags[i].Id] = &tags[i]
	}
	tagCache := models.TagCache{
		List:  tags,
		IdMap: tagIdMap,
	}
	return &tagCache
}

func (db DB) createTagsTable() error {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "tags" (
			"id"	INTEGER,
			"name"	TEXT NOT NULL UNIQUE,
			"icon"	TEXT NOT NULL,
			PRIMARY KEY("id" AUTOINCREMENT)
		)`,
	)
	return err
}

func (db DB) createArticleTagsTable() error {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "article_tags" (
			"article_id"	INTEGER NOT NULL,
			"tag_id"		INTEGER NOT NULL,
			PRIMARY KEY("article_id","tag_id"),
			FOREIGN KEY("article_id") REFERENCES "articles"("id"),
			FOREIGN KEY("tag_id") REFERENCES "tags"("id")
		)`,
	)
	return err
}

func (db DB) GetTags() ([]models.Tag, error) {
	var tags []models.Tag
	var err error
	rows, err := db.Query(
		`SELECT 
			t.id,
			t.name,
			t.icon
		FROM tags AS t`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t models.Tag
		err = rows.Scan(
			&t.Id,
			&t.Name,
			&t.Icon,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, err
}

func (db DB) GetTag(id int) (*models.Tag, error) {
	var tag models.Tag
	var err error
	stmt :=
		`SELECT 
			t.id, 
			t.name, 
			t.icon
		FROM tags as t
		WHERE id = ?`
	row := db.QueryRow(stmt, id)
	err = row.Scan(
		&tag.Id,
		&tag.Name,
		&tag.Icon,
	)
	return &tag, err
}

func (db DB) NewTag(name, icon string) (*models.Tag, error) {
	var tag *models.Tag
	tx, err := db.Begin()
	if err != nil {
		return tag, err
	}
	defer tx.Rollback()
	stmt :=
		`INSERT INTO tags (name, icon) 
		VALUES (?, ?)`
	result, err := tx.Exec(stmt, name, icon)
	if err != nil {
		return tag, err
	}
	resultId, err := result.LastInsertId()
	if err != nil {
		return tag, err
	}
	tag, err = tx.GetTag(int(resultId))
	if err != nil {
		return tag, err
	}
	err = tx.Commit()
	return tag, err
}

func (db DB) UpdateTag(tag *models.Tag) error {
	// Create transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Run Update
	stmt :=
		`UPDATE tags SET
			name = ?,
			icon = ?
		WHERE id = ?`
	_, err = db.Exec(stmt, tag.Name, tag.Icon, tag.Id)
	if err == nil {
		return err
	}
	// Read back tag from DB
	updatedTag, err := tx.GetTag(tag.Id)
	if err != nil {
		return err
	}
	if updatedTag.Id == 0 {
		err = &InvalidRowUpdate{
			Table:      "tags",
			PrimaryKey: tag.Id,
			Message:    "Unable to read updated row back from database.",
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	// Update the tag pointer
	tag.Name = updatedTag.Name
	tag.Icon = updatedTag.Icon
	return err
}
