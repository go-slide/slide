package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ferry-go/ferry"

	"github.com/go-playground/validator/v10"
)

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}


func main() {
	validate := validator.New()
	config := ferry.Config{
		Validator: validate,
	}
	app := ferry.InitServer(&config)
	app.Get("/", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "Hello, World")
	})

	app.Post("/login", func(ctx *ferry.Ctx) error {
		login := Login{}
		err := ctx.Bind(&login)
		if err != nil {
			return err
		}
		return ctx.Json(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Welcome %s", login.Username),
		})
	})

	log.Fatal(app.Listen("localhost:4321"))
}
