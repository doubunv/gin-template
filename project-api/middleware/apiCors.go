package middleware

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"project-api/common/tool"
)

func ApiCors() gin.HandlerFunc {
	defer func() {
		if err := recover(); err != nil {
			tool.Error(nil, fmt.Sprintf("%v", err))
		}
	}()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}                                       // 允许所有来源访问
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // 允许的请求方法
	corsConfig.AllowHeaders = []string{"*"}                                       // 允许的请求头
	corsConfig.AllowCredentials = true
	return cors.New(corsConfig)
}
