# 概述

本仓库为个人博客网站的前后端代码仓库

# 博客部署

```bash
# 1. 将前后端编译并打包至 blog.tar.gz
make

# 2. 将 blog.tar.gz 传到服务器上

# 3. 将 blog.tar.gz 解压
tar --warning=no-unknown-keyword -zxf blog.tar.gz

# 4. 启动服务器
nohup ./blog_server -c config/config.yaml > blog_server.log 2>&1 &
```

Nginx 配置放在 /etc/nginx/conf.d/blog.conf

```
server {
    listen 80;
    server_name huanglianjing.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name huanglianjing.com;

    # Let's Encrypt 证书路径
    ssl_certificate     /etc/nginx/huanglianjing.com_bundle.crt;
    ssl_certificate_key /etc/nginx/huanglianjing.com.key;

    # SSL 安全参数（推荐）
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache   shared:SSL:10m;
    ssl_session_timeout 1d;

    # HSTS：强制浏览器后续只用 HTTPS（确认站点稳定跑 HTTPS 后再开启）
    add_header Strict-Transport-Security "max-age=31536000" always;

    root /root/blog/dist;
    index index.html;

    # 后端 API：与 vite.config.js 里代理的路径保持一致
    location ~ ^/(article/list|article/detail|category/overview|category/list|tag/overview|tag/list) {
        proxy_pass http://127.0.0.1:6000;
        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 前端 SPA：其余路径回退到 index.html
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

初次部署需要重新加载 Nginx 配置

```bash
nginx -t
nginx -s reload
```

# 文章部署

```bash
# 1. 将文章打包至 article.tar.gz
rm -f article.tar.gz
tar --exclude=article/.git -zcf article.tar.gz article

# 2. 将 article.tar.gz 传到服务器上

# 3. 将 article.tar.gz 解压
tar --warning=no-unknown-keyword -zxf article.tar.gz

# 4. 将 markdown 转为 html，同时写入数据库
cd blog
./article_converter -src ../article -db db/blog.db -out article_html
```
