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
	env := &config.Env{DB: db}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.POST("/login", controllers.LoginHandler(env))
	g := e.Group("/todos")
	g.Use(middleware.JWT([]byte(config.Secret)))
	g.GET("", controllers.GetTodosHandler(env))
	g.POST("", controllers.CreateTodoHandler(env))
	g.GET("/:id", controllers.GetTodoHandler(env))
	g.PUT("/:id", controllers.UpdateTodoHandler(env))
	g.DELETE("/:id", controllers.DeleteTodoHandler(env))
	e.Logger.Fatal(e.Start(":1323"))
}
