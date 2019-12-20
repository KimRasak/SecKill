package engine

import (
	"SecKill/api"
	"SecKill/api/controller"
	"SecKill/data"
	"SecKill/middleware/jwt"
	"SecKill/model"
	"encoding/gob"
	"github.com/gin-gonic/gin"
)

func SeckillEngine() *gin.Engine {
	router := gin.New()
	gob.Register(&model.User{})

	// 设置路由
	userRouter := router.Group("/api/users")
	userRouter.POST("", controller.RegisterUser)
	userRouter.Use(jwt.JWTAuth())
	{
		userRouter.PATCH("/:username/coupons/:name", controller.FetchCoupon)
		userRouter.GET("/:username/coupons", controller.GetCoupons)
		userRouter.POST("/:username/coupons", controller.AddCoupon)
	}

	authRouter := router.Group("/api/auth")
	{
		authRouter.POST("", controller.LoginAuth)
		authRouter.POST("/logout", controller.Logout)
	}

	testRouter := router.Group("/testing")
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
