package dbService

import (
	"SecKill/data"
	"SecKill/model"
)


func GetAllUsers() ([]model.User, error) {
	var users []model.User
	result := data.Db.Find(&users)
	return users, result.Error
}

func GetUser(userName string) (model.User, error) {
	user := model.User{}
	operation := data.Db.Where("username = ?", userName).First(&user)
	return user, operation.Error
}