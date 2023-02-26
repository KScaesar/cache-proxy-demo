package api

import (
	"net/http"

	"github.com/KScaesar/cache-proxy-demo/ddd/application"
)

type UserHandler struct {
	userService application.UserService
}

func (h *UserHandler) SignInUser(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	account := request.PostFormValue("account")
	password := request.PostFormValue("password")

	h.userService.SignInUser(ctx, account, password)
}
