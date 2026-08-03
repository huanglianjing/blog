package common

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
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
			// 优先级小于默认 html renderer 的 1000，故 img 节点由它接管
			renderer.WithNodeRenderers(util.Prioritized(newLazyImageRenderer(), 500)),
		),
	)
}

// lazyImageRenderer 只接管 img 节点的渲染：先给节点补上懒加载属性，
// 再交回 goldmark 默认的 image 渲染逻辑，避免自行拼接 img 标签导致
// alt / title 转义、Unsafe / XHTML 等行为与默认实现不一致。
type lazyImageRenderer struct {
	// 内嵌默认 renderer，使其 SetOption 方法被提升上来，
	// goldmark 才能把 WithUnsafe / WithXHTML 等选项同步给它。
	*html.Renderer
	renderImage renderer.NodeRendererFunc
}

func newLazyImageRenderer() renderer.NodeRenderer {
	inner := html.NewRenderer().(*html.Renderer)
	var reg imageFuncRegisterer
	inner.RegisterFuncs(&reg)
	return &lazyImageRenderer{Renderer: inner, renderImage: reg.fn}
}

// RegisterFuncs 实现 renderer.NodeRenderer，只注册 img 节点。
func (r *lazyImageRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.renderLazyImage)
}

func (r *lazyImageRenderer) renderLazyImage(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// loading / decoding 均在 html.ImageAttributeFilter 白名单内，会被渲染出来
		node.SetAttributeString("loading", []byte("lazy"))
		node.SetAttributeString("decoding", []byte("async"))
	}
	return r.renderImage(w, source, node, entering)
}

// imageFuncRegisterer 用于从默认 renderer 中取出 img 节点的渲染函数。
type imageFuncRegisterer struct {
	fn renderer.NodeRendererFunc
}

func (c *imageFuncRegisterer) Register(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
	if kind == ast.KindImage {
		c.fn = fn
	}
}

// MarkdownToHTML 将 markdown 源内容转换为 HTML 字节。
func MarkdownToHTML(source []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := newMarkdown().Convert(source, &buf); err != nil {
		return nil, fmt.Errorf("convert markdown: %w", err)
	}
	return buf.Bytes(), nil
}

// MarkdownFileToHTMLFile 读取 srcPath 的 markdown 文件，转换后写入 dstPath，
// 并返回转换出的 html 内容，供调用方复用（如提取搜索用的正文纯文本）。
func MarkdownFileToHTMLFile(srcPath, dstPath string) ([]byte, error) {
	source, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("read markdown file %q: %w", srcPath, err)
	}

	out, err := MarkdownToHTML(source)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(dstPath, out, 0644); err != nil {
		return nil, fmt.Errorf("write html file %q: %w", dstPath, err)
	}
	return out, nil
}
