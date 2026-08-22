package database

import (
	"Blog/models"
	"errors"
)

//创建文章
func CreatePost(p *models.Postdb) error {
	sqlstr := "insert into tb_post (id , title , slug  , summary ,content, cover_url,status) values (?,?,?,?,?,?,?)"
	_ , err := DB.Exec(sqlstr,p.ID,p.Title,p.Slug,p.Summary,p.Content,p.CoverURL,p.Status)
	if err != nil {
		return  errors.New("创建失败")
	}
	return nil
}