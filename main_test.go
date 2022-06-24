package main

import (
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"gin_websocket/router"

	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
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

func TestAddAdmin(t *testing.T) {
	routerTest := router.InitRouter()
	tests := []struct {
		testName string
		method   string
		url      string
		param    map[string]interface{}
	}{
		{
			"add_admin_test",
			"POST",
			"/admin/user/add",
			map[string]interface{}{
				"username": "1s",
				"password": "password",
				"name":     "sta",
				"role":     "1",
			},
		},
	}
	for _, test := range tests {
		var bodyStr string
		urlTest, _ := url.Parse(test.url)
		t.Run(test.testName, func(t *testing.T) {
			if test.method == http.MethodPost {
				var r http.Request
				_ = r.ParseForm()
				for paramKey, paramVal := range test.param {
					r.Form.Add(paramKey, paramVal.(string))
				}
				bodyStr = r.Form.Encode()
			} else {
				urlString := url.Values{}
				for paramKey, paramVal := range test.param {
					urlString.Set(paramKey, paramVal.(string))
				}
				urlTest.RawQuery = urlString.Encode()
			}
			req := httptest.NewRequest(test.method, urlTest.String(), strings.NewReader(bodyStr))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			writer := httptest.NewRecorder()
			routerTest.ServeHTTP(writer, req)
			//code 判断
			Convey("status", t, func() {
				So(writer.Code, ShouldEqual, http.StatusOK)
			})
			//msg 判断
			var resp map[string]interface{}
			err := jsoniter.Unmarshal([]byte(writer.Body.String()), &resp)
			Convey("err", t, func() {
				So(err, ShouldBeNil)
			})
			Convey("msg", t, func() {
				So(resp["message"], ShouldEqual, "成功")
			})

		})
	}
}
