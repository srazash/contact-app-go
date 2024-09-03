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

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
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

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		counter.Increment()
		return c.Redirect(http.StatusFound, "/contacts")
	})

	e.GET("/contacts", func(c echo.Context) error {
		term := c.QueryParam("q")
		page := func() int {
			if c.QueryParam("page") == "" {
				return 1
			}
			p, err := strconv.Atoi(c.QueryParam("page"))
			if err != nil {
				panic(err)
			}
			return p
		}()
		items := 5
		hasNext := contact.NextPage(page, items)

		contacts := func() []contact.Contact {
			if term != "" {
				return contact.Search(term)
			}
			return contact.PaginatedContacts(page, items)
		}()

		if c.Request().Header.Get("HX-Trigger") == "search" {
			data := map[string]interface{}{
				"Contacts":    contacts,
				"HasNextPage": hasNext,
			}
			return c.Render(http.StatusOK, "rows", data)
		}

		title := func() string {
			if term != "" {
				return fmt.Sprintf("search results for \"%s\"", term)
			}
			return "all contacts"
		}()

		message := func() string {
			if term != "" {
				return fmt.Sprintf("Showing results for search term: \"%s\", %d found", term, len(contacts))
			}
			return ""
		}()

		reset := func() string {
			if term != "" {
				return "Reset"
			}
			return ""
		}()

		data := map[string]interface{}{
			"Title":         title,
			"Term":          term,
			"Message":       message,
			"Reset":         reset,
			"Contacts":      contacts,
			"Counter":       counter.PaddedCount(),
			"HasNextPage":   hasNext,
			"NextPage":      page + 1,
			"ContactsCount": contact.ContactsCount(),
			"Template":      "index",
		}

		return c.Render(http.StatusOK, "layout", data)
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
		values := make(map[string]string)

		values["First"] = c.FormValue("first")
		values["Last"] = c.FormValue("last")
		values["Email"] = c.FormValue("email")
		values["Phone"] = c.FormValue("phone")

		errors := contact.ValidateForm(&values)

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

	e.GET("/contacts/new/email", func(c echo.Context) error {
		email := c.FormValue("email")
		email_error := contact.ValidateEmail(email)
		return c.String(http.StatusOK, email_error)
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
			"Title":   fmt.Sprintf("view contact %d", con.Id),
			"Contact": con,
			"Counter": counter.PaddedCount(),
		}
		return c.Render(http.StatusOK, "view", data)
	})

	e.GET("/contacts/:contact_id/edit", func(c echo.Context) error {
		contact_id, err := strconv.Atoi(c.Param("contact_id"))
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/contacts")
		}
		con, err := contact.Find(contact_id)
		if err != nil {
			return c.Redirect(http.StatusNotFound, "/contacts")
		}
		data := map[string]interface{}{
			"Title":   fmt.Sprintf("edit contact %d", con.Id),
			"Contact": con,
			"Counter": counter.PaddedCount(),
		}
		return c.Render(http.StatusOK, "edit", data)
	})

	e.POST("/contacts/:contact_id/edit", func(c echo.Context) error {
		contact_id, err := strconv.Atoi(c.Param("contact_id"))
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/contacts")
		}

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
			return c.Render(http.StatusOK, "edit", data)
		}

		contact.Update(contact_id, values["First"], values["Last"], values["Email"], values["Phone"])
		path := fmt.Sprintf("/contacts/%d", contact_id)
		return c.Redirect(http.StatusFound, path)
	})

	e.DELETE("/contacts/:contact_id", func(c echo.Context) error {
		contact_id, err := strconv.Atoi(c.Param("contact_id"))
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/contacts")
		}

		contact.Delete(contact_id)
		return c.Redirect(http.StatusSeeOther, "/contacts")
	})

	e.Logger.Fatal(e.Start(":3000"))
}
