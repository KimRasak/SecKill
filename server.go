package main

import (
	"SecKill/api"
	"SecKill/model"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	router := gin.Default()

	// 设置Redis存储
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("seckill"))
	router.Use(sessions.Sessions("mysession", store))
	gob.Register(&model.User{})

	// 设置路由
	v1 := router.Group("/api/users")
	{
		v1.PATCH("/:username/coupons/:name", api.FetchCoupon)
		v1.GET("/:username/coupons", api.GetCoupons)
		v1.POST("/:username/coupons", api.AddCoupon)
		v1.POST("/", api.RegisterUser)
	}

	v2 := router.Group("/api/auth")
	{
		v2.POST("/", api.LoginAuth)
	}
	router.Run(":8000")

}

