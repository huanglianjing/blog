// Command article_converter 扫描博客源目录，将 blog_meta.yaml 中登记
// 且 markdown 文件真实存在的文章转换为 html，按「分类/文章.html」的结构
// 输出到目标目录，同时把分类、标签、文章信息全量写入 sqlite（清理旧数据）。
//
// 用法:
//
//	article_converter -src <源目录> -out <html输出目录> -db <sqlite文件>
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"

	"github.com/huanglianjing/blog/server/internal/common"
	"github.com/huanglianjing/blog/server/internal/model"
)

// metaCategory 对应 blog_meta.yaml 中的一个分类条目。
type metaCategory struct {
	Category string     `yaml:"category"`
	Files    []metaFile `yaml:"files"`
}

// metaFile 对应一篇文章的元信息。
type metaFile struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
}

func main() {
	var (
		srcDir = flag.String("src", "", "博客源目录（内含 blog_meta.yaml 和各分类子目录）")
		outDir = flag.String("out", "", "html 输出目录")
		dbPath = flag.String("db", "", "sqlite 数据库文件路径")
	)
	flag.Parse()

	if *srcDir == "" || *outDir == "" || *dbPath == "" {
		fmt.Fprintln(os.Stderr, "用法: article_converter -src <源目录> -out <html输出目录> -db <sqlite文件>")
		flag.PrintDefaults()
		os.Exit(2)
	}

	if err := run(*srcDir, *outDir, *dbPath); err != nil {
		log.Fatalf("转换失败: %v", err)
	}
}

func run(srcDir, outDir, dbPath string) error {
	// 1. 读取并解析 blog_meta.yaml。
	metas, err := loadMeta(filepath.Join(srcDir, "blog_meta.yaml"))
	if err != nil {
		return err
	}

	// 2. 初始化数据库（建表）。
	if err := model.InitDB(dbPath); err != nil {
		return fmt.Errorf("init db: %w", err)
	}

	// 3. 遍历 meta，转换存在的 markdown，收集要写库的数据。
	b := newBuilder()
	converted, skipped := 0, 0
	for _, mc := range metas {
		categoryID := b.categoryID(mc.Category)

		for _, mf := range mc.Files {
			srcPath := filepath.Join(srcDir, mc.Category, mf.Title+".md")
			if _, statErr := os.Stat(srcPath); statErr != nil {
				// meta 有记录但 md 文件不存在，跳过。
				log.Printf("跳过（缺少 md 文件）: %s", srcPath)
				skipped++
				continue
			}

			// 转换 html，输出到 out/分类/标题.html，保留一级分类目录。
			relPath := filepath.Join(mc.Category, mf.Title+".html")
			dstPath := filepath.Join(outDir, relPath)
			if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
				return fmt.Errorf("create out dir: %w", err)
			}
			if err := common.MarkdownFileToHTMLFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("convert %q: %w", srcPath, err)
			}

			// db 中 path 存 html 文件的绝对路径，便于服务端直接读取。
			absPath, err := filepath.Abs(dstPath)
			if err != nil {
				return fmt.Errorf("abs path %q: %w", dstPath, err)
			}

			b.addArticle(mf, categoryID, absPath)
			converted++
		}
	}

	// 4. 全量写入数据库（事务内清理旧数据）。
	if err := model.SyncBlogData(b.categories, b.tags, b.articles, b.relations); err != nil {
		return fmt.Errorf("sync db: %w", err)
	}

	log.Printf("完成: 转换 %d 篇, 跳过 %d 篇; 分类 %d, 标签 %d",
		converted, skipped, len(b.categories), len(b.tags))
	return nil
}

// loadMeta 读取并解析 blog_meta.yaml。
func loadMeta(path string) ([]metaCategory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read meta %q: %w", path, err)
	}
	var metas []metaCategory
	if err := yaml.Unmarshal(data, &metas); err != nil {
		return nil, fmt.Errorf("parse meta %q: %w", path, err)
	}
	return metas, nil
}

// builder 负责把 meta 数据组装成待写库的各表记录，并统一分配 id、去重。
type builder struct {
	categories []model.Category
	tags       []model.Tag
	articles   []model.Article
	relations  []model.ArticleTag

	categoryIDs map[string]int64 // 分类名 -> id
	tagIDs      map[string]int64 // 标签名 -> id
	nextArticle int64
}

func newBuilder() *builder {
	return &builder{
		categoryIDs: make(map[string]int64),
		tagIDs:      make(map[string]int64),
	}
}

// categoryID 返回分类名对应的 id，不存在则新建。
func (b *builder) categoryID(name string) int64 {
	if id, ok := b.categoryIDs[name]; ok {
		return id
	}
	id := int64(len(b.categories) + 1)
	b.categoryIDs[name] = id
	b.categories = append(b.categories, model.Category{ID: id, Name: name})
	return id
}

// tagID 返回标签名对应的 id，不存在则新建。
func (b *builder) tagID(name string) int64 {
	if id, ok := b.tagIDs[name]; ok {
		return id
	}
	id := int64(len(b.tags) + 1)
	b.tagIDs[name] = id
	b.tags = append(b.tags, model.Tag{ID: id, Name: name})
	return id
}

// addArticle 记录一篇文章及其标签关联。path 是 html 文件的绝对路径。
func (b *builder) addArticle(mf metaFile, categoryID int64, path string) {
	b.nextArticle++
	articleID := b.nextArticle
	b.articles = append(b.articles, model.Article{
		ID:         articleID,
		Title:      mf.Title,
		Date:       mf.Date,
		Path:       path,
		CategoryID: categoryID,
	})

	// 同一篇文章内标签去重，避免关联表唯一索引冲突。
	seen := make(map[int64]bool)
	for _, tagName := range mf.Tags {
		tagID := b.tagID(tagName)
		if seen[tagID] {
			continue
		}
		seen[tagID] = true
		b.relations = append(b.relations, model.ArticleTag{
			ArticleID: articleID,
			TagID:     tagID,
		})
	}
}
