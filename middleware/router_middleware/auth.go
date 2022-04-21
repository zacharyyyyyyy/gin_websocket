package router_middleware

import (
	"gin_websocket/controller"
	"gin_websocket/lib/session"
	"net/http"
	"strconv"
	"time"

	"gin_websocket/dao"

	"github.com/gin-gonic/gin"
	"github.com/storyicon/grbac"
)

func loadAuthorizationRules() (rules grbac.Rules, err error) {
	rules = make(grbac.Rules, 0)
	result, err := dao.GetAllAuthByEnable()
	for _, auth := range result {
		method := auth.Method
		if method != "*" {
			method = "{" + method + "}"
		}
		authorizedRoles := make([]string, 0)
		roleMaps, _ := dao.GetRoleByAuth(auth.Id)
		for _, roleMap := range roleMaps {
			authorizedRoles = append(authorizedRoles, strconv.Itoa(roleMap.Role))
		}
		if len(authorizedRoles) > 0 {
			rules = append(rules, &grbac.Rule{
				ID: auth.Id,
				Resource: &grbac.Resource{
					Host:   "*",
					Path:   auth.Path,
					Method: method,
				},
				Permission: &grbac.Permission{
					AuthorizedRoles: authorizedRoles,
					ForbiddenRoles:  []string{},
					AllowAnyone:     false,
				},
			})
		}
	}
	return
}

func QueryRoles(cRequest *http.Request, cResponse gin.ResponseWriter) (roles []string, err error) {
	userSession := session.NewSession(cRequest, cResponse)
	roleString, err := userSession.GetString("role")
	if err != nil {
		return nil, err
	}
	roles = append(roles, roleString)
	return
}

func AdminAuthentication() gin.HandlerFunc {
	rbac, err := grbac.New(grbac.WithLoader(loadAuthorizationRules, time.Minute))
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		roles, err := QueryRoles(c.Request, c.Writer)
		if err != nil {
			controller.PanicResponse(c, err, http.StatusUnauthorized, "")
			c.Abort()
			return
		}

		state, err := rbac.IsRequestGranted(c.Request, roles)
		if err != nil {
			controller.PanicResponse(c, err, http.StatusUnauthorized, "")
			c.Abort()
			return
		}

		if !state.IsGranted() {
			controller.PanicResponse(c, err, http.StatusUnauthorized, "")
			c.Abort()
			return
		}

	}

}
