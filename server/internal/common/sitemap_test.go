package common

import (
	"strings"
	"testing"
)

func TestEncodeURIComponent(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Go", "Go"},
		{"Claude Code", "Claude%20Code"},
		{"a-_.!~*'()", "a-_.!~*'()"},
		// 汉字与全角标点按 UTF-8 逐字节编码。
		{"分布式", "%E5%88%86%E5%B8%83%E5%BC%8F"},
		{"GORM：Go", "GORM%EF%BC%9AGo"},
		// url.PathEscape 会原样保留这些字符，encodeURIComponent 不会。
		{"a&b=c+d/e:f", "a%26b%3Dc%2Bd%2Fe%3Af"},
		{"", ""},
	}
	for _, c := range cases {
		if got := EncodeURIComponent(c.in); got != c.want {
			t.Errorf("EncodeURIComponent(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestBuildSitemap(t *testing.T) {
	entries := []SitemapEntry{
		{Segments: nil, LastMod: "2026-08-08"},
		{Segments: []string{"article"}, LastMod: "2026-08-08"},
		{Segments: []string{"article", "Claude Code使用指南"}, LastMod: "2026-05-05"},
		{Segments: []string{"tag", "Go"}}, // 无 lastmod
	}
	// base_url 末尾多余的斜杠应被去掉，不能出现 // 。
	got, err := BuildSitemap("https://example.com/", entries)
	if err != nil {
		t.Fatalf("BuildSitemap: %v", err)
	}
	out := string(got)

	for _, want := range []string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`,
		"<loc>https://example.com</loc>",
		"<loc>https://example.com/article</loc>",
		"<loc>https://example.com/article/Claude%20Code%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97</loc>",
		"<loc>https://example.com/tag/Go</loc>",
		"<lastmod>2026-08-08</lastmod>",
		"</urlset>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("sitemap 缺少 %q，实际内容:\n%s", want, out)
		}
	}
	if strings.Contains(out, "com//") {
		t.Errorf("base_url 末尾斜杠未去掉:\n%s", out)
	}
	// 无 lastmod 的条目不应输出空的 lastmod 标签。
	if strings.Contains(out, "<lastmod></lastmod>") {
		t.Errorf("出现空 lastmod 标签:\n%s", out)
	}
	if n := strings.Count(out, "<url>"); n != len(entries) {
		t.Errorf("url 条数 = %d, want %d", n, len(entries))
	}
}

func TestBuildSitemapEmptyBaseURL(t *testing.T) {
	if _, err := BuildSitemap("  ", nil); err == nil {
		t.Error("base_url 为空时应报错")
	}
}
