CREATE TABLE tb_post (
    id bigint UNSIGNED NOT NULL UNIQUE COMMENT '文章ID',
    title VARCHAR(200) NOT NULL COMMENT '文章标题',
    slug VARCHAR(200) NOT NULL COMMENT '文章URL标识',
    summary VARCHAR(500) DEFAULT NULL COMMENT '文章摘要',
    content LONGTEXT NOT NULL COMMENT 'Markdown正文',
    cover_url VARCHAR(500) DEFAULT NULL COMMENT '封面图片地址',
    category_id BIGINT UNSIGNED DEFAULT NULL COMMENT '分类ID',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '文章状态：0草稿，1已发布',
    view_count BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '阅读量',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    PRIMARY KEY (id),
    UNIQUE KEY uk_slug (slug),
    KEY idx_category_id (category_id),
    KEY idx_status_created_at (status, created_at)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='博客文章表';