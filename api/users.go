package api

import (
	"SecKill/api/redisService"
	"SecKill/data"
	"SecKill/model"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

const errMsgKey = "errMsg"


func FetchCoupon(ctx *gin.Context)  {
	// 登陆检查
	session := sessions.Default(ctx)
	sessionUser := session.Get("user")
	if sessionUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{errMsgKey: "Not Logged in."})
		return
	}
	user := sessionUser.(*model.User)
	if user.IsSeller() {
		ctx.JSON(http.StatusUnauthorized, gin.H{errMsgKey: "Sellers aren't allowed to get coupons."})
		return
	}

	paramSellerName := ctx.Param("username")
	paramCouponName := ctx.Param("name")

	// ---用户抢优惠券。后面需要高并发处理---

	// 先在缓存执行原子性的秒杀操作。将原子性地完成"判断能否秒杀-执行秒杀"的步骤
	secKillRes, err := redisService.CacheAtomicSecKill(user.Username, paramSellerName, paramCouponName)
	if err == nil {
		print(fmt.Sprintf("result: %d", secKillRes))
		coupon := redisService.GetCoupon(paramSellerName, paramCouponName)
		// 交给[协程]完成数据库写入操作
		SecKillChannel <- secKillMessage{user.Username, coupon}
		// TODO:
		// 1. 把用户拥有优惠券的行为存到数据库
		// 2. 将数据库里优惠券的库存-1
		// 可以建立一个带缓冲的channel
		// 传输的信息要包含user.Username, paramSellerName, paramCouponName
		ctx.JSON(http.StatusOK, gin.H{errMsgKey: ""})
	} else {
		println("Cache secKill error. " + err.Error())
		// 可在此将err输出到log.
	}

	// ----此后是数据库实现的抢购----
	// 查看是否有该优惠券，若没有，直接返回错误
	/* coupon := model.Coupon{Username: paramSellerName, CouponName: paramCouponName}
	err := data.Db.Where(&coupon).First(&coupon).Error
	if err != nil {
		outputQueryError(ctx, err)
		return
	}

	// 用户没抢过，才能加入抢购
	hasCoupon := model.Coupon{Username: user.Username, CouponName:paramCouponName}
	err = data.Db.Where(&hasCoupon).First(&hasCoupon).Error
	if err == nil || !gorm.IsRecordNotFoundError(err) {
		ctx.JSON(http.StatusUnauthorized, gin.H{errMsgKey: "Customer has got the coupon."})
		return
	}

	// 抢优惠券
	updateExec := data.Db.Exec(fmt.Sprintf("UPDATE ginhello.coupons c SET c.left=c.left-1 WHERE " +
		"c.username='%s' AND c.coupon_name='%s' AND c.left>0", paramSellerName, paramCouponName))
	updateErr := updateExec.Error
	updatedNum := updateExec.RowsAffected
	if updateErr != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{errMsgKey: "Get coupon failed. Db error."})
		return
	}
	if updatedNum == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{errMsgKey: "Get coupon failed. No such coupon or no coupon rest."})
		return
	}

	// 先直接返回
	ctx.JSON(http.StatusUnauthorized, gin.H{errMsgKey: "Please wait for the result."})

	// 抢到后加入用户拥有的优惠券
	userCoupon := model.Coupon{Username: user.Username, CouponName:paramCouponName,
		Amount:1, Left:1, Stock:coupon.Stock, Description:coupon.Description}
	err = data.Db.Create(&userCoupon).Error
	if err != nil {
		// 应再次尝试创建
	} */
}

const (
	couponPageSize = 20
)

// 工具函数 输出查询错误
func outputQueryError(ctx *gin.Context, err error) {
	if gorm.IsRecordNotFoundError(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Record not found."})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Query error."})
	}
}

