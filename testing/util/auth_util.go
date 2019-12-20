package util

import (
	"SecKill/api"
	"SecKill/model"
	"github.com/gavv/httpexpect"
	"net/http"
)

// 登录成功时header返回的凭证, 同时也是带身份操作时需要在header附上的凭证
const authorizationKey = "Authorization"

// 正确的注册消息体
type stdRegisterBody struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Kind     string  `json:"kind"`
}

// 正确的登录消息体
type stdLoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Convert user to register body.
func userToRegBody(user model.User) stdRegisterBody {
	return stdRegisterBody{
		Username: user.Username,
		Password: user.Password,
		Kind:     user.Kind,
	}
}

// Convert user to login body.
func userToLoginBody(user model.User) stdLoginBody {
	return stdLoginBody{
		Username: user.Username,
		Password: user.Password,
	}
}

// 测试注册用户成功
func regSuccess(e *httpexpect.Expect, user model.User) {
	regBody := userToRegBody(user)
	e.POST(defaultPath.Register).
		WithJSON(regBody).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")
}

// 测试用户登录成功
// 返回: 登录凭证
func loginSuccess(e *httpexpect.Expect, user model.User) string {
	loginBody := userToLoginBody(user)

	expect := e.POST(defaultPath.Register).
		WithJSON(loginBody).
		Expect()

	expect.
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")

	authorization := expect.JSON().Object().Value(authorizationKey).String().Raw()
	return authorization
}