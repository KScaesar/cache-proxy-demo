package application

import (
	"context"

	"github.com/KScaesar/cache-proxy-demo/ddd/domain"
)

type UserService interface {
	SignInUser(ctx context.Context, account, password string) error
	QueryUserList(ctx context.Context, dto *domain.DtoQryUserOption) (domain.DtoUserListResponse, error)
}
