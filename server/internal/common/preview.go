package common

import (
	"strings"

	"golang.org/x/net/html"
)

// 提取预览时跳过的块级元素（标题、图片、表格、代码块、图形等）。
var previewSkipTags = map[string]bool{
	"h1":       true,
	"h2":       true,
	"h3":       true,
	"h4":       true,
	"h5":       true,
	"h6":       true,
	"img":      true,
	"table":    true,
	"pre":      true,
	"code":     true,
	"figure":   true,
	"svg":      true,
	"picture":  true,
	"video":    true,
	"audio":    true,
	"iframe":   true,
	"script":   true,
	"style":    true,
	"blockquote": true,
}

// HTMLToPreview 从 html 正文中提取纯文本预览，
// 跳过图片、表格、代码块等非文字内容，最多返回 maxRunes 个字符。
func HTMLToPreview(source []byte, maxRunes int) string {
	doc, err := html.Parse(strings.NewReader(string(source)))
	if err != nil {
		return ""
	}

	var b strings.Builder
	collectPreviewText(doc, &b, maxRunes)

	// 归一化空白：把连续空白折叠成单个空格。
	text := strings.Join(strings.Fields(b.String()), " ")

	// 按字符数截断。
	runes := []rune(text)
	if maxRunes > 0 && len(runes) > maxRunes {
		return string(runes[:maxRunes])
	}
	return text
}

// collectPreviewText 深度遍历 DOM，把文本节点写入 b，
// 遇到需要跳过的元素则整棵子树略过；累计长度达到 maxRunes 即停止。
func collectPreviewText(n *html.Node, b *strings.Builder, maxRunes int) {
	if b.Len() >= maxRunes*4 { // 按 utf-8 最坏情况粗略限制，避免拼过长字符串
		return
	}

	if n.Type == html.ElementNode && previewSkipTags[n.Data] {
		return
	}

	if n.Type == html.TextNode {
		b.WriteString(n.Data)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectPreviewText(c, b, maxRunes)
	}
}
