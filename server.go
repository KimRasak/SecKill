package main

import (
	"SecKill/api"
	"SecKill/model"
	"SecKill/conf"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	router := gin.Default()

	// 设置Redis存储
	config, err := conf.GetAppConfig()
	if err != nil {
		panic("failed to load redis config" + err.Error())
	}
	store, _ := redis.NewStore(config.App.Redis.MaxIdle, config.App.Redis.Network,
		config.App.Redis.Address, config.App.Redis.Password, []byte("seckill"))
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

