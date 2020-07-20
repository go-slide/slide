package ferry

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RouterSuit struct {
	suite.Suite
	Ferry *Ferry
}

type testGroupRoutes struct {
	path   string
	method string
}

func (suite *RouterSuit) SetupTest() {
	config := &Config{}
	app := InitServer(config)
	suite.Ferry = app
}

func (suite *RouterSuit) TestGetMethod() {
	groupPath := "/group"
	routes := []testGroupRoutes{
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
	group := suite.Ferry.Group(groupPath)
	for _, testRoute := range routes {
		if testRoute.method == GET {
			group.Get(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
		if testRoute.method == POST {
			group.Post(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
		if testRoute.method == PUT {
			group.Put(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
		if testRoute.method == DELETE {
			group.Delete(testRoute.path, func(ctx *Ctx) error {
				return ctx.Send(http.StatusOK, "hey")
			})
		}
	}
	for _, testRoute := range routes {
		h := suite.Ferry.routerMap[testRoute.method]
		if assert.NotNil(suite.T(), h) {
			assert.Equal(suite.T(), h[0].routerPath, groupPath+testRoute.path, "router path should match")
			regexPath := findAndReplace(groupPath + testRoute.path)
			assert.Equal(suite.T(), h[0].regexPath, regexPath, "regex path should match")
		}
	}
}

func (suite *RouterSuit) TestGetMethodResponse() {
	path := "/hey"
	response := "hello, world!"
	group := suite.Ferry.Group("/group")
	group.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	r, err := http.NewRequest(GET, "http://test/group"+path, nil)
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

func (suite *RouterSuit) Test404() {
	path := "/hey"
	response := "hello, world!"
	group := suite.Ferry.Group("/group")
	group.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	// here we are giving a wrong URL
	r, err := http.NewRequest(GET, "http://test/groups"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Ferry)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Error(err)
			}
			assert.Equal(suite.T(), res.StatusCode, http.StatusNotFound)
			assert.Equal(suite.T(), string(body), NotFoundMessage)
		}
	}
}

func (suite *RouterSuit) TestCustom404Handler() {
	notFoundMessage := "check url"
	suite.Ferry.HandleNotFound(func(ctx *Ctx) error {
		return ctx.Send(http.StatusNotFound, notFoundMessage)
	})
	// here we are giving a random path
	r, err := http.NewRequest(GET, "http://test/random", nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Ferry)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Error(err)
			}
			assert.Equal(suite.T(), res.StatusCode, http.StatusNotFound)
			assert.Equal(suite.T(), string(body), notFoundMessage)
		}
	}
}

func TestRoutes(t *testing.T) {
	suite.Run(t, new(RouterSuit))
}
