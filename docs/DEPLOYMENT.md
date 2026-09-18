# 🍵 UK Tea House — 小白级部署指南

> 你手里有一台 Ubuntu/Debian 服务器，你知道服务器的公网 IP，你知道 root 密码，你会用 SSH。
> 复制粘贴 → 等它跑完 → 复制粘贴 → 等它跑完 → 部署完成。
> **不需要懂 Docker，不需要懂 Postgres，不需要懂 Linux 命令。**

---

## 📋 开始前 — 你需要准备好这些

| # | 需要什么 | 去哪拿 | 准备好没？ |
|---|---------|--------|-----------|
| 1 | 服务器公网 IP | 阿里云/腾讯云/AWS 控制台，或者登录服务器后敲 `curl ifconfig.me` | ☐ |
| 2 | 服务器 root 密码 | 你买服务器时设置的 | ☐ |
| 3 | GitHub 账号密码 | https://github.com（没有就注册一个） | ☐ |
| 4 | GitHub PAT（个人访问令牌） | 见 Step 2，有图教你一步步点 | ☐ |
| 5 | 一个域名（可选） | 万网/GoDaddy/Cloudflare 随便买 | ☐ |

---

## 🔧 Step 1 — 登录服务器，装 Docker

### 1.1 SSH 登录（本地电脑上打开终端）

```bash
ssh root@你的服务器IP
# 例如: ssh root@124.174.33

# 第一次会问 Are you sure you want to continue? 敲 yes 回车
# 然后输入 root 密码，敲回车
# 登录成功后会看到 root@xxx:~# 这样的提示符
```

### 1.2 更新系统

```bash
apt update && apt upgrade -y
```

**✅ 成功标志**：最后一行是 `Reading package lists... Done`，然后没红色报错。

### 1.3 安装 Docker（30 秒搞定）

```bash
curl -fsSL https://get.docker.com | sh
```

**✅ 成功标志**：看到 `Docker Engine has been installed successfully` 绿色提示。

### 1.4 验证 Docker 能跑

```bash
docker run hello-world
```

**✅ 成功标志**：看到 `Hello from Docker! This message shows...`。

### 1.5 验证 Docker Compose

```bash
docker compose version
```

**✅ 成功标志**：看到 `Docker Compose version v2.x.x` 版本号。

### 1.6 如果上面报错 "docker: command not found" 或者 compose 不行

那就直接**跳到本节最后一行**执行这个：

```bash
echo "alias dc='docker compose'" >> ~/.bashrc && source ~/.bashrc
```

---

## 🔧 Step 2 — 获取 GitHub PAT（一次性，5 分钟）

PAT 就是 GitHub 给你的一个"密码"，让服务器有权限从 GitHub 拉预构建好的 Docker 镜像。

### 2.1 打开浏览器访问

👉 https://github.com/settings/tokens/new

### 2.2 按以下配置填写

| 填什么 | 填成 |
|--------|------|
| **Note**（名字） | `tea-deploy`（随便起，自己认得就行） |
| **Expiration**（有效期） | 选 `No expiration`（或者选 90 天，过期了再来生成一次） |
| **勾选权限** | 只勾一个：**`read:packages`** — 向下滚动找 `packages` 分类，勾 `read:packages` 那一行 |

### 2.3 点最下面绿色按钮 Generate token

**⚠️ 生成后立刻复制保存到记事本！它只显示一次，关掉就没了。**

### 2.4 拿到了 PAT 之后

继续在服务器上操作：

```bash
# 登录 GitHub Container Registry
docker login ghcr.io -u jinheqing

# 它会问你 Username 和 Password
# Username:  jinheqing
# Password:  把 PAT 粘贴进去（粘贴时终端不会显示任何字符，正常现象），回车
```

**✅ 成功标志**：看到 `Login Succeeded`。

---

## 🔧 Step 3 — 下载代码

### 3.1 创建目录

```bash
mkdir -p /opt/tea && cd /opt/tea
```

### 3.2 从 GitHub 下载全部代码

```bash
git clone https://github.com/jinheqing/shop.git .
```

**✅ 成功标志**：看到 `Cloning into '.'... done`。

### 3.3 查看目录里有啥（确认一下）

```bash
ls -la
```

**✅ 成功标志**：应该能看到 `docker-compose.yml`、`.env.example`、`tea-system`、`tea-translate`、`tea-frontend` 这些文件夹。

---

## 🔧 Step 4 — 一键生成密钥 + 写配置（一条命令搞定）

项目里自带了 `setup.sh` 脚本，**只需要跑一条命令**就能自动：
- 自动获取服务器公网 IP
- 生成全部 7 组安全密钥（密码/Token/JWT）
- 写好完整的 `.env` 配置文件
- 打印一张密钥清单让你保存

