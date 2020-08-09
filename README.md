# Slide, a Go web framework for Building API(s)

[![codecov](https://codecov.io/gh/go-slide/slide/branch/master/graph/badge.svg)](https://codecov.io/gh/go-slide/slide)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-slide/slide)](https://goreportcard.com/report/github.com/go-slide/slide)

###### tags: `Go` `Express`

> Slide is one of the fastest frameworks, built on top of fasthttp.
> People coming form express will feel at home. 


## Motivation
While playing around with the Go’s net/http, one thing that we missed most was the lack of middleware support and the problems it solves. After a little reading, we decided to write our own web framework which would solve the issues like middleware support with next, handle a wide range of files, Upload, Download, etc.

**:bulb: **Note:** We are still in experimental stage, would love to hear feedback.**


### Installation
```cmd
go get -u github.com/go-slide/slide
```

### :rocket:  Example

For more API information check [Docs](https://goslide-framework.gitbook.io/slide/)

```go
package main

import (
	"log"
	"github.com/go-slide/slide"
	"github.com/go-playground/validator/v10"
)

func main() {
    validate := validator.New()
	config := slide.Config{
		Validator: validate,
	}
    app := slide.InitServer(&config)
    app.Get("/", func(ctx *slide.Ctx) error {
        return ctx.Send(http.StatusOK, "Hello, World")
    })
    log.Fatal(app.Listen("localhost:4321"))
}
```

## Routing

Slide supports multilevel routing.

```go
app := slide.InitServer(&config)

app.Get("/", func(ctx *slide.Ctx) error {
    return ctx.Send(http.StatusOK, "Hello, World")
})

// Grouping your route
auth := app.Group("/auth")
auth.Get("/login", func(ctx *slide.Ctx) error {
    return ctx.Send(http.StatusOK, "Hello, World")
})

```

## Middleware
Slide supports wide range of middlewares. 
1. Application Level
2. Group Level
3. Route Level

```go
app := slide.InitServer(&config)

## Application level
// you can multiple middlewares also
app.Use(func(ctx *slide.Ctx) error {
    fmt.Println("this will run for all URL(s)")
    return ctx.Next()
})

//Group Level

auth := app.Group("/auth")
auth.Use(func(ctx *slide.Ctx) error {
    fmt.Println("this will run for all /auth URL(s)")
    return ctx.Next()
})
auth.Get("/login", func(ctx *slide.Ctx) error {
    return ctx.Send(http.StatusOK, "Hello, World")
})

// Route level
// you can have router level middleware which works in Right -> Left or Bottom to Top
app.Get("/routermiddleware", func(ctx *slide.Ctx) error {
    return ctx.Send(http.StatusOK, "hola!")
}, func(ctx *slide.Ctx) error {
    fmt.Println("this prints second", ctx.RequestCtx.UserValue("lol"))
    return ctx.Next()
}, func(ctx *slide.Ctx) error {
    fmt.Println("this prints first")
    return ctx.Next()
})

```


## Benchmark

```cmd
autocannon -c 100 -d 40 -p http://localhost:4321/
```

| Framework | No of requests |
| -------- | -------- |
| Slide     | 2765K     |




**:computer: Core Contributors**
[Sai Umesh](https://twitter.com/saiumesh)
[Madhuri](https://twitter.com/pittalamadhuri)