// 查询用户的优惠券
func GetCoupons(ctx *gin.Context) {
	// 登陆检查
	session := sessions.Default(ctx)
	sessionUser := session.Get("user")
	if sessionUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not Logged in."})
		return
	}
	user := sessionUser.(*model.User)

	queryUserName, queryPage := ctx.Param("username"), ctx.Query("page")

	// 检查page参数
	var page int64
	if queryPage == "" {
		page = 0
	} else {
		var err error
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Wrong format of page."})
			return
		}
	}

	fmt.Printf("Querying coupon with name %s, page %d\n", queryUserName, page)
	// 查找对应用户
	queryUser := model.User{Username:queryUserName}
	queryErr := data.Db.Where(&queryUser).
		First(&queryUser).Error
	if queryErr != nil {
		outputQueryError(ctx, queryErr)
		return
	}

	// 根据用户名查找他的优惠券
	if queryUserName == user.Username {
		// 查询名与用户名相同，返回查询名用户的优惠券
		var coupons []model.Coupon
		data.Db.Offset(page *couponPageSize).Limit(couponPageSize).
			Find(&coupons, model.Coupon{Username:queryUserName})

		if queryUser.IsSeller() {
			sellerCoupons := model.ParseSellerResCoupons(coupons)
			ctx.JSON(http.StatusOK, sellerCoupons)
			return
		} else if queryUser.IsCustomer() {
			customerCoupons := model.ParseCustomerResCoupons(coupons)
			ctx.JSON(http.StatusOK, customerCoupons)
			return
		}
	} else {
		// 查询名与用户名不同
		if queryUser.IsCustomer() {
			// 不可查询其它顾客
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Cannot check other customer."})
			return
		} else if queryUser.IsSeller() {
			// 可以查询其它商家
			var coupons []model.Coupon
			data.Db.Offset(page *couponPageSize).Limit(couponPageSize).
				Find(&coupons, model.Coupon{Username:queryUserName})
			sellerCoupons := model.ParseSellerResCoupons(coupons)
			ctx.JSON(http.StatusOK, sellerCoupons)
			return
		}
	}

}

// 商家添加优惠券
func AddCoupon(ctx *gin.Context) {
	// 登陆检查
	session := sessions.Default(ctx)
	sessionUser := session.Get("user")
	if sessionUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{errMsgKey: "Not Logged in."})
		return
	}
	user := sessionUser.(*model.User)
	if !user.IsSeller() {
		ctx.JSON(http.StatusUnauthorized, gin.H{errMsgKey: "Only sellers can create coupons."})
		return
	}

	// 检查参数
	paramUserName := ctx.Param("username")
	couponName := ctx.PostForm("name")
	formAmount := ctx.PostForm("amount")
	description := ctx.PostForm("description")
	formStock := ctx.PostForm("stock")
	if user.Username != paramUserName {
		ctx.JSON(http.StatusUnauthorized, gin.H{errMsgKey: "Cannot create coupons for other users."})
		return
	}
	amount, amountErr := strconv.ParseInt(formAmount, 10, 64)
	if amountErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Amount field wrong format."})
		return
	}
	stock, stockErr := strconv.ParseInt(formStock, 10, 64)
	if stockErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Stock field wrong format."})
		return
	}
	if len(couponName) == 0 || len(description) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Coupon name or description should not be empty."})
		return
	}

	// 在数据库添加优惠券
	coupon := model.Coupon{
		Username:    user.Username,
		CouponName:  couponName,
		Amount:      amount,
		Left:        amount,
		Stock:       stock,
		Description: description,
	}
	err := data.Db.Create(&coupon).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Create failed. Maybe (username,coupon name) duplicates"})
		return
	}

	// 在Redis添加优惠券
	res, err := redisService.CacheCoupon(coupon)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Create Cache failed. result: " + res})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{errMsgKey: ""})
	return

}

// 用户注册
func RegisterUser(ctx *gin.Context) {
	postUserName, postKind, postPassword := ctx.PostForm("username"), ctx.PostForm("kind"), ctx.PostForm("password")

	// 查看参数长度、是否为空、格式
	var kind int64
	if len(postUserName) < model.MinUserNameLen {
		ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "User name too short."})
		return
	} else if len(postPassword) < model.MinPasswordLen {
		ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Password too short."})
		return
	} else if postKind == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Empty field of kind."})
		return
	} else {
		var err error
		kind, err = strconv.ParseInt(postKind, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Bad form of kind."})
			return
		} else if !model.IsValidKind(kind) {
			ctx.JSON(http.StatusBadRequest, gin.H{errMsgKey: "Unexpected value of kind, " + postKind})
			return
		}
	}

	// 插入用户
	username, password := postUserName, model.GetMD5(postPassword)
	user := model.User{Username: username, Kind: kind, Password: password}
	err := data.Db.Create(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{errMsgKey: "Insert user failed. Maybe user name duplicates."})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{errMsgKey: ""})
		return
	}
}
