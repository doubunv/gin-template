package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"project-api/common/tool"
	"strconv"
)

func checkToken(token string) int64 {
	return 0
}

func AuthToken(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			tool.Error(ctx, fmt.Sprintf("%v", err))
		}
	}()

	token := ctx.GetHeader("Token")
	if token == "" {
		return
	}

	tokenUid := checkToken(token)
	if tokenUid <= 0 {
		return
	}

	dd := ctx.Value("metadata")
	md, bl := metadata.FromOutgoingContext(dd.(context.Context))
	if !bl {
		return
	}
	md.Set("tokenUid", strconv.FormatInt(tokenUid, 10))
	ctx.Set("metadata", metadata.NewOutgoingContext(context.Background(), md))
}
