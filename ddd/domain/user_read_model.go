package domain

import "fmt"

type DtoQryUserOption struct {
	UserId string
	Name   string
}

func (dto *DtoQryUserOption) String() string {
	return fmt.Sprintf("user:user_id=%v&name=%v", dto.UserId, dto.Name)
}

func (dto *DtoQryUserOption) Validate() error { return nil }

type DtoUserResponse struct{}

type DtoUserListResponse struct {
	List []DtoUserResponse
}
