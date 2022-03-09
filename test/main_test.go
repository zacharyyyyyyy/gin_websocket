package test

import (
	"gin_websocket/router"
	jsoniter "github.com/json-iterator/go"
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSyncTestHandler(t *testing.T) {
	router := router.InitRouter()
	urlString := url.Values{}
	urlTest, _ := url.Parse("/admin/ping")
	//urlString.Set("")
	urlTest.RawQuery = urlString.Encode()
	t.Run("ping_test", func(t *testing.T) {
		req := httptest.NewRequest("GET", urlTest.String(), strings.NewReader(""))
		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, req)
		//code 判断
		assert.Equal(t, http.StatusOK, writer.Code)
		//msg 判断
		var resp map[string]string
		err := jsoniter.Unmarshal([]byte(writer.Body.String()), &resp)
		assert.NilError(t, err)
		assert.Equal(t, "success", resp["message"])
	})

}
