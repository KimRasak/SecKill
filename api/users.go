package api

import (
	"SecKill/api/redisService"
	"SecKill/data"
	"SecKill/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	myjwt "SecKill/middleware/jwt"
)

// Visible for testing
const ErrMsgKey = "errMsg"
const DataKey = "data"


// 秒杀优惠券
func FetchCoupon(ctx *gin.Context)  	{
	// 登陆检查token
	claims := ctx.MustGet("claims").(*myjwt.CustomClaims)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{ErrMsgKey: "Not Authorized."})
		return
	}

	if claims.Kind == "saler"{//user.IsSeller() {
		ctx.JSON(http.StatusUnauthorized, gin.H{ErrMsgKey: "Sellers aren't allowed to get coupons."})
		return
	}

	paramSellerName := ctx.Param("username")
	paramCouponName := ctx.Param("name")

	// ---用户抢优惠券。后面需要高并发处理---
	// 先在缓存执行原子性的秒杀操作。将原子性地完成"判断能否秒杀-执行秒杀"的步骤
	secKillRes, err := redisService.CacheAtomicSecKill(claims.Username, paramSellerName, paramCouponName)
	if err == nil {
		log.Println(fmt.Sprintf("result: %d", secKillRes))
		coupon := redisService.GetCoupon(paramCouponName)
		// 交给[协程]完成数据库写入操作
		SecKillChannel <- secKillMessage{claims.Username, coupon}
		ctx.JSON(http.StatusCreated, gin.H{ErrMsgKey: ""})
		return
	} else {
		if redisService.IsRedisEvalError(err) {
			log.Println("Server error" + err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{ErrMsgKey: err.Error()})
			return
		} else {
			log.Println("Fail to fetch coupon. " + err.Error())
			ctx.JSON(http.StatusNoContent, gin.H{})
			return
		}
		// 可在此将err输出到log.
	}
}

const (
	couponPageSize int64 = 20
)

// 工具函数 输出查询错误
func outputQueryError(ctx *gin.Context, err error) {
	if gorm.IsRecordNotFoundError(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Record not found."})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Query error."})
	}
}

// 取得合理切片范围的coupons
func getValidCouponSlice(allCoupons []model.Coupon, page int64) []model.Coupon {
	if len(allCoupons) == 0 {
		return allCoupons
	}
	couponLen := int64(len(allCoupons))
	startIndex := page * couponPageSize
	endIndex := page * couponPageSize + couponPageSize
	if startIndex < 0 {
		startIndex = 0
	} else if startIndex > couponLen {
		startIndex = couponLen
	}
	if endIndex < 1 {
		if couponLen < couponPageSize {
			endIndex = couponLen
		} else {
			endIndex = couponPageSize
		}
	} else if endIndex > couponLen {
		endIndex = couponLen
	}
	return allCoupons[startIndex:endIndex]
}

// 数据长度为空则返回204,否则返回200
func getDataStatusCode(len int) int {
	if len == 0 {
		return http.StatusNoContent
	} else {
		return http.StatusOK
	}
}

// 查询优惠券
func GetCoupons(ctx *gin.Context) {
	// 登陆检查token
	claims := ctx.MustGet("claims").(*myjwt.CustomClaims)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{ErrMsgKey: "Not Authorized."})
		return
	}

	queryUserName, queryPage := ctx.Param("username"), ctx.Query("page")

	// 检查page参数, TODO：全部下标改为从1开始
	var page int64
	var tmpPage int64
	if queryPage == "" {
		tmpPage = 1
	} else {
		var err error
		tmpPage, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Wrong format of page."})
			return
		}
	}

	// 数据库从0开始，但是查找从1开始
	page = tmpPage - 1

	fmt.Printf("Querying coupon with name %s, page %d\n", queryUserName, page)
	// TODO: 查询用户需要从缓存查
	// 查找对应用户
	queryUser := model.User{Username:queryUserName}
	queryErr := data.Db.Where(&queryUser).
		First(&queryUser).Error
	if queryErr != nil {
		outputQueryError(ctx, queryErr)
		return
	}

	// 根据用户名查找其拥有/创建的优惠券
	if queryUserName == claims.Username {
		// 查询名与用户名相同，返回查询名用户拥有的优惠券
		var allCoupons []model.Coupon
		var err error
		if allCoupons, err = redisService.GetCoupons(claims.Username); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{ErrMsgKey: "Server error"})
			return
		}

		coupons := getValidCouponSlice(allCoupons, page)

		if queryUser.IsSeller() {
			sellerCoupons := model.ParseSellerResCoupons(coupons)
			statusCode := getDataStatusCode(len(sellerCoupons))
			ctx.JSON(statusCode, gin.H{ErrMsgKey: "", DataKey: sellerCoupons})
			return
		} else if queryUser.IsCustomer() {
			customerCoupons := model.ParseCustomerResCoupons(coupons)
			statusCode := getDataStatusCode(len(customerCoupons))
			ctx.JSON(statusCode, gin.H{ErrMsgKey: "", DataKey: customerCoupons})
			return
		}
	} else {
		// 查询名与用户名不同
		if queryUser.IsCustomer() {
			// 不可查询其它顾客的优惠券
			ctx.JSON(http.StatusUnauthorized, gin.H{ErrMsgKey: "Cannot check other customer.", DataKey: []model.Coupon{}})
			return
		} else if queryUser.IsSeller() {
			// 可以查询其它商家拥有的优惠券
			var allCoupons []model.Coupon
			var err error
			if allCoupons, err = redisService.GetCoupons(queryUserName); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{ErrMsgKey: "Error when getting seller's coupons.", DataKey: allCoupons})
				return
			}
			coupons := getValidCouponSlice(allCoupons, page)

			sellerCoupons := model.ParseSellerResCoupons(coupons)
			statusCode := getDataStatusCode(len(sellerCoupons))
			ctx.JSON(statusCode, gin.H{ErrMsgKey: "", DataKey: sellerCoupons})
			return
		}
	}

}

