package httptest

/* 该文件下依赖于注册过的demo用户，需要先调用registerDemoUsers */

// 定义了添加优惠券、查看优惠券的表格的表格
const addCouponPath = "api/users/{username}/coupons"
type AddCouponForm struct {
	Username    string `form:"username"`
	Name        string `form:"name"`
	Amount      string `form:"amount"`      // 应当int
	Description string `form:"description"`
	Stock       string `form:"stock"`       // 应当int
}

// 定义了demo优惠券
demoCouponName := ""

// 测试添加优惠券时的表格格式
func testAddCouponWrongFormat() {
	amountNotNumberForm := AddCouponForm{
		Username:    demoSellerName,
		Name:        ,
		Amount:      "",
		Description: "",
		Stock:       "",
	}
}

// 测试非商家添加优惠券或为其它用户添加优惠券

// 测试未登录添加优惠券