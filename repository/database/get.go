package database

import (
	"Blog/models"
	"errors"
)

// 获取单篇文章
func GetPost(id int64) (*models.Postdb, error) {
	var p = new(models.Postdb)
	sqlstr := "select * from tb_post where id = ?"
	err := DB.Get(p, sqlstr, id)
	if err != nil {
		return nil, errors.New("获取文章失败")
	}
	return p, nil
}

// 获取文章列表（分页）
func GetPostList(page, pageSize int) ([]*models.Postdb, error) {
	var p = []*models.Postdb{}
	offset := (page - 1) * pageSize
	sqlstr := "select * from tb_post order by created_at desc limit ? offset ?"
	err := DB.Select(&p, sqlstr, pageSize, offset)
	if err != nil {
		return nil, errors.New("获取文章列表失败")
	}
	return p, nil
}

// 用户端：通过 slug 获取已发布文章
func GetPostBySlug(slug string) (*models.Postdb, error) {
	var p = new(models.Postdb)
	sqlstr := "select * from tb_post where slug = ? and status = 1"
	err := DB.Get(p, sqlstr, slug)
	if err != nil {
		return nil, errors.New("文章不存在或未发布")
	}
	return p, nil
}

// 用户端：已发布文章分页列表
func GetPublishedPostList(page, pageSize int) ([]*models.Postdb, error) {
	var p = []*models.Postdb{}
	offset := (page - 1) * pageSize
	sqlstr := "select * from tb_post where status = 1 order by created_at desc limit ? offset ?"
	err := DB.Select(&p, sqlstr, pageSize, offset)
	if err != nil {
		return nil, errors.New("获取文章列表失败")
	}
	return p, nil
}

// 阅读量 +1（打开文章详情时调用，失败不影响文章展示）
func IncrementViewCount(id int64) error {
	sqlstr := "update tb_post set view_count = view_count + 1 where id = ?"
	_, err := DB.Exec(sqlstr, id)
	if err != nil {
		return errors.New("阅读量更新失败")
	}
	return nil
}
