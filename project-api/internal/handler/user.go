package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project-api/api/user/proto"
	"project-api/common"
	"project-api/internal/logic"
)

func UserList(ctx *gin.Context) {
	result := &common.Result{}

	//step 1
	var req = &userProto.UserListReq{}
	err := ctx.BindJSON(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, result.Fail(ctx, err.Error()))
		return
	}

	//step 2
	s := logic.NewUserService()
	list, err := s.UserList(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, result.Fail(ctx, err.Error()))
		return
	}

	//step 3
	ctx.JSON(http.StatusOK, result.Success(ctx, list))
}
