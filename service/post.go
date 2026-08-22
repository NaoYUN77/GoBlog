package service

import (

	//json -> db

	"Blog/models"
	"Blog/pkg/snowflake"
	"Blog/repository/database"
)

// json  -> db
func CreatePost(cpr *models.CreatePostRequest) (err error) {
	var p = &models.Postdb{
		ID:       snowflake.GenId(),
		Title:    cpr.Title,
		Slug:     cpr.Slug,
		Summary:  cpr.Summary,
		Content:  cpr.Content,
		CoverURL: cpr.CoverURL,
		Status:   cpr.Status,
	}

	//插入数据

	if err = database.CreatePost(p); err != nil {
		return err
	}

	return

}

func UpdatePost(id int64, upr *models.UpdatePostRequest) error {
	//获取从handler传来的数据
	var p = &models.UpdateDb{
		ID:       id,
		Title:    upr.Title,
		Slug:     upr.Slug,
		Summary:  upr.Summary,
		Content:  upr.Content,
		CoverURL: upr.CoverURL,

		//patch之后status 变为1,
		Status: upr.Status,
	}

	//操作数据库

	if err := database.UpdatePost(p); err != nil {
		return err
	}
	return nil

}
