package common

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/metadata"
	"project-api/common/headInfo"
	"project-api/common/tool"
	"strings"
)

type BusinessCode int

const (
	Success BusinessCode = 0
	Error   BusinessCode = 500
	Author  BusinessCode = 401
)

type Result struct {
	Code    BusinessCode `json:"code"`
	Msg     string       `json:"msg"`
	Data    any          `json:"data"`
	TraceId string       `json:"traceId"`
}

func (r *Result) Success(ctx *gin.Context, data any) *Result {
	dd, _ := ctx.Get("metadata")
	md, _ := metadata.FromOutgoingContext(dd.(context.Context))
	r.TraceId = strings.Join(md.Get("traceId"), "")
	r.Code = Success
	r.Msg = "success"
	r.Data = data
	t, _ := jsoniter.Marshal(r)
	tool.Info(ctx, fmt.Sprintf("API Response:%s,%s", ctx.Request.RequestURI, string(t)))
	return r
}

func (r *Result) Fail(ctx *gin.Context, msg string) *Result {
	r.TraceId = headInfo.GetTraceId(ctx)
	r.Code = Error
	r.Msg = msg
	t, _ := jsoniter.Marshal(r)
	tool.Error(ctx, fmt.Sprintf("API Response:%s,%s", ctx.Request.RequestURI, string(t)))
	return r
}

func (r *Result) FailCode(ctx *gin.Context, msg string, code BusinessCode) *Result {
	r.TraceId = headInfo.GetTraceId(ctx)
	r.Code = code
	r.Msg = msg
	t, _ := jsoniter.Marshal(r)
	tool.Error(ctx, fmt.Sprintf("API Response:%s,%s", ctx.Request.RequestURI, string(t)))
	return r
}
