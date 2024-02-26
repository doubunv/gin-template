package headInfo

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"strconv"
	"strings"
)

type headInfo struct {
	Token    string
	TokenUid int64
	ClientIp string
	DeviceId string
}

func GetToken(ctx *gin.Context) string {
	dd := ctx.Value("metadata")
	md, bl := metadata.FromOutgoingContext(dd.(context.Context))
	if !bl {
		return ""
	}
	return strings.Join(md.Get("token"), "")
}

func GetTokenUid(ctx *gin.Context) int64 {
	dd := ctx.Value("metadata")
	md, bl := metadata.FromOutgoingContext(dd.(context.Context))
	if !bl {
		return 0
	}
	parseInt, err := strconv.ParseInt(strings.Join(md.Get("tokenUid"), ""), 10, 64)
	if err != nil {
		return 0
	}
	return parseInt
}

func GetClientIp(ctx *gin.Context) string {
	dd := ctx.Value("metadata")
	md, bl := metadata.FromOutgoingContext(dd.(context.Context))
	if !bl {
		return ""
	}
	return strings.Join(md.Get("clientIp"), "")
}

func GetAppVersion(ctx *gin.Context) string {
	dd := ctx.Value("metadata")
	md, bl := metadata.FromOutgoingContext(dd.(context.Context))
	if !bl {
		return ""
	}
	return strings.Join(md.Get("appVersion"), "")
}

func GetDeviceId(ctx *gin.Context) string {
	return ""
}

func GetRpcContext(ctx *gin.Context) context.Context {
	if _, ok := ctx.Value("metadata").(context.Context); ok {
		return ctx.Value("metadata").(context.Context)
	}
	return context.Background()
}

func GetTraceId(ctx *gin.Context) string {
	dd, _ := ctx.Get("metadata")
	md, _ := metadata.FromOutgoingContext(dd.(context.Context))
	return strings.Join(md.Get("traceId"), "")
}
