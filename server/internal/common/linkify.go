package common

import (
	"bytes"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// bareURLPattern 匹配裸链接。字符集与 goldmark linkify 扩展内置的
// urlRegexp / wwwURLRegxp 保持一致，区别只有一处：合并成一条正则
// （副作用是 www. 开头的地址也支持端口号了，属于放宽而非收紧）。
const bareURLPattern = `(?:https?://|ftp://|www\.)` +
	`[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]+(?::\d+)?` +
	`(?:[/#?][-a-zA-Z0-9@:%_+.~#$!?&/=\(\);,'">\^{}\[\]` + "`" + `]*)?`

var (
	// 在一段文本中间找裸链接。左边界不靠 ^ 锚定，改由 hasURLBoundaryBefore 判断。
	bareURLRegexp = regexp.MustCompile(bareURLPattern)
	// 锚定版，用于从某个已知起点量出 URL 的真实长度，见 urlExtendsBeyond
	anchoredURLRegexp = regexp.MustCompile(`^(?:` + bareURLPattern + `)`)
)

// urlPrefixRejectChars 是那些出现在 URL 之前，就说明这段文本只是更长 token
// 的一部分（而非一个独立链接）的 ASCII 标点，如 `foo=https://a.com` 中的 `=`。
// 字母数字另行判断。
const urlPrefixRejectChars = `-._~:/@%+=#?&`

// linkifySkipKinds 是不做裸链接识别的节点：代码、已有链接、原始 HTML。
// 遍历时整棵子树跳过，所以行内代码与围栏代码块（含 mermaid）里的 URL
// 会保持纯文本，显式 markdown 链接的 label 里也不会再嵌一层 a 标签。
var linkifySkipKinds = map[ast.NodeKind]bool{
	ast.KindCodeSpan:        true,
	ast.KindLink:            true,
	ast.KindAutoLink:        true,
	ast.KindRawHTML:         true,
	ast.KindHTMLBlock:       true,
	ast.KindCodeBlock:       true,
	ast.KindFencedCodeBlock: true,
}

// bareLinkTransformer 把 goldmark 内置 linkify 漏掉的裸链接补成 AutoLink 节点。
//
// 内置 linkify（extension.GFM 的一部分）是个 inline parser，只在
// ' ' '*' '_' '~' '(' 五个字符处和行首触发（见 extension/linkify.go 的 Trigger），
// 所以「中文紧贴https://example.com」「（https://example.com）」这类中文写法
// 识别不到。而且这个缺口没法靠再写一个 inline parser 补上：goldmark 只在
// ASCII 标点、空白和行首处派发 inline parser（见 parser.parseBlock 里的
// util.IsPunct 判断，其查表对 >= 0x80 的字节一律返回 false），中文字符处
// 根本不会调用任何 inline parser，改 Trigger 也没用。
//
// 所以只能等 inline 解析完再遍历一遍 AST 补。已被内置 linkify 处理掉的部分
// 此时已经是 AutoLink 节点、不再是 Text，故两者天然互补，不会重复识别。
//
// 不识别 example.com 这类不带协议也不带 www 的裸域名：它与 config.io、
// 句子里的 "etc.com" 无法区分，误判代价大于收益。
type bareLinkTransformer struct{}

func newBareLinkTransformer() parser.ASTTransformer {
	return &bareLinkTransformer{}
}

// Transform 实现 parser.ASTTransformer。
func (t *bareLinkTransformer) Transform(doc *ast.Document, reader text.Reader, _ parser.Context) {
	source := reader.Source()

	// 先收集再改写：ast.Walk 靠 NextSibling 推进，边遍历边把节点
	// 从父节点上摘下来会让它的兄弟指针被清空，遍历就此中断。
	var runs []textRun
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if linkifySkipKinds[n.Kind()] {
			return ast.WalkSkipChildren, nil
		}
		runs = append(runs, collectTextRuns(n)...)
		return ast.WalkContinue, nil
	})

	for _, run := range runs {
		linkifyTextRun(run, source)
	}
}

