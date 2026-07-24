# blog 目标的目标平台，默认为 Linux amd64（部署到服务器）。
# 本机运行可覆盖，例如：make blog GOOS=darwin GOARCH=arm64
GOOS  ?= linux
GOARCH ?= amd64
# sqlite 驱动为纯 Go 实现（modernc.org/sqlite），关闭 CGO 以生成静态二进制、支持交叉编译。
GOBUILD := CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build

.PHONY: release dev clean

# 构建生产环境
release:
	rm -rf blog blog.tar.gz
	mkdir -p blog/config
	# 后端：交叉编译二进制到 blog/（默认 linux/amd64）
	cd server && $(GOBUILD) -o ../blog/blog_server ./cmd/blog_server
	cd server && $(GOBUILD) -o ../blog/article_converter ./cmd/article_converter
	cp ./server/config/config.yaml ./blog/config/config.yaml
	# 前端：构建 dist 并放入 blog/
	cd web && npm install && npm run build
	cp -r ./web/dist ./blog/dist
	# 压缩整个 blog/ 目录
	tar zcf blog.tar.gz blog

# 构建开发环境
dev:
	mkdir -p blog_dev
	mkdir -p blog_dev/config
	mkdir -p data
	cd server && go build -o ../blog_dev/blog_server ./cmd/blog_server
	cd server && go build -o ../blog_dev/article_converter ./cmd/article_converter
	./blog_dev/article_converter -src ../article -db data/db/blog.db -out data/article_html
	cp ./server/config/config.yaml ./blog_dev/config/config.yaml

# 清理构建产物
clean:
	rm -rf blog blog.tar.gz blog_dev data
