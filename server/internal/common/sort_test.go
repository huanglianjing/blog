package common

import (
	"sort"
	"testing"
)

// TestNameSortKey 验证排序规则：数字 -> 字母 -> 汉字，
// 字母不区分大小写，汉字按拼音升序。
func TestNameSortKey(t *testing.T) {
	input := []string{"苹果", "Golang", "阿里", "123", "banana", "3D", "香蕉", "Apple"}
	want := []string{"123", "3D", "Apple", "banana", "Golang", "阿里", "苹果", "香蕉"}

	got := make([]string, len(input))
	copy(got, input)
	sort.SliceStable(got, func(i, j int) bool {
		return NameSortKey(got[i]) < NameSortKey(got[j])
	})

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("排序结果不符\n got:  %v\n want: %v", got, want)
		}
	}
}
