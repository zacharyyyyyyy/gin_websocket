package global_middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gin_websocket/service/tracer"
	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
)

func HttpTrace(c *gin.Context) {
	var (
		paramString       string
		contentTypeString string
		param             = make([]string, 0)
		contentTypeSlice  = make([]string, 0)
	)

	httpTracer := &tracer.Tracer{}
	httpTracer.AddTag(tracer.Tag{Key: tracer.TagHttpURL, Value: c.Request.URL.Path})
	httpTracer.AddTag(tracer.Tag{Key: tracer.TagHttpRawURL, Value: c.Request.URL.String()})
	httpTracer.AddTag(tracer.Tag{Key: tracer.TagHttpMethod, Value: c.Request.Method})
	httpTracer.AddTag(tracer.Tag{Key: tracer.TagHttpClientURL, Value: c.ClientIP()})

	contentType := c.Request.Header.Get("Content-Type")
	contentTypeSlice = strings.Split(contentType, ";")
	if len(contentTypeSlice) > 0 {
		contentTypeString = contentTypeSlice[0]
	} else {
		contentTypeString = ""
	}
	httpTracer.AddTag(tracer.Tag{Key: tracer.TagHttpContentType, Value: contentTypeString})

	switch c.Request.Method {
	case http.MethodGet:
		for queryKey, queryVal := range c.Request.URL.Query() {
			param = append(param, fmt.Sprintf("%s:%s", queryKey, queryVal))
		}
	case http.MethodDelete:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodPost:
		switch contentTypeString {
		//form 由于可能存在文件流 不记录数据
		case "multipart/form-data":
		case "application/x-www-form-urlencoded":
			_ = c.Request.ParseForm()
			for postKey, postVal := range c.Request.PostForm {
				param = append(param, fmt.Sprintf("%s:%s", postKey, postVal))
			}
		case "application/json":
			jsonData := make(map[string]interface{}, 0)
			data, _ := ioutil.ReadAll(c.Request.Body)
			_ = jsoniter.Unmarshal(data, &jsonData)
			for jsonKey, jsonValue := range jsonData {
				param = append(param, fmt.Sprintf("%s:%v", jsonKey, jsonValue))
			}
		case "":
		}
	}
	paramString = "[" + strings.Join(param, ",") + "]"
	httpTracer.AddTag(tracer.Tag{Key: tracer.TagHttpParam, Value: paramString})

	c.Next()
	httpTracer.AddTag(tracer.Tag{Key: tracer.TagHttpStatusCode, Value: strconv.Itoa(c.Writer.Status())})
	httpTracer.Finish()
}
