package common

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// newMarkdown 构造一个带常用扩展的 goldmark 实例。
// 只负责生成结构化 HTML（GFM 表格/删除线/任务列表、脚注、标题锚点），
// 代码高亮与 mermaid 图表等外观交由前端 JS + CSS 处理。
func newMarkdown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,      // 表格、删除线、自动链接、任务列表
			extension.Footnote, // 脚注
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // 给标题自动生成 id，便于锚点跳转
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			// 保留 markdown 中的原始 HTML；文章来源可信时开启
			html.WithUnsafe(),
		),
	)
}

// MarkdownToHTML 将 markdown 源内容转换为 HTML 字节。
func MarkdownToHTML(source []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := newMarkdown().Convert(source, &buf); err != nil {
		return nil, fmt.Errorf("convert markdown: %w", err)
	}
	return buf.Bytes(), nil
}

// MarkdownFileToHTMLFile 读取 srcPath 的 markdown 文件，转换后写入 dstPath。
func MarkdownFileToHTMLFile(srcPath, dstPath string) error {
	source, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("read markdown file %q: %w", srcPath, err)
	}

	out, err := MarkdownToHTML(source)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dstPath, out, 0644); err != nil {
		return fmt.Errorf("write html file %q: %w", dstPath, err)
	}
	return nil
}
