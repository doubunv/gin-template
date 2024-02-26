package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"project-api/common/tool"
)

func CheckHeader(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			tool.Error(ctx, fmt.Sprintf("%v", err))
		}
	}()
	headData := map[string]string{
		"tokenUid": "0",
		"traceId":  tool.GenerateUUID(),
		"clientIp": ctx.ClientIP(),
		"version":  ctx.GetHeader("version"),
		"deviceId": ctx.GetHeader("deviceId"),
		"token":    ctx.GetHeader("AuthorizationJwt"),
		"apiUrl":   ctx.Request.RequestURI,
	}
	c := metadata.NewOutgoingContext(context.Background(), metadata.New(headData))
	ctx.Set("metadata", c)
}
