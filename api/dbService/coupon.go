package dbService

import (
	"SecKill/data"
	"SecKill/model"
	"fmt"
)

func GetAllCoupons() ([]model.Coupon, error) {
	var coupons []model.Coupon
	result := data.Db.Find(&coupons)
	return coupons, result.Error
}

// 插入用户拥有优惠券的数据
func UserHasCoupon(userName string, coupon model.Coupon) error {
	return data.Db.Exec(fmt.Sprintf("INSERT IGNORE INTO coupons " +
		"(`username`,`coupon_name`,`amount`,`left`,`stock`,`description`) " +
		"values('%s', '%s', %d, %d, %f, '%s')",
		userName, coupon.CouponName, 1, 1, coupon.Stock, coupon.Description)).Error
}

// 优惠券库存自减1
func DecreaseOneCouponLeft(sellerName string, couponName string) error {
	return data.Db.Exec(fmt.Sprintf("UPDATE coupons c SET c.left=c.left-1 WHERE " +
		"c.username='%s' AND c.coupon_name='%s' AND c.left>0", sellerName, couponName)).Error
}