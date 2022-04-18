package controller

import (
	"gin_websocket/lib/logger"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

const (
	UnexpectedSystemErrorMsg  = "此服务忙碌中，请稍后重试"
	UnauthorizedErrorMsg      = "权限校验失败"
	ThirdPartyServiceErrorMsg = "服务暂不可用，有可能正在维护"
)

var DefaultErrorMsgMap = map[int]string{
	http.StatusInternalServerError: UnexpectedSystemErrorMsg,
	http.StatusUnauthorized:        UnauthorizedErrorMsg,
}

type ResponseStruct struct {
	C       *gin.Context
	Data    interface{}
	message string
	Code    int
}

type jsonResponseStruct struct {
	JsonData    interface{} `json:"data"`
	JsonMessage string      `json:"message"`
	JsonCode    int         `json:"code"`
}

func PanicResponse(c *gin.Context, err error, code int) {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	logger.Api.Error(err.Error())
	//TODO
	baseController := ResponseStruct{Code: code, C: c}
	baseController.JsonResponse()
}

//成功且无需返回任何信息时使用
func QuickSuccessResponse(c *gin.Context) {
	baseController := ResponseStruct{Code: http.StatusOK, C: c}
	baseController.JsonResponse()
}

func (resp *ResponseStruct) JsonResponse() {
	resp.C.Writer.Header().Set("Content-Type", "application/json")
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
	resp.C.String(resp.Code, string(jsonStr))
}

func (resp *ResponseStruct) SetMessage(msg string) {
	resp.message = msg
}

func (resp *ResponseStruct) SetHeaders(headers map[string]string) {
	for key, value := range headers {
		resp.C.Writer.Header().Set(key, value)
	}
}

func (resp *ResponseStruct) setMessageByCode() {
	_, ok := DefaultErrorMsgMap[resp.Code]
	if !ok {
		resp.SetMessage(UnexpectedSystemErrorMsg)
	} else {
		resp.SetMessage(DefaultErrorMsgMap[resp.Code])
	}
}
