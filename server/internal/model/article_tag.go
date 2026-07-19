package model

// ArticleTag 文章标签关联表
type ArticleTag struct {
	ID        int64 `json:"id"         gorm:"column:id;primaryKey;autoIncrement"`
	ArticleID int64 `json:"article_id" gorm:"column:article_id;type:bigint;not null;default:0;uniqueIndex:uk_article_tag,priority:1"`
	TagID     int64 `json:"tag_id"     gorm:"column:tag_id;type:bigint;not null;default:0;uniqueIndex:uk_article_tag,priority:2;index"`
}

// TableName 表名
func (ArticleTag) TableName() string {
	return "article_tag"
}

// ListArticleTagsByArticleIDs 按文章 id 批量查询关联记录。
func ListArticleTagsByArticleIDs(articleIDs []int64) ([]ArticleTag, error) {
	relations := make([]ArticleTag, 0)
	if len(articleIDs) == 0 {
		return relations, nil
	}
	if err := DB.Where("article_id IN ?", articleIDs).Find(&relations).Error; err != nil {
		return nil, err
	}
	return relations, nil
}

// ListArticlesByTag 按标签 id 分页查询文章，日期倒序。
// 通过 article_tag 关联表连接 article 表。
func ListArticlesByTag(tagID int64, offset, limit int) ([]Article, error) {
	articles := make([]Article, 0)
	err := DB.Model(&Article{}).
		Joins("JOIN article_tag ON article_tag.article_id = article.id").
		Where("article_tag.tag_id = ?", tagID).
		Order("article.date DESC, article.id DESC").
		Offset(offset).Limit(limit).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// CountArticlesByTag 返回某标签下的文章总数。
func CountArticlesByTag(tagID int64) (int64, error) {
	var count int64
	err := DB.Model(&ArticleTag{}).Where("tag_id = ?", tagID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