### 4.1 执行一键脚本

```bash
cd /opt/tea && bash setup.sh
```

### 4.2 执行完后会看到这样的输出

```
╔═══════════════════════════════════════════════════════════╗
║  ⚠️  重要! 下面这张清单请全部复制保存到你的记事本!          ║
║     只显示这一次, 关掉终端就没了!                           ║
╚═══════════════════════════════════════════════════════════╝

────────────────────────────────────────────────────────────
  服务器 IP           :  124.174.33.195
  管理后台            :  http://124.174.33.195/admin/

  管理员 Email        :  admin@ukteahouse.co.uk
  管理员密码          :  aB3xK9mQ2nF7

  数据库密码          :  cD4yL8pR1qE6
  审计库密码          :  gH5zN7sT3wJ4
  RabbitMQ 密码       :  kL2mP9vB8cX3
  JWT Secret          :  dGhpc2lzYXZlcnlsb25nc2VjcmV0...
  LiveKit API Key     :  3a7f9e5c1d2b4f6a8c0e
  LiveKit API Secret  :  eU3xQ9mR5kT7vB2nC4yP6zW8dF1h
────────────────────────────────────────────────────────────
```

**✅ 你需要做的：立刻把上面整张表里的内容全部复制到记事本！**

### 4.3 验证 .env 写对了

```bash
cat /opt/tea/.env | head -20
```

**✅ 成功标志**：看到 `REGISTRY=ghcr.io`、`DB_PASSWORD=XXXX...` 这些值已经填好了，没有 `【替换】` 字样。

---

## 🔧 Step 5 — 启动系统

所有配置已经就绪，直接启动就行。

### 5.1 预下载翻译模型（Whisper + NLLB-200，约 2GB）

这一步最耗时，等它跑完再继续。

```bash
cd /opt/tea
docker compose pull
```

**⏳ 预计耗时**：3-10 分钟，取决于网速。

**✅ 成功标志**：最后一行是 `Pull complete` 或者所有镜像都显示 `Downloaded newer image` / `Image is up to date`。

### 5.2 启动所有 10 个服务

```bash
cd /opt/tea
docker compose up -d --no-build
```

**⏳ 预计耗时**：2-3 分钟（Postgres 首次启动 + 翻译模型首次加载）。

**✅ 成功标志**：每个服务显示 `Started`，最后一行退出码 0。

### 5.3 盯着看哪些服务在"等模型"

```bash
docker compose ps
```

**正常状态说明**：
| 容器名 | 状态 | 说明 |
|--------|------|------|
| tea-postgres | Up (healthy) | ✅ 好 |
| tea-postgres-audit | Up (healthy) | ✅ 好 |
| tea-redis | Up (healthy) | ✅ 好 |
| tea-rabbitmq | Up (healthy) | ✅ 好 |
| tea-translate | Up | ⏳ 首次加载模型可能还没 healthy，**等 3-5 分钟**再看 |
| tea-mediamtx | Up | ✅ 好 |
| tea-livekit | Up | ✅ 好 |
| tea-gin | Up | ⏳ 等 translate，可能还没 healthy |
| tea-frontend | Up (healthy) | ✅ 好 |
| tea-nginx | Up | ✅ 好 |

### 5.4 等 translate 加载完成

```bash
# 每隔 15 秒刷新一次，直到 translate 和 gin 都变 healthy
watch -n 15 docker compose ps
# 按 Ctrl+C 退出 watch
```

**⏳ 等待时间**：首次大约 3-5 分钟。翻译模型从缓存加载 ~2GB 到内存。

**✅ 完成标志**：所有容器的 State 列都显示 `Up (healthy)`。

### 5.5 查看翻译加载日志（如果想确认）

```bash
docker compose logs translate --tail=30
```

**✅ 好的标志**：看到最后两行是 `whisper_loaded: true` 和 `nllb_loaded: true`。

---

## 🔧 Step 6 — 打开浏览器，用起来

### 6.1 先看你的服务器 IP

```bash
echo "你的服务器 IP 是: $(curl -s ifconfig.me)"
```

### 6.2 打开你的本地浏览器（Chrome/Safari/Edge），输入

**管理后台**:
```
http://你的服务器IP/admin/
# 例如 http://124.174.33.195/admin/
```

**公共站点**:
```
http://你的服务器IP/
# 例如 http://124.174.33.195/
```

### 6.3 登录管理后台

| 项 | 值 |
|----|-----|
| Email | `.env` 里的 `ADMIN_SEED_EMAIL`，默认是 `admin@ukteahouse.co.uk` |
| 密码 | **Step 4 保存的 ADMIN_SEED_PASSWORD** |

### 6.4 冒烟测试（验证真的跑起来了）

