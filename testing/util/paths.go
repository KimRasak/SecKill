package util

type path struct {
	FetchCoupon   string
	GetCoupons    string
	AddCoupon     string
	LoginAuth     string
	Register      string
}


var defaultPath path = path{
	FetchCoupon: "/api/users/{username}/coupons/{name}",
	GetCoupons:  "/api/users/{username}/coupons",
	AddCoupon:   "/api/users/{username}/coupons",
	Register:    "/api/users",
	LoginAuth:   "/api/auth",
}
