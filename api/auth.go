package api

import (
	"SecKill/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"google.golang.org/appengine/log"
	"net/http"
)

import "SecKill/data"

const kindKey = "kind"
// 用户登录
func LoginAuth(ctx *gin.Context)  {
	username := ctx.PostForm("username")
	password := model.GetMD5(ctx.PostForm("password"))

	// 查找该用户
	queryUser := model.User{Username: username}
	err := data.Db.Where(&queryUser).
		First(&queryUser).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		ctx.JSON(http.StatusUnauthorized, gin.H{kindKey: "", ErrMsgKey: "No such queryUser."})
		return
	}

	// 匹配密码
	if queryUser.Password != password {
		ctx.JSON(http.StatusUnauthorized, gin.H{kindKey: queryUser.Kind, ErrMsgKey: "Password mismatched."})
		return
	}

	// 保存Session
	user := queryUser
	session := sessions.Default(ctx)
	session.Set("user", user)
	err = session.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{kindKey: queryUser.Kind, ErrMsgKey: "Save session failed."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{kindKey: queryUser.Kind, ErrMsgKey: ""})
	return
}

// 用户登出
func Logout(ctx *gin.Context)  {
	session := sessions.Default(ctx)
	session.Delete("user")
	if err := session.Save(); err != nil {
		log.Warningf(ctx, "Error when save deleted session. %v", err.Error())
	}


	ctx.JSON(http.StatusOK, gin.H{ErrMsgKey: "log out."})
	return
}