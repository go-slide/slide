package slide

import "strings"

func appLevelMiddleware(ctx *Ctx, slide *Slide) {
	if len(slide.middleware) > 0 {
		ctx.appMiddlewareIndex = 0
		var next func() error
		next = func() error {
			ctx.appMiddlewareIndex = ctx.appMiddlewareIndex + 1
			if ctx.appMiddlewareIndex != len(slide.middleware) {
				handler := slide.middleware[ctx.appMiddlewareIndex]
				if err := handler(ctx); err != nil {
					handlerRouterError(err, ctx, slide)
				}
			} else {
				// handling request route
				handleRouting(slide, ctx)
			}
			return nil
		}
		handler := slide.middleware[ctx.appMiddlewareIndex]
		ctx.Next = next
		if err := handler(ctx); err != nil {
			handlerRouterError(err, ctx, slide)
		}
	} else {
		handleRouting(slide, ctx)
	}
}

func groupLevelMiddleware(ctx *Ctx, slide *Slide, routers []router) {
	path := string(ctx.RequestCtx.Path())
	// check if path is available in Group middleware
	if len(slide.groupMiddlewareMap) == 0 {
		handleRouter(ctx, slide, routers)
		return
	}
	for groupPath, groupMiddlewares := range slide.groupMiddlewareMap {
		// replace this with wild card
		if strings.Contains(path, groupPath) && len(groupMiddlewares) > 0 {
			var next func() error
			next = func() error {
				ctx.groupMiddlewareIndex = ctx.groupMiddlewareIndex + 1
				if ctx.groupMiddlewareIndex != len(groupMiddlewares) {
					handler := groupMiddlewares[ctx.groupMiddlewareIndex]
					if err := handler(ctx); err != nil {
						handlerRouterError(err, ctx, slide)
					}
				}
				return nil
			}
			ctx.groupMiddlewareIndex = 0
			ctx.Next = next
			handler := groupMiddlewares[ctx.groupMiddlewareIndex]
			if err := handler(ctx); err != nil {
				handlerRouterError(err, ctx, slide)
			}
		}
	}
	handleRouter(ctx, slide, routers)
}
