package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/dirges/todo/config"
	"github.com/dirges/todo/models"
)

func GetTodoHandler(env *config.Env) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		todo := models.Todo{}
		err := env.DB.QueryRow("select id, content from todos where id = $1", id).Scan(&todo.ID, &todo.Content)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, todo)
	})
}

func GetTodosHandler(env *config.Env) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		rows, err := env.DB.Query("select * from todos")
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		defer rows.Close()
		var todos []models.Todo
		for rows.Next() {
			todo := models.Todo{}
			if err := rows.Scan(&todo.ID, &todo.Content); err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			todos = append(todos, todo)
		}
		return c.JSON(http.StatusOK, &todos)
	})
}

type CreateTodoRequest struct {
	Content string `json:"content"`
}

func CreateTodoHandler(env *config.Env) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		req := &CreateTodoRequest{}
		if err := c.Bind(req); err != nil {
			return err
		}
		tx, err := env.DB.Begin()
		if err != nil {
			return err
		}
		var id int
		err = tx.QueryRow("insert into todos (content) values ($1) returning id", req.Content).Scan(&id)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}
		return c.JSON(http.StatusCreated, &models.Todo{id, req.Content, false})
	})
}

type UpdateTodoRequest struct {
	Content string `json:"content"`
}

func UpdateTodoHandler(env *config.Env) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		req := &UpdateTodoRequest{}
		if err := c.Bind(req); err != nil {
			return err
		}
		tx, err := env.DB.Begin()
		if err != nil {
			return err
		}
		_, err = tx.Exec("update todos set content = $1 where id = $2", req.Content, id)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}
		return c.JSON(http.StatusOK, &models.Todo{id, req.Content, false})
	})
}

func DeleteTodoHandler(env *config.Env) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		tx, err := env.DB.Begin()
		if err != nil {
			return err
		}
		_, err = tx.Exec("delete from todos where id = $1", id)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}
		return c.NoContent(http.StatusNoContent)
	})
}