// textRun 是同一父节点下一串在源文里首尾相接的文本节点，当成一整段来匹配 URL。
//
// 必须合并，因为 markdown 的行内语法会把文本切碎：URL 里带 _ 时
// （真实例子 ?version=1.3&_preview=1），emphasis 的分隔符处理会在 _ 处
// 断开文本节点，单看头一个节点只能匹配到 ...&，生成的 href 就是错的。
type textRun []*ast.Text

// collectTextRuns 把 parent 的直接子节点按「连续的文本节点」切成若干段。
// 相邻文本节点只有源区间首尾相接、且前一段不以换行结尾才算同一段；
// 中间隔了换行或别的行内节点就必须断开，否则会把不相干的文字拼进 URL。
func collectTextRuns(parent ast.Node) []textRun {
	var runs []textRun
	var cur textRun
	flush := func() {
		if len(cur) > 0 {
			runs = append(runs, cur)
			cur = nil
		}
	}

	for c := parent.FirstChild(); c != nil; c = c.NextSibling() {
		node, ok := c.(*ast.Text)
		// raw text 是行内代码的内容，正常已被 linkifySkipKinds 挡住，这里兜底
		if !ok || node.IsRaw() {
			flush()
			continue
		}
		if len(cur) > 0 {
			prev := cur[len(cur)-1]
			if prev.Segment.Stop != node.Segment.Start ||
				prev.SoftLineBreak() || prev.HardLineBreak() {
				flush()
			}
		}
		cur = append(cur, node)
	}
	flush()
	return runs
}

// linkifyTextRun 把 run 中的裸链接换成 AutoLink 节点：
// 命中时用「文本 / 链接 / 文本 …」若干节点原地替换掉整个 run，未命中则不动。
func linkifyTextRun(run textRun, source []byte) {
	first, last := run[0], run[len(run)-1]
	parent := first.Parent()
	if parent == nil {
		return
	}

	lo, hi := first.Segment.Start, last.Segment.Stop
	// 用 source 切片而不是 Segment.Value()：后者会在头部补 Padding 个空格、
	// 尾部按 ForceNewline 补换行，匹配到的下标就对不上源文位置了。
	value := source[lo:hi]

	var replacements []ast.Node
	cursor := 0
	for _, m := range bareURLRegexp.FindAllIndex(value, -1) {
		start, stop := m[0], m[1]
		if !hasURLBoundaryBefore(value, start) {
			continue
		}
		stop = start + trimURLTrailing(value[start:stop])
		if stop <= start {
			continue
		}
		// URL 一直顶到这段文本的末尾、而源文里它其实还没结束时，说明真正的
		// 地址被行内语法切断在这里了（URL 中成对的 _ 或 * 被解析成了
		// emphasis，那部分已经是独立节点，拼不回来）。这时宁可整段留作纯文本，
		// 也不要生成一个指向错误地址的链接。
		if stop == len(value) && urlExtendsBeyond(source, lo+start, hi) {
			continue
		}

		if start > cursor {
			replacements = append(replacements, runTextNode(run, cursor, start))
		}
		replacements = append(replacements, autoLinkNode(value[start:stop], lo+start, lo+stop))
		cursor = stop
	}
	if len(replacements) == 0 {
		return
	}

	if cursor < len(value) {
		replacements = append(replacements, runTextNode(run, cursor, len(value)))
	} else if last.SoftLineBreak() || last.HardLineBreak() || last.Segment.ForceNewline {
		// URL 顶到末尾时没有尾段文本可以承载换行标记，补一个空文本节点，
		// 否则原本的软换行会丢，下一行会被直接接到 URL 后面。
		replacements = append(replacements, runTextNode(run, len(value), len(value)))
	}

	for _, r := range replacements {
		parent.InsertBefore(parent, first, r)
	}
	for _, node := range run {
		parent.RemoveChild(parent, node)
	}
}

