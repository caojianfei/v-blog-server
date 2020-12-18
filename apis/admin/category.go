package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"v-blog/consts"
	"v-blog/databases"
	"v-blog/helpers"
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
		if err := c.ShouldBind(&form); err != nil {
			helpers.ResponseValidateError(c, err)
			return
		}
		existCategory := &models.Category{Name: form.Name}
		// 去重查询
		databases.DB.Where(existCategory).First(existCategory)
		if existCategory.ID > 0 {
			helpers.ResponseError(c, helpers.RecordExist, "分类名称已存在")
			return
		}

		// 创建分类
		newCategory := models.Category{Name: form.Name, Description: form.Description}
		if err := databases.DB.Create(&newCategory).Error; err != nil {
			helpers.ResponseError(c, helpers.RecordCreatedFail, "分类添加失败")
			return
		}

		helpers.ResponseOk(c, "添加成功", &gin.H{
			"id": newCategory.ID,
		})
		return
	}
}

// 更新分类
func (c CategoryController) Edit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form CategoryEditForm
		if err := c.ShouldBind(&form); err != nil {
			helpers.ResponseValidateError(c, err)
			return
		}

		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}

		category := &models.Category{}
		if err := databases.DB.First(&category, id).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseError(c, helpers.RecordNotFound, "分类不存在或已经被删除")
			} else {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "查询失败")
			}
			return
		}

		// 去重查询
		existCategory := &models.Category{Name: form.Name}
		databases.DB.Where(existCategory).Where("id <> ?", category.ID).First(&existCategory)
		if existCategory.ID > 0 {
			helpers.ResponseError(c, helpers.RecordExist, "该分类名称已存在")
			return
		}

		// 更新分类
		category.Name = form.Name
		category.Description = form.Description
		if err := databases.DB.Save(&category).Error; err != nil {
			helpers.ResponseError(c, helpers.RecordUpdateFail, "分类更新失败")
			return
		}

		helpers.ResponseOk(c, "更新成功", &gin.H{})
		return
	}
}

// 获取分类详情
func (c CategoryController) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}
		category := &models.Category{}
		if err := databases.DB.First(&category, id).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseError(c, helpers.RecordNotFound, "分类不存在或已经被删除")
			} else {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "查询失败")
			}
			return
		}

		helpers.ResponseOk(c, "success", &gin.H{
			"id":          category.ID,
			"name":        category.Name,
			"description": category.Description,
			"createdAt":   category.CreatedAt.Format(consts.DefaultTimeFormat),
			"updatedAt":   category.UpdatedAt.Format(consts.DefaultTimeFormat),
		})

		return
	}
}

// 获取分类列表
func (c CategoryController) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.DefaultQuery("name", "")
		pageSizeStr := c.DefaultQuery("pageSize", "20")
		pageStr := c.DefaultQuery("page", "1")

		query := databases.DB
		if name != "" {
			query = query.Where("name like ?", "%"+name+"%")
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			pageSize = 20
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
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
				"createdAt":   item.CreatedAt.Format(consts.DefaultTimeFormat),
				"updatedAt":   item.UpdatedAt.Format(consts.DefaultTimeFormat),
			}
		}

		helpers.ResponseOk(c, "success", &gin.H{
			"list":     result,
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
			"isEnd":    len(result) < pageSize,
		})
		return
	}
}

// 删除分类
func (c CategoryController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}

		category := models.Category{}
		if err := databases.DB.First(&category, id).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseError(c, helpers.RecordNotFound, "分类不存在或已经被删除")
			} else {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "查询失败")
			}
			return
		}

		err = databases.DB.Delete(&category).Error
		if err != nil {
			helpers.ResponseError(c, helpers.RecordDeleteFail, "分类删除失败")
			return
		}

		helpers.ResponseOk(c, "删除成功", &gin.H{})
		return
	}
}

// 根据标题查询分类列表
func (c CategoryController) QueryByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("title")
		if name == "" {
			helpers.ResponseOk(c, "success", &gin.H{
				"list": make(gin.H),
			})
			return
		}

		var categories []models.Category

		err := databases.DB.Where("name like ?", fmt.Sprintf("%%%s%%", name)).Find(&categories).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseOk(c, "success", &gin.H{
					"list": make(gin.H),
				})
				return
			} else {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "查询失败")
				return
			}
		}

		list := make([]gin.H, len(categories))
		for index, item := range categories {
			list[index] = gin.H{
				"value": item.ID,
				"label": item.Name,
			}
		}

		helpers.ResponseOk(c, "success", &gin.H{
			"list": list,
		})
		return
	}
}
