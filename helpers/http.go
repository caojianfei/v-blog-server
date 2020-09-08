package helpers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"time"
	"unicode"
)

const (
	Success = 0

	// 数据库相关
	RecordNotFound     = 100
	RecordExist        = 101
	RecordCreatedFail  = 102
	RecordUpdateFail   = 103
	RecordDeleteFail   = 104
	DatabaseUnknownErr = 110

	// 请求参数相关
	RequestParamError = 200

	// 文件
	PathBaseError      = 300
	PathCreateFail     = 301
	FileUploadFail     = 302
	UploadParamInvalid = 303
	UploadFileEmpty    = 304
	UploadFileInvalid  = 305
)

type HttpResponse struct {
	Code    int
	Message string
	Data    *gin.H
}

// http 响应
func (resp *HttpResponse) Json(ctx *gin.Context) {
	result := gin.H{}
	result["code"] = resp.Code
	result["data"] = resp.Data
	if resp.Message == "" {
		if resp.Code == 0 {
			result["message"] = "success"
		} else {
			result["message"] = "error"
		}
	} else {
		result["message"] = resp.Message
	}
	if resp.Data == nil {
		result["data"] = gin.H{}
	} else {
		result["data"] = resp.Data
	}

	result["time"] = time.Now().Unix()

	ctx.JSON(http.StatusOK, result)
}

// 业务正常响应
func ResponseOk(ctx *gin.Context, message string, data *gin.H) {
	(&HttpResponse{Code: Success, Message: message, Data: data}).Json(ctx)
}

// 业务正常响应，data 数据为空
func ResponseOkWithoutData(ctx *gin.Context, message string) {
	(&HttpResponse{Code: Success, Message: message}).Json(ctx)
}

// 业务异常响应
func ResponseError(ctx *gin.Context, code int, message string) {
	(&HttpResponse{Code: code, Message: message}).Json(ctx)
}

// 业务异常响应，data 数据不为空
func ResponseErrorWithData(ctx *gin.Context, code int, message string, data *gin.H) {
	(&HttpResponse{Code: code, Message: message, Data: data}).Json(ctx)
}

func ResponseValidateError(ctx *gin.Context, err error) {
	errorList := gin.H{}
	switch err.(type) {
	case validator.ValidationErrors:
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
				//fmt.Println(err.Translate(Trans))
				errorList[ToLowerCase(err.Field())] = err.Translate(Trans)
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errorList, "message": "参数错误"})
		}
	case *json.UnmarshalTypeError:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": gin.H{}, "message": "参数错误"})
	default:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": gin.H{}, "message": "参数错误"})
	}
}

// 首字母小写
func ToLowerCase(s string) string {
	if s == "" {
		return s
	}
	sRune := []rune(s)
	firstLetter := sRune[0]
	if unicode.IsUpper(firstLetter) {
		sRune[0] += 32
	}

	return string(sRune)
}

// 从请求参数中获取 id
func GetIdFromParam(ctx *gin.Context) (int, error) {
	idStr := ctx.Param("id")
	if idStr == "" {
		return 0, errors.New("no id param")
	}

	return strconv.Atoi(idStr)
}
