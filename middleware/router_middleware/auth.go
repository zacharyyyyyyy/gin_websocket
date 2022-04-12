package router_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/storyicon/grbac"
	"net/http"
	"time"
)

func loadAuthorizationRules() (rules grbac.Rules, err error) {
	//TODO
	return
}

func QueryRolesByHeaders(header http.Header) (roles []string, err error) {
	//TODO
	return roles, err
}

func AdminAuthentication() gin.HandlerFunc {
	rbac, err := grbac.New(grbac.WithLoader(loadAuthorizationRules, time.Minute))
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		roles, err := QueryRolesByHeaders(c.Request.Header)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		state, err := rbac.IsRequestGranted(c.Request, roles)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !state.IsGranted() {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	}

}
