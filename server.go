package slide

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/valyala/fasthttp"
)

// Slide -- Slide config
type Slide struct {
	config             *Config
	routerMap          map[string][]router
	middleware         []handler
	groupMiddlewareMap map[string][]handler
	urlNotFoundHandler handler
	errorHandler       errHandler
}

// InitServer -- initializing server with slide config
func InitServer(config *Config) *Slide {
	return &Slide{
		config:             config,
		routerMap:          map[string][]router{},
		middleware:         []handler{},
		groupMiddlewareMap: map[string][]handler{},
	}
}

func requestHandler(c *fasthttp.RequestCtx, slide *Slide) {
	ctx := getRouterContext(c, slide)
	appLevelMiddleware(ctx, slide)
}

// Listen -- starting server with given host
func (slide *Slide) Listen(host string) error {
	handler := func(c *fasthttp.RequestCtx) {
		requestHandler(c, slide)
	}
	server := &fasthttp.Server{
		NoDefaultServerHeader: true,
		Handler:               handler,
		ErrorHandler: func(r *fasthttp.RequestCtx, err error) {
			if slide.errorHandler != nil {
				ctx := getRouterContext(r, slide)
				_ = slide.errorHandler(ctx, err)
			} else {
				// TODO replace it with logger
				fmt.Println(err.Error())
			}

		},
	}
	return server.ListenAndServe(host)
}

func (slide *Slide) addRoute(method, path string, h []handler) {
	pathWithRegex := findAndReplace(path)
	slide.routerMap[method] = append(slide.routerMap[method], router{
		routerPath: path,
		regexPath:  pathWithRegex,
		handlers:   h,
	})
}

// Use -- application level middleware
func (slide *Slide) Use(h handler) {
	slide.middleware = append(slide.middleware, h)
}

// Get method of slide
func (slide *Slide) Get(path string, h ...handler) {
	slide.addRoute(GET, path, h)
}

// Post method of slide
func (slide *Slide) Post(path string, h ...handler) {
	slide.addRoute(POST, path, h)
}

// Put method of slide
func (slide *Slide) Put(path string, h ...handler) {
	slide.addRoute(PUT, path, h)
}

// Delete method of slide
func (slide *Slide) Delete(path string, h ...handler) {
	slide.addRoute(DELETE, path, h)
}

// Group method
func (slide *Slide) Group(path string) *Group {
	return &Group{
		path:  path,
		slide: slide,
	}
}

// HandleNotFound custom 404 handler
func (slide *Slide) HandleNotFound(h handler) {
	slide.urlNotFoundHandler = h
}

// HandleErrors Handling errors at application level
func (slide *Slide) HandleErrors(h errHandler) {
	slide.errorHandler = h
}

// Serving
func (slide *Slide) serveFile(path, filePath, contentType string) {
	slide.Get(path, func(ctx *Ctx) error {
		ctx.RequestCtx.Response.Header.Set(ContentType, contentType)
		return ctx.RequestCtx.Response.SendFile(filePath)
	})
}

// ServerDir files as routes
func (slide *Slide) ServerDir(path, dir string) {
	var paths []string
	if err := getAllPaths(dir, &paths); err != nil {
		panic(err)
	}
	indexFile := fmt.Sprintf("%s%s", dir, "/index.html")
	indexFileContentType, err := getFileContentType(indexFile)
	if err != nil {
		panic(err)
	}
	slide.serveFile(path, indexFile, indexFileContentType)
	for _, p := range paths {
		// replace dir name
		fileRoutePath := strings.Replace(p, dir, "", 1)
		contentType, err := getFileContentType(p)
		if err != nil {
			panic(err)
		}
		slide.serveFile(fileRoutePath, p, contentType)
	}
}

// ServeFile -- serving single file
func (slide *Slide) ServeFile(path, filePath string) {
	contentType, err := getFileContentType(filePath)
	if err != nil {
		panic(err)
	}
	slide.serveFile(path, filePath, contentType)
}

func testServer(req *http.Request, slide *Slide) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()
	go func() {
		handler := func(c *fasthttp.RequestCtx) {
			requestHandler(c, slide)
		}
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()
	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}
	return client.Do(req)
}
