package redisService

import (
	"SecKill/api/dbService"
	"SecKill/data"
	"SecKill/model"
	"fmt"
	"strconv"
)


// 获取"用户持有优惠券"的key
func getHasCouponKeyByName(userName string, sellerName string, couponName string) string {
	return fmt.Sprintf("%s-has-(%s,%s)", userName, sellerName, couponName)
}

// 获取"优惠券"的key
func getCouponKeyByCoupon(coupon model.Coupon) string {
	return getCouponKeyByName(coupon.Username, coupon.CouponName)
}
func getCouponKeyByName(sellerName string, couponName string) string {
	return fmt.Sprintf("(%s,%s)", sellerName, couponName)
}

func CacheHasCoupon(username string, coupon model.Coupon) (string, error) {
	key := getHasCouponKeyByName(username, coupon.Username, coupon.CouponName)
	val, err := data.SetForever(key, 1)
	return val, err
}

func CacheFullCoupon(coupon model.Coupon) (string, error) {
	key := getCouponKeyByCoupon(coupon)
	fields := map[string]interface{}{
		"id": coupon.Id,
		"username": coupon.Username,
		"couponName": coupon.CouponName,
		"amount": coupon.Amount,
		"left": coupon.Left,
		"stock": coupon.Stock,
		"description": coupon.Description,
	}
	val, err := data.SetMapForever(key, fields)
	return val, err
}

// 缓存优惠券
func CacheCoupon(coupon model.Coupon) (string, error) {
	user, err := dbService.GetUser(coupon.Username)
	if err != nil {
		println("Error when getting user.")
	}
	if user.IsCustomer() {
		// 缓存拥有优惠券消息
		return CacheHasCoupon(user.Username, coupon)
	} else if user.IsSeller() {
		// 缓存完整优惠券消息
		return CacheFullCoupon(coupon)
	}
	panic(fmt.Sprintf("Wrong type of user. %s %d", user.Username, user.Kind))
}

func GetCoupon(sellerName string, couponName string) model.Coupon {
	key := getCouponKeyByName(sellerName, couponName)
	values, err := data.GetMap(key, "id", "username", "couponName", "amount", "left", "stock", "description")
	if err != nil {
		println("Error on getting coupon. " + err.Error())
	}

	id, err := strconv.ParseInt(values[0].(string), 10, 64)
	if err != nil {
		println("Wrong type of id. " + err.Error())
	}
	amount, err := strconv.ParseInt(values[3].(string), 10, 64)
	if err != nil {
		println("Wrong type of id. " + err.Error())
	}
	left, err := strconv.ParseInt(values[4].(string), 10, 64)
	if err != nil {
		println("Wrong type of id. " + err.Error())
	}
	stock, err := strconv.ParseInt(values[5].(string), 10, 64)
	if err != nil {
		println("Wrong type of id. " + err.Error())
	}
	return model.Coupon{
		Id:          id,
		Username:    values[1].(string),
		CouponName:  values[2].(string),
		Amount:      amount,
		Left:        left,
		Stock:       stock,
		Description: values[6].(string),
	}

}