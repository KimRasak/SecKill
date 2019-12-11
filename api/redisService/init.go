package redisService

import (
	"SecKill/api/dbService"
	"SecKill/data"
)

const secKillScript = `
    --- Check if User has coupon ---
    --- KEYS[1]: hasCouponKey "{username}-has"
    --- KEYS[2]: couponName   "{couponName}"
    --- KEYS[3]: couponKey    "{couponName}-info"
	local userHasCoupon = redis.call("SISMEMBER", KEYS[1], KEYS[2]);
	if (userHasCoupon == 1)
	then
		return -1;
	end

    --- Check if coupon exists and is cached ---
	local couponLeft = redis.call("hget", KEYS[3], "left");
	if (couponLeft == false)
	then
		return -2;  --- No such coupon
	end
	if (couponLeft == 0)
    then
		return -3;  ---  No Coupon Left.
	end
	
    --- User gets the coupon ---
	redis.call("hset", KEYS[3], "left", couponLeft - 1);
	redis.call("SADD", KEYS[1], KEYS[2]);
	return 1;
`
var secKillSHA string  // SHA expression of secKillScript



// 将数据加载到缓存预热，防止缓存穿透
// 预热加载了商品库存key
func preHeatKeys()  {
	coupons, err := dbService.GetAllCoupons()
	if err != nil {
		panic("Error when getting all coupons." + err.Error())
	}

	for _, coupon := range coupons {
		err := CacheCouponAndHasCoupon(coupon)
		if err != nil {
			panic("Error while setting redis keys of coupons. " + err.Error())
		}
	}
	print("Set redis keys of coupons success.")
}

func init() {
	// 让redis加载秒杀的lua脚本
	secKillSHA = data.PrepareScript(secKillScript)

	// 预热
	preHeatKeys()
}
