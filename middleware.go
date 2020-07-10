package ferry

import "strings"

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
				// handling request route
				handleRouting(ferry, ctx)
			}
			return nil
		}
		handler := ferry.middleware[ctx.appMiddlewareIndex]
		ctx.Next = next
		if err := handler(ctx); err != nil {
			handlerRouterError(err, ctx.Writer)
		}
	} else {
		handleRouting(ferry, ctx)
	}
}

func groupLevelMiddleware(ctx *Ctx, ferry *Ferry, routers []router) {
	path := ctx.Request.URL.Path
	// check if path is available in group middleware
	if len(ferry.groupMiddlewareMap) == 0 {
		handleRouter(ctx, ferry, routers)
		return
	}
	for groupPath, groupMiddlewares := range ferry.groupMiddlewareMap {
		// replace this with wild card
		if strings.Contains(path, groupPath) && len(groupMiddlewares) > 0 {
			var next func() error
			next = func() error {
				ctx.groupMiddlewareIndex = ctx.groupMiddlewareIndex + 1
				if ctx.groupMiddlewareIndex != len(groupMiddlewares) {
					handler := groupMiddlewares[ctx.groupMiddlewareIndex]
					if err := handler(ctx); err != nil {
						handlerRouterError(err, ctx.Writer)
					}
				}
				return nil
			}
			ctx.groupMiddlewareIndex = 0
			ctx.Next = next
			handler := groupMiddlewares[ctx.groupMiddlewareIndex]
			if err := handler(ctx); err != nil {
				handlerRouterError(err, ctx.Writer)
			}
		}
	}
	handleRouter(ctx, ferry, routers)
}
