package common

import "strings"

// 需要从名称中剔除的不可见字符（零宽字符 / BOM）。
// 用数值码点转义而非字面量，避免源码里嵌入不可见字符本身。
const (
	runeZWSP = '\u200b' // 零宽空格 ZERO WIDTH SPACE
	runeZWNJ = '\u200c' // 零宽非连接符 ZERO WIDTH NON-JOINER
	runeZWJ  = '\u200d' // 零宽连接符 ZERO WIDTH JOINER
	runeBOM  = '\ufeff' // 字节序标记 / 零宽不换行空格
)

// CleanName 规范化名称文本：去除零宽字符与首尾空白。
// 用于清洗 meta 中的分类名、标题、标签名，避免不可见字符（如 U+200B
// 零宽空格）导致本应相同的名称被当成两个，进而破坏去重与排序。
func CleanName(s string) string {
	cleaned := strings.Map(func(r rune) rune {
		switch r {
		case runeZWSP, runeZWNJ, runeZWJ, runeBOM:
			return -1 // 丢弃
		}
		return r
	}, s)
	return strings.TrimSpace(cleaned)
}
