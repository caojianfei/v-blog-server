package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/fatih/set.v0"
	"strconv"
	"time"
	"v-blog/consts"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/models"
)

type ArticleController struct {
}

var Article ArticleController

type ArticleEditForm struct {
	Title      string `json:"title" binding:"required"`
	HeadImage  string `form:"headImage" binding:"required"`
	Content    string `form:"content" binding:"required"`
	Intro      string `form:"intro"`
	CategoryId uint   `form:"categoryId" binding:"required"`
	IsDraft    int    `form:"isDraft"`
	//PublishedAt time.Time `form:"publishedAt" time_format:"2006-01-02 15:04:05"` //# https://github.com/gin-gonic/gin/issues/1193
	PublishedAt string `form:"publishedAt" binding:"formatData=2006-01-02 15:04:05"`
	Tags        []uint `form:"tags"`
}

// 创建文章
func (c ArticleController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form ArticleEditForm
		if err := c.ShouldBind(&form); err != nil {
			helpers.ResponseValidateError(c, err)
			return
		}

		// 分类检查
		category := &models.Category{}
		if databases.DB.First(category, form.CategoryId).RecordNotFound() {
			helpers.ResponseError(c, helpers.RecordNotFound, "文章分类不存在")
			return
		}

		articlePublishedAt := form.PublishedAt
		var publishedAt time.Time
		current := time.Now()
		if articlePublishedAt == "" {
			publishedAt = current
		} else {
			var err error
			publishedAt, err = time.Parse(consts.DefaultTimeFormat, articlePublishedAt)
			if err != nil {
				helpers.ResponseError(c, helpers.RequestParamError, "文章发布时间错误")
				return
			}
			if publishedAt.Before(current) {
				publishedAt = current
			}
		}

		// 入库
		article := models.Article{
			Title:       form.Title,
			HeadImage:   form.HeadImage,
			Content:     form.Content,
			Intro:       form.Intro,
			CategoryId:  form.CategoryId,
			IsDraft:     form.IsDraft,
			PublishedAt: publishedAt,
		}

		tags := make([]models.Tag, len(form.Tags))
		// 关联标签
		if len(form.Tags) > 0 {
			databases.DB.Where(form.Tags).Find(&tags)
			if len(tags) == 0 {
				helpers.ResponseError(c, helpers.RecordNotFound, "文章标签有误")
				return
			} else {
				article.Tags = tags
			}
		}

		if databases.DB.Create(&article).Error != nil {
			helpers.ResponseError(c, helpers.RecordCreatedFail, "文章创建失败")
			return
		}

		helpers.ResponseOk(c, "文章创建成功", &gin.H{"id": article.ID})

		// 标签下文章数量更新
		for _, tag := range tags {
			go func(tag models.Tag) {
				tag.IncreaseArticleCount()
			}(tag)
		}
		return
	}
}

