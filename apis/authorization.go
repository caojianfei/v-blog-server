package apis

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/models"
)

type FormatValidateError map[string]string

// 400 响应
func ResponseFormValidateError(ctx *gin.Context, err error) {
	validateErrors := make(FormatValidateError)
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, err := range err {
			//fmt.Println("Tag ", err.Tag())
			//fmt.Println("ActualTag ", err.ActualTag())
			//fmt.Println("Namespace ", err.Namespace())
			//fmt.Println("StructNamespace ", err.StructNamespace())
			//fmt.Println("Field ", err.Field())
			//fmt.Println("StructField ", err.StructField())
			//fmt.Println("Value ", err.Value())
			//fmt.Println("Param ", err.Param())
			//fmt.Println("Kind ", err.Kind())
			//fmt.Println("Type ", err.Type())
			//fmt.Println(err)
			validateErrors[err.Field()] = err.Translate(helpers.Trans)
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, validateErrors)
		return
	} else {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{},
		)
		return
	}
}


type ResponseBody struct {
	Code int
	Message string
	Data *gin.H
}
// 200 响应
func Response(ctx *gin.Context, body ResponseBody) {
	responseBody := gin.H{}
	responseBody["code"] = body.Code
	if body.Message == "" {
		responseBody["message"] = "success"
	} else {
		responseBody["message"] = body.Message
	}
	if body.Data == nil {
		responseBody["data"] = gin.H{}
	} else {
		responseBody["data"] = body.Data
	}
	responseBody["time"] = time.Now().Unix()

	ctx.JSON(http.StatusOK, responseBody)
}

type LoginForm struct {
	Email string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

var Login = func(ctx *gin.Context) {
	var form LoginForm
	if err := ctx.ShouldBind(&form); err != nil {
		validateErrors := make(FormatValidateError)
		for _, err := range err.(validator.ValidationErrors) {
			//fmt.Println("Tag ", err.Tag())
			//fmt.Println("ActualTag ", err.ActualTag())
			//fmt.Println("Namespace ", err.Namespace())
			//fmt.Println("StructNamespace ", err.StructNamespace())
			//fmt.Println("Field ", err.Field())
			//fmt.Println("StructField ", err.StructField())
			//fmt.Println("Value ", err.Value())
			//fmt.Println("Param ", err.Param())
			//fmt.Println("Kind ", err.Kind())
			//fmt.Println("Type ", err.Type())
			//fmt.Println(err)
			validateErrors[err.Field()] = err.Translate(helpers.Trans)
		}
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": validateErrors})
		return
	}

	// 通过 email 查询用户信息
	user := &models.User{}
	if databases.DB.Where(&models.User{Email:form.Email}).First(user).RecordNotFound() {
		ctx.JSON(http.StatusOK, gin.H{
			"code": RecordNotFound,
			"message": "账号或密码错误",
			"time": time.Now().Unix(),
			"data": gin.H{},
		})
		return
	}

	encryptPassword, _ := helpers.EncryptPassword(form.Password)
	if user.Password != encryptPassword {
		ctx.JSON(http.StatusOK, gin.H{
			"code": RecordNotFound,
			"message": "账号或密码错误",
			"time": time.Now().Unix(),
			"data": gin.H{},
		})
		return
	}

	// 返回 token
	token, _ := helpers.GenerateJWTByUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"code": Success,
		"message": "登录成功",
		"time": time.Now().Unix(),
		"data": gin.H{
			"token": token,
		},
	})

	return
}


