package ferry

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"os"
)

type Ctx struct {
	GzipWriter           *gzip.Writer
	RequestCtx           *fasthttp.RequestCtx
	Next                 func() error
	config               *Config
	appMiddlewareIndex   int
	groupMiddlewareIndex int
	routerPath           string
}

// Sending application/json response
func (ctx *Ctx) Json(statusCode int, payload interface{}) error {
	ctx.RequestCtx.Response.Header.Set("Content-Type", "application/json")
	ctx.RequestCtx.SetStatusCode(statusCode)
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
	ctx.RequestCtx.SetBody(response)
	return nil
}

// Sending a text response
func (ctx *Ctx) Send(statusCode int, payload string) error {
	ctx.RequestCtx.SetStatusCode(statusCode)
	ctx.RequestCtx.SetBody([]byte(payload))
	return nil
}

// redirect to new url
// reference https://developer.mozilla.org/en-US/docs/Web/HTTP/Redirections#Temporary_redirections
// status codes between 300-308
func (ctx *Ctx) Redirect(statusCode int, url string) error {
	ctx.RequestCtx.Redirect(url, statusCode)
	return nil
}

// Sending attachment
// filePath
func (ctx *Ctx) SendAttachment(filePath, fileName string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	headerValue := fmt.Sprintf("attachment; filename=%s", fileName)
	ctx.RequestCtx.Response.Header.Set("Content-Disposition", headerValue)
	ctx.RequestCtx.SetBody(content)
	return nil
}

// uploads file to given path
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

// Deserialize body to struct
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

func (ctx *Ctx) GetParam(name string) string {
	return extractParamFromPath(ctx.routerPath, string(ctx.RequestCtx.Path()), name)
}

func (ctx *Ctx) GetParams() map[string]string {
	return getParamsFromPath(ctx.routerPath, string(ctx.RequestCtx.Path()))
}

func getRouterContext(r *fasthttp.RequestCtx, ferry *Ferry) *Ctx {
	return &Ctx{
		RequestCtx: r,
		config:     ferry.config,
	}
}
