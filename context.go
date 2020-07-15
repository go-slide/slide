package ferry

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/valyala/fasthttp"
)

type Ctx struct {
	RequestCtx           *fasthttp.RequestCtx
	Next                 func() error
	config               *Config
	appMiddlewareIndex   int
	groupMiddlewareIndex int
	routerPath           string
	queryPath            string
}

// Json Sending application/json response
func (ctx *Ctx) Json(statusCode int, payload interface{}) error {
	ctx.RequestCtx.Response.Header.Set("Content-Type", "application/json")
	ctx.RequestCtx.SetStatusCode(statusCode)
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	ctx.RequestCtx.SetBody(response)
	return nil
}

// Send Sending a text response
func (ctx *Ctx) Send(statusCode int, payload string) error {
	ctx.RequestCtx.SetStatusCode(statusCode)
	ctx.RequestCtx.SetBody([]byte(payload))
	return nil
}

// Redirect to new url
// reference https://developer.mozilla.org/en-US/docs/Web/HTTP/Redirections#Temporary_redirections
// status codes between 300-308
func (ctx *Ctx) Redirect(statusCode int, url string) error {
	ctx.RequestCtx.Redirect(url, statusCode)
	return nil
}

// SendAttachment Sending attachment
func (ctx *Ctx) SendAttachment(filePath, fileName string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	contentType, err := getFileContentType(filePath)
	if err != nil {
		panic(err)
	}
	ctx.RequestCtx.Response.Header.Set("Content-Type", contentType)
	headerValue := fmt.Sprintf("attachment; filename=%s", fileName)
	ctx.RequestCtx.Response.Header.Set("Content-Disposition", headerValue)
	ctx.RequestCtx.SetBody(content)
	return nil
}

// UploadFile uploads file to given path
func (ctx *Ctx) UploadFile(filePath, fileName string) error {
	form, err := ctx.RequestCtx.FormFile(fileName)
	if err != nil {
		return err
	}
	file, err := form.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Bind Deserialize body to struct
func (ctx *Ctx) Bind(input interface{}) error {
	body := ctx.RequestCtx.Request.Body()
	if err := json.Unmarshal(body, input); err != nil {
		return err
	}
	if ctx.config.Validator != nil {
		if err := ctx.config.Validator.Struct(input); err != nil {
			return err
		}
	}
	return nil
}

// GetParam - Getting path param
//
// /name/:name
//
// /name/ferry
//
// name := ctx.GetParam("name")
//
// name = ferry
//
func (ctx *Ctx) GetParam(name string) string {
	return extractParamFromPath(ctx.routerPath, string(ctx.RequestCtx.Path()), name)
}

// GetParams returns map of path params
//
// routerPath /auth/:name/:age
//
// requestPath /auth/madhuri/32
//
// returns { name: madhuri, age: 32 }
//
func (ctx *Ctx) GetParams() map[string]string {
	return getParamsFromPath(ctx.routerPath, string(ctx.RequestCtx.Path()))
}

//	returns value of a single query Param
//
//	route path /hello?key=test&value=bbp
//
//	keyValue = GetQueryParam(key)
//
//	keyValue = test

func (ctx *Ctx) GetQueryParam(name string) string {
	return getQueryParam(ctx.queryPath, name)
}

//	returns map of query Params
//
//	route path /hello?key=test&value=bbp
//
//	returns {key : test, value : bbp}

func (ctx *Ctx) GetQueryParams() map[string]string {
	return getAllQueryParams(ctx.queryPath)
}

// ServeFile serving file as response
func (ctx *Ctx) ServeFile(filePath string) error {
	contentType, err := getFileContentType(filePath)
	if err != nil {
		return err
	}
	ctx.RequestCtx.Response.Header.Set("Content-Type", contentType)
	return ctx.RequestCtx.Response.SendFile(filePath)
}

func getRouterContext(r *fasthttp.RequestCtx, ferry *Ferry) *Ctx {
	return &Ctx{
		RequestCtx: r,
		config:     ferry.config,
	}
}
