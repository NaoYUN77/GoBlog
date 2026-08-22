package router

import (
	"Blog/handler"
	"Blog/middleware"

	"Blog/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	//定义默认路由引擎

	//现在定义的是一个admin的api
	r := gin.New()
	r.Use(middleware.LoggerMiddleWare(logger.Logger), gin.Recovery())
	r.Group("/admin")
	//创建文章

	//定义admin后端api
	admin := r.Group("/admin")
	{
		admin.POST("/login", handler.LoginHandler)
		admin.POST("/post", middleware.ParseTokenMiddleware(), handler.PostHandler)
		admin.PATCH("/patch/:id", middleware.ParseTokenMiddleware(), handler.UpdatePostHandler)
		admin.DELETE("/delete/:id", middleware.ParseTokenMiddleware(), handler.DeletePostHandler)

		//获取单个文章
		admin.GET("/get/:id", middleware.ParseTokenMiddleware(), handler.GetPostSigleHandler)
		//获取文章列表
		admin.GET("/list", middleware.ParseTokenMiddleware(), handler.GetPostListHandler)
	}

	// 用户端公开接口（无需 JWT）：通过 slug 打开文章 / 文章列表
	r.GET("/post/:slug", handler.UserGetPostHandler)
	r.GET("/posts", handler.UserGetPostListHandler)

	return r
}
