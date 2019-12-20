package util

import (
	"SecKill/api"
	"SecKill/api/controller"
	"SecKill/middleware/jwt"
	"fmt"
	"github.com/gavv/httpexpect"
)

var customerSchema = fmt.Sprintf(`{
	"type": "object",
	"properties": {
		"%s": {
				"type": "string"
			},
        "%s": {
				"type": "array",
				"items": {
					"type":        "object",
					"name":        "string",
					"amount":      "integer",
					"left":        "integer",
					"stock":       "integer",
					"description": "string"
				}
			}
	}
}`, api.ErrMsgKey, controller.DataKey)

var sellerSchema = fmt.Sprintf(`{
	"type": "object",
	"properties": {
		"%s": {
				"type": "string"
			},
        "%s": {
				"type": "array",
				"items": {
					"type":        "object",
					"name":        "string",
					"stock":       "integer",
					"description": "string"
				}
			}
	}
}`, api.ErrMsgKey, controller.DataKey)

// 验证符合顾客的格式
func MatchCustomerSchema(e *httpexpect.Expect, authorization string, username string, page int) {
	e.GET(defaultPath.GetCoupons, username).
		WithHeader(jwt.AuthorizationKey, authorization).
		WithQuery(controller.PageKey, page).
		Expect().JSON().Schema(customerSchema)

}

// 验证符合商家的格式
func MatchSellerSchema(e *httpexpect.Expect, authorization string, username string, page int) {
	e.GET(defaultPath.GetCoupons, username).
		WithHeader(jwt.AuthorizationKey, authorization).
		WithQuery(controller.PageKey, page).
		Expect().JSON().Schema(sellerSchema)
}

// 验证优惠券的剩余量与预期一致
func isCouponExpectedLeft(e *httpexpect.Expect, username string, page int, index int, expectedLeft int)  {
	e.GET(defaultPath.GetCoupons, username).WithQuery(controller.PageKey, page).
		Expect().JSON().Object().Value(controller.DataKey).Array().
		Element(index).Object().Value("left").Equal(expectedLeft)
}
