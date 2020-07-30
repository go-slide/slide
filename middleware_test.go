package slide

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MiddlewareSuite struct {
	suite.Suite
	Slide *Slide
}

func (suite *MiddlewareSuite) SetupTest() {
	config := &Config{}
	app := InitServer(config)
	suite.Slide = app
}

func (suite *MiddlewareSuite) TestAppLevelMiddleware() {
	path := "/hey"
	response := "hello, world!"
	suite.Slide.Use(func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, "response from middleware")
	})
	suite.Slide.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	// first send early response
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
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
	suite.Slide.Use(func(ctx *Ctx) error {
		ctx.RequestCtx.Response.Header.Set("hey1", "hey1")
		return ctx.Next()
	})
	suite.Slide.Use(func(ctx *Ctx) error {
		ctx.RequestCtx.Response.Header.Set("hey2", "hey2")
		return ctx.Next()
	})
	suite.Slide.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	// first send early response
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
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
	group := suite.Slide.Group(groupPath)
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
		res, err := testServer(r, suite.Slide)
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

func (suite *MiddlewareSuite) TestAppLevelMiddlewareError() {
	path := "/hey"
	response := "hello, world!"
	suite.Slide.Use(func(ctx *Ctx) error {
		return errors.New("error from middleware")
	})
	suite.Slide.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	// first send early response
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Error(err)
			}
			assert.Equal(suite.T(), res.StatusCode, http.StatusInternalServerError)
			assert.Equal(suite.T(), string(body), "error from middleware")
		}
	}
}

func (suite *MiddlewareSuite) TestAppLevelMiddlewareMultiError() {
	path := "/hey"
	response := "hello, world!"
	suite.Slide.Use(func(ctx *Ctx) error {
		return ctx.Next()
	})
	suite.Slide.Use(func(ctx *Ctx) error {
		return errors.New("error from middleware")
	})
	suite.Slide.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	// first send early response
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Error(err)
			}
			assert.Equal(suite.T(), res.StatusCode, http.StatusInternalServerError)
			assert.Equal(suite.T(), string(body), "error from middleware")
		}
	}
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}