```bash
# 在服务器上敲这几条命令，看返回值对不对
curl http://localhost/health             # 期望: {"service":"tea-system"...}
curl http://localhost:8090/health        # 期望: {"status":"ok","whisper_loaded":true,...}
curl http://localhost/api/health          # 期望: {"service":"tea-system",...}
curl -I http://localhost/                 # 期望: 200 OK
curl -I http://localhost/admin/           # 期望: 200 OK
```

**✅ 全绿就是系统完全正常了。**

---

## 🔧 Step 7 — HTTPS（生产必做，大约 5 分钟）

### 前提：你有一个域名，并且 DNS 已经指向你的服务器 IP

| 域名 | 操作 |
|------|------|
| `@` | A 记录 → 你的服务器 IP |
| `www` | A 记录 → 你的服务器 IP |

### 7.1 改 .env 里的域名

```bash
# 把下面的 YOUR_DOMAIN 换成你的域名（不带 https://）
sed -i 's|APP_DOMAIN=http://.*|APP_DOMAIN=https://YOUR_DOMAIN|' /opt/tea/.env

# 确认一下
grep APP_DOMAIN /opt/tea/.env
```

### 7.2 装 certbot 申请免费证书

```bash
apt install -y certbot python3-certbot-nginx
certbot certonly --standalone -d 你的域名 -d www.你的域名
```

申请过程中：
- 它会问邮箱，填一个能收信的
- 它会问同意不同意 EULA → 敲 Y
- 它会问要不要强制 HTTPS → 敲 2（Redirect）

**✅ 成功标志**：看到 `Successfully received certificate`。

### 7.3 证书位置

```bash
# 确认证书在
ls /etc/letsencrypt/live/你的域名/fullchain.pem
ls /etc/letsencrypt/live/你的域名/privkey.pem
```

### 7.4 修改 nginx 配置启用 HTTPS

```bash
# 直接粘贴执行（把 YOUR_DOMAIN 换成你的域名）
DOMAIN="你的域名"

cat > /opt/tea/docker/nginx.conf << NGEOF
# HTTPS 配置
upstream tea_backend   { server gin:8080; }
upstream tea_translate { server translate:8090; }
upstream tea_frontend  { server frontend:80; }

# HTTP → HTTPS 跳转
server {
    listen 80 default_server;
    server_name _;
    return 301 https://\$host\$request_uri;
}

# HTTPS 主服务
server {
    listen 443 ssl http2 default_server;
    server_name _;

    ssl_certificate     /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         HIGH:!aNULL:!MD5;

    # / → 公共站点
    location / {
        proxy_pass http://tea_frontend;
        proxy_set_header Host              \$host;
        proxy_set_header X-Real-IP         \$remote_addr;
        proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    # /admin/ → 管理后台
    location /admin/ {
        proxy_pass http://tea_frontend;
        proxy_set_header Host              \$host;
        proxy_set_header X-Real-IP         \$remote_addr;
        proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    # /api/ → Go 后端
    location /api/ {
        proxy_pass http://tea_backend;
        proxy_set_header Host              \$host;
        proxy_set_header X-Real-IP         \$remote_addr;
        proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        client_max_body_size 20m;
    }

    # /webhooks/ → Go 后端回调
    location /webhooks/ {
        proxy_pass http://tea_backend;
        proxy_set_header Host              \$host;
        proxy_set_header X-Real-IP         \$remote_addr;
    }

    # /ws/ → IM WebSocket
    location /ws/ {
        proxy_pass http://tea_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade    \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host       \$host;
        proxy_read_timeout          3600s;
    }

    # /translate/ → 翻译服务
    location /translate/ {
        proxy_pass http://tea_translate;
        proxy_http_version 1.1;
        proxy_set_header Upgrade    \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host       \$host;
        proxy_read_timeout          3600s;
        client_max_body_size        10m;
    }
}
NGEOF
```

### 7.5 把证书挂进 nginx 容器

```bash
# 先看 certbot 证书实际在哪（certbot 可能用了软链接）
ls /etc/letsencrypt/live/${DOMAIN}/fullchain.pem
# 上面路径存在的话，直接用；不存在就手动指定实际路径

# 修改 docker-compose.yml 给 nginx 加证书挂载
# 这一条可能报 "File edited"，那是正常的
sed -i "s|#    volumes:|    volumes:\n      - /etc/letsencrypt/live/${DOMAIN}/fullchain.pem:/etc/nginx/ssl/cert.pem:ro\n      - /etc/letsencrypt/live/${DOMAIN}/privkey.pem:/etc/nginx/ssl/key.pem:ro\n#    volumes (original disabled)|" /opt/tea/docker-compose.yml 2>/dev/null || true
```

