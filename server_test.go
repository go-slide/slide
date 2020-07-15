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

func (suite *ServerSuite) SetupTest() {
	config := &Config{}
	app := InitServer(config)
	suite.Ferry = app
}

func (suite *ServerSuite) TestGetMethod() {
	path := "/hey"
	suite.Ferry.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, "hey")
	})
	h := suite.Ferry.routerMap[get]
	if assert.NotNil(suite.T(), h) {
		assert.Equal(suite.T(), h[0].routerPath, path, "router path should match")
		regexPath := findAndReplace(path)
		assert.Equal(suite.T(), h[0].regexPath, regexPath, "regex path should match")
	}
}

func (suite *ServerSuite) TestGetMethodResponse() {
	path := "/hey"
	response := "hello, world!"
	suite.Ferry.Get(path, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, response)
	})
	r, err := http.NewRequest(get, "http://test"+path, nil)
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

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
