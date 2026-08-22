package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//自定义recovery和logger日志中间件

func LoggerMiddleWare(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		//记录请求发送时间
		start := time.Now()
		c.Next()

		//请求花费的时间
		cost := time.Since(start)

		//定义日志格式
		logger.Info(
			"HTTP Request",
			zap.String(
				"method",
				c.Request.Method,
			),

			zap.String(
				"path",
				c.Request.URL.Path,
			),

			zap.Int(
				"status",
				c.Writer.Status(),
			),

			zap.Duration(
				"cost",
				cost,
			),
		)
	}
}

// 定义recovery, 不让程序panic
//func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//
//		defer func() {
//			//recover会捕捉panic ,
//			if err := recover(); err != nil {
//
//				logger.Error(
//					"panic recovered",
//					zap.Any(
//						"error",
//						err,
//					),
//					zap.String(
//						"path",
//						c.Request.URL.Path,
//					),
//				)
//				//返回一个json响应
//				c.JSON(
//					http.StatusInternalServerError, gin.H{
//						"error": "Internal Server Error",
//					},
//				)
//				//组织
//				c.Abort()
//
//			}
//
//		}()
//
//		c.Next()
//
//	}
//
//}
