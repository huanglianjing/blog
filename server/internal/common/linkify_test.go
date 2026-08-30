package common

import (
	"strings"
	"testing"
)

// TestLinkifyBareURL 验证裸链接识别：内置 linkify 覆盖的场景要保持不变，
// 它漏掉的中文紧贴、全角标点场景要补上，代码与已有链接内不得改动。
func TestLinkifyBareURL(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want string
	}{
		// 内置 linkify 本就支持的场景，不能被改坏
		{
			name: "空格分隔",
			src:  "访问 https://example.com/a?b=1 看看",
			want: `访问 <a href="https://example.com/a?b=1" target="_blank" rel="noopener noreferrer">https://example.com/a?b=1</a> 看看`,
		},
		{
			name: "www 补协议",
			src:  "访问 www.example.com 看看",
			want: `访问 <a href="http://www.example.com" target="_blank" rel="noopener noreferrer">www.example.com</a> 看看`,
		},

		// 本次要补上的场景
		{
			name: "中文紧贴",
			src:  "见https://example.com结尾",
			want: `见<a href="https://example.com" target="_blank" rel="noopener noreferrer">https://example.com</a>结尾`,
		},
		{
			name: "全角括号包裹",
			src:  "括号（https://example.com）",
			want: `括号（<a href="https://example.com" target="_blank" rel="noopener noreferrer">https://example.com</a>）`,
		},
		{
			name: "中文书名号包裹",
			src:  "「https://example.com/x」",
			want: `「<a href="https://example.com/x" target="_blank" rel="noopener noreferrer">https://example.com/x</a>」`,
		},
		{
			name: "中文紧贴 www",
			src:  "官网www.example.com上",
			want: `官网<a href="http://www.example.com" target="_blank" rel="noopener noreferrer">www.example.com</a>上`,
		},
		{
			name: "中文紧贴 ftp",
			src:  "镜像ftp://example.com/pub里",
			want: `镜像<a href="ftp://example.com/pub" target="_blank" rel="noopener noreferrer">ftp://example.com/pub</a>里`,
		},
		{
			name: "一行内两个",
			src:  "从https://a.example.com到https://b.example.com",
			want: `从<a href="https://a.example.com" target="_blank" rel="noopener noreferrer">https://a.example.com</a>到` +
				`<a href="https://b.example.com" target="_blank" rel="noopener noreferrer">https://b.example.com</a>`,
		},
		{
			name: "中文句号不算进链接",
			src:  "见https://example.com。",
			want: `见<a href="https://example.com" target="_blank" rel="noopener noreferrer">https://example.com</a>。`,
		},
		{
			name: "半角句号剥离",
			src:  "见https://example.com/a.",
			want: `见<a href="https://example.com/a" target="_blank" rel="noopener noreferrer">https://example.com/a</a>.`,
		},
		{
			name: "尾随逗号剥离",
			src:  "见https://example.com/a,然后",
			want: `见<a href="https://example.com/a" target="_blank" rel="noopener noreferrer">https://example.com/a</a>,然后`,
		},

		// URL 里带 _ 时 emphasis 的分隔符处理会把文本节点切碎，
		// 必须跨节点拼回完整地址，否则 href 会截断在 _ 之前
		{
			name: "URL 含未配对的下划线",
			src:  "返回码：https://a.example.com/x?version=1.3&_preview=1",
			want: `返回码：<a href="https://a.example.com/x?version=1.3&amp;_preview=1" target="_blank" rel="noopener noreferrer">` +
				`https://a.example.com/x?version=1.3&amp;_preview=1</a>`,
		},
		{
			name: "路径含未配对的下划线",
			src:  "见https://a.example.com/a_b_c",
			want: `见<a href="https://a.example.com/a_b_c" target="_blank" rel="noopener noreferrer">https://a.example.com/a_b_c</a>`,
		},
		// 不该识别的场景
		{
			// URL 里的下划线成对配上了，那几段已经是 emphasis 节点、拼不回来。
			// 这时有意放弃识别：截断的 href 会指向错误地址，纯文本至少还是对的。
			// 想要链接就用 []() 显式写，或在 URL 前留个半角空格。
			name: "下划线配对成 emphasis 时放弃",
			src:  "见https://a.example.com/x?a=_1_&b=_2_",
			want: `见https://a.example.com/x?a=<em>1</em>&amp;b=<em>2</em>`,
		},
		{
			name: "无协议裸域名",
			src:  "配置 example.com 即可",
			want: "配置 example.com 即可",
		},
		{
			name: "紧跟在字母后",
			src:  "见xhttps://example.com",
			want: "见xhttps://example.com",
		},
		{
			name: "已是 query 参数的一部分",
			src:  "见https://a.example.com/x?u=https://b.example.com",
			want: `见<a href="https://a.example.com/x?u=https://b.example.com" target="_blank" rel="noopener noreferrer">` +
				`https://a.example.com/x?u=https://b.example.com</a>`,
		},
		{
			name: "行内代码",
			src:  "值为`https://example.com`时",
			want: "值为<code>https://example.com</code>时",
		},
		{
			name: "显式链接的 label 内不再嵌套",
			src:  "[见https://a.example.com](https://b.example.com)",
			want: `<a href="https://b.example.com" target="_blank" rel="noopener noreferrer">见https://a.example.com</a>`,
		},
		{
			name: "站内相对链接不加 target",
			src:  "[另一篇](/article/foo)",
			want: `<a href="/article/foo">另一篇</a>`,
		},
		{
			name: "锚点不加 target",
			src:  "[小节](#anchor)",
			want: `<a href="#anchor">小节</a>`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := MarkdownToHTML([]byte(c.src))
			if err != nil {
				t.Fatalf("MarkdownToHTML() error = %v", err)
			}
			want := "<p>" + c.want + "</p>\n"
			if string(got) != want {
				t.Errorf("MarkdownToHTML(%q)\n got = %q\nwant = %q", c.src, got, want)
			}
		})
	}
}

