package common

import (
	"sort"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// 名称排序的分组序号：数字 < 字母 < 汉字 < 其它。
const (
	sortRankDigit   = 0
	sortRankLetter  = 1
	sortRankChinese = 2
	sortRankOther   = 3
)

var pinyinArg = pinyin.NewArgs()

// NameSortKey 生成用于名称排序的比较键，排序规则为：
// 数字 -> 字母 -> 汉字 -> 其它；字母不区分大小写，汉字按拼音升序。
// 逐字符转换：为每个字符加上分组序号前缀，字母转小写，汉字转拼音，
// 用分隔符拼接，返回的字符串可直接用 strings 字典序比较。
func NameSortKey(name string) string {
	var b strings.Builder
	for _, r := range name {
		var rank int
		var norm string
		switch {
		case unicode.IsDigit(r):
			rank = sortRankDigit
			norm = string(r)
		case r <= unicode.MaxASCII && unicode.IsLetter(r):
			rank = sortRankLetter
			norm = string(unicode.ToLower(r))
		case unicode.Is(unicode.Han, r):
			rank = sortRankChinese
			if py := pinyin.SinglePinyin(r, pinyinArg); len(py) > 0 {
				norm = py[0]
			} else {
				norm = string(r)
			}
		default:
			rank = sortRankOther
			norm = string(unicode.ToLower(r))
		}
		// 每个字符写成 "<rank><norm>\x00"，rank 保证分组优先，\x00 作为字符边界，
		// 避免前缀相同的短名排在长名之后时的错乱。
		b.WriteByte(byte('0' + rank))
		b.WriteString(norm)
		b.WriteByte(0)
	}
	return b.String()
}

// SortNames 按 NameSortKey 的规则就地排序名称列表，
// 与标签页面、分类页面的名称排序规则一致。
func SortNames(names []string) {
	keys := make(map[string]string, len(names))
	for _, n := range names {
		keys[n] = NameSortKey(n)
	}
	sort.SliceStable(names, func(i, j int) bool {
		return keys[names[i]] < keys[names[j]]
	})
}
