package base

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

const (
	UnexpectedSystemErrorMsg  = "此服务维护中，暂时不可用"
	UnauthorizedErrorMsg      = "权限校验失败"
	ThirdPartyServiceErrorMsg = "服务暂不可用，有可能正在维护"
)

var DefaultErrorMsgMap = map[int]string{
	http.StatusInternalServerError: UnexpectedSystemErrorMsg,
	http.StatusUnauthorized:        UnauthorizedErrorMsg,
}

type ResponseStruct struct {
	Data    interface{}
	message string
	Code    int
}

type jsonResponseStruct struct {
	JsonData    interface{} `json:"data"`
	JsonMessage string      `json:"message"`
	JsonCode    int         `json:"code"`
}

func (resp *ResponseStruct) JsonResponse(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	if resp.Code == 0 {
		resp.Code = http.StatusOK
	}
	if resp.message == "" && resp.Code != http.StatusOK {
		resp.setMessageByCode()
	} else if resp.message == "" && resp.Code == http.StatusOK {
		resp.SetMessage("成功")
	}
	jsonResp := jsonResponseStruct{
		JsonData:    resp.Data,
		JsonMessage: resp.message,
		JsonCode:    resp.Code,
	}
	jsonStr, _ := jsoniter.Marshal(jsonResp)
	c.String(resp.Code, string(jsonStr))
}

func (resp *ResponseStruct) SetMessage(msg string) {
	resp.message = msg
}

func (resp *ResponseStruct) setMessageByCode() {
	_, ok := DefaultErrorMsgMap[resp.Code]
	if !ok {
		resp.SetMessage(UnexpectedSystemErrorMsg)
	} else {
		resp.SetMessage(DefaultErrorMsgMap[resp.Code])
	}
}
