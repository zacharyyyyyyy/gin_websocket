package admin

import (
	"gin_websocket/dao"
	"gin_websocket/lib/session"
	"github.com/gin-gonic/gin"
)

func Login(username, password string) error {
	adminDao, err := dao.SelectOneByUsername(username)
	if err != nil {
		return err
	}
	//Todo
}

func Logout(c *gin.Context) error {
	sessionCtl := session.NewSession(c)
	return sessionCtl.Del()
}
