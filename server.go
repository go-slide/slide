package ferry

import (
	"net/http"
)

type Ferry struct {
	config    *Config
	routerMap map[string][]router
	middleware []handler
}

// Init server
func InitServer(config *Config) *Ferry {
	return &Ferry{
		config:    config,
		routerMap: map[string][]router{},
		middleware: []handler{},
	}
}

func (ferry *Ferry) Listen(host string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := getRouterContext(w, r, ferry)
		appLevelMiddleware(ctx, ferry)
	})
	return http.ListenAndServe(host, nil)
}

func (ferry *Ferry) addRoute(method,path string, h handler) {
	ferry.routerMap[method] = append(ferry.routerMap[method], router{
		path:    path,
		handler: h,
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

// Group method
func (ferry *Ferry) Group (path string) *group {
	return &group{
		path: path,
		ferry: ferry,
	}
}