// 编辑文章
func (c ArticleController) Edit() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		form := ArticleEditForm{}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "文章id错误")
			return
		}
		if err := c.ShouldBind(&form); err != nil {
			helpers.ResponseValidateError(c, err)
			return
		}

		// 查询文章
		article := models.Article{}
		if databases.DB.First(&article, id).RecordNotFound() {
			Response(c, ResponseBody{
				Code:    RecordNotFound,
				Message: "文章不存在或已经被删除",
			})
			helpers.ResponseError(c, helpers.RecordNotFound, "文章不存在或已经被删除")
			return
		}

		// 分类检查
		category := &models.Category{}
		if databases.DB.First(category, form.CategoryId).RecordNotFound() {
			helpers.ResponseError(c, helpers.RecordNotFound, "文章分类不存在")
			return
		}

		// 关联标签
		oldTags := make([]models.Tag, 0)
		databases.DB.Model(&article).Related(&oldTags, "Tags")
		for _, oldTag := range oldTags {
			go func(tag models.Tag) {
				oldTag.DecreaseArticleCount()
			}(oldTag)
		}

		databases.DB.Model(&article).Association("Tags").Clear()
		tags := make([]models.Tag, len(form.Tags))
		if len(form.Tags) > 0 {
			databases.DB.Where(form.Tags).Find(&tags)
			if len(tags) == 0 {
				Response(c, ResponseBody{
					Code:    RecordNotFound,
					Message: "文章标签有误",
				})
				return
			} else {
				article.Tags = tags
			}
		}

		article.Title = form.Title
		article.HeadImage = form.HeadImage
		article.Content = form.Content
		article.Intro = form.Intro
		article.CategoryId = form.CategoryId
		article.IsDraft = form.IsDraft
		if form.PublishedAt != "" {
			publishAt, err := time.Parse(consts.DefaultTimeFormat, form.PublishedAt)
			if err != nil {
				Response(c, ResponseBody{
					Code:    RequestParamError,
					Message: "发布时间填写错误",
				})
				return
			}
			now := time.Now()
			if publishAt.After(now) {
				article.PublishedAt = publishAt
			}
		}

		if err := databases.DB.Save(&article).Error; err != nil {
			Response(c, ResponseBody{
				Code:    RecordUpdateFail,
				Message: "文章更新失败",
			})
			return
		}

		Response(c, ResponseBody{
			Code:    Success,
			Message: "文章更新成功",
			Data: &gin.H{
				"id": article.ID,
			},
		})

		for _, tag := range tags {
			go func(tag models.Tag) {
				tag.IncreaseArticleCount()
			}(tag)
		}
		return
	}
}

// 文章详情
func (c ArticleController) Show() gin.HandlerFunc {
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
		tags := []models.Tag{{}}
		category := models.Category{}
		var responseBody ResponseBody
		article := models.Article{}
		if databases.DB.First(&article, id).RecordNotFound() {
			responseBody.Code = RecordNotFound
			responseBody.Message = "文章不存在或已经被删除"
			Response(c, responseBody)
			return
		}

		databases.DB.Model(&article).Related(&tags, "Tags")
		databases.DB.Model(&article).Related(&category)

		formatTags := make([]gin.H, len(tags))
		for index, tag := range tags {
			item := gin.H{
				"id":   tag.ID,
				"name": tag.Name,
			}
			formatTags[index] = item
		}

		responseBody.Data = &gin.H{
			"id":        article.ID,
			"title":     article.Title,
			"headImage": article.HeadImage,
			"content":   article.Content,
			"intro":     article.Intro,
			"category": gin.H{
				"id":   category.ID,
				"name": category.Name,
			},
			"tags": formatTags,
		}
		Response(c, responseBody)
		return
	}
}

