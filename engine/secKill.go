package engine

import (
	"SecKill/api"
	"SecKill/conf"
	"SecKill/data"
	"SecKill/middleware/jwt"
	"SecKill/model"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// Visible for test
const SessionHeaderKey =  "Authorization"

func SeckillEngine() *gin.Engine {
	router := gin.Default()

	// 设置session为Redis存储
	config, err := conf.GetAppConfig()
	if err != nil {
		panic("failed to load redisService config" + err.Error())
	}
	store, _ := redis.NewStore(config.App.Redis.MaxIdle, config.App.Redis.Network,
		config.App.Redis.Address, config.App.Redis.Password, []byte("seckill"))
	router.Use(sessions.Sessions(SessionHeaderKey, store))
	gob.Register(&model.User{})

	// 设置路由
	userRouter := router.Group("/api/users")
	userRouter.POST("/", api.RegisterUser)
	userRouter.Use(jwt.JWTAuth())
	{
		userRouter.PATCH("/:username/coupons/:name", api.FetchCoupon)
		userRouter.GET("/:username/coupons", api.GetCoupons)
		userRouter.POST("/:username/coupons", api.AddCoupon)
	}

	authRouter := router.Group("/api/auth")
	{
		authRouter.POST("/", api.LoginAuth)
		authRouter.POST("/logout", api.Logout)
	}

	testRouter := router.Group("/test")
	{
		testRouter.GET("/", api.Welcome)
		testRouter.GET("/flush", func(context *gin.Context) {
			if _, err := data.FlushAll(); err != nil {
				println("Error when flushAll. " + err.Error())
			} else {
				println("Flushall succeed.")
			}
		})
	}

	// 启动秒杀功能的消费者
	api.RunSecKillConsumer()

	return router
}
