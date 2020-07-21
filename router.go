package ferry

import (
	"fmt"
	"net/http"
	"regexp"
)

// Middleware/Route Handler
type handler func(ctx *Ctx) error

// error handler
type errHandler func(ctx *Ctx, err error) error

type router struct {
	routerPath string
	regexPath  string
	handlers   []handler
}

// Group -- router group
type Group struct {
	path                   string
	ferry                  *Ferry
	middleware             []handler
	middlewareCurrentIndex int
}

func (g *Group) addRoute(method, path string, h ...handler) {
	groupPath := fmt.Sprintf("%s%s", g.path, path)
	pathWithRegex := findAndReplace(groupPath)
	g.ferry.routerMap[method] = append(g.ferry.routerMap[method], router{
		routerPath: groupPath,
		regexPath:  pathWithRegex,
		handlers:   h,
	})
}

// Get method of ferry
func (g *Group) Get(path string, h handler) {
	g.addRoute(GET, path, h)
}

// Post method of ferry
func (g *Group) Post(path string, h handler) {
	g.addRoute(POST, path, h)
}

// Put method of ferry
func (g *Group) Put(path string, h handler) {
	g.addRoute(PUT, path, h)
}

// Delete method of ferry
func (g *Group) Delete(path string, h handler) {
	g.addRoute(DELETE, path, h)
}

// Use -- group level middleware
func (g *Group) Use(h handler) {
	g.ferry.groupMiddlewareMap[g.path] = append(g.ferry.groupMiddlewareMap[g.path], h)
}

// Group method
func (g *Group) Group(path string) *Group {
	return &Group{
		path:       fmt.Sprintf("%s%s", g.path, path),
		ferry:      g.ferry,
		middleware: []handler{},
	}
}

// handler 404
func handle404(ferry *Ferry, ctx *Ctx) {
	if ferry.urlNotFoundHandler != nil {
		if err := ferry.urlNotFoundHandler(ctx); err != nil {
			handlerRouterError(err, ctx, ferry)
		}
		return
	}
	ctx.RequestCtx.Response.SetStatusCode(http.StatusNotFound)
	ctx.RequestCtx.Response.SetBody([]byte(NotFoundMessage))
}

func handlerRouterError(err error, ctx *Ctx, ferry *Ferry) {
	if ferry.errorHandler != nil {
		if handlerError := ferry.errorHandler(ctx, err); handlerError != nil {
			ctx.RequestCtx.Response.SetStatusCode(http.StatusInternalServerError)
			ctx.RequestCtx.Response.SetBody([]byte(handlerError.Error()))
		}
		return
	}
	ctx.RequestCtx.Response.SetStatusCode(http.StatusInternalServerError)
	ctx.RequestCtx.Response.SetBody([]byte(err.Error()))
}

func handleRouting(ferry *Ferry, ctx *Ctx) {
	// first GET handler by method
	routesByMethod := ferry.routerMap[string(ctx.RequestCtx.Method())]
	if routesByMethod != nil {
		groupLevelMiddleware(ctx, ferry, routesByMethod)
	} else {
		// run 404
		handle404(ferry, ctx)
	}

}

// calls actual handler
func handleRouter(ctx *Ctx, ferry *Ferry, routers []router) {
	urlPath := string(ctx.RequestCtx.Path())
	query := ctx.RequestCtx.QueryArgs()
	var route *router
	for _, r := range routers {
		match, _ := regexp.MatchString(r.regexPath, urlPath)
		if match {
			route = &r
			break
		}
	}
	if route != nil {
		ctx.routerPath = route.routerPath
		ctx.queryPath = query.String()
		index := 0
		var next func() error
		next = func() error {
			index = index + 1
			if index <= len(route.handlers)-1 {
				if err := route.handlers[len(route.handlers)-1-index](ctx); err != nil {
					handlerRouterError(err, ctx, ferry)
				}
			}
			return nil
		}
		ctx.Next = next
		if err := route.handlers[len(route.handlers)-1-index](ctx); err != nil {
			handlerRouterError(err, ctx, ferry)
		}
	} else {
		handle404(ferry, ctx)
	}
}
