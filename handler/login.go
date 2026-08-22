package handler

import (
	"Blog/models"
	"Blog/pkg/jwt"
	"Blog/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	//1.获取参数
	var pa models.ParamsAdmin
	if err := c.ShouldBind(&pa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "字段为必填字段",
		})
		return
	}

	//2. 查找用户并校验密码
	err := service.GetAdminUser(pa.Username, pa.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	var token string
	//生成access token
	token, err = jwt.GenJWT(pa.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "生成token失败,用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"token": token,
	})

}
