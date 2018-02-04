package main

import (
	"database/sql"
	"html/template"
	"io"
	"log"

	_ "github.com/lib/pq"

	"github.com/dirges/todo/config"
	"github.com/dirges/todo/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	connStr := "postgres://postgres:example@db/todos?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	env := &config.Env{db}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.GET("/todos", controllers.GetTodosHandler(env))
	e.POST("/todos", controllers.CreateTodoHandler(env))
	e.GET("/todos/:id", controllers.GetTodoHandler(env))
	e.PUT("/todos/:id", controllers.UpdateTodoHandler(env))
	e.DELETE("/todos/:id", controllers.DeleteTodoHandler(env))

	e.POST("/login", controllers.LoginHandler(env))
	e.Logger.Fatal(e.Start(":1323"))
}
