package redisService

import (
	"SecKill/data"
	"SecKill/model"
	"fmt"
)


// 获取"用户持有优惠券"的key
func getUserHasCouponKey(userName string, sellerName string, couponName string) string {
	return fmt.Sprintf("%s-has-(%s,%s)", userName, sellerName, couponName)
}

// 获取"优惠券库存"的key
func getCouponLeftKey(sellerName string, couponName string) string {
	return fmt.Sprintf("(%s,%s)-left", sellerName, couponName)
}


func CacheCoupon(coupon model.Coupon) (string, error) {
	sellerName := coupon.Username
	couponName := coupon.CouponName
	left := coupon.Left
	key := getCouponLeftKey(sellerName, couponName)
	res, err := data.SetForever(key, left)
	return res, err
}
