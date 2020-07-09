# ferry-go

## Example

```go
package main

import (
	"ferry"
	"fmt"
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

	app.Use(func(ctx *ferry.Ctx) error {
		fmt.Println("hey!, this is middleware")
		return ctx.Next()
	})
	// uncomment below code to see early response from middleware
	//app.Use(func(ctx *ferry.Ctx) error {
	//	fmt.Println("Early response from middleware")
	//	return ctx.Send(http.StatusOK, "From app middleware")
	//})

	app.Get("/", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "Hello, World!")
	})

	app.Get("/login", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "Hello, World!, This is login")
	})

	app.Get("/hey", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "Heey!")
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

	// group routing
	auth := app.Group("/auth")
	auth.Get("/signup", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "Registered")
	})

	// you can make nested groups as well
	auth.Group("/doubleauth").Get("/signup", func(ctx *ferry.Ctx) error {
		return ctx.Send(http.StatusOK, "double")
	})

	log.Fatal(app.Listen("localhost:3000"))
}

```