// 商家添加优惠券
func AddCoupon(ctx *gin.Context) {
	// 登陆检查token
	claims := ctx.MustGet("claims").(*myjwt.CustomClaims)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{ErrMsgKey: "Not Authorized."})
		return
	}


	if claims.Kind == "customer"{//!user.IsSeller() {
		ctx.JSON(http.StatusUnauthorized, gin.H{ErrMsgKey: "Only sellers can create coupons."})
		return
	}

	// 检查参数

	paramUserName := ctx.Param("username")  // 注意: 该参数是网址路径参数
	var postCoupon model.ReqCoupon
	if err := ctx.BindJSON(&postCoupon); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Only receive JSON format."})
		return
	}
	couponName := postCoupon.Name
	formAmount := postCoupon.Amount
	description := postCoupon.Description
	formStock := postCoupon.Stock
	if claims.Username != paramUserName {
		ctx.JSON(http.StatusUnauthorized, gin.H{ErrMsgKey: "Cannot create coupons for other users."})
		return
	}
	amount := formAmount
	// TODO: 检查int的范围合法性
	//amount, amountErr := strconv.ParseInt(formAmount, 10, 64)
	//if amountErr != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Amount field wrong format."})
	//	return
	//}
	stock := formStock
	// TODO: 检查stock的范围合法性
	//stock, stockErr := strconv.ParseInt(formStock, 10, 64)
	//if stockErr != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Stock field wrong format."})
	//	return
	//}

	// 优惠券描述可以为空的，不需要检查长度
	//if len(couponName) == 0 || len(description) == 0 {
	//	ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Coupon name or description should not be empty."})
	//	return
	//}

	// 在数据库添加优惠券
	coupon := model.Coupon{
		Username:    claims.Username,
		CouponName:  couponName,
		Amount:      amount,
		Left:        amount,
		Stock:       stock,
		Description: description,
	}
	var err error
	err = data.Db.Create(&coupon).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Create failed. Maybe (username,coupon name) duplicates"})
		return
	}

	// 在Redis添加优惠券
	if err = redisService.CacheCouponAndHasCoupon(coupon); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{ErrMsgKey: "Create Cache failed. " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{ErrMsgKey: ""})
	return

}

// 用户注册
func RegisterUser(ctx *gin.Context) {
	var postUser model.RegisterUser
	if err := ctx.BindJSON(&postUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Only receive JSON format."})
		return
	}
	// 查看参数长度、是否为空、格式
	if len(postUser.Username) < model.MinUserNameLen {
		ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "User name too short."})
		return
	} else if len(postUser.Password) < model.MinPasswordLen {
		ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Password too short."})
		return
	} else if postUser.Kind == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Empty field of kind."})
		return
	} else if !model.IsValidKind(postUser.Kind) {
			ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Unexpected value of kind, " + postUser.Kind})
			return
	}

	// 插入用户
	user := model.User{Username: postUser.Username, Kind: postUser.Kind, Password: model.GetMD5(postUser.Password)}
	err := data.Db.Create(&user).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ErrMsgKey: "Insert user failed. Maybe user name duplicates."})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{ErrMsgKey: ""})
		return
	}
}
