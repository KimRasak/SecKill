package httptest

import (
	"SecKill/api"
	"SecKill/data"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

const loginPath = "/api/users/"

// 测试登录不存在的用户或错误的密码
func testWrongLogin(e *httpexpect.Expect) {
	wrongUserName := "wrongUserName"
	e.POST(loginPath).
		WithForm(LoginForm{wrongUserName, "whatever_pw"}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueNotEqual(api.ErrMsgKey, "No such queryUser.")

	wrongPassword := "sysucs515"
	e.POST(loginPath).
		WithForm(LoginForm{demoSellerName, wrongPassword}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueNotEqual(api.ErrMsgKey, "Password mismatched.")
}

// 测试登录demo商家和demo顾客
func testUsersLogin(e *httpexpect.Expect) {
	demoSellerLogin(e)
	demoCustomerLogin(e)
}

func testAddAndGetCoupon(e *httpexpect.Expect) {
	veryLargePage := 10000

	// 作为商家和用户分别获取一次优惠券信息
	// 顾客查询顾客/商家的优惠券
	demoCustomerLogin(e)
	// 自己没抢过优惠券，查询不到
	isEmptyBody(e, demoCustomerName, -1)
	isEmptyBody(e, demoCustomerName, 0)
	isEmptyBody(e, demoCustomerName, veryLargePage)
	// 商家没创建过优惠券，查询不到
	isEmptyBody(e, demoSellerName, -1)
	isEmptyBody(e, demoSellerName, 0)
	isEmptyBody(e, demoSellerName, veryLargePage)
	// 商家查询商家的优惠券
	demoSellerLogin(e)
	isEmptyBody(e, demoSellerName, -1)
	isEmptyBody(e, demoSellerName, 0)
	isEmptyBody(e, demoSellerName, veryLargePage)

	// 创建demo优惠券
	demoAddCoupon(e)

	// 顾客查询该商家创建的优惠券信息
	demoCustomerLogin(e)
	isNonEmptyCoupons(e, demoSellerName, -1)
	isNonEmptyCoupons(e, demoSellerName, 0)
	isEmptyBody(e, demoSellerName, veryLargePage)
	isCustomerSchema(e, demoSellerName, 0)

	// 自己没抢过优惠券，查询不到
	isEmptyBody(e, demoCustomerName, -1)
	isEmptyBody(e, demoCustomerName, 0)
	isEmptyBody(e, demoCustomerName, veryLargePage)
	// 商家查询到自己创建的优惠券信息
	demoSellerLogin(e)
	isNonEmptyCoupons(e, demoSellerName, -1)
	isNonEmptyCoupons(e, demoSellerName, 0)
	isEmptyBody(e, demoSellerName, veryLargePage)
	isSellerSchema(e, demoSellerName, 0)
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

	// 测试查看、添加优惠券功能
	testAddAndGetCoupon(e)
}
