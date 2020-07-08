package main

import (
	"ferry"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

func main() {
	validate := validator.New()
	config := ferry.Config{
		Validator: validate,
	}
	app := ferry.InitServer(&config)

	app.Get("/", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "Hello, World!")
	})

	app.Get("/login", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "Hello, World!")
	})

	app.Get("/hey", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "Heey!")
	})

	log.Fatal(app.Listen("localhost:3000"))
}
