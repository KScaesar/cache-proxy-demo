package application

import (
	"context"

	"github.com/KScaesar/cache-proxy-demo/ddd/domain"
)

type UserUserCase struct {
	userRepo domain.UserRepository
}

func (uc *UserUserCase) SignInUser(ctx context.Context, account, password string) error {
	// TODO implement me
	panic("implement me")
}

func (uc *UserUserCase) QueryUserList(ctx context.Context, dto *domain.DtoQryUserOption) (domain.DtoUserListResponse, error) {
	err := dto.Validate()
	if err != nil {
		return domain.DtoUserListResponse{}, err
	}

	return uc.userRepo.QueryUserList(ctx, dto)
}
