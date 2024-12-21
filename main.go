package main

import (
	"berlang/terminal"
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
    e := echo.New()
    e.Use(middleware.Logger())

    terminal := terminal.NewTerminal()
    e.Renderer = newTemplate()

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index.html", nil)
    })

e.POST("/execute", func(c echo.Context) error {
    command := c.FormValue("command")
    result := terminal.ExecuteCommand(command)
    return c.Render(200, "terminal_output.html", result)
})

    e.Logger.Fatal(e.Start(":3000"))
}
