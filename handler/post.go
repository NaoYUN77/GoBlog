package handler

import (
	"Blog/models"
	"net/http"
	"strconv"

	"Blog/service"

	"github.com/gin-gonic/gin"
)

// 创建文章
func PostHandler(c *gin.Context) {
	//验证access token
	_, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无法解析token",
		})
		return
	}
	//1 获取数据
	var cpr models.CreatePostRequest

	if err := c.ShouldBind(&cpr); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "需要填必要字段",
		})
		return

	}

	//逻辑处理

	if err := service.CreatePost(&cpr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "创建文章失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "创建文章成功",
	})

}

//定义admin更新api

func UpdatePostHandler(c *gin.Context) {
	_, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无法解析token",
		})
		return
	}
	//获取数据
	var upr models.UpdatePostRequest

	//从Url中文章id
	idstr := c.Param("id")

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "文章格式id错误",
		})

		return
	}

	if err := c.ShouldBind(&upr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "需要填充必要的字段",
		})
		return
	}
	//逻辑处理
	if err := service.UpdatePost(id, &upr); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":   "更新失败",
			"error": err,
		})
		return
	}

	//
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})

}

//定规删除文章api

func DeletePostHandler(c *gin.Context) {
	_, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无法解析token",
		})
		return
	}

	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "id错误",
		})
	}

	//执行删除逻辑
	if err := service.DeletePost(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "删除文章失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})

}

func GetPostSigleHandler(c *gin.Context) {
	_, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无法解析token",
		})
		return
	}

	//获取文章Id
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取id失败",
		})
		return
	}
	//执行逻辑返回一个post结构体

	p, err := service.GetPost(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取文章失败",
		})
		return
	}

	//构造一个返回响应
	c.JSON(
		http.StatusOK,
		gin.H{
			"msg": "ok",
			"data": gin.H{
				"id":        p.ID,
				"title":     p.Title,
				"slug":      p.Slug,
				"content":   p.Content,
				"summary":   p.Summary,
				"cover_url": p.CoverURL,
				"status":    p.Status,
				"create_at": p.CreatedAt,
				"update_at": p.UpdatedAt,
			},
		},
	)
}

// 管理员获取文章列表（支持分页，默认第1页每页10条）
func GetPostListHandler(c *gin.Context) {
	_, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无法解析token",
		})
		return
	}

	var req models.GetPostListRequest
	if err := c.ShouldBindQuery(&req); err != nil || req.Page <= 0 || req.PageSize <= 0 {
		req.Page = 1
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	pl, err := service.GetPostList(req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":       "ok",
		"page":      req.Page,
		"page_size": req.PageSize,
		"data":      pl,
	})
}

// 用户端文章视图：不暴露内部文章 id，统一 snake_case 字段
func postView(p *models.Postdb) gin.H {
	return gin.H{
		"title":       p.Title,
		"slug":        p.Slug,
		"summary":     p.Summary,
		"content":     p.Content,
		"cover_url":   p.CoverURL,
		"category_id": p.CategoryID,
		"status":      p.Status,
		"view_count":  p.ViewCount,
		"created_at":  p.CreatedAt,
		"updated_at":  p.UpdatedAt,
	}
}

// 用户获取文章列表
//
// 获取单个文章
func UserGetPostHandler(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "slug不能为空",
		})
		return
	}

	p, err := service.GetPostBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": postView(p),
	})
}

func UserGetPostListHandler(c *gin.Context) {
	var req models.GetPostListRequest

	//从url上获取参数 ,用于分页
	if err := c.ShouldBindQuery(&req); err != nil || req.Page <= 0 || req.PageSize <= 0 {
		req.Page = 1
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	pl, err := service.GetPublishedPostList(req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 用户端不暴露文章 id
	var data = make([]gin.H, 0, len(pl))
	for _, p := range pl {
		data = append(data, postView(p))
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":       "ok",
		"page":      req.Page,
		"page_size": req.PageSize,
		"data":      data,
	})
}
