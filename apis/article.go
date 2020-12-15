package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"v-blog/consts"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/libs/slice"
	"v-blog/models"
)

type ArticleController struct {
}

var Article ArticleController

func (c ArticleController) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		category := c.Query("category")
		tag := c.Query("tag")
		title := c.Query("title")
		page := c.Query("page")
		pageSize := c.Query("pageSize")

		categoryId, _ := strconv.Atoi(category)
		tagId, _ := strconv.Atoi(tag)
		currentPage, _ := strconv.Atoi(page)
		count, _ := strconv.Atoi(pageSize)

		if count <= 0 {
			count = 10
		}
		if currentPage < 1  {
			 currentPage = 1
		}

		var articleIds []int
		if tagId > 0 {
			databases.DB.Table("article_tags").Where("tag_id = ?", tagId).Pluck("article_id", &articleIds)
		}

		query := databases.DB.Model(&models.Article{}).Where("is_draft = ?", 0).Order("comment_count desc, views desc, id desc")
		if categoryId > 0 {
			query = query.Where("category_id = ?", categoryId)
		}
		if len(articleIds) > 0 {
			query = query.Where("id IN (?)", articleIds)
		}

		if title != "" {
			query = query.Where("title like ?", fmt.Sprintf("%%%s%%", title))
		}

		var articles []models.Article
		var total int
		query.Count(&total)
		query.Offset((currentPage - 1) * count).Limit(count).Find(&articles)

		// 文章图片 md5 切片
		articleImageMd5, err := slice.ToSlice(articles).Column("HeadImage").Unique().CovertToString()
		// 文章分类 id 切片
		categoryIdArr := make([]uint, 0, len(articles))
		// 文章 id 切片
		articleIdArr := make([]uint, 0, len(articles))

		for _, a := range articles {
			if a.HeadImage != "" {
				md5Repeat := false
				for _, v := range articleImageMd5 {
					if v == a.HeadImage {
						md5Repeat = true
					}
				}
				if !md5Repeat {
					articleImageMd5 = append(articleImageMd5, a.HeadImage)
				}
			}
			cateIdRepeat := false
			for _, v := range categoryIdArr {
				if a.CategoryId == v {
					cateIdRepeat = true
				}
			}
			if !cateIdRepeat {
				categoryIdArr = append(categoryIdArr, a.CategoryId)
			}
			articleIdArr = append(articleIdArr, a.ID)
		}

		imgMap := make(map[string]string)
		imgMap, _ = helpers.BatchGetImageUrlsByMd5(articleImageMd5)

		categoryMap := make(map[uint]models.Category)
		categories := make([]models.Category, 0, len(articles))
		tagMap := make(map[uint]models.Tag)

		databases.DB.Where("id In (?)", categoryIdArr).Find(&categories)

		tagIdArr := make([]uint, 0, len(articles))
		articleTagMap := make(map[uint][]uint)

		rows, err := databases.DB.Table("article_tags").Where("article_id IN (?)", articleIdArr).Rows()
		if err == nil {
			var articleId, tagId uint
			for rows.Next() {
				err = rows.Scan(&articleId, &tagId)
				if err != nil {
					continue
				}
				repeat := false
				for _, id := range tagIdArr {
					if id == tagId {
						repeat = true
					}
				}
				if !repeat {
					tagIdArr = append(tagIdArr, tagId)
				}
				articleTagMap[articleId] = append(articleTagMap[articleId], tagId)
			}
		}

		tags := make([]models.Tag, 0, len(tagIdArr))
		databases.DB.Where("id in (?)", tagIdArr).Find(&tags)
		for _, tag := range tags {
			tagMap[tag.ID] = tag
		}

		for _, category := range categories {
			categoryMap[category.ID] = category
		}

		list := make([]gin.H, len(articles))
		for k, article := range articles {
			headImageUrl := ""
			if article.HeadImage != "" {
				if url, ok := imgMap[article.HeadImage]; ok {
					headImageUrl = url
				}
			}
			cateInfo := make(gin.H)
			if category, ok := categoryMap[article.CategoryId]; ok {
				cateInfo["label"] = category.Name
				cateInfo["value"] = category.ID
			}

			articleTags := make([]gin.H, 0)
			// 标签
			if articleTagIds, ok := articleTagMap[article.ID]; ok {
				var tags []models.Tag
				databases.DB.Where("id in (?)", articleTagIds).Find(&tags)
				for _, tag := range tags {
					articleTags = append(articleTags, gin.H{
						"value": tag.ID,
						"label": tag.Name,
					})
				}
			}

			list[k] = gin.H{
				"id":           article.ID,
				"title":        article.Title,
				"headImageUrl": headImageUrl,
				"intro":        article.Intro,
				"views":        article.Views,
				"commentCount": article.CommentCount,
				"category":     cateInfo,
				"publishedAt":  article.PublishedAt.Format(consts.DefaultTimeFormat),
				"tags":         articleTags,
			}
		}

		helpers.ResponseOk(c, "success", &gin.H{
			"list": list,
			"total": total,
			"currentPage": currentPage,
		})
	}
}

func (c ArticleController) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
		}

		var article models.Article
		err = databases.DB.Preload("Category").Preload("Tags").Where("is_draft = ?", 0).First(&article, id).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseError(c, helpers.RecordNotFound, "文章不存在")
				return
			} else {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "查询失败")
				return
			}
		}

		result := gin.H{
			"id":      article.ID,
			"title":   article.Title,
			"headImageUrl": helpers.SingleGetImageUrlByMd5(article.HeadImage),
			"content": article.Content,
			"intro":   article.Intro,
			"category": gin.H{
				"value": article.Category.ID,
				"label": article.Category.Name,
			},
			"views":        article.Views,
			"commentCount": article.CommentCount,
			"publishedAt":  article.PublishedAt.Format(consts.DefaultTimeFormat),
		}

		articleTags := make([]gin.H, 0, 5)
		for _, tag := range article.Tags {
			articleTags = append(articleTags, gin.H{
				"value": tag.ID,
				"label": tag.Name,
			})
		}

		result["tags"] = articleTags

		helpers.ResponseOk(c, "success", &result)
	}
}
