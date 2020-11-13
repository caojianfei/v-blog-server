package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/models"
)

type FormatValidateError map[string]string

type LoginForm struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

var Login = func(ctx *gin.Context) {
	var form LoginForm
	if err := ctx.ShouldBind(&form); err != nil {
		helpers.ResponseValidateError(ctx, err)
		return
	}

	// 通过 email 查询用户信息
	user := &models.User{Email: form.Email}
	if err := databases.DB.Where(user).First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			helpers.ResponseError(ctx, helpers.RecordNotFound, "账号或密码错误")
		} else {
			helpers.ResponseError(ctx, helpers.DatabaseUnknownErr, "账号查询失败")
		}
		return
	}

	encryptPassword, _ := helpers.EncryptPassword(form.Password)
	if user.Password != encryptPassword {
		helpers.ResponseError(ctx, helpers.RecordNotFound, "账号或密码错误")
		return
	}

	// 返回 token
	token, _ := helpers.GenerateJWTByUser(user)
	helpers.ResponseOk(ctx, "登录成功", &gin.H{
		"token": token,
	})
	return
}
