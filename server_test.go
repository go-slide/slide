package ferry

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	Ferry *Ferry
}

type testRoutes struct {
	path   string
	method string
}

func (suite *ServerSuite) SetupTest() {
	config := &Config{}
	app := InitServer(config)
	suite.Ferry = app
}

func (suite *ServerSuite) TestGetMethod() {
	routes := []testRoutes{
		{
			path:   "/hey",
			method: GET,
		},
		{
			path:   "/hey/:name",
			method: POST,
		},
		{
			path:   "/hey/:name",
			method: PUT,
		},
		{
			path:   "/hey/:name",
			method: DELETE,
		},
	}
	for _, testRoute := range routes {
		if testRoute.method == GET {
			suite.Ferry.Get(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
		if testRoute.method == POST {
			suite.Ferry.Post(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
		if testRoute.method == PUT {
			suite.Ferry.Put(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
		if testRoute.method == DELETE {
			suite.Ferry.Delete(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
	}
	for _, testRoute := range routes {
		h := suite.Ferry.routerMap[testRoute.method]
		if assert.NotNil(suite.T(), h) {
			assert.Equal(suite.T(), h[0].routerPath, testRoute.path, "router path should match")
			regexPath := findAndReplace(testRoute.path)
			assert.Equal(suite.T(), h[0].regexPath, regexPath, "regex path should match")
		}
	}
}

func (suite *ServerSuite) TestGetMethodResponse() {
	path := "/hey"
	response := "hello, world!"
	suite.Ferry.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Ferry)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Error(err)
			}
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
			assert.Equal(suite.T(), string(body), response)
		}
	}
}

func (suite *ServerSuite) TestServeDir() {
	suite.Ferry.ServerDir("/", "example")
	r, err := http.NewRequest(GET, "http://test/main.go", nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Ferry)
		if assert.Nil(suite.T(), err) {
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
		}
	}
}

func (suite *ServerSuite) TestServeFile() {
	suite.Ferry.ServeFile("/main", "example/main.go")
	r, err := http.NewRequest(GET, "http://test/main", nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Ferry)
		if assert.Nil(suite.T(), err) {
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
		}
	}
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
