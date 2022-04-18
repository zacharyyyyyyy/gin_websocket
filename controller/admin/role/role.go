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
		return
	}
	fmt.Println(result)

}
