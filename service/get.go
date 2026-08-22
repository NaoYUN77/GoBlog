package service

import (
	"Blog/models"
	"Blog/repository/database"
)

// 获取单篇文章
func GetPost(id int64) (*models.Postdb, error) {
	return database.GetPost(id)
}

// 获取文章列表（分页）
func GetPostList(page, pageSize int) ([]*models.Postdb, error) {
	return database.GetPostList(page, pageSize)
}

// 用户端：通过 slug 获取已发布文章（并累加阅读量）
func GetPostBySlug(slug string) (*models.Postdb, error) {
	p, err := database.GetPostBySlug(slug)
	if err != nil {
		return nil, err
	}
	// 阅读量自增失败不影响文章展示
	_ = database.IncrementViewCount(p.ID)
	return p, nil
}

// 用户端：已发布文章分页列表
func GetPublishedPostList(page, pageSize int) ([]*models.Postdb, error) {
	return database.GetPublishedPostList(page, pageSize)
}
