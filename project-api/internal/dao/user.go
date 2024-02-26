package dao

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"project-api/common/tool"
	"project-api/config"
	"project-api/internal/model/entity"
)

type UserDao struct {
	Db       *gorm.DB
	RedisCli *redis.Client
}

func NewUserDao() *UserDao {
	return &UserDao{
		Db:       config.MysqlClient,
		RedisCli: config.RedisClient,
	}
}

func (u *UserDao) UserQuery(ctx context.Context, filter *entity.User, limit, offset int) (user []*entity.User, err error) {
	err = u.Db.Model(&entity.User{}).Where("name like ?", "%"+filter.Name+"%").Find(&user).Limit(limit).Offset(offset).Scan(&user).Error
	if err != nil {
		tool.Error(ctx, "[dao|user] FindUser sql err:"+err.Error())
	}
	return
}

func (u *UserDao) UserCount(ctx context.Context, filter *entity.User) (count int64, err error) {
	err = u.Db.Find(&entity.User{}).Where("name like ?", "%"+filter.Name+"%").Count(&count).Error
	if err != nil {
		tool.Error(ctx, "[dao|user] CountUser sql err:"+err.Error())
	}
	return
}

func (u *UserDao) UserUpdate(ctx context.Context, filter *entity.User) (res bool, err error) {
	err = u.Db.Model(&entity.User{}).Where("id = ?", filter.Id).Updates(filter).Error
	if err != nil {
		tool.Error(ctx, "[dao|user] UpdateUser sql err:"+err.Error())
	}
	return true, err
}

func (u *UserDao) UserHardDelete(ctx context.Context, filter *entity.User) (res bool, err error) {
	err = u.Db.Find(&entity.User{}).Delete(filter).Error
	if err != nil {
		tool.Error(ctx, "[dao|user] HardDeleteUser sql err:"+err.Error())
	}
	return true, err
}
