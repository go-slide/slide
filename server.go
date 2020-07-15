package ferry

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp/fasthttputil"
	"net"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

type Ferry struct {
	config             *Config
	routerMap          map[string][]router
	middleware         []handler
	groupMiddlewareMap map[string][]handler
	urlNotFoundHandler handler
}

// Init server
func InitServer(config *Config) *Ferry {
	return &Ferry{
		config:             config,
		routerMap:          map[string][]router{},
		middleware:         []handler{},
		groupMiddlewareMap: map[string][]handler{},
	}
}

func requestHandler(c *fasthttp.RequestCtx, ferry *Ferry) {
	ctx := getRouterContext(c, ferry)
	appLevelMiddleware(ctx, ferry)
}

func (ferry *Ferry) Listen(host string) error {
	handler := func(c *fasthttp.RequestCtx) {
		requestHandler(c, ferry)
	}
	server := &fasthttp.Server{
		NoDefaultServerHeader: true,
		Handler:               handler,
	}
	return server.ListenAndServe(host)
}

func (ferry *Ferry) addRoute(method, path string, h handler) {
	pathWithRegex := findAndReplace(path)
	ferry.routerMap[method] = append(ferry.routerMap[method], router{
		routerPath: path,
		regexPath:  pathWithRegex,
		handler:    h,
	})
}

// application level middleware
func (ferry *Ferry) Use(h handler) {
	ferry.middleware = append(ferry.middleware, h)
}

// Get method of ferry
func (ferry *Ferry) Get(path string, h handler) {
	ferry.addRoute(get, path, h)
}

// Post method of ferry
func (ferry *Ferry) Post(path string, h handler) {
	ferry.addRoute(post, path, h)
}

// Put method of ferry
func (ferry *Ferry) Put(path string, h handler) {
	ferry.addRoute(put, path, h)
}

// Delete method of ferry
func (ferry *Ferry) Delete(path string, h handler) {
	ferry.addRoute(delete, path, h)
}

// Group method
func (ferry *Ferry) Group(path string) *group {
	return &group{
		path:  path,
		ferry: ferry,
	}
}

// custom 404 handler
func (ferry *Ferry) HandleNotFound(h handler) {
	ferry.urlNotFoundHandler = h
}

// Serving
func (ferry *Ferry) serveFile(path, filePath, contentType string) {
	ferry.Get(path, func(ctx *Ctx) error {
		ctx.RequestCtx.Response.Header.Set("Content-Type", contentType)
		return ctx.RequestCtx.Response.SendFile(filePath)
	})
}

func (ferry *Ferry) ServerDir(path, dir string) {
	var paths []string
	if err := getAllPaths(dir, &paths); err != nil {
		panic(err)
	}
	indexFile := fmt.Sprintf("%s%s", dir, "/index.html")
	indexFileContentType, err := getFileContentType(indexFile)
	if err != nil {
		panic(err)
	}
	ferry.serveFile(path, indexFile, indexFileContentType)
	for _, p := range paths {
		// replace dir name
		fileRoutePath := strings.Replace(p, dir, "", 1)
		contentType, err := getFileContentType(p)
		if err != nil {
			panic(err)
		}
		ferry.serveFile(fileRoutePath, p, contentType)
	}
}

// ServeFile -- serving single file
func (ferry *Ferry) ServeFile(path, filePath string) {
	contentType, err := getFileContentType(filePath)
	if err != nil {
		panic(err)
	}
	ferry.serveFile(path, filePath, contentType)
}

func testServer(req *http.Request, ferry *Ferry) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()
	go func() {
		handler := func(c *fasthttp.RequestCtx) {
			requestHandler(c, ferry)
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
