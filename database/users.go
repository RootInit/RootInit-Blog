package database

import (
	"blog/models"
)

func createUserCache(users []models.User) *models.UserCache {
	userIdMap := make(map[int]*models.User, len(users))
	userNameMap := make(map[string]*models.User, len(users))
	for i := 0; i < len(users); i++ {
		userIdMap[users[i].Id] = &users[i]
		userNameMap[users[i].Username] = &users[i]
	}
	userCache := models.UserCache{
		List:    users,
		IdMap:   userIdMap,
		NameMap: userNameMap,
	}
	return &userCache
}

func (db DB) createUserTable() error {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "users" (
			"id"		INTEGER,
			"username"	VARCHAR(50)  NOT NULL UNIQUE,
			"password"	VARCHAR(50)	 NOT NULL DEFAULT "",
			"email"		VARCHAR(100) NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		)
	`)
	return err
}

func (db DB) GetUser(userId int) (*models.User, error) {
	var user models.User
	var err error
	stmt :=
		`SELECT 
			u.id,
			u.username,
			u.password,
			u.email
		FROM users AS u
		WHERE u.id = ?`
	row := db.QueryRow(stmt, userId)
	err = row.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Email,
	)
	return &user, err
}

func (db DB) GetUsers() ([]models.User, error) {
	var users []models.User
	var err error
	stmt :=
		`SELECT 
			u.id,
			u.username,
			u.password,
			u.email
		FROM users AS u`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		err = rows.Scan(
			&u.Id,
			&u.Username,
			&u.Password,
			&u.Email,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, err
}

func (db DB) NewUser(username, email, password string) (*models.User, error) {
	var user *models.User
	tx, err := db.Begin()
	if err != nil {
		return user, err
	}
	defer tx.Rollback()
	stmt := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	result, err := tx.Exec(stmt, username, email, password)
	if err != nil {
		return user, err
	}
	resultId, err := result.LastInsertId()
	if err != nil {
		return user, err
	}
	user, err = tx.GetUser(int(resultId))
	if err != nil {
		return user, err
	}
	err = tx.Commit()
	return user, err
}

func (db DB) UpdateUser(user *models.User) error {
	// Create transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Run Update
	stmt :=
		`UPDATE tags SET
			username = ?,
			password = ?,
			email = ?
		WHERE id = ?`
	_, err = db.Exec(stmt, user.Username, user.Password, user.Email, user.Id)
	if err == nil {
		return err
	}
	// Read back user from DB
	updatedUser, err := tx.GetUser(user.Id)
	if err != nil {
		return err
	}
	if updatedUser.Id == 0 {
		err = &InvalidRowUpdate{
			Table:      "articles",
			PrimaryKey: user.Id,
			Message:    "Unable to read updated row back from database.",
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	// Update the user pointer
	user.Username = updatedUser.Username
	user.Password = updatedUser.Password
	user.Email = updatedUser.Email
	return err
}
