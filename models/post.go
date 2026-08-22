package models

import (
	"time"
)

type Postdb struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
	Slug  string `db:"slug"`

	//*string的控制是Nil
	Summary    *string `db:"summary"`
	Content    string  `db:"content"`
	CoverURL   *string `db:"cover_url"`
	CategoryID *uint64 `db:"category_id"`
	Status     int8    `db:"status"`
	//viewcount不使用外来参数
	//每get一次，viewcount加1
	ViewCount uint64    `db:"view_count"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// 前端发送的json数据
type CreatePostRequest struct {
	Title string `json:"title" binding:"required,max=200"`

	//文章叙述,后续会映射到url
	Slug string `json:"slug" binding:"required,max=200"`

	//简要说明可为空
	Summary    *string `json:"summary"`
	Content    string  `json:"content" binding:"required"`
	CoverURL   *string `json:"cover_url"`
	CategoryID *uint64 `json:"category_id"`
	Status     int8    `json:"status"`
}

type UpdateDb struct {
	ID         int64   `db:"id"`
	Title      *string `db:"title"`
	Slug       *string `db:"slug"`
	Summary    *string `db:"summary"`
	Content    *string `db:"content"`
	CoverURL   *string `db:"cover_url"`
	CategoryID *uint64 `db:"category_id"`
	Status     *int8   `db:"status"`
}
type UpdatePostRequest struct {

	//防止文章是一个空字符串
	Title *string `json:"title"`
	Slug  *string `json:"slug"`
	//简要说明可为空
	Summary    *string `json:"summary"`
	Content    *string `json:"content"`
	CoverURL   *string `json:"cover_url"`
	CategoryID *uint64 `json:"category_id"`
	Status     *int8   `json:"status"`
}
