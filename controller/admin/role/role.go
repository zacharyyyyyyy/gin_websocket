package role

import (
	"fmt"
	"gin_websocket/controller"
	"gin_websocket/dao"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllAdminAuth(c *gin.Context) {
	result, err := dao.GetAllAuthByEnable()
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError)
	}
	result1, err := dao.GetRoleByAuth(result[0].Id)
	fmt.Println(result[0].Id)
	fmt.Println(result1[0].Role)

}
