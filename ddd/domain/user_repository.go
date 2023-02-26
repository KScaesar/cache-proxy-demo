package domain

import "context"

type UserRepository interface {
	QueryUserList(ctx context.Context, dto *DtoQryUserOption) (DtoUserListResponse, error)
	QueryUserByAccount(ctx context.Context, account string) (DtoUserResponse, error)
	GetUserByAccount(ctx context.Context, account string) (User, error)
}
