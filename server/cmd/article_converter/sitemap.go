package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/huanglianjing/blog/server/internal/common"
	"github.com/huanglianjing/blog/server/internal/model"
)

// writeSitemap 依据本次转换收集到的分类、标签、文章，生成 sitemap.xml 写入 path。
// baseURL 为站点根地址，来自配置文件的 site.base_url。
// 收录的是前端 SPA 真实存在的路由：首页、文章 / 分类 / 标签三个列表页，
// 以及每篇文章、每个分类、每个标签的详情页；
// 列表页的翻页只改组件内部状态、不改 URL，搜索页也无固定内容，故都不收录。
func writeSitemap(path, baseURL string, b *builder) error {
	content, err := common.BuildSitemap(baseURL, sitemapEntries(b))
	if err != nil {
		return err
	}

	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create sitemap dir %q: %w", dir, err)
		}
	}
	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("write sitemap %q: %w", path, err)
	}
	return nil
}

// sitemapEntries 组装 sitemap 的全部 URL，顺序固定以保证多次运行结果一致：
// 首页与三个列表页 -> 文章（日期倒序，与列表页一致）-> 分类 -> 标签（均按名称排序）。
func sitemapEntries(b *builder) []common.SitemapEntry {
	// 文章按日期倒序，与前端列表页顺序一致。
	articles := make([]model.Article, len(b.articles))
	copy(articles, b.articles)
	sort.SliceStable(articles, func(i, j int) bool {
		if articles[i].Date != articles[j].Date {
			return articles[i].Date > articles[j].Date
		}
		return articles[i].ID > articles[j].ID
	})

	// 分类 / 标签页的更新时间取其下文章的最新日期。
	categoryLastMod := make(map[int64]string, len(b.categories))
	for _, a := range articles {
		if d := sitemapDate(a.Date); d > categoryLastMod[a.CategoryID] {
			categoryLastMod[a.CategoryID] = d
		}
	}
	articleDates := make(map[int64]string, len(articles))
	for _, a := range articles {
		articleDates[a.ID] = sitemapDate(a.Date)
	}
	tagLastMod := make(map[int64]string, len(b.tags))
	for _, r := range b.relations {
		if d := articleDates[r.ArticleID]; d > tagLastMod[r.TagID] {
			tagLastMod[r.TagID] = d
		}
	}

	// 全站最新日期，作为首页与各列表页的更新时间。
	var latest string
	for _, d := range articleDates {
		if d > latest {
			latest = d
		}
	}

	entries := []common.SitemapEntry{
		{Segments: nil, LastMod: latest},
		{Segments: []string{"article"}, LastMod: latest},
		{Segments: []string{"category"}, LastMod: latest},
		{Segments: []string{"tag"}, LastMod: latest},
	}
	for _, a := range articles {
		entries = append(entries, common.SitemapEntry{
			Segments: []string{"article", a.Title},
			LastMod:  sitemapDate(a.Date),
		})
	}
	for _, c := range sortedByName(b.categories, func(c model.Category) string { return c.Name }) {
		entries = append(entries, common.SitemapEntry{
			Segments: []string{"category", c.Name},
			LastMod:  categoryLastMod[c.ID],
		})
	}
	for _, t := range sortedByName(b.tags, func(t model.Tag) string { return t.Name }) {
		entries = append(entries, common.SitemapEntry{
			Segments: []string{"tag", t.Name},
			LastMod:  tagLastMod[t.ID],
		})
	}
	return entries
}

// sortedByName 返回按名称排序（规则同标签概览接口）的副本，不改动入参切片。
func sortedByName[T any](items []T, name func(T) string) []T {
	sorted := make([]T, len(items))
	copy(sorted, items)
	keys := make(map[string]string, len(items))
	key := func(s string) string {
		if k, ok := keys[s]; ok {
			return k
		}
		k := common.NameSortKey(s)
		keys[s] = k
		return k
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		return key(name(sorted[i])) < key(name(sorted[j]))
	})
	return sorted
}

// sitemapDate 把 meta 中的 "YYYY-MM-DD HH:MM:SS" 截成 sitemap 的 lastmod
// 所需的 YYYY-MM-DD；格式不符（长度不足或含非法字符）时返回空串，
// 该条 URL 就不带 lastmod，而不是写入一个无效日期。
func sitemapDate(date string) string {
	if len(date) < 10 {
		return ""
	}
	d := date[:10]
	for i := 0; i < len(d); i++ {
		expectDash := i == 4 || i == 7
		isDigit := d[i] >= '0' && d[i] <= '9'
		if expectDash && d[i] != '-' {
			return ""
		}
		if !expectDash && !isDigit {
			return ""
		}
	}
	if strings.Count(d, "-") != 2 {
		return ""
	}
	return d
}
