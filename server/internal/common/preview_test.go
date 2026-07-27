package common

import "testing"

const testHTML = `<h2>安装步骤</h2>
<p>先执行下面的命令。</p>
<pre><code>go install example.com/tool@latest</code></pre>
<script>var secret = 1;</script>
<p>然后重启服务。</p>`

// TestHTMLToPreview 验证预览跳过标题与代码块，并按字符数截断。
func TestHTMLToPreview(t *testing.T) {
	got := HTMLToPreview([]byte(testHTML), 200)
	want := "先执行下面的命令。 然后重启服务。"
	if got != want {
		t.Errorf("HTMLToPreview() = %q, want %q", got, want)
	}

	if got := HTMLToPreview([]byte(testHTML), 3); len([]rune(got)) != 3 {
		t.Errorf("HTMLToPreview() 截断到 3 字符, got %q", got)
	}
}

// TestHTMLToPlainText 验证搜索用全文保留标题与代码块、排除脚本，且不截断。
func TestHTMLToPlainText(t *testing.T) {
	got := HTMLToPlainText([]byte(testHTML))
	want := "安装步骤 先执行下面的命令。 go install example.com/tool@latest 然后重启服务。"
	if got != want {
		t.Errorf("HTMLToPlainText() = %q, want %q", got, want)
	}
}
