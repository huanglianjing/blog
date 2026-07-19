package model

import "gorm.io/gorm"

// SyncBlogData 全量重建博客数据：在一个事务内先清空各表，
// 再写入传入的数据，从而清理掉配置中已不存在的旧数据。
// 传入的各切片需自行分配好 id（含关联表的 article_id / tag_id）。
func SyncBlogData(categories []Category, tags []Tag, articles []Article, relations []ArticleTag) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 清空旧数据（顺序：先删关联表，再删主表）。
		for _, table := range []string{"article_tag", "article", "tag", "category"} {
			if err := tx.Exec("DELETE FROM " + table).Error; err != nil {
				return err
			}
		}

		// 写入新数据。
		if len(categories) > 0 {
			if err := tx.Create(&categories).Error; err != nil {
				return err
			}
		}
		if len(tags) > 0 {
			if err := tx.Create(&tags).Error; err != nil {
				return err
			}
		}
		if len(articles) > 0 {
			if err := tx.Create(&articles).Error; err != nil {
				return err
			}
		}
		if len(relations) > 0 {
			if err := tx.Create(&relations).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
