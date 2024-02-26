package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"project-api/api/user/proto"
	"project-api/config"
)

type UserService struct {
	redisClient *redis.Client
}

func NewUserService() *UserService {
	return &UserService{
		redisClient: config.RedisClient,
	}
}

func (u *UserService) UserList(ctx *gin.Context, req *userProto.UserListReq) (rsp *userProto.UserListResp, err error) {
	return
}
