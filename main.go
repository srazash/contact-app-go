package main

import (
	"contactapp/controllers/counter"
	"contactapp/models/contact"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	contact.Load()
	counter.Load()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Static("/", "static")

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/contacts")
	})

	e.GET("/contacts", func(c echo.Context) error {
		return c.String(http.StatusFound, "I live!")
	})

	e.Logger.Fatal(e.Start(":3000"))
}
