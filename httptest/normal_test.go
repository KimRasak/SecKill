package httptest

import (
	"SecKill/api"
	"SecKill/data"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

func testWrongLogin(e *httpexpect.Expect) {
	wrongUserName := "wrongUserName"
	e.POST("/api/users/").
		WithForm(LoginForm{wrongUserName, "whatever_pw"}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueNotEqual(api.ErrMsgKey, "No such queryUser.")

	wrongPassword := "sysucs515"
	e.POST("/api/users/").
		WithForm(LoginForm{demoSellerName, wrongPassword}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueNotEqual(api.ErrMsgKey, "Password mismatched.")
}

func testUsersLogin(e *httpexpect.Expect) {
	demoSellerLogin(e)
	demoCustomerLogin(e)
}

// 进行普通的测试，用户注册、登录后进行常规操作
func TestNormal(t *testing.T) {
	_, e := startServer(t)
	defer data.Close()

	// 注册用户,商家
	registerDemoUsers(e)

	// 用户登录错误
	testWrongLogin(e)

	// 用户登录
	testUsersLogin(e)

	// 用户创建优惠券

	//
}
