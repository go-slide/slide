# slide-go

## Installation
```cmd
go get -u github.com/go-slide/slide

```
**Not Production Ready**

[![codecov](https://codecov.io/gh/go-slide/slide/branch/master/graph/badge.svg)](https://codecov.io/gh/go-slide/slide)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-slide/slide)](https://goreportcard.com/report/github.com/go-slide/slide)

## Example

```go
package main

import (
	"fmt"
	"github.com/slide-go/slide/middleware"
	"log"
	"net/http"

	"github.com/slide-go/slide"

	"github.com/go-playground/validator/v10"
)

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}


func main() {
	validate := validator.New()
	config := slide.Config{
		Validator: validate,
	}

	app := slide.InitServer(&config)

	// compression middleware
	app.Use(middleware.Compress())

	// you can multiple middlewares also
	app.Use(func(ctx *slide.Ctx) error {
		fmt.Println("this will run for all URL(s)")
		return ctx.Next()
	})

	app.Get("/", func(ctx *slide.Ctx) error {
		return ctx.Send(http.StatusOK, "Hello, World")
	})

	// redirect to new url
	app.Get("/redirect", func(ctx *slide.Ctx) error {
		return ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/static")
	})

	app.Get("/name/:name", func(ctx *slide.Ctx) error {
		name := ctx.GetParam("name")
		return ctx.Send(http.StatusOK, fmt.Sprintf("hello, %s", name))
	})

	app.Post("/login", func(ctx *slide.Ctx) error {
		login := Login{}
		err := ctx.Bind(&login)
		if err != nil {
			return err
		}
		return ctx.Json(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Welcome %s", login.Username),
		})
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
	app.ServerDir("/static", "static")

	// single file
	app.ServeFile("/js", "static/login.js")

	// Downloading file
	app.Get("/download", func(ctx *slide.Ctx) error {
		return ctx.SendAttachment("static/login.js", "login.js")
	})

	// uploading file
	app.Post("/upload", func(ctx *slide.Ctx) error {
		return ctx.UploadFile("static/login.js", "login.js")
	})

	log.Fatal(app.Listen("localhost:3000"))
}
```

## Benchmarks
![](https://i.ibb.co/TWdgzB8/slide-benchmark.png)
