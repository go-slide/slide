package ferry

import (
	"fmt"
	"net/http"
	"strings"
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

func (ferry *Ferry) Listen(host string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := getRouterContext(w, r, ferry)
		appLevelMiddleware(ctx, ferry)
	})
	return http.ListenAndServe(host, nil)
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
func (ferry *Ferry) ServeFile(path, fileName string) {
	ferry.Get(path, func(ctx *Ctx) error {
		http.ServeFile(ctx.Writer, ctx.Request, fileName)
		return nil
	})
}

func (ferry *Ferry) ServerDir(path, dir string) {
	var paths []string
	if err := getAllPaths(dir, &paths); err != nil {
		panic(err)
	}
	ferry.Get(path, func(ctx *Ctx) error {
		indexFile := fmt.Sprintf("%s%s", dir, "/index.html")
		http.ServeFile(ctx.Writer, ctx.Request, indexFile)
		return nil
	})
	for _, p := range paths {
		// replace dir name
		filePath := strings.Replace(p, dir, "", 1)
		ferry.ServeFile(filePath, p)
	}
}
