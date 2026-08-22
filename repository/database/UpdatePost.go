package database

import (
	"Blog/models"
	"errors"
	"fmt"
	"strings"
)

// 文章部分patch
func UpdatePost(p *models.UpdateDb) error {
	//定义一个string类型切片用来存储占位符
	var fields []string
	//存储数据库字段
	var args []interface{}

	if p.Title != nil {
		fields = append(fields, "title = ?")
		args = append(args, *p.Title)

	}

	if p.Slug != nil {
		fields = append(fields, "slug = ?")
		args = append(args, *p.Slug)
	}
	if p.Summary != nil {
		fields = append(fields, "summary = ?")
		args = append(args, *p.Summary)
	}

	if p.CategoryID != nil {
		fields = append(fields, "category_id = ?")
		args = append(args, *p.CategoryID)
	}

	if p.Content != nil {
		fields = append(fields, "content = ?")
		args = append(args, *p.Content)

	}

	if p.CoverURL != nil {
		fields = append(fields, "cover_url = ?")
		args = append(args, *p.CoverURL)

	}

	if p.Status != nil {
		fields = append(fields, "status = ?")
		args = append(args, *p.Status)
	}

	if len(fields) == 0 {
		return errors.New("没有可更新的字段")
	}

	query := fmt.Sprintf(
		"UPDATE tb_post SET %s, updated_at = NOW() WHERE id = ?",
		strings.Join(fields, ", "),
	)
	args = append(args, p.ID)

	_, err := DB.Exec(query, args...)
	if err != nil {
		return errors.New("更新失败")
	}
	return nil
}

// POST   /api/admin/posts
// GET    /api/admin/posts/:id
// PATCH  /api/admin/posts/:id
// DELETE /api/admin/posts/:id
