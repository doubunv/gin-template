package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project-api/common"
	"project-api/common/headInfo"
)

var mapPath = map[string]int{
	"/api/user/v1/list": 1,
}

func AuthWhitePath(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	if _, ok := mapPath[path]; ok == false && headInfo.GetTokenUid(ctx) == 0 {
		res := common.Result{}
		ctx.AbortWithStatusJSON(http.StatusOK, res.FailCode(ctx, "visitor not done！", 401))
	}
}
