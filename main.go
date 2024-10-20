package main

import (
	"contactapp/controllers/archiver"
	"contactapp/controllers/counter"
	"contactapp/models/contact"
	"fmt"
	"html/template"
	"io"
	"log"
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

	a := archiver.Get()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Static("/", "static")

	e.GET("/", func(c echo.Context) error {
		counter.Increment()
		return c.Redirect(http.StatusFound, "/contacts")
	})

	e.GET("/contacts", func(c echo.Context) error {
		t := &Template{
			templates: template.Must(template.ParseFiles("views/layout.html",
				"views/index.html", "views/rows.html")),
		}
		e.Renderer = t

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
				"HasNextPage": false,
				"Term":        term,
				"Search":      true,
			}
			return c.Render(http.StatusOK, "rows", data)
		}

		title := func() string {
			if term != "" {
				return fmt.Sprintf("search results for \"%s\"", term)
			}
			return "all contacts"
		}()

		data := map[string]interface{}{
			"Title":       title,
			"Term":        term,
			"Search":      false,
			"Contacts":    contacts,
			"Counter":     counter.PaddedCount(),
			"HasNextPage": hasNext,
			"NextPage":    page + 1,
			"Archive":     a,
		}

		return c.Render(http.StatusOK, "index", data)
	})

	e.GET("/contacts/count", func(c echo.Context) error {
		cc := contact.ContactsCount()
		s := func() string {
			if cc == 1 {
				return fmt.Sprintf(`(%d total contact)`, cc)
			}
			return fmt.Sprintf(`(%d total contacts)`, cc)
		}()
		return c.String(http.StatusOK, s)
	})

	e.GET("/contacts/new", func(c echo.Context) error {
		t := &Template{
			templates: template.Must(template.ParseFiles("views/layout.html",
				"views/new.html")),
		}
		e.Renderer = t

		data := map[string]interface{}{
			"Title":   "new contact",
			"Values":  make(map[string]string),
			"Errors":  make(map[string]string),
			"Counter": counter.PaddedCount(),
		}
		return c.Render(http.StatusOK, "new", data)
	})

	e.POST("/contacts/new", func(c echo.Context) error {
		t := &Template{
			templates: template.Must(template.ParseFiles("views/layout.html",
				"views/new.html")),
		}
		e.Renderer = t

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
		t := &Template{
			templates: template.Must(template.ParseFiles("views/layout.html",
				"views/view.html")),
		}
		e.Renderer = t

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
		t := &Template{
			templates: template.Must(template.ParseFiles("views/layout.html",
				"views/edit.html")),
		}
		e.Renderer = t

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
		t := &Template{
			templates: template.Must(template.ParseFiles("views/layout.html",
				"views/edit.html")),
		}
		e.Renderer = t

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

	e.DELETE("/contacts", func(c echo.Context) error {
		contact_ids := []int{}

		c.FormParams()
		for _, idstr := range c.Request().Form["selected_contact_ids"] {
			log.Println(idstr)
			id, err := strconv.Atoi(idstr)
			if err != nil {
				continue
			}
			contact_ids = append(contact_ids, id)
		}

		contact.MultiDelete(contact_ids)

		t := &Template{
			templates: template.Must(template.ParseFiles("views/layout.html",
				"views/index.html", "views/rows.html")),
		}
		e.Renderer = t

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

		data := map[string]interface{}{
			"Term":        term,
			"Search":      term != "",
			"Contacts":    contacts,
			"Counter":     counter.PaddedCount(),
			"HasNextPage": hasNext,
			"NextPage":    page + 1,
		}

		return c.Render(http.StatusOK, "index", data)
	})

	e.DELETE("/contacts/:contact_id", func(c echo.Context) error {
		contact_id, err := strconv.Atoi(c.Param("contact_id"))
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/contacts")
		}

		contact.Delete(contact_id)

		if c.Request().Header.Get("HX-Trigger") == "delete-button" {
			return c.Redirect(http.StatusSeeOther, "/contacts")
		}

		return c.String(http.StatusOK, "")
	})

	e.GET("/contacts/archive", func(c echo.Context) error {
		if a.Status == archiver.Running {
			return c.Render(http.StatusOK, "archive-running", a)
		}
		return c.Render(http.StatusOK, "archive-complete", nil)
	})

	e.POST("/contacts/archive", func(c echo.Context) error {
		a.Run()
		return c.Render(http.StatusOK, "archive-running", a)
	})

	e.GET("/contacts/archive/file", func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, a.ArchiveFile())
	})

	e.DELETE("/contacts/archive", func(c echo.Context) error {
		a.Reset()
		return c.Render(http.StatusSeeOther, "archive", nil)
	})

	// Scripting Examples

	e.GET("/counter/js", func(c echo.Context) error {
		return c.Render(http.StatusOK, "js-counter", nil)
	})

	e.GET("/counter/alpine", func(c echo.Context) error {
		return c.Render(http.StatusOK, "alpine-counter", nil)
	})

	e.GET("/counter/hs", func(c echo.Context) error {
		return c.Render(http.StatusOK, "hs-counter", nil)
	})

	// JSON Data APIs
	e.GET("/api/v1/contacts", func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, contact.JsonAllContacts())
	})

	e.POST("/api/v1/contacts", func(c echo.Context) error {
		values := make(map[string]string)

		values["First"] = c.FormValue("first")
		values["Last"] = c.FormValue("last")
		values["Email"] = c.FormValue("email")
		values["Phone"] = c.FormValue("phone")

		errors := contact.ValidateForm(&values)

		if len(errors) != 0 {
			return c.NoContent(http.StatusBadRequest)
		}

		_ = contact.Create(values["First"], values["Last"], values["Email"], values["Phone"])

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
