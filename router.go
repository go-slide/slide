package ferry

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Middleware/Route Handler
type handler func(ctx *Ctx) error

type router struct {
	path    string
	handler handler
}


type group struct {
	path string
	ferry *Ferry
	middleware []handler
	middlewareCurrentIndex int
}

var (
	get  = "GET"
	post = "POST"
)

var routerRegexReplace = "[a-zA-Z0-9]*"

func (g *group) addRoute(method,path string, h handler) {
	groupPath := fmt.Sprintf("%s%s", g.path, path)
	pathWithRegex := findAndReplace(groupPath)
	g.ferry.routerMap[method] = append(g.ferry.routerMap[method], router{
		path:    pathWithRegex,
		handler: h,
	})
}

// Get method of ferry
func (g *group) Get(path string, h handler) {
	g.addRoute(get, path, h)
}

// Post method of ferry
func (g *group) Post(path string, h handler) {
	g.addRoute(post, path, h)
}

func (g *group) Use (h handler) {
	g.ferry.groupMiddlewareMap[g.path] = append(g.ferry.groupMiddlewareMap[g.path], h)
}

// Group method
func (g *group) Group (path string) *group {
	return &group{
		path: fmt.Sprintf("%s%s", g.path, path),
		ferry: g.ferry,
		middleware: []handler{},
	}
}

// handler 404
func handle404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(w, "Check URL")
}

func handlerRouterError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = fmt.Fprint(w, err.Error())
}


func handleRouting(ferry *Ferry, ctx *Ctx) {
	// first get handler by method
	routesByMethod := ferry.routerMap[ctx.Request.Method]
	if routesByMethod != nil {
		groupLevelMiddleware(ctx, ferry, routesByMethod)
	} else {
		// run 404
		handle404(ctx.Writer)
	}

}

// Finds wild card in URL and replace them with a regex for,
// ex if path is /auth/:name -> /auth/[a-zA-Z0-9]*
// ex if path is /auth/name -> /auth/name
func findAndReplace(path string) string {
	if !strings.Contains(path, ":") {
		return fmt.Sprintf("%s%s%s", "^", path, "$")
	}
	result := ""
	slitted := strings.Split(path, "/")
	for _, v := range slitted {
		if v == "" {
			continue
		}
		if strings.Contains(v, ":") {
			result = fmt.Sprintf("%s/%s", result, routerRegexReplace)
			continue
		}
		result = fmt.Sprintf("%s/%s", result, v)
	}
	// replace slashes
	result = strings.ReplaceAll(result, "/", "\\/")
	result = fmt.Sprintf("%s%s%s", "^", result, "$")
	return result
}


// calls actual handler
func handleRouter(ctx *Ctx, ferry *Ferry, routers []router)  {
	urlPath := ctx.Request.URL.Path
	for _, route := range routers {
		match, _ := regexp.MatchString(route.path, urlPath)
		if match {
			if err := route.handler(ctx); err != nil {
				handlerRouterError(err, ctx.Writer)
			}
			return
		}
	}
	handle404(ctx.Writer)
}