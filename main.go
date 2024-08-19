package main

import (
	"contactapp/controllers/counter"
	"contactapp/models/contact"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRegistry struct {
	templates map[string]*template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates[name].Execute(w, data)
}

func main() {
	contact.Load()
	counter.Load()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Static("/", "static")

	templates := make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles("views/layout.html", "views/index.html"))

	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	e.GET("/", func(c echo.Context) error {
		counter.Increment()
		return c.Redirect(http.StatusFound, "/contacts")
	})

	e.GET("/contacts", func(c echo.Context) error {
		data := map[string]interface{}{
			"Title":    "contacts.app",
			"Term":     "",
			"Contacts": *contact.Ptr(),
			"Counter":  counter.PaddedCount(),
		}
		return c.Render(http.StatusOK, "index", data)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
