package models

import (
	"database/sql"
)

type Todo struct {
	ID      int    `json:"id"`
	UserID  int    `json:"userId"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

func FindTodoByID(db *sql.DB, id int) (*Todo, error) {
	todo := &Todo{}
	err := db.QueryRow("select id, user_id, content from todos where id = $1", id).Scan(&todo.ID, &todo.UserID, &todo.Content)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func FindTodosByUserID(db *sql.DB, userId int) ([]*Todo, error) {
	rows, err := db.Query("select id, user_id, content from todos where user_id = $1", userId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		if err := rows.Scan(&todo.ID, &todo.UserID, &todo.Content); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func SaveTodo(tx *sql.Tx, todo *Todo) (*Todo, error) {
	var err error
	if todo.ID <= 0 {
		var id int
		err = tx.QueryRow("insert into todos (user_id, content) values ($1, $2) returning id", todo.UserID, todo.Content).Scan(&id)
		todo.ID = id
	} else {
		_, err = tx.Exec("update todos set content = $1 where id = $2", todo.Content, todo.ID)
	}
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func DeleteTodo(tx *sql.Tx, id int) error {
	_, err := tx.Exec("delete rom todos where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
