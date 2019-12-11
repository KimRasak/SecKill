package main

import (
	"SecKill/api"
	"SecKill/conf"
	"SecKill/data"
	"SecKill/model"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	router := gin.Default()

	// 设置session为Redis存储
	config, err := conf.GetAppConfig()
	if err != nil {
		panic("failed to load redisService config" + err.Error())
	}
	store, _ := redis.NewStore(config.App.Redis.MaxIdle, config.App.Redis.Network,
		config.App.Redis.Address, config.App.Redis.Password, []byte("seckill"))
	router.Use(sessions.Sessions("mysession", store))
	gob.Register(&model.User{})

	// 设置路由
	userRouter := router.Group("/api/users")
	{
		userRouter.PATCH("/:username/coupons/:name", api.FetchCoupon)
		userRouter.GET("/:username/coupons", api.GetCoupons)
		userRouter.POST("/:username/coupons", api.AddCoupon)
		userRouter.POST("/", api.RegisterUser)
	}

	authRouter := router.Group("/api/auth")
	{
		authRouter.POST("/", api.LoginAuth)
	}

	testRouter := router.Group("/test")
	{
		testRouter.GET("/", api.Welcome)
	}

	// 启动秒杀功能的消费者
	api.RunSecKillConsumer()

	router.Run(":8000")
	defer data.Close()
}

