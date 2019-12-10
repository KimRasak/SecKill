package redisService

import "SecKill/data"

const secKillScript = `
	local userHasCoupon = redis.call("get", KEYS[1]);
	if (userHasCoupon ~= false)
	then
		return -1;  --- user has had the coupon
	end

	local couponLeft = redis.call("get", KEYS[2]);
	if (couponLeft == false)
	then
		return -2;  --- No such coupon
	end
	if (couponLeft == 0)
    then
		return -3;  ---  No Coupon Left.
	end
	
   ---/ User gets the coupon
	redis.call("set", KEYS[2], couponLeft - 1);
	redis.call("set", KEYS[1], 1);
	return 1;
`
var secKillSHA string  // SHA expression of secKillScript



// 将数据加载到缓存预热，防止缓存穿透
// 预热加载了商品库存key
func preHeatKeys()  {
	coupons, err := data.GetAllCoupons()
	if err != nil {
		panic("Error when getting all coupons." + err.Error())
	}

	printResultOnce := true
	for _, coupon := range coupons {
		res, err := CacheCoupon(coupon)
		if err != nil {
			panic("Error while setting redisService keys of coupons left " + err.Error())
		}
		if printResultOnce {
			print("Set redis keys of coupons left success. " + res)
			printResultOnce = false
		}
	}
}

func init() {
	// 让redis加载秒杀的lua脚本
	secKillSHA = data.PrepareScript(secKillScript)

	// 预热
	preHeatKeys()
}
