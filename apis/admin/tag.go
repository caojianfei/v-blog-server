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
			helpers.ResponseValidateError(c, err)
			return
		}

		// 去重查询
		existTag := &models.Tag{Name: form.Name}
		affected := databases.DB.Where(existTag).First(existTag).RowsAffected
		if affected > 0 {
			helpers.ResponseError(c, helpers.RecordExist, "标签名称已经存在")
			return
		}

		// 创建标签
		newTag := models.Tag{Name: form.Name, Description: form.Description}
		if err := databases.DB.Create(&newTag).Error; err != nil {
			helpers.ResponseError(c, helpers.RecordCreatedFail, "标签创建失败")
			return
		}

		helpers.ResponseOk(c, "标签创建成功", &gin.H{"id": newTag.ID})
		return
	}
}

// 更新标签
func (c TagController) Edit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form TagEditForm
		if err := c.ShouldBind(&form); err != nil {
			helpers.ResponseValidateError(c, err)
			return
		}
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}

		// 查找标签
		tag := &models.Tag{}
		if err := databases.DB.First(tag, id).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseError(c, helpers.RecordNotFound, "标签不存在或已经被删除")
				return
			}

			helpers.ResponseError(c, helpers.DatabaseUnknownErr, "数据库查询出错")
			return
		}

		// 去重查询
		existTag := &models.Tag{}
		databases.DB.Where("name = ?", form.Name).Where("id <> ?", tag.ID).First(existTag)
		if existTag.ID > 0 {
			helpers.ResponseError(c, helpers.RecordExist, "标签名称已经存在")
			return
		}

		// 创建标签
		tag.Name = form.Name
		tag.Description = form.Description
		if err := databases.DB.Save(&tag).Error; err != nil {
			helpers.ResponseError(c, helpers.RecordUpdateFail, "标签更新失败")
			return
		}

		helpers.ResponseOk(c, "更新成功", &gin.H{"id": tag.ID})
		return
	}
}

// 获取标签详情
func (c TagController) Show() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}

		tag := &models.Tag{}
		if databases.DB.First(tag, id).Error != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseError(c, helpers.RecordNotFound, "标签不存在或已经被删除")
			} else {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "查询出错")
			}
			return
		}

		helpers.ResponseOk(c, "查询成功", &gin.H{
			"id":           tag.ID,
			"name":         tag.Name,
			"description":  tag.Description,
			"articleCount": tag.ArticleCount,
			"createdAt":    tag.CreatedAt.Format(consts.DefaultTimeFormat),
			"updatedAt":    tag.UpdatedAt.Format(consts.DefaultTimeFormat),
		})
		return
	}
}

// 获取标签列表
func (c TagController) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.DefaultQuery("name", "")
		pageSizeStr := c.DefaultQuery("pageSize", "20")
		pageStr := c.DefaultQuery("page", "1")

		query := databases.DB
		if name != "" {
			query = query.Where("name like ?", "%"+name+"%")
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
			pageSize = consts.DefaultPageSize
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

// 删除标签
func (c TagController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}

		tag := &models.Tag{}
		if err := databases.DB.First(tag, id).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseError(c, helpers.RecordNotFound, "标签不存在或已经被删除")
			} else {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "查询失败")
			}
			return
		}

		if databases.DB.Delete(tag).Error != nil {
			helpers.ResponseError(c, helpers.RecordDeleteFail, "删除失败")
			return
		}

		helpers.ResponseOk(c, "删除成功", &gin.H{})
		return
	}
}

func (c TagController) QueryByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			helpers.ResponseOk(c, "success", &gin.H{
				"list": make(gin.H),
			})
			return
		}

		var tags []models.Tag
		err := databases.DB.Where("name like ?", fmt.Sprintf("%%%s%%", name)).Find(&tags).Error
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

		list := make([]gin.H, len(tags))
		for index, item := range tags {
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
