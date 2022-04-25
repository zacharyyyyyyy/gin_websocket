package admin

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"

	"gin_websocket/dao"
	"gin_websocket/lib/session"
	"gin_websocket/model"

	"github.com/gin-gonic/gin"
)

var (
	VerifyUsernameAndPasswordFailErr = errors.New("账号或密码错误,请重试!")
	UnKnownErr                       = errors.New("服务忙碌，请稍后重试")
	AdminNotFoundErr                 = errors.New("未登录或登录过期，请重新登录")
)

func verifyPassword(username, password string) (*model.Admin, error) {
	adminDao, err := dao.SelectOneByUsername(username)
	if err != nil {
		return nil, VerifyUsernameAndPasswordFailErr
	}
	saltingPassword := `3'2W4E($*^%*URFY7"&HEASfa<@#RCVSATY4590-GA` + password + `%9da%$^#'saT"HS>fdhgashs#@fA`
	sha1Handle := sha1.New()
	sha1Handle.Write([]byte(saltingPassword))
	hexString := hex.EncodeToString(sha1Handle.Sum([]byte(nil)))
	if hexString != adminDao.Password {
		return nil, VerifyUsernameAndPasswordFailErr
	}
	return adminDao, nil
}

func ChangePassword(password string) string {
	saltingPassword := `3'2W4E($*^%*URFY7"&HEASfa<@#RCVSATY4590-GA` + password + `%9da%$^#'saT"HS>fdhgashs#@fA`
	sha1Handle := sha1.New()
	sha1Handle.Write([]byte(saltingPassword))
	hexString := hex.EncodeToString(sha1Handle.Sum([]byte(nil)))
	return hexString
}

func Login(username, password string, cRequest *http.Request, cResponse gin.ResponseWriter) error {
	adminId, err := verifyPassword(username, password)
	if err != nil {
		return VerifyUsernameAndPasswordFailErr
	}
	sessionCtl := session.NewSession(cRequest, cResponse)
	if err = sessionCtl.Set("role", strconv.Itoa(adminId.Role)); err != nil {
		return UnKnownErr
	}
	if err = sessionCtl.Set("admin", strconv.Itoa(adminId.Id)); err != nil {
		return UnKnownErr
	}
	return nil
}

func Logout(cRequest *http.Request, cResponse gin.ResponseWriter) error {
	sessionCtl := session.NewSession(cRequest, cResponse)
	return sessionCtl.Del()
}
