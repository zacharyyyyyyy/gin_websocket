package test

import (
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"gin_websocket/router"
	jsoniter "github.com/json-iterator/go"
)

func TestSyncTestHandler(t *testing.T) {
	routerTest := router.InitRouter()
	tests := []struct {
		testName string
		method   string
		url      string
	}{
		{
			"ping_test",
			"POST",
			"/admin/info",
		},
	}

	for _, test := range tests {
		urlString := url.Values{}
		urlTest, _ := url.Parse(test.url)
		//urlString.Set("")
		urlTest.RawQuery = urlString.Encode()
		t.Run(test.testName, func(t *testing.T) {
			req := httptest.NewRequest(test.method, urlTest.String(), strings.NewReader(""))
			writer := httptest.NewRecorder()
			routerTest.ServeHTTP(writer, req)
			//code 判断
			assert.Equal(t, http.StatusOK, writer.Code)
			//msg 判断
			var resp map[string]interface{}
			err := jsoniter.Unmarshal([]byte(writer.Body.String()), &resp)
			assert.NilError(t, err)
			assert.Equal(t, "成功", resp["message"])
		})
	}

}
