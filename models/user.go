package models

import (
	"database/sql"
)

type User struct {
	ID       int
	Username string
	Password string
}

func FindUserByID(db *sql.DB, id int) (*User, error) {
	user := &User{}
	err := db.QueryRow("select id, username, password from users where id = $1", id).Scan(&user.ID, &user.Username,
		&user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindUserByUsername(db *sql.DB, username string) (*User, error) {
	user := &User{}
	err := db.QueryRow("select id, username, password from users where username = $1", username).Scan(&user.ID,
		&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
