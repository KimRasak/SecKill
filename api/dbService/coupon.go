package dbService

import (
	"SecKill/data"
	"SecKill/model"
)

func GetAllCoupons() ([]model.Coupon, error) {
	var coupons []model.Coupon
	result := data.Db.Find(&coupons)
	return coupons, result.Error
}

// 插入用户拥有优惠券的数据
func UserHasCoupon(userName string, coupon model.Coupon) error {
	operation := data.Db.Raw("INSERT IGNORE INTO coupons " +
		"(username, coupon_name, amount, left, stock, description) " +
		"values(?, ?, ?, ?, ?, ?)",
		userName, coupon.CouponName, 1, 1, coupon.Stock, coupon.Description)
	return operation.Error
}

// 优惠券库存自减1
func DecreaseOneCouponLeft(sellerName string, couponName string) error {
	coupon := model.Coupon{Username:sellerName, CouponName:couponName}
	return data.Db.Model(&coupon).
		Where("left > 0").
		UpdateColumn("left", "left-1").Error
}