// runTextNode 从 run 中切出 [from, to) 这段（下标相对 run 的整段文本）构造文本节点。
// Padding 是整段文本前的缩进空格，只有首段该带；ForceNewline 与换行标记表示
// 「本节点以换行结尾」，只有末段该带。
func runTextNode(run textRun, from, to int) *ast.Text {
	first, last := run[0], run[len(run)-1]
	total := last.Segment.Stop - first.Segment.Start
	atEnd := to == total

	sub := first.Segment
	sub.Start = first.Segment.Start + from
	sub.Stop = first.Segment.Start + to
	if from != 0 {
		sub.Padding = 0
	}
	sub.ForceNewline = atEnd && last.Segment.ForceNewline

	out := ast.NewTextSegment(sub)
	if atEnd {
		out.SetSoftLineBreak(last.SoftLineBreak())
		out.SetHardLineBreak(last.HardLineBreak())
	}
	return out
}

// autoLinkNode 用 [start, stop) 这段源文本构造一个 URL 型 AutoLink 节点。
func autoLinkNode(url []byte, start, stop int) *ast.AutoLink {
	link := ast.NewAutoLink(ast.AutoLinkURL, ast.NewTextSegment(text.NewSegment(start, stop)))
	// 没写协议的 www.xxx 要补上，否则渲染出的 href 会是站内相对路径
	if bytes.HasPrefix(url, []byte("www.")) {
		link.Protocol = []byte("http")
	}
	return link
}

// urlExtendsBeyond 判断源文 start 处的 URL 是否越过了 hi。
// 越过说明 hi 并非 URL 的真正结尾，而是被行内语法切断的位置。
func urlExtendsBeyond(source []byte, start, hi int) bool {
	m := anchoredURLRegexp.FindIndex(source[start:])
	return m != nil && start+m[1] > hi
}

// hasURLBoundaryBefore 判断 s 中 i 处的 URL 左侧是否是一个合法的起始边界。
// 前一个字节属于多字节 UTF-8（>= 0x80，即中文、全角标点等）时一律放行 ——
// 这正是 goldmark 内置 linkify 覆盖不到、本文件要补上的场景。
func hasURLBoundaryBefore(s []byte, i int) bool {
	if i == 0 {
		return true
	}
	c := s[i-1]
	if c >= utf8.RuneSelf {
		return true
	}
	if util.IsAlphaNumeric(c) {
		return false
	}
	return strings.IndexByte(urlPrefixRejectChars, c) < 0
}

// trimURLTrailing 剥掉 url 尾部不该算进链接的字符，返回剩下的长度。
// 规则照搬 goldmark linkify 扩展 Parse 方法的尾部处理，
// 免得同一篇文章里两套 linkify 对尾随标点的处理不一致。
func trimURLTrailing(url []byte) int {
	stop := len(url)
	if stop == 0 {
		return 0
	}

	switch url[stop-1] {
	case '.':
		stop--
	case ')':
		// 右括号多于左括号时，多出来的是包裹用的，如「(见 https://a.com/b)」
		closing := 0
		for i := stop - 1; i >= 0; i-- {
			switch url[i] {
			case ')':
				closing++
			case '(':
				closing--
			}
		}
		if closing > 0 {
			stop -= closing
		}
	case ';':
		// 结尾是 HTML 实体（&amp; 之类）时整个实体都不算
		i := stop - 2
		for ; i >= 0 && util.IsAlphaNumeric(url[i]); i-- {
		}
		if i >= 0 && i != stop-2 && url[i] == '&' {
			stop = i
		}
	}

	// URL 后面直接跟句读符号是常见写法，一并剥掉
	for stop > 0 {
		switch url[stop-1] {
		case '?', '!', '.', ',', ':', '*', '_', '~':
			stop--
			continue
		}
		break
	}
	return stop
}
