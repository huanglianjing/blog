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
			// 补上 GFM linkify 覆盖不到的裸链接（如中文紧贴 URL），见 linkify.go
			parser.WithASTTransformers(util.Prioritized(newBareLinkTransformer(), 998)),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			// 保留 markdown 中的原始 HTML；文章来源可信时开启
			html.WithUnsafe(),
			// 优先级小于默认 html renderer 的 1000，故 img / a 节点由它接管
			renderer.WithNodeRenderers(util.Prioritized(newAttributeRenderer(), 500)),
		),
	)
}

// attributeRenderer 只接管需要补属性的节点（img 懒加载、外链 target / rel）：
// 先给节点补上属性，再交回 goldmark 默认的渲染逻辑，避免自行拼接标签导致
// alt / title 转义、Unsafe / XHTML 等行为与默认实现不一致。
type attributeRenderer struct {
	// 内嵌默认 renderer，使其 SetOption 方法被提升上来，
	// goldmark 才能把 WithUnsafe / WithXHTML 等选项同步给它。
	*html.Renderer
	// 默认 renderer 中各节点原本的渲染函数，补完属性后回落到它
	inner map[ast.NodeKind]renderer.NodeRendererFunc
}

func newAttributeRenderer() renderer.NodeRenderer {
	base := html.NewRenderer().(*html.Renderer)
	collector := &nodeFuncCollector{fns: map[ast.NodeKind]renderer.NodeRendererFunc{}}
	base.RegisterFuncs(collector)
	return &attributeRenderer{Renderer: base, inner: collector.fns}
}

// RegisterFuncs 实现 renderer.NodeRenderer，只注册要补属性的那几种节点。
func (r *attributeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
}

func (r *attributeRenderer) renderImage(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// loading / decoding 均在 html.ImageAttributeFilter 白名单内，会被渲染出来
		node.SetAttributeString("loading", []byte("lazy"))
		node.SetAttributeString("decoding", []byte("async"))
	}
	return r.inner[ast.KindImage](w, source, node, entering)
}

// renderLink 处理 markdown 显式写的 []() 链接。
func (r *attributeRenderer) renderLink(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		markExternalLink(node, node.(*ast.Link).Destination)
	}
	return r.inner[ast.KindLink](w, source, node, entering)
}

// renderAutoLink 处理 <> 尖括号链接与 linkify 出来的裸链接。
func (r *attributeRenderer) renderAutoLink(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// 邮箱型 autolink 渲染成 mailto:，不是页面跳转，无需 target
		if n := node.(*ast.AutoLink); n.AutoLinkType == ast.AutoLinkURL {
			markExternalLink(node, n.URL(source))
		}
	}
	return r.inner[ast.KindAutoLink](w, source, node, entering)
}

// externalLinkPrefixes 判定外链的前缀。站内跳转在文章里写成相对路径或 #锚点，
// 不会命中；本站的绝对地址会被当成外链，但另开标签页也无妨，
// 故这里不为了比对 site.base_url 而引入配置依赖。
var externalLinkPrefixes = [][]byte{
	[]byte("http://"),
	[]byte("https://"),
	[]byte("ftp://"),
	[]byte("//"),
}

// markExternalLink 给指向站外的链接补上 target / rel，
// 避免点击后当前标签页直接跳走、丢掉文章的目录与滚动位置。
// target / rel 均在 html.LinkAttributeFilter 白名单内，会被渲染出来。
func markExternalLink(node ast.Node, dest []byte) {
	if !isExternalURL(dest) {
		return
	}
	node.SetAttributeString("target", []byte("_blank"))
	// noreferrer 一并给上：新标签页不应能通过 window.opener 反向操作本页
	node.SetAttributeString("rel", []byte("noopener noreferrer"))
}

func isExternalURL(dest []byte) bool {
	for _, prefix := range externalLinkPrefixes {
		// 协议部分大小写不敏感，作者手写的 Destination 可能是 HTTPS://
		if len(dest) >= len(prefix) && bytes.EqualFold(dest[:len(prefix)], prefix) {
			return true
		}
	}
	return false
}

// nodeFuncCollector 用于从默认 renderer 中取出各节点原本的渲染函数。
type nodeFuncCollector struct {
	fns map[ast.NodeKind]renderer.NodeRendererFunc
}

func (c *nodeFuncCollector) Register(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
	c.fns[kind] = fn
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
