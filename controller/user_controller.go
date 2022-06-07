package controller

import (
	"github.com/gin-gonic/gin"
	"go_demo/common"
	"go_demo/model"
	"go_demo/utils"
	"golang.org/x/crypto/bcrypt"
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

	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "加密错误"})
	}

	log.Println(name, password, telephone)
	//创建用户
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(fromPassword),
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

func Login(context *gin.Context) {
	db := common.InitDB()
	//获取参数
	password := context.PostForm("password")
	telephone := context.PostForm("telephone")

	if len(telephone) != 11 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}

	if len(password) != 6 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}

	var user model.User
	db.Where("telephone = ?", telephone).First(&user)

	if user.ID == 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}

	//密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码错误"})
		return
	}

	token, err := common.ReleaseToken(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("token generate error: %V", err)
		return
	}

	context.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})

}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": user}})
}
