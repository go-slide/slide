package main

import (
	"ferry"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
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
		return ctx.Send(http.StatusOK, "Hello")
	})

	//app.Use(middleware.Compress())
	//
	//app.Get("/hello", func(ctx *ferry.Ctx) error {
	//	return ctx.Send(http.StatusOK, "Hello")
	//})
	//
	//app.Get("/download", func(ctx *ferry.Ctx) error {
	//	return ctx.SendAttachment("static/js/login.js", "login.js")
	//})
	//
	//app.Post("/upload", func(ctx *ferry.Ctx) error {
	//	return ctx.UploadFile("static/index.html", "index.html")
	//})
	//
	//app.Get("/redirect", func(ctx *ferry.Ctx) error {
	//	return ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000")
	//})
	//
	//app.Post("/login", func(ctx *ferry.Ctx) error {
	//	login := Login{}
	//	err := ctx.Bind(&login)
	//	if err != nil {
	//		return err
	//	}
	//	return ctx.Json(http.StatusOK, map[string]string{
	//		"message": fmt.Sprintf("Welcome %s", login.Username),
	//	})
	//})
	//
	//app.ServerDir("/", "static")
	log.Fatal(app.Listen("localhost:3000"))
}
