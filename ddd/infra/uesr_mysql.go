package infra

import (
	"context"

	"github.com/KScaesar/cache-proxy-demo/ddd/domain"
)

func NewUserMysql() *UserMysql {
	return &UserMysql{}
}

type UserMysql struct{}

func (u UserMysql) QueryUserList(ctx context.Context, dto *domain.DtoQryUserOption) (domain.DtoUserListResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (u UserMysql) QueryUserByAccount(ctx context.Context, account string) (domain.DtoUserResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (u UserMysql) GetUserByAccount(ctx context.Context, account string) (domain.User, error) {
	// TODO implement me
	panic("implement me")
}
