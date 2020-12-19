package admin

import (
	"github.com/gin-gonic/gin"
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
	HeadImage  string `form:"headImage"`
	Content    string `form:"content" binding:"required"`
	Intro      string `form:"intro"`
	CategoryId uint   `form:"categoryId" binding:"required"`
	IsDraft    int    `form:"isDraft"`
	// PublishedAt time.Time `form:"publishedAt" time_format:"2006-01-02 15:04:05"` //# https://github.com/gin-gonic/gin/issues/1193
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
		// 关联标签 todo 会更新 tags 表update_at字段，不知道有什么作用
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
		return
	}
}

// 编辑文章
func (c ArticleController) Edit() gin.HandlerFunc {
	return func(c *gin.Context) {
		form := ArticleEditForm{}
		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}
		if err := c.ShouldBind(&form); err != nil {
			helpers.ResponseValidateError(c, err)
			return
		}

		// 查询文章
		article := &models.Article{}
		if databases.DB.First(article, id).RecordNotFound() {
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
		databases.DB.Model(article).Related(&oldTags, "Tags")
		databases.DB.Model(article).Association("Tags").Clear()

		tags := make([]models.Tag, len(form.Tags))
		if len(form.Tags) > 0 {
			databases.DB.Where(form.Tags).Find(&tags)
			if len(tags) == 0 {
				helpers.ResponseError(c, helpers.RecordNotFound, "文章标签有误")
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
				helpers.ResponseError(c, helpers.RequestParamError, "发布时间填写错误")
				return
			}
			now := time.Now()
			if publishAt.After(now) {
				article.PublishedAt = publishAt
			}
		}

		if err := databases.DB.Save(article).Error; err != nil {
			helpers.ResponseError(c, helpers.RecordUpdateFail, "文章更新失败")
			return
		}

		helpers.ResponseOk(c, "文章更新成功", &gin.H{
			"id": article.ID,
		})
		return
	}
}

// 文章详情
func (c ArticleController) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}

		article := models.Article{}
		if databases.DB.Preload("Category").Preload("Tags").First(&article, id).RecordNotFound() {
			helpers.ResponseError(c, helpers.RecordNotFound, "文章不存在或已经被删除")
			return
		}

		headImageFile := &models.File{}
		databases.DB.Where("md5 = ?", article.HeadImage).First(headImageFile)

		formatTags := make([]gin.H, 0)
		for _, tag := range article.Tags {
			formatTags = append(formatTags, gin.H{
				"id":   tag.ID,
				"name": tag.Name,
			})
		}

		res := gin.H{
			"id":            article.ID,
			"title":         article.Title,
			"headImage":     article.HeadImage,
			"headImageFile": gin.H{},
			"content":       article.Content,
			"intro":         article.Intro,
			"category": gin.H{
				"id":   article.Category.ID,
				"name": article.Category.Name,
			},
			"tags":        formatTags,
			"isDraft":     article.IsDraft,
			"publishedAt": article.PublishedAt.Format(consts.DefaultTimeFormat),
		}

		if headImageFile.ID > 0 {
			url, _ := headImageFile.Url()
			res["headImageFile"] = gin.H{
				"id":   headImageFile.ID,
				"md5":  headImageFile.Md5,
				"name": headImageFile.Name,
				"url":  url,
			}
		}

		helpers.ResponseOk(c, "success", &res)
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
		pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "20"))
		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		isDraftStr := c.DefaultQuery("isDraft", "")
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = consts.DefaultPageSize
		}

		query := databases.DB
		if categoryId > 0 {
			query = query.Where("category_id = ?", categoryId)
		}
		if isDraftStr != "" {
			isDraft, err = strconv.Atoi(isDraftStr)
			if err == nil && isDraft == 0 || isDraft == 1 {
				query = query.Where("is_draft = ?", isDraft)
			}
		}
		if title != "" {
			query = query.Where("title like ?", "%"+title+"%")
		}
		query.Model(&models.Article{}).Count(&total)
		query.Order("id desc").Preload("Category").Preload("Tags").Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles)

		formatArticles := make([]gin.H, len(articles))
		for index, article := range articles {
			formatArticle := map[string]interface{}{}
			formatArticle["id"] = article.ID
			formatArticle["title"] = article.Title
			formatArticle["headImage"] = article.HeadImage
			formatArticle["views"] = article.Views
			formatArticle["commentCount"] = article.CommentCount
			formatArticle["isDraft"] = article.IsDraft
			formatArticle["publishedAt"] = article.PublishedAt.Format(consts.DefaultTimeFormat)
			formatArticle["createdAt"] = article.CreatedAt.Format(consts.DefaultTimeFormat)
			formatArticle["updatedAt"] = article.UpdatedAt.Format(consts.DefaultTimeFormat)
			formatArticle["category"] = gin.H{
				"id":   article.Category.ID,
				"name": article.Category.Name,
			}
			// 标签
			tags := make([]gin.H, 0)
			for _, tag := range article.Tags {
				tags = append(tags, gin.H{
					"id":   tag.ID,
					"name": tag.Name,
				})
			}
			formatArticle["tags"] = tags
			formatArticles[index] = formatArticle
		}

		helpers.ResponseOk(c, "success", &gin.H{
			"list":     formatArticles,
			"pageSize": pageSize,
			"page":     page,
			"total":    total,
			"isEnd":    len(articles) < pageSize,
		})
		return
	}
}

// 删除文章
func (c ArticleController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}
		// 查询文章
		article := models.Article{}
		if databases.DB.First(&article, id).RecordNotFound() {
			helpers.ResponseError(c, helpers.RecordNotFound, "文章不存在或已经被删除")
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
			helpers.ResponseError(c, helpers.RecordDeleteFail, "删除失败")
			return
		}

		helpers.ResponseOk(c, "删除成功", &gin.H{})
		return
	}
}
