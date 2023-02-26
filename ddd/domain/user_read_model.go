package domain

type DtoQryUserOption struct {
	UserId string
	Name   string
}

func (dto *DtoQryUserOption) Validate() error { return nil }

type DtoUserResponse struct{}

type DtoUserListResponse struct {
	List []DtoUserResponse
}
