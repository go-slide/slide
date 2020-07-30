package slide

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ContextSuite struct {
	suite.Suite
	Slide *Slide
}

type login struct {
	Username string `json:"username"`
}

func (suite *ContextSuite) SetupTest() {
	validate := validator.New()
	config := &Config{
		Validator: validate,
	}
	app := InitServer(config)
	suite.Slide = app
}

func (suite *ContextSuite) TestJsonResponse() {
	path := "/hey"
	response := map[string]string{
		"user": "slide",
	}
	suite.Slide.Get(path, func(ctx *Ctx) error {
		return ctx.JSON(http.StatusOK, response)
	})
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				suite.T().Error(err)
			}
			bytes, err := json.Marshal(response)
			if err != nil {
				suite.T().Error(err)
			}
			assert.Equal(suite.T(), res.Header.Get(ContentType), ApplicationJSON)
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
			assert.Equal(suite.T(), body, bytes)
		}
	}
}

func (suite *ContextSuite) TestRedirect() {
	path := "/hey"
	redirectPath := "/redirect"
	suite.Slide.Get(path, func(ctx *Ctx) error {
		return ctx.Redirect(http.StatusTemporaryRedirect, redirectPath)
	})
	suite.Slide.Get(redirectPath, func(ctx *Ctx) error {
		return ctx.Send(http.StatusOK, "redirected")
	})
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			assert.Equal(suite.T(), res.StatusCode, http.StatusOK)
		}
	}
}

func (suite *ContextSuite) TestBind() {
	path := "/hey"
	suite.Slide.Post(path, func(ctx *Ctx) error {
		body := login{}
		_ = ctx.Bind(&body)
		return ctx.JSON(http.StatusOK, body)
	})
	postBody, _ := json.Marshal(login{
		Username: "Slide",
	})
	r, err := http.NewRequest(POST, "http://test"+path, strings.NewReader(string(postBody)))
	if assert.Nil(suite.T(), err) {
		res, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			body, err := ioutil.ReadAll(res.Body)
			if assert.Nil(suite.T(), err) {
				res.Body.Close()
				assert.Equal(suite.T(), body, postBody)
			}
		}
	}
}

func (suite *ContextSuite) TestParamsAndQueryParams() {
	path := "/hey/:name"
	name := "slide"
	paramsMap := map[string]string{
		"name": name,
	}
	queryMap := map[string]string{
		"key":  "value",
		"name": "slide",
	}
	requestPath := fmt.Sprintf("/hey/%s?key=value&name=slide", name)
	suite.Slide.Get(path, func(ctx *Ctx) error {
		name := ctx.GetParam("name")
		queryParamName := ctx.GetQueryParam("name")

		paramsMapExp := ctx.GetParams()
		queryMapExp := ctx.GetQueryParams()

		assert.Equal(suite.T(), name, name)
		assert.Equal(suite.T(), queryParamName, name)
		assert.Equal(suite.T(), paramsMapExp, paramsMap)
		assert.Equal(suite.T(), queryMapExp, queryMap)
		return ctx.SendStatusCode(http.StatusOK)
	})
	r, err := http.NewRequest(GET, "http://test"+requestPath, nil)
	if assert.Nil(suite.T(), err) {
		_, err := testServer(r, suite.Slide)
		assert.Nil(suite.T(), err)
	}
}

func (suite *ContextSuite) TestServeFiles() {
	path := "/hey"
	filePath := "server.go"
	suite.Slide.Get(path, func(ctx *Ctx) error {
		return ctx.ServeFile(filePath)
	})
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		_, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			fileType, _ := getFileContentType(filePath)
			assert.Equal(suite.T(), r.Header.Get(ContentType), fileType)
		}
		assert.Nil(suite.T(), err)
	}
}

func (suite *ContextSuite) TestSendAttachment() {
	path := "/hey"
	filePath := "server.go"
	suite.Slide.Get(path, func(ctx *Ctx) error {
		return ctx.SendAttachment(filePath, "server.go")
	})
	r, err := http.NewRequest(GET, "http://test"+path, nil)
	if assert.Nil(suite.T(), err) {
		resp, err := testServer(r, suite.Slide)
		if assert.Nil(suite.T(), err) {
			header := getAttachmentHeader("server.go")
			assert.Equal(suite.T(), resp.Header.Get(ContentDeposition), header)
		}
		assert.Nil(suite.T(), err)
	}
}

func createMultipartFormData(suite *ContextSuite, fieldName, filePath string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file, err := os.Open(filePath)
	if err != nil {
		suite.T().Errorf("Error while reading file: %v", err)
	}
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		suite.T().Errorf("Error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		suite.T().Errorf("Error with io.Copy: %v", err)
	}
	w.Close()
	return b, w
}

func (suite *ContextSuite) TestUploadFile() {
	path := "/hey"
	dirPath := "temp"
	fileName := "server.go"
	uploadPath := dirPath + "/" + fileName
	// first create a folder
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, os.ModeDir)
		suite.NotNil(suite.T(), err)
	}
	suite.Slide.Post(path, func(ctx *Ctx) error {
		return ctx.UploadFile(uploadPath, fileName)
	})
	buffer, multiWriter := createMultipartFormData(suite, fileName, fileName)
	r, err := http.NewRequest(POST, "http://test"+path, &buffer)
	if assert.Nil(suite.T(), err) {
		r.Header.Set("Content-Type", multiWriter.FormDataContentType())
		_, err := testServer(r, suite.Slide)
		assert.Nil(suite.T(), err)
		if assert.Nil(suite.T(), err) {
			if _, pathError := os.Stat(uploadPath); pathError != nil {
				if os.IsNotExist(pathError) {
					assert.Errorf(suite.T(), pathError, "error while upload file for path "+uploadPath)
				}
			}
		}
	}
	// delete that file
	deleteErr := os.RemoveAll(dirPath)
	if deleteErr != nil {
		assert.Fail(suite.T(), "error while deleting dir", deleteErr.Error())
	}
}

func TestContext(t *testing.T) {
	suite.Run(t, new(ContextSuite))
}
