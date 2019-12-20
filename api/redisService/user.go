package redisService

import (
	"SecKill/data"
	"SecKill/model"
	"errors"
	"fmt"
	"log"
	"strconv"
)

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


// 从缓存获取用户信息
func GetUserByName(userName string) (model.User, error) {
	userKey := getUserKeyByUserName(userName)
	userInfo, err := data.GetMap(userKey, "id", "username", "kind", "password")
	if err != nil {
		log.Println("Error when getting user info. " + err.Error())
		return model.User{}, err
	}

	// 有nil即redis没缓存到这个用户或出错
	for _, val := range userInfo {
		if val == nil {
			log.Println("Got nil user info from redis.")
			return model.User{}, errors.New("Got nil user info from redis.")
		}
	}

	id, err := strconv.ParseInt(userInfo[0].(string), 10, 64)
	if err != nil {
		log.Println("Wrong type of id. " + err.Error())
		return model.User{}, err
	}

	return model.User{
		Id:       id,
		Username: userInfo[1].(string),
		Kind:     userInfo[2].(string),
		Password: userInfo[3].(string),
	}, nil
}

// 缓存用户的注册信息
func CacheUser(user model.User) (string, error) {
	key := getUserKeyByUser(user)
	fields := map[string]interface{}{
		"id":          user.Id,
		"username":    user.Username,
		"kind":		   user.Kind,
		"password":	   user.Password,
	}
	val, err := data.SetMapForever(key, fields)
	return val, err
}

func getUserKeyByUser(user model.User) string {
	return getUserKeyByUserName(user.Username)
}

func getUserKeyByUserName(userName string) string {
	return fmt.Sprintf("user-%s-info", userName)
}
