package ferry

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MiddlewareSuite struct {
	suite.Suite
	Ferry *Ferry
}

func (suite *MiddlewareSuite) SetupTest() {
	config := &Config{}
	app := InitServer(config)
	suite.Ferry = app
}

func (suite *MiddlewareSuite) TestAppLevelMiddleware() {
	path := "/hey"
	response := "hello, world!"
	suite.Ferry.Use(func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, "response from middleware")
	})
	suite.Ferry.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	// first send early response
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Ferry)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Error(err)
			}
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
			assert.Equal(suite.T(), string(body), "response from middleware")
		}
	}
}

func (suite *MiddlewareSuite) TestAppLevelMultiMiddleware() {
	path := "/hey"
	response := "hello, world!"
	suite.Ferry.Use(func(ctx *Ctx) error {
		ctx.RequestCtx.Response.Header.Set("hey1", "hey1")
		return ctx.Next()
	})
	suite.Ferry.Use(func(ctx *Ctx) error {
		ctx.RequestCtx.Response.Header.Set("hey2", "hey2")
		return ctx.Next()
	})
	suite.Ferry.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	// first send early response
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Ferry)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Error(err)
			}
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
			assert.Equal(suite.T(), res.Header.Get("hey1"), "hey1")
			assert.Equal(suite.T(), res.Header.Get("hey2"), "hey2")
			assert.Equal(suite.T(), string(body), response)
		}
	}
}

func (suite *MiddlewareSuite) TestGroupMultiMiddleware() {
	path := "/hey"
	response := "hello, world!"
	groupPath := "/group"
	group := suite.Ferry.Group(groupPath)
	group.Use(func(ctx *Ctx) error {
		ctx.RequestCtx.Response.Header.Set("hey1", "hey1")
		return ctx.Next()
	})
	group.Use(func(ctx *Ctx) error {
		ctx.RequestCtx.Response.Header.Set("hey2", "hey2")
		return ctx.Next()
	})
	group.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	// first send early response
	r, err := http.NewRequest(GET, "http://test"+groupPath+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Ferry)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Error(err)
			}
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
			assert.Equal(suite.T(), res.Header.Get("hey1"), "hey1")
			assert.Equal(suite.T(), res.Header.Get("hey2"), "hey2")
			assert.Equal(suite.T(), string(body), response)
		}
	}
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}
