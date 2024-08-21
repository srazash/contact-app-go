package main

import (
	"contactapp/controllers/counter"
	"contactapp/models/contact"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

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
	templates["new"] = template.Must(template.ParseFiles("views/layout.html", "views/new.html"))
	templates["view"] = template.Must(template.ParseFiles("views/layout.html", "views/view.html"))

	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	e.GET("/", func(c echo.Context) error {
		counter.Increment()
		return c.Redirect(http.StatusFound, "/contacts")
	})

	e.GET("/contacts", func(c echo.Context) error {
		data := map[string]interface{}{
			"Title":    "all contacts",
			"Term":     "",
			"Contacts": *contact.Ptr(),
			"Counter":  counter.PaddedCount(),
		}
		return c.Render(http.StatusOK, "index", data)
	})

	e.GET("/contacts/new", func(c echo.Context) error {
		data := map[string]interface{}{
			"Title":   "new contact",
			"Values":  make(map[string]string),
			"Errors":  make(map[string]string),
			"Counter": counter.PaddedCount(),
		}
		return c.Render(http.StatusOK, "new", data)
	})

	e.POST("/contacts/new", func(c echo.Context) error {
		errors := make(map[string]string)
		values := make(map[string]string)

		values["First"] = c.FormValue("first")
		values["Last"] = c.FormValue("last")
		values["Email"] = c.FormValue("email")
		values["Phone"] = c.FormValue("phone")

		switch {
		case values["First"] == "":
			errors["First"] = "First name is required"
			fallthrough
		case values["Last"] == "":
			errors["Last"] = "Last name is required"
			fallthrough
		case values["Email"] == "":
			errors["Email"] = "Email is required"
			fallthrough
		case values["Phone"] == "":
			errors["Phone"] = "Phone number is required"
		}

		if len(errors) != 0 {
			data := map[string]interface{}{
				"Title":   "new contact",
				"Values":  values,
				"Errors":  errors,
				"Counter": counter.PaddedCount(),
			}
			return c.Render(http.StatusOK, "new", data)
		}

		contact_id := contact.Create(values["First"], values["Last"], values["Email"], values["Phone"])
		path := fmt.Sprintf("/contacts/%d", contact_id)
		return c.Redirect(http.StatusFound, path)
	})

	e.GET("/contacts/:contact_id", func(c echo.Context) error {
		contact_id, err := strconv.Atoi(c.Param("contact_id"))
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/contacts")
		}
		con, err := contact.Find(contact_id)
		if err != nil {
			return c.Redirect(http.StatusNotFound, "/contacts")
		}
		data := map[string]interface{}{
			"Title":   "all contacts",
			"Term":    "",
			"Contact": con,
			"Counter": counter.PaddedCount(),
		}
		return c.Render(http.StatusOK, "view", data)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
