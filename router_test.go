package slide

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RouterSuit struct {
	suite.Suite
	Slide *Slide
}

type testGroupRoutes struct {
	path   string
	method string
}

func (suite *RouterSuit) SetupTest() {
	config := &Config{}
	app := InitServer(config)
	suite.Slide = app
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
	group := suite.Slide.Group(groupPath)
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
		h := suite.Slide.routerMap[testRoute.method]
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
	group := suite.Slide.Group("/group")
	group.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	r, err := http.NewRequest(GET, "http://test/group"+path, nil)
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

func (suite *RouterSuit) Test404() {
	path := "/hey"
	response := "hello, world!"
	group := suite.Slide.Group("/group")
	group.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	// here we are giving a wrong URL
	r, err := http.NewRequest(GET, "http://test/groups"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
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
	suite.Slide.HandleNotFound(func(ctx *Ctx) error {
		return ctx.Send(http.StatusNotFound, notFoundMessage)
	})
	// here we are giving a random path
	r, err := http.NewRequest(GET, "http://test/random", nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
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

func (suite *RouterSuit) TestNestedGroup() {
	group := suite.Slide.Group("/hey")
	nestedGroup := group.Group("/hello")
	nestedGroup.Get("/slide", func(ctx *Ctx) error {
		return ctx.SendStatusCode(http.StatusOK)
	})
	r, err := http.NewRequest(GET, "http://test/hey/hello/slide", nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
		}
	}
}

func (suite *RouterSuit) TestRouterError() {
	suite.Slide.Get("/hey", func(ctx *Ctx) error {
		return errors.New("test error")
	})
	r, err := http.NewRequest(GET, "http://test/hey", nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Errorf("error while readingn body")
			}
			defer res.Body.Close()
			assert.Equal(suite.T(), res.StatusCode, http.StatusInternalServerError)
			assert.Equal(suite.T(), string(body), "test error")
		}
	}
}

func (suite *RouterSuit) TestCustomErrorRouter() {
	suite.Slide.Get("/hey", func(ctx *Ctx) error {
		return errors.New("test error")
	})
	suite.Slide.HandleErrors(func(ctx *Ctx, err error) error {
		return ctx.Send(http.StatusInternalServerError, fmt.Sprintf("%s from custom handler", err.Error()))
	})
	r, err := http.NewRequest(GET, "http://test/hey", nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Errorf("error while readingn body")
			}
			defer res.Body.Close()
			assert.Equal(suite.T(), res.StatusCode, http.StatusInternalServerError)
			assert.Equal(suite.T(), string(body), "test error from custom handler")
		}
	}
}

func (suite *RouterSuit) TestRouteMiddleware() {
	suite.Slide.Get("/hey", func(ctx *Ctx) error {
		return ctx.SendStatusCode(http.StatusOK)
	}, func(ctx *Ctx) error {
		// early response from middleware
		return ctx.Send(http.StatusOK, "response from middleware")
	}, func(ctx *Ctx) error {
		ctx.RequestCtx.Response.Header.Set("server", "slide")
		return ctx.Next()
	})
	r, err := http.NewRequest(GET, "http://test/hey", nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Errorf("error while readingn body")
			}
			defer res.Body.Close()
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
			assert.Equal(suite.T(), res.Header.Get("server"), "slide")
			assert.Equal(suite.T(), string(body), "response from middleware")
		}
	}
}

func TestRoutes(t *testing.T) {
	suite.Run(t, new(RouterSuit))
}
