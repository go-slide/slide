package ferry

import (
	"fmt"
	"net/http"
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

func (g *group) addRoute(method,path string, h handler) {
	groupPath := fmt.Sprintf("%s%s", g.path, path)
	g.ferry.routerMap[method] = append(g.ferry.routerMap[method], router{
		path:    groupPath,
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
