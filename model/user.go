package model

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	MinUserNameLen = 4
	MinPasswordLen = 4
	NormalCustomer = "customer"
	NormalSeller   = "saler"
)
type User struct {
	Id       int     `gorm:"primary_key;auto_increment"`
	Username string  `gorm:"type:varchar(20)"`
	Kind     string  `gorm:"type:varchar(20)"`
	Password string  `gorm:"type:varchar(32)"`
}

func (user User)IsCustomer() bool {
	return user.Kind == NormalCustomer
}

func (user User)IsSeller() bool {
	return user.Kind == NormalSeller
}

func IsValidKind(kind string) bool {
	return kind == NormalCustomer || kind == NormalSeller
}

func GetMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

