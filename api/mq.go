package api

import (
	"SecKill/api/dbService"
	"SecKill/model"
	"log"
)

type SecKillMessage struct {
	Username string
	Coupon   model.Coupon
}

const maxMessageNum = 20000
var SecKillChannel = make(chan SecKillMessage, maxMessageNum)

func seckillConsumer() {
	for {
		message := <- SecKillChannel
		log.Println("Got one message: " + message.Username)

		username := message.Username
		sellerName := message.Coupon.Username
		couponName := message.Coupon.CouponName

		var err error
		err = dbService.UserHasCoupon(username, message.Coupon)
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