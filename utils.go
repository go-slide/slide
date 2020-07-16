package ferry

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	GET    = http.MethodGet
	POST   = http.MethodPost
	PUT    = http.MethodPut
	DELETE = http.MethodDelete

	// headers
	ContentType       = "Content-Type"
	ContentDeposition = "Content-Disposition"
	ApplicationJson   = "application/json"
	Attachment        = "attachment"

	routerRegexReplace = "[a-zA-Z0-9_-]*"

	// routing error messages
	NotFoundMessage = "Not Found, Check URL"
)

// Finds wild card in URL and replace them with a regex for,
// ex if path is /auth/:name -> /auth/[a-zA-Z0-9]*
// ex if path is /auth/name -> /auth/name
func findAndReplace(path string) string {
	if !strings.Contains(path, ":") {
		return fmt.Sprintf("%s%s%s", "^", path, "$")
	}
	result := ""
	slitted := strings.Split(path, "/")
	for _, v := range slitted {
		if v == "" {
			continue
		}
		if strings.Contains(v, ":") {
			result = fmt.Sprintf("%s/%s", result, routerRegexReplace)
			continue
		}
		result = fmt.Sprintf("%s/%s", result, v)
	}
	// replace slashes
	result = strings.ReplaceAll(result, "/", "\\/")
	result = fmt.Sprintf("%s%s%s", "^", result, "$")
	return result
}

// routerPath /auth/:name
// requestPath /auth/madhuri
// paramName name
// returns madhuri
func extractParamFromPath(routerPath, requestPath, paramName string) string {
	routerSplit := strings.Split(routerPath, "/")
	requestSplit := strings.Split(requestPath, "/")
	if len(routerSplit) != len(requestSplit) {
		return ""
	}
	paramWithWildCard := fmt.Sprintf(":%s", paramName)
	for k, v := range routerSplit {
		if v == paramWithWildCard {
			return requestSplit[k]
		}
	}
	return ""
}

// routerPath /auth/:name/:age
// requestPath /auth/madhuri/32
// returns { name: madhuri, age: 32 }
func getParamsFromPath(routerPath, requestPath string) map[string]string {
	paramsMap := map[string]string{}
	routerSplit := strings.Split(routerPath, "/")
	requestSplit := strings.Split(requestPath, "/")
	for k, v := range routerSplit {
		if strings.Contains(v, ":") {
			key := strings.ReplaceAll(v, ":", "")
			paramsMap[key] = requestSplit[k]
		}
	}
	return paramsMap
}

//	returns value of a single query Param
//
//	route path /hello?key=test&value=bbp
//
//	keyValue = GetQueryParam(key)
//
//	keyValue = test
func getAllQueryParams(queryPath string) map[string]string {
	queryParamsMap := map[string]string{}
	params := strings.Split(queryPath, "&")
	for _, v := range params {
		if strings.Contains(v, "=") {
			pair := strings.Split(v, "=")
			queryParamsMap[pair[0]] = pair[1]
		}
	}
	return queryParamsMap
}

//	returns map of query Params
//
//	route path /hello?key=test&value=bbp
//
//	returns {key : test, value : bbp}
func getQueryParam(querypath string, name string) string {
	params := strings.Split(querypath, "&")
	for _, v := range params {
		if strings.Contains(v, "=") {
			pair := strings.Split(v, "=")
			if pair[0] == name {
				return pair[1]
			}
		}
	}
	return ""
}

func getAttachmentHeader(fileName string) string {
	if fileName == "" {
		return fmt.Sprintf(Attachment)
	}
	return fmt.Sprintf("%s; filename=%s", Attachment, fileName)
}
