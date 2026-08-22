package models

// 获取文章列表的分页请求参数
type GetPostListRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}