// TestLinkifySkipCodeBlock 验证围栏代码块（含 mermaid）里的 URL 保持纯文本。
func TestLinkifySkipCodeBlock(t *testing.T) {
	src := "```\n见https://example.com\n```\n"
	got, err := MarkdownToHTML([]byte(src))
	if err != nil {
		t.Fatalf("MarkdownToHTML() error = %v", err)
	}
	if strings.Contains(string(got), "<a") {
		t.Errorf("代码块内的 URL 不应被转成链接, got %q", got)
	}
}

// TestLinkifyKeepsSoftLineBreak 验证 URL 顶到行尾时软换行不丢，
// 否则下一行会被直接接到 URL 后面。
func TestLinkifyKeepsSoftLineBreak(t *testing.T) {
	src := "见https://example.com\n下一行"
	got, err := MarkdownToHTML([]byte(src))
	if err != nil {
		t.Fatalf("MarkdownToHTML() error = %v", err)
	}
	if !strings.Contains(string(got), "</a>\n下一行") {
		t.Errorf("URL 在行尾时应保留换行, got %q", got)
	}
}

// TestLinkifyInListItem 验证列表项内（Segment 带 Padding）的裸链接
// 识别后缩进不错位、文字不丢。
func TestLinkifyInListItem(t *testing.T) {
	src := "- 参考https://example.com/a 说明\n- 第二项\n"
	got, err := MarkdownToHTML([]byte(src))
	if err != nil {
		t.Fatalf("MarkdownToHTML() error = %v", err)
	}
	out := string(got)
	if !strings.Contains(out, `>https://example.com/a</a> 说明`) {
		t.Errorf("列表项内链接与后续文字应完整, got %q", out)
	}
	if !strings.Contains(out, "第二项") {
		t.Errorf("列表其余项不应受影响, got %q", out)
	}
}

// TestLinkifyInTableCell 验证表格单元格内的裸链接也能识别。
func TestLinkifyInTableCell(t *testing.T) {
	src := "| 名称 | 地址 |\n| --- | --- |\n| 示例 | 见https://example.com |\n"
	got, err := MarkdownToHTML([]byte(src))
	if err != nil {
		t.Fatalf("MarkdownToHTML() error = %v", err)
	}
	if !strings.Contains(string(got), `<a href="https://example.com"`) {
		t.Errorf("表格单元格内的 URL 应转成链接, got %q", got)
	}
}

// TestLinkifyEmailUnaffected 验证邮箱仍由内置 linkify 处理、且不加 target。
func TestLinkifyEmailUnaffected(t *testing.T) {
	got, err := MarkdownToHTML([]byte("邮箱 foo@example.com 联系"))
	if err != nil {
		t.Fatalf("MarkdownToHTML() error = %v", err)
	}
	want := `<p>邮箱 <a href="mailto:foo@example.com">foo@example.com</a> 联系</p>` + "\n"
	if string(got) != want {
		t.Errorf("MarkdownToHTML()\n got = %q\nwant = %q", got, want)
	}
}

// TestLinkifyRawHTMLUnaffected 验证文章内联的原始 HTML 不被改写
// （WithUnsafe 开着，属性里的 URL 不能被当成文本 linkify）。
func TestLinkifyRawHTMLUnaffected(t *testing.T) {
	src := `<div data-src="https://example.com/x">内容</div>`
	got, err := MarkdownToHTML([]byte(src))
	if err != nil {
		t.Fatalf("MarkdownToHTML() error = %v", err)
	}
	if !strings.Contains(string(got), src) {
		t.Errorf("原始 HTML 应原样保留, got %q", got)
	}
}
