package ferry

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"regexp"
	"testing"
)

func TestAttachmentHeader(t *testing.T) {
	fileName := "hey"
	attachment := getAttachmentHeader(fileName)
	attachmentExpected := fmt.Sprintf("%s; filename=%s", Attachment, fileName)
	assert.Equal(t, attachment, attachmentExpected)
}

func TestPathParamRegex(t *testing.T) {
	pathParam := "/auth/:name"
	routerParam := "/auth/ferry"
	regexParam := findAndReplace(pathParam)
	assert.MatchRegex(t, routerParam, regexParam)
}

func TestPathParamRegexAbs(t *testing.T) {
	pathParam := "/auth/ferry"
	routerParam := "/auth/ferry"
	regexParam := findAndReplace(pathParam)
	assert.MatchRegex(t, routerParam, regexParam)
}

func TestPathParamRegexFail(t *testing.T) {
	pathParam := "/auth/:name"
	routerParam := "/auth/ferry/ss"
	regexParam := findAndReplace(pathParam)
	fail, _ := regexp.MatchString(routerParam, regexParam)
	assert.Equal(t, false, fail)
}

func TestGetParam(t *testing.T) {
	name := "ferry"
	pathParam := "/auth/:name/hello/:age"
	routerParam := fmt.Sprintf("/auth/%s/hello/%d", name, 1)
	wantedName := extractParamFromPath(pathParam, routerParam, "name")
	assert.Equal(t, wantedName, name)
}

func TestGetParamEmpty(t *testing.T) {
	name := ""
	pathParam := "/auth/:name"
	routerParam := "/auth/" + name
	wantedName := extractParamFromPath(pathParam, routerParam, "names")
	assert.Equal(t, wantedName, name)
}

func TestGetParams(t *testing.T) {
	name := "ferry"
	pathParam := "/auth/:name/hello/:age"
	routerParam := fmt.Sprintf("/auth/%s/hello/%d", name, 1)
	paramsMap := getParamsFromPath(pathParam, routerParam)
	wantedParamsMap := map[string]string{
		"name": name,
		"age":  "1",
	}
	assert.Equal(t, paramsMap, wantedParamsMap)
}
