package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/dirges/todo/config"
	"github.com/dirges/todo/models"
)

func GetTodoHandler(env *config.Env) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		todo, err := models.FindTodoByID(env.DB, id)
		switch {
		case err == sql.ErrNoRows:
			return echo.ErrNotFound
		case err != nil:
			return err
		}
		if todo.UserID != env.GetCurrentUserID(c) {
			return echo.ErrForbidden
		}
		return c.JSON(http.StatusOK, todo)
	})
}

func GetTodosHandler(env *config.Env) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		todos, err := models.FindTodosByUserID(env.DB, env.GetCurrentUserID(c))
		if err != nil {
			return err
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
		userID := env.GetCurrentUserID(c)
		todo, err := models.SaveTodo(tx, &models.Todo{ID: 0, UserID: userID, Content: req.Content, Done: false})
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}
		return c.JSON(http.StatusCreated, todo)
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
		todo, err := models.FindTodoByID(env.DB, id)
		if err != nil {
			return err
		}
		tx, err := env.DB.Begin()
		if err != nil {
			return err
		}
		todo.Content = req.Content
		todo, err = models.SaveTodo(tx, todo)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}
		return c.JSON(http.StatusOK, todo)
	})
}

func DeleteTodoHandler(env *config.Env) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		tx, err := env.DB.Begin()
		if err != nil {
			return err
		}
		err = models.DeleteTodo(tx, id)
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
