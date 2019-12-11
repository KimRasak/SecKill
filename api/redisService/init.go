package redisService

import (
	"SecKill/api/dbService"
	"SecKill/data"
)

const secKillScript = `
    --- Check if User has coupon ---
	local userHasCoupon = redis.call("get", KEYS[1]);
	if (userHasCoupon ~= false)
	then
		return -1;
	end

    --- Check if coupon exists and is cached ---
	local couponLeft = redis.call("hget", KEYS[2], "left");
	if (couponLeft == false)
	then
		return -2;  --- No such coupon
	end
	if (couponLeft == 0)
    then
		return -3;  ---  No Coupon Left.
	end
	
    --- User gets the coupon ---
	redis.call("hset", KEYS[2], "left", couponLeft - 1);
	redis.call("set", KEYS[1], 1);
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
		_, err := CacheCoupon(coupon)
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
