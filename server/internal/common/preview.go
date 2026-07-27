package common

import (
	"strings"

	"golang.org/x/net/html"
)

// 提取预览时跳过的块级元素（标题、图片、表格、代码块、图形等）。
var previewSkipTags = map[string]bool{
	"h1":         true,
	"h2":         true,
	"h3":         true,
	"h4":         true,
	"h5":         true,
	"h6":         true,
	"img":        true,
	"table":      true,
	"pre":        true,
	"code":       true,
	"figure":     true,
	"svg":        true,
	"picture":    true,
	"video":      true,
	"audio":      true,
	"iframe":     true,
	"script":     true,
	"style":      true,
	"blockquote": true,
}

// 提取搜索用全文时跳过的元素：只排除非文字内容，
// 标题、代码块、表格、引用中的文字都应可被搜索到。
var plainTextSkipTags = map[string]bool{
	"svg":    true,
	"video":  true,
	"audio":  true,
	"iframe": true,
	"script": true,
	"style":  true,
}

// HTMLToPreview 从 html 正文中提取纯文本预览，
// 跳过图片、表格、代码块等非文字内容，最多返回 maxRunes 个字符。
func HTMLToPreview(source []byte, maxRunes int) string {
	return htmlToText(source, previewSkipTags, maxRunes)
}

// HTMLToPlainText 提取 html 的全部正文纯文本，用于搜索匹配，不截断。
func HTMLToPlainText(source []byte) string {
	return htmlToText(source, plainTextSkipTags, 0)
}

// htmlToText 遍历 html 提取纯文本，跳过 skip 中的元素，
// 归一化空白后按字符数截断；maxRunes <= 0 表示不限长度。
func htmlToText(source []byte, skip map[string]bool, maxRunes int) string {
	doc, err := html.Parse(strings.NewReader(string(source)))
	if err != nil {
		return ""
	}

	var b strings.Builder
	collectText(doc, &b, skip, maxRunes)

	// 归一化空白：把连续空白折叠成单个空格。
	text := strings.Join(strings.Fields(b.String()), " ")

	// 按字符数截断。
	runes := []rune(text)
	if maxRunes > 0 && len(runes) > maxRunes {
		return string(runes[:maxRunes])
	}
	return text
}

// collectText 深度遍历 DOM，把文本节点写入 b，
// 遇到 skip 中的元素则整棵子树略过。
// maxRunes > 0 时累计长度达到上限即停止，<= 0 表示不限长度。
func collectText(n *html.Node, b *strings.Builder, skip map[string]bool, maxRunes int) {
	if maxRunes > 0 && b.Len() >= maxRunes*4 { // 按 utf-8 最坏情况粗略限制，避免拼过长字符串
		return
	}

	if n.Type == html.ElementNode && skip[n.Data] {
		return
	}

	if n.Type == html.TextNode {
		b.WriteString(n.Data)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, b, skip, maxRunes)
	}
}