// 文章列表
func (c ArticleController) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			categoryId, isDraft, pageSize, page, total int
			err                                        error
			articles                                   []models.Article
		)

		title := c.DefaultQuery("title", "")
		categoryId, _ = strconv.Atoi(c.DefaultQuery("categoryId", ""))
		pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", ""))
		page, _ = strconv.Atoi(c.DefaultQuery("page", ""))
		isDraftStr := c.DefaultQuery("isDraft", "")

		query := databases.DB.Model(&models.Article{})
		if categoryId > 0 {
			query.Where("category_id > ?", categoryId)
		}
		if isDraftStr != "" {
			isDraft, err = strconv.Atoi(isDraftStr)
			if err == nil && isDraft == 0 || isDraft == 1 {
				query.Where("is_draft = ?", isDraft)
			}
		}
		if title != "" {
			query.Where("title like ?", "%"+title+"%")
		}
		query.Model(&models.Article{}).Count(&total)
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = consts.DefaultPageSize
		}
		databases.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles)
		categoryIds := set.New(set.ThreadSafe)
		articleIds := make([]uint, len(articles))
		for index, article := range articles {
			articleIds[index] = article.ID
			if article.CategoryId > 0 {
				categoryIds.Add(article.CategoryId)
			}
		}

		var categories []models.Category
		if categoryIds.Size() > 0 {
			databases.DB.Model(&models.Category{}).
				Select("id, name").
				Where(categoryIds.List()).
				Find(&categories)
		}

		categoryMap := make(map[uint]models.Category)
		for _, category := range categories {
			categoryMap[category.ID] = category
		}

		tagIds := set.New(set.ThreadSafe)
		articleTagsMap := make(map[uint][]uint)
		rows, err := databases.DB.Table("article_tags").Where("article_id in (?)", articleIds).Select("article_id, tag_id").Rows()
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var articleId, tagId uint
			if err := rows.Scan(&articleId, &tagId); err != nil {
				fmt.Println("scan error", err)
				continue
			}
			tagIds.Add(tagId)
			articleTagsMap[articleId] = append(articleTagsMap[articleId], tagId)
		}

		var tags []models.Tag
		tagsMap := make(map[uint]models.Tag)
		if tagIds.Size() > 0 {
			err := databases.DB.Find(&tags, tagIds.List()).Error
			if err != nil {
				panic(err)
			}
		}
		for _, tag := range tags {
			tagsMap[tag.ID] = tag
		}

		formatArticles := make([]gin.H, len(articles))
		for index, article := range articles {
			formatArticle := map[string]interface{}{}
			formatArticle["id"] = article.ID
			formatArticle["title"] = article.Title
			formatArticle["headImage"] = article.HeadImage
			formatArticle["views"] = article.Views
			formatArticle["commentCount"] = article.CommentCount
			formatArticle["is_draft"] = article.IsDraft
			formatArticle["publishedAt"] = article.PublishedAt.Format(consts.DefaultTimeFormat)
			formatArticle["createdAt"] = article.CreatedAt.Format(consts.DefaultTimeFormat)
			formatArticle["updatedAt"] = article.UpdatedAt.Format(consts.DefaultTimeFormat)

			var articleCategory gin.H
			// 分类
			if category, ok := categoryMap[article.CategoryId]; ok {
				articleCategory = gin.H{
					"id":   category.ID,
					"name": category.Name,
				}
			}
			formatArticle["category"] = articleCategory
			// 标签
			tags := make([]gin.H, 0)
			if tagIds, ok := articleTagsMap[article.ID]; ok {
				for _, tagId := range tagIds {
					if tag, ok := tagsMap[tagId]; ok {
						tags = append(tags, gin.H{
							"id":   tag.ID,
							"name": tag.Name,
						})
					}
				}
			}
			formatArticle["tags"] = tags
			formatArticles[index] = formatArticle
		}

		Response(c, ResponseBody{
			Data: &gin.H{
				"list":     formatArticles,
				"pageSize": pageSize,
				"page":     page,
				"total":    total,
				"isEnd":    len(articles) < pageSize,
			},
		})
		return
	}
}

// 删除文章
func (c ArticleController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Query("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			Response(c, ResponseBody{
				Code:    RequestParamError,
				Message: "参数错误",
			})
			return
		}
		// 查询文章
		article := models.Article{}
		if databases.DB.First(&article, id).RecordNotFound() {
			Response(c, ResponseBody{
				Code:    RecordNotFound,
				Message: "文章不存在或已经被删除",
			})
			return
		}

		tags := make([]models.Tag, 0)
		databases.DB.Model(&article).Related(&tags, "Tags")
		for _, tag := range tags {
			go func(tag models.Tag) {
				tag.DecreaseArticleCount()
			}(tag)
		}

		err = databases.DB.Delete(&article).Error
		if err != nil {
			Response(c, ResponseBody{
				Code:    RecordDeleteFail,
				Message: "文章删除失败",
			})
			return
		}

		Response(c, ResponseBody{
			Code:    Success,
			Message: "文章删除成功",
		})
		return
	}
}
