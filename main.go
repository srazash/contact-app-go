package main

import (
	"contactapp/controllers/counter"
	"contactapp/models/contact"
	"fmt"
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
	err := t.template.ExecuteTemplate(w, "layout", data)
	if err != nil {
		c.Logger().Errorf("template rendering error: %v", err)
	}
	return err
}

func main() {
	contact.Load()
	counter.Load()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Static("/static", "static")

	templates, err := template.ParseGlob("views/*.html")
	if err != nil {
		fmt.Printf("Error parsing templates: %v\n", err)
		return // or handle the error appropriately
	}

	renderer := &Templates{
		template: templates,
	}

	fmt.Printf("Number of templates loaded: %d\n", len(templates.Templates()))

	e.Renderer = renderer

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/contacts")
	})

	e.GET("/contacts", func(c echo.Context) error {
		data := map[string]interface{}{
			"Term":     "",
			"Contacts": *contact.Ptr(),
			"Counter":  counter.Count,
			"Debug":    "This is a debug message",
		}
		err := c.Render(http.StatusOK, "index.html", data)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Template error: "+err.Error())
		}
		return nil
	})

	e.GET("/debug", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Debug Page</h1>")
	})

	e.Logger.Fatal(e.Start(":3000"))
}
