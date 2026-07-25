package common

import "testing"

// TestCleanName 验证零宽字符与首尾空白被正确剔除。
func TestCleanName(t *testing.T) {
	cases := map[string]string{
		"\u200b关系型数据库": "关系型数据库", // 前缀零宽空格
		"关系型数据库":       "关系型数据库", // 干净的应保持不变
		"  Go  ":         "Go",     // 首尾空白
		"a\u200cb\u200dc": "abc",    // 中间零宽字符
		"\ufeffMySQL":     "MySQL",  // BOM 前缀
	}
	for in, want := range cases {
		if got := CleanName(in); got != want {
			t.Errorf("CleanName(%q) = %q, want %q", in, got, want)
		}
	}
}
