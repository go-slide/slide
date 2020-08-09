package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-slide/slide/middleware"

	"github.com/go-slide/slide"

	"github.com/go-playground/validator/v10"
)

type login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func main() {
	validate := validator.New()
	config := slide.Config{
		Validator: validate,
	}

	app := slide.InitServer(&config)

	// this is with config
	app.Use(middleware.CorsWithConfig(middleware.CorsConfig{
		AllowOrigins: []string{"https://www.postgresqltutorial.com"},
	}))

	// without config, which uses default config
	// app.Use(middleware.Cors())

	app.HandleNotFound(func(ctx *slide.Ctx) error {
		return ctx.JSON(http.StatusNotFound, "check url idiot")
	})
	app.HandleErrors(func(ctx *slide.Ctx, err error) error {
		return ctx.Send(http.StatusInternalServerError, err.Error())
	})

	// compression middleware
	app.Use(middleware.Compress())

	// you can multiple middlewares also
	app.Use(func(ctx *slide.Ctx) error {
		fmt.Println("this will run for all URL(s)", string(ctx.RequestCtx.Path()))
		return ctx.Next()
	})

	app.Get("/", func(ctx *slide.Ctx) error {
		return ctx.SendStatusCode(http.StatusOK)
	})

	// redirect to new url
	app.Get("/redirect", func(ctx *slide.Ctx) error {
		return ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/static")
	})

	app.Get("/name/:name", func(ctx *slide.Ctx) error {
		name := ctx.GetParam("name")
		key := ctx.GetQueryParam("ss")
		return ctx.Send(http.StatusOK, fmt.Sprintf("hello, %s %s", name, key))
	})

	app.Post("/login", func(ctx *slide.Ctx) error {
		login := login{}
		err := ctx.Bind(&login)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Welcome %s", login.Username),
		})
	})

	app.Get("/hello", func(ctx *slide.Ctx) error {
		params := ctx.GetQueryParams()
		return ctx.JSON(http.StatusOK, params)
	})

	app.Get("/hello/single", func(ctx *slide.Ctx) error {
		params := ctx.GetQueryParam("key")
		return ctx.Send(http.StatusOK, fmt.Sprintf("key is %s", params))
	})

	// Grouping your route
	auth := app.Group("/auth")
	// you can multiple middlewares also
	auth.Use(func(ctx *slide.Ctx) error {
		fmt.Println("this will run for all urls with /auth")
		return ctx.Next()
	})
	auth.Get("/login", func(ctx *slide.Ctx) error {
		return ctx.Send(http.StatusOK, "Hello, World")
	})

	// path and dir name
	// app.ServerDir("/static", "static")

	// single file
	// app.ServeFile("/js", "static/login.js")

	// Downloading file
	app.Get("/download", func(ctx *slide.Ctx) error {
		return ctx.SendAttachment("static/login.js", "login.js")
	})

	app.Get("/servefile", func(ctx *slide.Ctx) error {
		return ctx.ServeFile("static/index.html")
	})

	// uploading file
	app.Post("/upload", func(ctx *slide.Ctx) error {
		return ctx.UploadFile("static/login.js", "login.js")
	})

	// you can have router level middleware which works in reverse way
	app.Get("/routermiddleware", func(ctx *slide.Ctx) error {
		return ctx.Send(http.StatusOK, "hola!")
	}, func(ctx *slide.Ctx) error {
		fmt.Println("this prints second", ctx.RequestCtx.UserValue("lol"))
		return ctx.Next()
	}, func(ctx *slide.Ctx) error {
		fmt.Println("this prints first")
		return ctx.Next()
	})

	log.Fatal(app.Listen("localhost:3000"))
}
