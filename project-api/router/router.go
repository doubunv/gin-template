package router

import "github.com/gin-gonic/gin"

type Router interface {
	Router(r *gin.Engine)
}

type RegisterRouter struct {
}

func New() *RegisterRouter {
	return &RegisterRouter{}
}

func (*RegisterRouter) Router(ro Router, r *gin.Engine) {
	ro.Router(r)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	for _, ro := range routers {
		ro.Router(r)
	}
}

func RegisterApi(ro ...Router) {
	routers = append(routers, ro...)
}
