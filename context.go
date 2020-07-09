package ferry

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Ctx struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Context context.Context
	Next func() error
	config  *Config
	appMiddlewareIndex int
	groupMiddlewareIndex int
}

// Sending application/json response
func (ctx *Ctx) Json(statusCode int, payload interface{}) error {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(statusCode)
	lol, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(ctx.Writer, string(lol))
	return err
}

// Sending a text response
func (ctx *Ctx) Send(statusCode int, payload string) error {
	ctx.Writer.WriteHeader(statusCode)
	_, err := fmt.Fprint(ctx.Writer, payload)
	return err
}

// Deserialize body to struct
func (ctx *Ctx) Bind(input interface{}) error {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, input); err != nil {
		return err
	}
	if ctx.config.Validator != nil {
		if err := ctx.config.Validator.Struct(input); err != nil {
			return err
		}
	}
	return nil
}

func getRouterContext(w http.ResponseWriter, r *http.Request, ferry *Ferry) *Ctx {
	return &Ctx{
		Writer:  w,
		Request: r,
		Context: r.Context(),
		config:  ferry.config,
	}
}
