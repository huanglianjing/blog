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
	cd web && npm install --no-audit --no-fund && npm run build
	cp -r ./web/dist ./blog/dist
	# 压缩整个 blog/ 目录
	tar zcf blog.tar.gz blog

# 开发环境后台进程的 pid 与日志文件
BACKEND_PID  := $(CURDIR)/blog_dev/blog_server.pid
FRONTEND_PID := $(CURDIR)/blog_dev/vite.pid

# 构建开发环境，并在后台启动前后端服务
dev:
	mkdir -p blog_dev
	mkdir -p blog_dev/config
	mkdir -p data
	cd server && go build -o ../blog_dev/blog_server ./cmd/blog_server
	cd server && go build -o ../blog_dev/article_converter ./cmd/article_converter
	./blog_dev/article_converter -src ../article -db data/db/blog.db -out data/article_html \
		-sitemap data/sitemap.xml -c ./server/config/config.yaml
	cp ./server/config/config.yaml ./blog_dev/config/config.yaml
	@if [ -f $(BACKEND_PID) ] && kill -0 $$(cat $(BACKEND_PID)) 2>/dev/null; then \
		echo "后端已在运行（pid $$(cat $(BACKEND_PID))），跳过"; \
	else \
		cd blog_dev; nohup ./blog_server > blog_server.log 2>&1 & echo $$! > $(BACKEND_PID); \
		echo "后端已启动（pid $$(cat $(BACKEND_PID))），日志 blog_dev/blog_server.log"; \
	fi
	@if [ -f $(FRONTEND_PID) ] && kill -0 $$(cat $(FRONTEND_PID)) 2>/dev/null; then \
		echo "前端已在运行（pid $$(cat $(FRONTEND_PID))），跳过"; \
	else \
		[ -d web/node_modules ] || (cd web && npm install --no-audit --no-fund); \
		cd web; nohup ./node_modules/.bin/vite > ../blog_dev/vite.log 2>&1 & echo $$! > $(FRONTEND_PID); \
		echo "前端已启动（pid $$(cat $(FRONTEND_PID))），日志 blog_dev/vite.log"; \
	fi
	@echo "浏览器打开 http://localhost:5173/"

# 停止后台的前后端进程，清理构建产物
clean:
	@for f in $(BACKEND_PID) $(FRONTEND_PID); do \
		if [ -f $$f ] && kill $$(cat $$f) 2>/dev/null; then \
			echo "已停止 pid $$(cat $$f)"; \
		fi; \
	done
	rm -rf blog blog.tar.gz blog_dev data
