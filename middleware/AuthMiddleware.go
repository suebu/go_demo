package middleware

import (
	"github.com/gin-gonic/gin"
	"go_demo/common"
	"go_demo/model"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			context.Abort()
			return
		}

		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)

		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			context.Abort()
			return
		}

		//获取userID
		userId := claims.UserId
		db := common.InitDB()
		var user model.User
		db.First(&user, userId)

		if user.ID == 0 {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			context.Abort()
			return
		}

		//user存在，写入上下文
		context.Set("user", user)
		context.Next()

	}
}
