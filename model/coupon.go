package model

const (
	USER_MAX_AMOUNT = 1
)

// 数据库实体
type Coupon struct {
	Id          int       `gorm:"primary_key;auto_increment"`
	Username    string    `gorm:"type:varchar(20); not null"` // 用户名
	CouponName  string    `gorm:"type:varchar(60); not null"` // 优惠券名称
	Amount      int64     `gorm:"type:varchar(32)"`           // 最大优惠券数
	Left        int64								          // 剩余优惠券数
	Stock       int64                                         // 面额
	Description string  `gorm:"type:varchar(60)"`             // 优惠券描述信息
}

type ResCoupon struct {
	Name            string  `json:"name"`
	Stock           int64     `json:"stock"`
	Description     string  `json:"description"`
}

// 商家查询优惠券时，返回的数据结构
type SellerResCoupon struct {
	ResCoupon
	TotalAmount int64  `json:"total_amount"`
	Left        int64  `json:"left"`
}

// 顾客查询优惠券时，返回的数据结构
type CustomerResCoupon struct {
	ResCoupon
}

func ParseSellerResCoupons(coupons []Coupon) []SellerResCoupon {
	sellerCoupons := []SellerResCoupon{}
	for _, coupon := range coupons {
		sellerCoupons = append(sellerCoupons,
			SellerResCoupon{ResCoupon{coupon.CouponName, coupon.Amount, coupon.Description},
				coupon.Amount, coupon.Left})
	}
	return sellerCoupons
}

func ParseCustomerResCoupons(coupons []Coupon) []CustomerResCoupon {
	var sellerCoupons []CustomerResCoupon
	for _, coupon := range coupons {
		sellerCoupons = append(sellerCoupons,
			CustomerResCoupon{ResCoupon{coupon.CouponName, coupon.Amount, coupon.Description}})
	}
	return sellerCoupons
}