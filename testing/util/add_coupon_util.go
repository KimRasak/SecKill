package util

import (
	"SecKill/api"
	"SecKill/middleware/jwt"
	"SecKill/model"
	"github.com/gavv/httpexpect"
	"net/http"
)

type stdAddCouponBody struct {
	Name        string  `form:"name"`
	Amount      int64   `form:"amount"`
	Description string  `form:"description"`
	Stock       int64   `form:"stock"`
}

// convert coupon to "add coupon" body
func couponToAddCouponBody(coupon model.Coupon) stdAddCouponBody {
	return stdAddCouponBody{
		Name:        coupon.CouponName,
		Amount:      coupon.Amount,
		Description: coupon.Description,
		Stock:       coupon.Stock,
	}
}

// 测试添加demo优惠券成功
func AddCouponSucess(e *httpexpect.Expect, authorization string, coupon model.Coupon)  {
	addCouponBody := couponToAddCouponBody(coupon)

	e.POST(defaultPath.AddCoupon, coupon.Username).
		WithJSON(addCouponBody).
		WithHeader(jwt.AuthorizationKey, authorization).
		Expect().
		Status(http.StatusCreated).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")
}