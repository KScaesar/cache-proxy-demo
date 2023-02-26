package infra

import (
	"context"
	"errors"

	"github.com/KScaesar/cache-proxy-demo/ddd/domain"
)

type UserCacheProxy struct {
	db    *UserMysql
	cache *CacheRedis
}

func (proxy *UserCacheProxy) QueryUserList(ctx context.Context, dto *domain.DtoQryUserOption) (domain.DtoUserListResponse, error) {
	queryUserList := &CacheProxyImpl{
		Cache: proxy.cache,

		TransformReadOption: func(readDtoOption any) (key string) {
			return readDtoOption.(*domain.DtoQryUserOption).String()
		},
		ReadDataSource: func(ctx context.Context, readDtoOption any) (readModel any, err error) {
			return proxy.db.QueryUserList(ctx, readDtoOption.(*domain.DtoQryUserOption))
		},

		IsAnNotFoundError:                func(err error) bool { return errors.Is(err, ErrNotFound) },
		CanIgnoreCacheError:              false,
		CanIgnoreReadSourceErrorNotFound: true,
		CacheTTL:                         0,
	}

	readModel, err := queryUserList.ReadValue(ctx, dto)
	if err != nil {
		return domain.DtoUserListResponse{}, err
	}
	return readModel.(domain.DtoUserListResponse), nil
}
func (proxy *UserCacheProxy) QueryUserByAccount(ctx context.Context, account string) (domain.DtoUserResponse, error) {
	queryUserByAccount := &CacheProxyImpl{
		Cache: proxy.cache,

		TransformReadOption: func(readDtoOption any) (key string) {
			return readDtoOption.(string)
		},
		ReadDataSource: func(ctx context.Context, readDtoOption any) (readModel any, err error) {
			return proxy.db.QueryUserByAccount(ctx, readDtoOption.(string))
		},

		IsAnNotFoundError:                func(err error) bool { return errors.Is(err, ErrNotFound) },
		CanIgnoreCacheError:              false,
		CanIgnoreReadSourceErrorNotFound: true,
		CacheTTL:                         0,
	}

	readModel, err := queryUserByAccount.ReadValue(ctx, account)
	if err != nil {
		return domain.DtoUserResponse{}, err
	}
	return readModel.(domain.DtoUserResponse), nil
}

func (proxy *UserCacheProxy) GetUserByAccount(ctx context.Context, account string) (domain.User, error) {
	// TODO implement me
	panic("implement me")
}
