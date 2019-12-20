package testing

import (
	"SecKill/data"
	"SecKill/engine"
	"SecKill/model"
	"SecKill/testing/util"
	"github.com/gavv/httpexpect"
	"net/http/httptest"
	"testing"
)

// 顾客 jzl
var jzl = model.User{
	Username: "jzl",
	Kind:     model.NormalCustomer,
	Password: "shen6508",
}

// 商家 京东
var jd = model.User{
	Username: "jd",
	Kind:     model.NormalSeller,
	Password: "jingdong",
}

// 京东特产 奶茶
var milkTea = model.Coupon{
	Username:    jd.Username,
	CouponName:  "milkTea",
	Amount:      3,
	Left:        3,
	Stock:       50,
	Description: "jingdong's milk tea.",
}


func TestNormal(t *testing.T)  {
	// 启动服务器
	server := httptest.NewServer(engine.SeckillEngine())
	e := httpexpect.New(t, server.URL)
	defer data.Close()

	util.RegSuccess(e, jzl)
	util.RegSuccess(e, jd)

	jdAuth := util.LoginSuccess(e, jd)
	util.AddCouponSucess(e, jdAuth, milkTea)
	util.MatchCustomerSchema(e, jdAuth, jd.Username, 1)

	jzlAuth := util.LoginSuccess(e, jzl)
	util.FetchCouponSuccess(e, jzlAuth, milkTea)
}