package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ferry-go/ferry/middleware"

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

	app.HandleNotFound(func(ctx *ferry.Ctx) error {
		return ctx.Json(http.StatusNotFound, "check url idiot")
	})
	app.HandleErrors(func(ctx *ferry.Ctx, err error) error {
		return ctx.Send(http.StatusInternalServerError, err.Error())
	})

	// compression middleware
	app.Use(middleware.Compress())

	// you can multiple middlewares also
	app.Use(func(ctx *ferry.Ctx) error {
		fmt.Println("this will run for all URL(s)")
		return ctx.Next()
	})

	app.Get("/", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "Hello, World")
	})

	// redirect to new url
	app.Get("/redirect", func(ctx *ferry.Ctx) error {
		return ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/static")
	})

	app.Get("/name/:name", func(ctx *ferry.Ctx) error {
		name := ctx.GetParam("name")
		return ctx.Send(http.StatusOK, fmt.Sprintf("hello, %s", name))
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

	app.Get("/hello", func(ctx *ferry.Ctx) error {
		params := ctx.GetQueryParams()
		return ctx.Json(http.StatusOK, params)
	})

	app.Get("/hello/single", func(ctx *ferry.Ctx) error {
		params := ctx.GetQueryParam("key")
		return ctx.Send(http.StatusOK, fmt.Sprintf("key is %s", params))
	})

	// Grouping your route
	auth := app.Group("/auth")
	// you can multiple middlewares also
	auth.Use(func(ctx *ferry.Ctx) error {
		fmt.Println("this will run for all urls with /auth")
		return ctx.Next()
	})
	auth.Get("/login", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "Hello, World")
	})

	// path and dir name
	app.ServerDir("/static", "static")

	// single file
	app.ServeFile("/js", "static/login.js")

	// Downloading file
	app.Get("/download", func(ctx *ferry.Ctx) error {
		return ctx.SendAttachment("static/login.js", "login.js")
	})

	app.Get("/servefile", func(ctx *ferry.Ctx) error {
		return ctx.ServeFile("static/index.html")
	})

	// uploading file
	app.Post("/upload", func(ctx *ferry.Ctx) error {
		return ctx.UploadFile("static/login.js", "login.js")
	})

	// you can have router level middleware which works in reverse way
	app.Get("/routermiddleware", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "hola!")
	}, func(ctx *ferry.Ctx) error {
		fmt.Println("this prints second")
		return ctx.Next()
	}, func(ctx *ferry.Ctx) error {
		fmt.Println("this prints first")
		return ctx.Next()
	})

	log.Fatal(app.Listen("localhost:3000"))
}
