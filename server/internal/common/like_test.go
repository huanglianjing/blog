package common

import "testing"

// TestEscapeLike 验证 LIKE 通配符与转义符被正确转义。
func TestEscapeLike(t *testing.T) {
	cases := map[string]string{
		"go":      "go",       // 普通关键词不变
		"100%":    `100\%`,    // 百分号
		"a_b":     `a\_b`,     // 下划线
		`C:\path`: `C:\\path`, // 反斜杠自身
		`\%`:      `\\\%`,     // 反斜杠先转义，不会连锁
		"关系型数据库":  "关系型数据库",   // 中文不受影响
	}
	for in, want := range cases {
		if got := EscapeLike(in); got != want {
			t.Errorf("EscapeLike(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestLikePattern 验证关键词被包装成子串匹配模式。
func TestLikePattern(t *testing.T) {
	if got, want := LikePattern("go"), "%go%"; got != want {
		t.Errorf("LikePattern(%q) = %q, want %q", "go", got, want)
	}
	if got, want := LikePattern("50%"), `%50\%%`; got != want {
		t.Errorf("LikePattern(%q) = %q, want %q", "50%", got, want)
	}
}
