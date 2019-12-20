package util

import (
	"SecKill/api"
	"SecKill/middleware/jwt"
	"SecKill/model"
	"github.com/gavv/httpexpect"
	"net/http"
)

func FetchCouponSuccess(e *httpexpect.Expect, authorization string, coupon model.Coupon) {

	e.PATCH(defaultPath.FetchCoupon, coupon.Username, coupon.CouponName).
		WithHeader(jwt.AuthorizationKey, authorization).
		Expect().
		Status(http.StatusCreated).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")
}

func FetchCouponFail(e *httpexpect.Expect, authorization string, coupon model.Coupon) {
	e.PATCH(defaultPath.FetchCoupon, coupon.Username, coupon.CouponName).
		WithHeader(jwt.AuthorizationKey, authorization).
		Expect().
		Status(http.StatusNoContent).
		Body().Empty()
}