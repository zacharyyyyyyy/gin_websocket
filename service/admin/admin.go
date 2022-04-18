package admin

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"gin_websocket/dao"
	"gin_websocket/lib/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	VerifyUsernameAndPasswordFailErr = errors.New("账号或密码错误,请重试!")
	UnKnownErr                       = errors.New("服务忙碌，请稍后重试")
)

func verifyPassword(username, password string) (int, error) {
	adminDao, err := dao.SelectOneByUsername(username)
	if err != nil {
		return 0, VerifyUsernameAndPasswordFailErr
	}
	saltingPassword := `3'2W4E($*^%*URFY7"&HEASfa<@#RCVSATY4590-GA` + password + `%9da%$^#'saT"HS>fdhgashs#@fA`
	sha1Handle := sha1.New()
	sha1Handle.Write([]byte(saltingPassword))
	hexString := hex.EncodeToString(sha1Handle.Sum([]byte(nil)))
	if hexString != adminDao.Password {
		return 0, VerifyUsernameAndPasswordFailErr
	}
	return adminDao.Id, nil
}

func Login(username, password string, cRequest *http.Request, cResponse gin.ResponseWriter) error {
	adminId, err := verifyPassword(username, password)
	if err != nil {
		return VerifyUsernameAndPasswordFailErr
	}
	sessionCtl := session.NewSession(cRequest, cResponse)
	if err = sessionCtl.Set("role", strconv.Itoa(adminId)); err != nil {
		return UnKnownErr
	}
	return nil

}

func Logout(cRequest *http.Request, cResponse gin.ResponseWriter) error {
	sessionCtl := session.NewSession(cRequest, cResponse)
	return sessionCtl.Del()
}
