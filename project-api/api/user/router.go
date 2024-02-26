package user

import (
	"github.com/gin-gonic/gin"
	"project-api/internal/handler"
	"project-api/router"
)

func init() {
	router.RegisterApi(&RouterDemo{})
}

type RouterDemo struct {
}

func (*RouterDemo) Router(r *gin.Engine) {
	r.POST("/api/user/v1/list", handler.UserList)
}
