package ferry

import (
	"net/http"
)

type Ferry struct {
	config *Config
	routerMap map[string][]router
}

// Init server
func InitServer(config *Config) *Ferry {
	return &Ferry{
		config:    config,
		routerMap: map[string][]router{},
	}
}

func (ferry *Ferry) Listen(host string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		handleRouting(ferry, writer, request)
	})
	return http.ListenAndServe(host, nil)
}

// Get method of ferry
func (ferry *Ferry) Get(path string, h handler) {
	ferry.routerMap[get] = append(ferry.routerMap[get], router{
		path: path,
		handler: h,
	})
}

// Get method of ferry
func (ferry *Ferry) Post(path string, h handler) {
	ferry.routerMap[post] = append(ferry.routerMap[post], router{
		path: path,
		handler: h,
	})
}
