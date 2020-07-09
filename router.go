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

// Group method
func (g *group) Group (path string) *group {
	return &group{
		path: fmt.Sprintf("%s%s", g.path, path),
		ferry: g.ferry,
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

func appLevelMiddleware(ctx *Ctx, ferry *Ferry) {
	if len(ferry.middleware) > 0 {
		ctx.appMiddlewareIndex = 0
		var next func() error
		next = func() error {
			ctx.appMiddlewareIndex = ctx.appMiddlewareIndex + 1
			if ctx.appMiddlewareIndex != len(ferry.middleware) {
				handler := ferry.middleware[ctx.appMiddlewareIndex]
				if err := handler(ctx); err != nil {
					handlerRouterError(err, ctx.Writer)
				}
			} else {
				handleRouting(ferry, ctx)
			}
			return nil
		}
		handler := ferry.middleware[ctx.appMiddlewareIndex]
		ctx.Next = next
		if err := handler(ctx); err != nil {
			handlerRouterError(err, ctx.Writer)
		}
	}
}

func handleRouting(ferry *Ferry, ctx *Ctx) {
	// first get handler by method
	routesByMethod := ferry.routerMap[ctx.Request.Method]
	if routesByMethod != nil {
		// get handler by path
		for _, route := range routesByMethod {
			if route.path == ctx.Request.URL.Path {
				if err := route.handler(ctx); err != nil {
					handlerRouterError(err, ctx.Writer)
				}
				return
			}
		}
	}
	// run 404
	handle404(ctx.Writer)
}
