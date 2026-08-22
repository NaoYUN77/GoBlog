package middleware

import (
	"net/http"
	"strings"

	jwt1 "Blog/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func ParseTokenMiddleware() gin.HandlerFunc {
	//解析token
	//请求头中存储token的方式有很多种, 这里用bearer authorization


	return func(c *gin.Context) {

		//验证首部字段是不是Authorization
		Handler := c.Request.Header.Get("Authorization")
		if Handler == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "请求头错误",
			})
			c.Abort()
			return

		}
		//按照空格切开字段, 最多切成2段
		parts := strings.SplitN(Handler, " ", 2)

		//如果part不满足2且part0不是Bearer报错
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "请求头格式错误",
			})
			//阻止后续的handler
			c.Abort()
			return

		}
		//获取token
		tokenString := parts[1]
		token, err := jwt1.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "解析token失败",
			})
			c.Abort()
			return
		}

		c.Set("token", token.Username)

	}
}