**更保险的做法**：手动编辑 docker-compose.yml 找到 nginx 服务的 volumes 部分，加两行：

```yaml
  nginx:
    volumes:
      - ./docker/nginx.conf:/etc/nginx/conf.d/default.conf:ro
      - /etc/letsencrypt/live/你的域名/fullchain.pem:/etc/nginx/ssl/cert.pem:ro
      - /etc/letsencrypt/live/你的域名/privkey.pem:/etc/nginx/ssl/key.pem:ro
```

### 7.6 重启 nginx 容器

```bash
cd /opt/tea
docker compose up -d --no-deps nginx

# 验证
curl -I https://你的域名/ --insecure
# 看到 200 OK 就是成功了
```

---

## 🛠️ 日常运维（最常用的 5 条命令）

### 全部停止

```bash
cd /opt/tea && docker compose down
```

### 全部启动

```bash
cd /opt/tea && docker compose up -d --no-build
```

### 查看状态

```bash
cd /opt/tea && docker compose ps
```

### 实时看日志（找报错就看这个）

```bash
# 全部日志
cd /opt/tea && docker compose logs -f

# 只看 Go 后端
docker compose logs -f gin

# 只看翻译服务
docker compose logs -f translate
```

### 更新（GitHub 有新代码时）

```bash
cd /opt/tea
git pull origin main
docker compose pull gin frontend translate
docker compose up -d gin frontend translate
```

### 备份数据库（每天定时）

```bash
# 立即备份一次
cd /opt/tea
docker compose exec -T postgres pg_dump -U tea_system tea_system > backup_$(date +%Y%m%d).sql

# 设置每天凌晨 2 点自动备份
(crontab -l 2>/dev/null; echo "0 2 * * * cd /opt/tea && docker compose exec -T postgres pg_dump -U tea_system tea_system > backups/tea_\$(date +\%Y\%m\%d).sql 2>&1 && find backups/ -name '*.sql' -mtime +7 -delete") | crontab -
mkdir -p backups
```

---

## 🆘 出了问题怎么办

### 服务起不来

```bash
# 看具体哪个挂了
docker compose ps
# 看这个服务的日志（把 SERVICE_NAME 换成 gin / translate / nginx 等）
docker compose logs SERVICE_NAME --tail=50
```

### 连不上数据库

```bash
# 测试 gin 能不能 ping 到 postgres
docker compose exec gin bash -c "wget -qO- http://postgres:5432 2>&1 || echo 'postgres 不通'"
# 或者直接看 .env 里的 DSN 对不对
grep DSN /opt/tea/.env
```

### 前端页面 502

```bash
# 502 = nginx 想代理的后端没起来
docker compose ps frontend gin
# 确保它们是 Up (healthy)
```

### 翻译服务一直 loading

```bash
# translate 首次加载 Whisper + NLLB 需要 3-5 分钟
# 检查模型文件在不在
docker compose exec translate ls -la /app/models/
# 检查 health
curl http://localhost:8090/health
# 期望: {"status":"ok","whisper_loaded":true,"nllb_loaded":true}
```

### 服务器重启后服务都在吗？

**都在的。** 所有服务都配了自动重启：
- Docker daemon 开机自启 → 自动拉起所有容器
- 容器崩溃 → Docker 自动重拉
- 数据库/Redis/翻译模型都在 named volume 里 → 重启不丢数据

### 清空重来（⚠️ 危险，删所有数据）

```bash
cd /opt/tea
docker compose down -v
# 然后重新 up
docker compose pull && docker compose up -d --no-build
```

---

## 📌 快速对照表

| 想做什么 | 命令 |
|----------|------|
| 查看状态 | `docker compose ps` |
| 实时日志 | `docker compose logs -f` |
| 停止 | `docker compose down` |
| 启动 | `docker compose up -d --no-build` |
| 重启一个 | `docker compose restart gin` |
| 更新代码 | `git pull && docker compose pull && docker compose up -d` |
| 备份数据库 | `docker compose exec -T postgres pg_dump -U tea_system tea_system > backup.sql` |
| 恢复数据库 | `cat backup.sql \| docker compose exec -T postgres psql -U tea_system tea_system` |

---

## 🎉 恭喜！

系统已经跑起来了。现在你可以：

1. 打开管理后台 → 登录 → 改管理员密码（Admin 页面 → 找自己那条记录编辑）
2. 打开 System Settings → 配置 SMTP 邮件服务器（魔法链接登录用）
3. 配置支付网关凭证（2Checkout 或 PayPal）
4. 买一个域名，按 Step 8 配 HTTPS

有任何问题，`docker compose logs 服务名 --tail=50` 一看日志就知道哪错了。

---

*文档版本: v2.0 (小白复刻版) | 2026-09-18*
