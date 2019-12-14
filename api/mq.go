package api

import (
	"SecKill/api/dbService"
	"SecKill/model"
)

type secKillMessage struct {
	username string
	coupon model.Coupon
}

const maxMessageNum = 20000
var SecKillChannel = make(chan secKillMessage, maxMessageNum)

func seckillConsumer() {
	for {
		message := <- SecKillChannel
		println("Got one message: " + message.username)

		username := message.username
		sellerName := message.coupon.Username
		couponName := message.coupon.CouponName

		var err error
		err = dbService.UserHasCoupon(username, message.coupon)
		if err != nil {
			println("Error when inserting user's coupon. " + err.Error())
		}
		err = dbService.DecreaseOneCouponLeft(sellerName, couponName)
		if err != nil {
			println("Error when decreasing coupon left. " + err.Error())
		}
	}

}

var isConsumerRun = false
func RunSecKillConsumer() {
	// Only Run one consumer.
	if !isConsumerRun {
		go seckillConsumer()
		isConsumerRun = true
	}
}