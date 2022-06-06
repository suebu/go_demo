package controller

import (
	"github.com/gin-gonic/gin"
	"go_demo/common"
	"go_demo/model"
	"go_demo/utils"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(context *gin.Context) {
	db := common.InitDB()
	//获取参数
	name := context.PostForm("name")
	password := context.PostForm("password")
	telephone := context.PostForm("telephone")
	//校验数据
	if len(telephone) != 11 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}

	if len(password) != 6 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}

	if len(name) == 0 {
		name = utils.RandomString(10)
	}

	if isTelephoneExit(db, telephone) {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在"})
		return
	}

	log.Println(name, password, telephone)
	//创建用户
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}

	db.Create(&newUser)

	//返回结果

	context.JSON(http.StatusOK, gin.H{
		"msg": "注册成功",
	})
}

func isTelephoneExit(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)

	if user.ID != 0 {
		return true
	}
	return false
}
