package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"project-api/common/tool"
)

func ApiParams(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			tool.Error(ctx, fmt.Sprintf("%v", err))
		}
	}()
	data, _ := ctx.GetRawData()
	tool.Info(ctx, fmt.Sprintf("API Params：%s,%+v", ctx.Request.RequestURI, string(data)))
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	ctx.Next()
}
