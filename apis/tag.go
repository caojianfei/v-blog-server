package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"v-blog/consts"
	"v-blog/databases"
	"v-blog/models"
)

type TagController struct {
}

var Tag TagController

type TagEditForm struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description"`
}

// 创建标签
func (c TagController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form TagEditForm
		if err := c.ShouldBind(&form); err != nil {
			ResponseFormValidateError(c, err)
			return
		}

		// 去重查询
		existTag := models.Tag{}
		databases.DB.Where(&models.Tag{Name: form.Name}).First(&existTag)
		fmt.Println(existTag)
		fmt.Println(form.Name)
		var responseBody ResponseBody
		if existTag.ID > 0 {
			responseBody.Code = RecordExist
			responseBody.Message = "标签名称已经存在"
			Response(c, responseBody)
			return
		}

		// 创建标签
		newTag := models.Tag{Name: form.Name, Description: form.Description}
		if err := databases.DB.Create(&newTag).Error; err != nil {
			fmt.Println(err)
			responseBody.Code = RecordCreatedFail
			responseBody.Message = "标签创建失败"
			Response(c, responseBody)
			return
		}

		responseBody.Code = Success
		responseBody.Message = "标签创建成功"
		responseBody.Data = &gin.H{"id": newTag.ID}
		Response(c, responseBody)
		return
	}
}

// 更新标签
func (c TagController) Edit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form TagEditForm
		if err := c.ShouldBind(&form); err != nil {
			ResponseFormValidateError(c, err)
			return
		}
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			Response(c, ResponseBody{
				Code:    RequestParamError,
				Message: "参数错误",
			})
			return
		}
		var responseBody ResponseBody
		// 查找标签
		tag := models.Tag{}
		if databases.DB.First(&tag, id).RecordNotFound() {
			responseBody.Code = RecordNotFound
			responseBody.Message = "标签不存在或已经被删除"
			Response(c, responseBody)
			return
		}

		// 去重查询
		existTag := models.Tag{}
		databases.DB.Where("name = ?", form.Name).Where("id <> ?", tag.ID).First(&existTag)
		if existTag.ID > 0 {
			responseBody.Code = RecordExist
			responseBody.Message = "标签名称已经存在"
			Response(c, responseBody)
			return
		}

		// 创建标签
		tag.Name = form.Name
		tag.Description = form.Description
		if err := databases.DB.Save(&tag).Error; err != nil {
			responseBody.Code = RecordUpdateFail
			responseBody.Message = "标签更新失败"
			Response(c, responseBody)
			return
		}

		responseBody.Code = Success
		responseBody.Message = "标签更新成功"
		responseBody.Data = &gin.H{"id": tag.ID}
		Response(c, responseBody)
		return
	}
}

// 获取标签详情
func (c TagController) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			Response(c, ResponseBody{
				Code:    RequestParamError,
				Message: "参数错误",
			})
			return
		}
		var responseBody ResponseBody
		tag := models.Tag{}
		if databases.DB.First(&tag, id).RecordNotFound() {
			responseBody.Code = RecordNotFound
			responseBody.Message = "标签不存在或已经被删除"
			Response(c, responseBody)
			return
		}

		responseBody.Code = Success
		responseBody.Message = "success"
		responseBody.Data = &gin.H{
			"id": tag.ID,
			"name": tag.Name,
			"description": tag.Description,
			"articleCount": tag.ArticleCount,
			"createdAt": tag.CreatedAt.Format(consts.DefaultTimeFormat),
			"updatedAt": tag.UpdatedAt.Format(consts.DefaultTimeFormat),
		}

		Response(c, responseBody)
		return
	}
}

// 获取标签列表
func (c TagController) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.DefaultQuery("name", "")
		pageSizeStr := c.DefaultQuery("pageSize", "20")
		pageStr := c.DefaultQuery("page", "")

		query := databases.DB
		if name != "" {
			query.Where("name like ?", "%"+name+"%")
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			pageSize = 20
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}

		if page < 1 {
			page = 1
		}

		if pageSize < 1 {
			pageSize = 1
		}
		list := make([]models.Tag, pageSize)
		total := 0
		query.Model(&models.Tag{}).Count(&total)
		query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)

		result := make([]gin.H, len(list))
		for index, item := range list {
			result[index] = gin.H{
				"id":          item.ID,
				"name":        item.Name,
				"description": item.Description,
				"createdAt":   item.CreatedAt.Format(consts.DefaultTimeFormat),
				"updatedAt":   item.UpdatedAt.Format(consts.DefaultTimeFormat),
			}
		}

		Response(c, ResponseBody{
			Code:    Success,
			Message: "success",
			Data: &gin.H{
				"list":     result,
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
				"isEnd":    len(result) < pageSize,
			},
		})
	}
}

func (c TagController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			Response(c, ResponseBody{
				Code:    RequestParamError,
				Message: "参数错误",
			})
			return
		}

		tag := models.Tag{}
		if databases.DB.First(&tag, id).RecordNotFound() {
			Response(c, ResponseBody{
				Code:    RecordNotFound,
				Message: "分类不存在或已经被删除",
			})
			return
		}

		err = databases.DB.Delete(&tag).Error
		if err != nil {
			Response(c, ResponseBody{
				Code:    RecordDeleteFail,
				Message: "分类删除失败",
			})
			return
		}

		Response(c, ResponseBody{
			Code:    Success,
			Message: "删除成功",
		})
		return
	}
}
