# AGENTS.md

个人博客网站的前后端单体仓库。后端 Go + Gin + GORM(SQLite)，前端 Vue 3 + Vite。

## 仓库结构

```
├── server/              // 后端（Go）
│   ├── cmd/
│   │   ├── blog_server/       // HTTP 服务入口
│   │   └── article_converter/ // 文章转换 + 写库的命令行工具
│   ├── config/config.yaml     // 服务配置（端口、db 路径）
│   └── internal/
│       ├── controller/  // HTTP handler，解析参数、调 service、统一返回
│       ├── service/     // 业务逻辑，组装各 model 查询结果
│       ├── model/       // GORM 模型 + 数据库查询；db.go 负责连接/建表
│       ├── common/      // markdown 转换、正文预览提取、统一响应
│       ├── config/      // 配置结构与加载
│       └── router/      // 路由注册
├── web/                 // 前端（Vue 3 + Vite）
│   └── src/
│       ├── views/       // 页面级组件（每个路由一个）
│       ├── components/  // 复用组件（TopBar、Pagination）
│       ├── router/      // vue-router 配置
│       └── style.css    // 全局样式
├── Makefile             // 交叉编译 + 前端构建 + 打包发布
└── README.md            // 部署说明（打包、Nginx、文章部署）
```

## 数据流

内容不是运行时从 markdown 渲染的，而是**离线预生成**：

1. 源目录含 `blog_meta.yaml`（登记分类、文章标题、日期、标签）和各分类子目录下的 `.md` 文件。
2. `article_converter -src <源目录> -out <html输出目录> -db <sqlite文件>` 读取 meta，对**真实存在**的 md 文件用 goldmark 转成 html 写入 `<out>/分类/标题.html`（meta 有记录但 md 缺失的会跳过并打日志），同时把分类 / 标签 / 文章 / 文章-标签关联**全量重建**写入 SQLite（事务内先清空各表再写入，见 [server/internal/model/sync.go](server/internal/model/sync.go)）。
3. `blog_server` 提供 JSON 接口：文章元信息来自 SQLite，文章正文在请求时从 `article` 表记录的 html **绝对路径**读取。列表页摘要由 [server/internal/common/preview.go](server/internal/common/preview.go) 从 html 提取纯文本（跳过标题、图片、表格、代码块等）。

改动文章内容后需重新运行 `article_converter` 才会生效。

## 构建与运行

所有 Go 命令在 `server/` 目录下执行（module 根在此）。sqlite 驱动为纯 Go 实现（modernc.org/sqlite），构建时 `CGO_ENABLED=0`，支持交叉编译与静态二进制。

```bash
# 发布打包（在仓库根）：交叉编译后端 + 构建前端 dist，产出 blog.tar.gz
make                        # 等价于 make blog，默认目标平台 linux/amd64
make blog GOOS=darwin GOARCH=arm64   # 本机平台可覆盖 GOOS/GOARCH
make clean                  # 清理 blog/ 与 blog.tar.gz

# 本机开发：直接运行后端服务（默认读 config/config.yaml，端口 6000）
cd server && go run ./cmd/blog_server -c config/config.yaml

# 本机开发：转换文章并写库（-src/-out/-db 三个参数均必填）
cd server && go run ./cmd/article_converter -src ../../article -out output/html -db db/blog.db

# 前端（在 web/ 下）
cd web && npm install
npm run dev                 # 开发服务器，接口代理到 localhost:6000
npm run build               # 生产构建到 web/dist
```

> 部署流程（打包、上传、解压、Nginx 配置、文章部署）见 [README.md](README.md)。
> Makefile 只有 `blog`（默认）和 `clean` 两个实际可用目标，`markdown_to_html`/`article_converter`/`blog_server`/`convert` 已不存在。

## 接口约定

- 统一响应结构 `{"code":0,"msg":"","data":{...}}`，用 [server/internal/common/response.go](server/internal/common/response.go) 的 `common.OK` / `common.Fail`。`code == 0` 为成功。
- 现有接口：
  - `GET /article/list?page=` — 文章分页列表（page 从 0 开始）
  - `GET /article/detail?title=` — 文章详情（含 html 正文）
  - `GET /category/overview` — 各分类及文章数（按文章数降序）
  - `GET /category/list?name=&page=` — 某分类下的文章分页
  - `GET /tag/overview` — 各标签及文章数（按文章数降序）
  - `GET /tag/list?name=&page=` — 某标签下的文章分页
- 分页每页固定 `service.PageSize`（10 条），列表按日期倒序。
- 前端 dev 环境靠 [web/vite.config.js](web/vite.config.js) 的 proxy 把接口转发到后端；**新增接口路径时必须同步在此登记**，否则会被前端 SPA 路由拦截。

## 约定与风格

- **代码注释与文档用中文**（现有代码、README 均为中文注释），保持一致。
- 后端分层清晰：controller 只做参数校验和响应，业务逻辑放 service，SQL / GORM 查询放 model；不要跨层。
- service 层查询列表时用 `enrichArticles` 统一填充分类名、标签、摘要等关联字段，新的文章列表接口应复用它。
- 「未找到」在 model 层返回 `(nil, nil)`，由上层转成对应的业务错误码，不要直接返回 gorm 的 ErrRecordNotFound。
- 前端每个路由对应一个 `views/*.vue`；列表页共用 `Pagination` 组件（数字分页，`page` 从 0 开始，`@change` 回传目标页）。
- 前端整体为浅色设计、已做移动端适配（viewport、流式布局、表格横向滚动、长内容断行）；改样式时注意别破坏窄屏表现。

## 验证改动

- 后端：`cd server && go build ./...` 必须通过。若涉及接口行为，起服务后用 `curl` 打对应路径核对返回（本机 6000 端口可能已被占用，可临时用别的端口起一个测试实例）。
- 前端：`cd web && npm run build` 必须通过。
