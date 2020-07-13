package ferry

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Ctx struct {
	Writer               http.ResponseWriter
	GzipWriter           *gzip.Writer
	Request              *http.Request
	Context              context.Context
	Next                 func() error
	config               *Config
	appMiddlewareIndex   int
	groupMiddlewareIndex int
	routerPath           string
}

// Sending application/json response
func (ctx *Ctx) Json(statusCode int, payload interface{}) error {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(statusCode)
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if ctx.GzipWriter != nil {
		_, err = ctx.GzipWriter.Write(response)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = ctx.Writer.Write(response)
	return err
}

// Sending a text response
func (ctx *Ctx) Send(statusCode int, payload string) error {
	ctx.Writer.WriteHeader(statusCode)
	_, err := ctx.Writer.Write([]byte(payload))
	return err
}

// redirect to new urls
// reference https://developer.mozilla.org/en-US/docs/Web/HTTP/Redirections#Temporary_redirections
// status codes between 300-308
func (ctx *Ctx) Redirect(statusCode int, url string) error {
	http.Redirect(ctx.Writer, ctx.Request, url, statusCode)
	return nil
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

func (ctx *Ctx) GetParam(name string) string {
	return extractParamFromPath(ctx.routerPath, ctx.Request.URL.Path, name)
}

func (ctx *Ctx) GetParams() map[string]string {
	return getParamsFromPath(ctx.routerPath, ctx.Request.URL.Path)
}

func getRouterContext(w http.ResponseWriter, r *http.Request, ferry *Ferry) *Ctx {
	return &Ctx{
		Writer:  w,
		Request: r,
		Context: r.Context(),
		config:  ferry.config,
	}
}
