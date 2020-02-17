package apis

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"v-blog/databases"
	"v-blog/models"
)

type CategoryController struct {
}

var Category CategoryController

type CategoryEditForm struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description"`
}

// 新增分类
func (c CategoryController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {

		var form CategoryEditForm
		var responseBody ResponseBody
		if err := c.ShouldBind(&form); err != nil {
			ResponseFormValidateError(c, err)
			return
		}

		var existCategory models.Category
		// 去重查询
		databases.DB.Where(&models.Category{Name: form.Name}).First(&existCategory)
		if existCategory.ID > 0 {
			responseBody.Message = "该分类名称已存在"
			responseBody.Code = RecordExist
			Response(c, responseBody)
			return
		}

		// 创建分类
		newCategory := models.Category{Name: form.Name, Description: form.Description}
		if err := databases.DB.Create(&newCategory).Error; err != nil {
			responseBody.Message = "分类添加失败"
			responseBody.Code = RecordCreatedFail
			Response(c, responseBody)
			return
		}

		responseBody.Code = Success
		responseBody.Message = "分类添加成功"
		responseBody.Data = &gin.H{
			"id": newCategory.ID,
		}
		Response(c, responseBody)
		return
	}
}

// 更新分类
func (c CategoryController) Edit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form CategoryEditForm
		var responseBody ResponseBody
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
		category := models.Category{}
		if databases.DB.First(&category, id).RecordNotFound() {
			Response(c, ResponseBody{
				Code:    RecordNotFound,
				Message: "分类不存在或已经被删除",
			})
			return
		}

		var existCategory models.Category
		// 去重查询
		databases.DB.Where(&models.Category{Name: form.Name}).Where("id <> ?", category.ID).First(&existCategory)
		if existCategory.ID > 0 {
			responseBody.Message = "该分类名称已存在"
			responseBody.Code = RecordExist
			Response(c, responseBody)
			return
		}

		category.Name = form.Name
		category.Description = form.Description
		if err := databases.DB.Save(&category).Error; err != nil {
			responseBody.Code = RecordUpdateFail
			responseBody.Message = "分类更新失败"
			Response(c, responseBody)
			return
		}

		responseBody.Code = Success
		responseBody.Message = "分类更新成功"
		Response(c, responseBody)
	}
}

// 获取分类详情
func (c CategoryController) Show() gin.HandlerFunc {
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
		category := models.Category{}
		if databases.DB.First(&category, id).RecordNotFound() {
			Response(c, ResponseBody{
				Code:    RecordNotFound,
				Message: "分类不存在或已经被删除",
			})
			return
		}

		Response(c, ResponseBody{
			Code:    Success,
			Message: "success",
			Data: &gin.H{
				"id":          category.ID,
				"name":        category.Name,
				"description": category.Description,
				"createdAt":   category.CreatedAt,
				"updatedAt":   category.UpdatedAt,
			},
		})
	}
}

// 获取分类列表
func (c CategoryController) List() gin.HandlerFunc {
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
		list := make([]models.Category, pageSize)
		total := 0
		query.Model(&models.Category{}).Count(&total)
		query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)

		result := make([]gin.H, len(list))
		for index, item := range list {
			result[index] = gin.H{
				"id":          item.ID,
				"name":        item.Name,
				"description": item.Description,
				"createdAt":   item.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":   item.UpdatedAt.Format("2006-01-02 15:04:05"),
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

// 删除分类
func (c CategoryController) Delete() gin.HandlerFunc {
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

		category := models.Category{}
		if databases.DB.First(&category, id).RecordNotFound() {
			Response(c, ResponseBody{
				Code:    RecordNotFound,
				Message: "分类不存在或已经被删除",
			})
			return
		}

		err = databases.DB.Delete(&category).Error
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
