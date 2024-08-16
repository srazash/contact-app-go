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

type Templates struct {
	template *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.template.ExecuteTemplate(w, name, data)
}

func main() {
	contact.Load()
	counter.Load()

	e := echo.New()

	renderer := &Templates{
		template: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = renderer

	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/contacts")
	})

	e.GET("/contacts", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"Term":     "",
			"Contacts": *contact.Ptr(),
			"Counter":  counter.Count,
		})
	})

	e.Logger.Fatal(e.Start(":3000"))
}
