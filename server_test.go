package slide

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	Slide *Slide
}

type testRoutes struct {
	path   string
	method string
}

func (suite *ServerSuite) SetupTest() {
	config := &Config{}
	app := InitServer(config)
	suite.Slide = app
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
			suite.Slide.Get(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
		if testRoute.method == POST {
			suite.Slide.Post(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
		if testRoute.method == PUT {
			suite.Slide.Put(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
		if testRoute.method == DELETE {
			suite.Slide.Delete(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
	}
	for _, testRoute := range routes {
		h := suite.Slide.routerMap[testRoute.method]
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
	suite.Slide.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
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
	suite.Slide.ServerDir("/", "example")
	r, err := http.NewRequest(GET, "http://test/main.go", nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
		}
	}
}

func (suite *ServerSuite) TestServeFile() {
	suite.Slide.ServeFile("/main", "example/main.go")
	r, err := http.NewRequest(GET, "http://test/main", nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
		}
	}
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
