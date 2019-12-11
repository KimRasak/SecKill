package redisService

//// 缓存优惠券的完整信息
//func CacheUser(user model.User) (string, error) {
//	key := getCouponKeyByCoupon(user)
//	fields := map[string]interface{}{
//		"id":          user.Id,
//		"username":    user.Username,
//		"couponName":  user.CouponName,
//		"amount":      user.Amount,
//		"left":        user.Left,
//		"stock":       user.Stock,
//		"description": user.Description,
//	}
//	val, err := data.SetMapForever(key, fields)
//	return val, err
//}
