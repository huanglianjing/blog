package common

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// sitemapXMLNS 是 sitemap 协议 0.9 的命名空间。
const sitemapXMLNS = "http://www.sitemaps.org/schemas/sitemap/0.9"

// SitemapEntry 是 sitemap 中的一条 URL。
type SitemapEntry struct {
	// Segments 是站点内路径的各段原始文本（不含斜杠），
	// 如 []string{"article", "GORM：Go ORM框架"}；空切片表示站点首页。
	// 生成时逐段做百分号编码，规则与前端 encodeURIComponent 一致。
	Segments []string

	// LastMod 是最后修改日期（YYYY-MM-DD），为空则该条不输出 lastmod。
	LastMod string
}

// sitemapURLSet 对应 sitemap.xml 的根元素 urlset。
type sitemapURLSet struct {
	XMLName xml.Name         `xml:"urlset"`
	Xmlns   string           `xml:"xmlns,attr"`
	URLs    []sitemapURLNode `xml:"url"`
}

// sitemapURLNode 对应 sitemap.xml 中的一个 url 元素。
type sitemapURLNode struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

// BuildSitemap 按 sitemap 协议 0.9 生成 sitemap.xml 的完整内容。
// baseURL 为站点根地址（如 https://huanglianjing.com，末尾斜杠会被去掉），
// 各条 URL 由 baseURL 与编码后的路径拼成。XML 转义由 encoding/xml 负责。
func BuildSitemap(baseURL string, entries []SitemapEntry) ([]byte, error) {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		return nil, fmt.Errorf("站点根地址为空")
	}

	urls := make([]sitemapURLNode, 0, len(entries))
	for _, e := range entries {
		loc := base
		for _, seg := range e.Segments {
			loc += "/" + EncodeURIComponent(seg)
		}
		urls = append(urls, sitemapURLNode{Loc: loc, LastMod: e.LastMod})
	}

	body, err := xml.MarshalIndent(sitemapURLSet{Xmlns: sitemapXMLNS, URLs: urls}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal sitemap: %w", err)
	}

	var b strings.Builder
	b.WriteString(xml.Header) // 自带结尾换行
	b.Write(body)
	b.WriteByte('\n')
	return []byte(b.String()), nil
}

// EncodeURIComponent 对字符串做百分号编码，规则与浏览器 encodeURIComponent
// 完全一致：只保留 A-Za-z0-9 与 -_.!~*'() ，其余字节按 UTF-8 逐字节转成 %XX。
// 不能用 url.PathEscape 代替：它会保留 $&+,:;=@ 等字符，
// 与前端 router-link 生成的路径不一致，会让同一页面出现两种 URL 写法。
func EncodeURIComponent(s string) string {
	const upperhex = "0123456789ABCDEF"
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isURIComponentUnreserved(c) {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		b.WriteByte(upperhex[c>>4])
		b.WriteByte(upperhex[c&0x0f])
	}
	return b.String()
}

// isURIComponentUnreserved 判断字节是否属于 encodeURIComponent 不编码的字符。
func isURIComponentUnreserved(c byte) bool {
	switch {
	case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9':
		return true
	}
	switch c {
	case '-', '_', '.', '!', '~', '*', '\'', '(', ')':
		return true
	}
	return false
}
