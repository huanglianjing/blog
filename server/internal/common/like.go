package common

import "strings"

// LikeEscapeChar 是 LIKE 模式中使用的转义字符，
// SQL 需相应写成 `col LIKE ? ESCAPE '\'`。
const LikeEscapeChar = `\`

// EscapeLike 转义 LIKE 模式中的特殊字符，使关键词按字面量匹配。
// 反斜杠必须最先替换，否则会把后续新加的转义符再转义一遍。
func EscapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "%", `\%`)
	s = strings.ReplaceAll(s, "_", `\_`)
	return s
}

// LikePattern 把关键词包装成子串匹配的 LIKE 模式 %关键词%。
func LikePattern(keyword string) string {
	return "%" + EscapeLike(keyword) + "%"
}
