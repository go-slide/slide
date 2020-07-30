package slide

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
	slide                  *Slide
	middleware             []handler
	middlewareCurrentIndex int
}

func (g *Group) addRoute(method, path string, h ...handler) {
	groupPath := fmt.Sprintf("%s%s", g.path, path)
	pathWithRegex := findAndReplace(groupPath)
	g.slide.routerMap[method] = append(g.slide.routerMap[method], router{
		routerPath: groupPath,
		regexPath:  pathWithRegex,
		handlers:   h,
	})
}

// Get method of slide
func (g *Group) Get(path string, h handler) {
	g.addRoute(GET, path, h)
}

// Post method of slide
func (g *Group) Post(path string, h handler) {
	g.addRoute(POST, path, h)
}

// Put method of slide
func (g *Group) Put(path string, h handler) {
	g.addRoute(PUT, path, h)
}

// Delete method of slide
func (g *Group) Delete(path string, h handler) {
	g.addRoute(DELETE, path, h)
}

// Use -- group level middleware
func (g *Group) Use(h handler) {
	g.slide.groupMiddlewareMap[g.path] = append(g.slide.groupMiddlewareMap[g.path], h)
}

// Group method
func (g *Group) Group(path string) *Group {
	return &Group{
		path:       fmt.Sprintf("%s%s", g.path, path),
		slide:      g.slide,
		middleware: []handler{},
	}
}

// handler 404
func handle404(slide *Slide, ctx *Ctx) {
	if slide.urlNotFoundHandler != nil {
		if err := slide.urlNotFoundHandler(ctx); err != nil {
			handlerRouterError(err, ctx, slide)
		}
		return
	}
	ctx.RequestCtx.Response.SetStatusCode(http.StatusNotFound)
	ctx.RequestCtx.Response.SetBody([]byte(NotFoundMessage))
}

func handlerRouterError(err error, ctx *Ctx, slide *Slide) {
	if slide.errorHandler != nil {
		if handlerError := slide.errorHandler(ctx, err); handlerError != nil {
			ctx.RequestCtx.Response.SetStatusCode(http.StatusInternalServerError)
			ctx.RequestCtx.Response.SetBody([]byte(handlerError.Error()))
		}
		return
	}
	ctx.RequestCtx.Response.SetStatusCode(http.StatusInternalServerError)
	ctx.RequestCtx.Response.SetBody([]byte(err.Error()))
}

func handleRouting(slide *Slide, ctx *Ctx) {
	// first GET handler by method
	routesByMethod := slide.routerMap[string(ctx.RequestCtx.Method())]
	if routesByMethod != nil {
		groupLevelMiddleware(ctx, slide, routesByMethod)
	} else {
		// run 404
		handle404(slide, ctx)
	}

}

// calls actual handler
func handleRouter(ctx *Ctx, slide *Slide, routers []router) {
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
					handlerRouterError(err, ctx, slide)
				}
			}
			return nil
		}
		ctx.Next = next
		if err := route.handlers[len(route.handlers)-1-index](ctx); err != nil {
			handlerRouterError(err, ctx, slide)
		}
	} else {
		handle404(slide, ctx)
	}
}
