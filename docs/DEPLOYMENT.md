# 🍵 UK Tea House — 一键部署 & 完整配置文档

> 本文档覆盖 **裸服务器 → 生产可用** 的完整流程，所有命令均可直接复制粘贴执行。
> 技术栈: Ubuntu 22.04 LTS + Docker Compose v2 + GitHub Packages

---

## 目录

1. [服务器前置条件](#1-服务器前置条件)
2. [一键安装（复制粘贴即可）](#2-一键安装复制粘贴即可)
3. [配置 .env（密钥/密码/域名）](#3-配置-env)
4. [启动部署](#4-启动部署)
5. [服务对接配置](#5-服务对接配置)
6. [验证清单](#6-验证清单)
7. [HTTPS 配置（生产必做）](#7-https-配置生产必做)
8. [运维手册](#8-运维手册)
9. [故障排查](#9-故障排查)

---

## 1. 服务器前置条件

### 最低硬件要求

| 资源 | 要求 | 说明 |
|------|------|------|
| CPU | 4 核+ | 翻译服务 Whisper/NLLB 吃 CPU |
| 内存 | 8 GB+ | Postgres 2GB + Translate 3GB + 其他 2GB |
| 硬盘 | 40 GB+ | Docker 镜像 ~5GB + 翻译模型 ~3GB + 数据卷 |
| 带宽 | 10Mbps+ | 首次拉镜像 + 模型下载 |

### 操作系统

Ubuntu 22.04 LTS / Debian 12 / CentOS 9 / Rocky 9 / AlmaLinux 9

### 端口开放清单

| 端口 | 用途 | 是否对外 |
|------|------|----------|
| 80 | Nginx HTTP | ✅ 必须 |
| 443 | Nginx HTTPS | ✅ 生产必开 |
| 1935 | MediaMTX RTMP (OBS 推流) | ✅ 直播需要 |
| 7880 | LiveKit SFU (WebRTC 信令) | ✅ 直播需要 |
| 7882 | LiveKit SFU UDP (WebRTC 媒体) | ✅ 直播需要 |
| 15672 | RabbitMQ 管理界面 | ❌ 开发调试用，建议关闭 |
| 5432/5433 | PostgreSQL | ❌ 仅容器内部通信 |
| 6379 | Redis | ❌ 仅容器内部通信 |
| 5672 | RabbitMQ AMQP | ❌ 仅容器内部通信 |

### 域名 DNS 解析（如走 HTTPS）

```
@      A    你的服务器公网IP
www    A    你的服务器公网IP
live   A    你的服务器公网IP    # 直播专用子域（可选）
```

---

## 2. 一键安装（复制粘贴即可）

以下命令按顺序执行，一条不能少。

### Step 1 — 系统基础更新

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl wget git vim ca-certificates gnupg lsb-release
```

### Step 2 — 安装 Docker & Docker Compose v2

```bash
# 官方一键安装脚本
curl -fsSL https://get.docker.com | sudo sh

# 当前用户加入 docker 组（免 sudo 跑 docker）
sudo usermod -aG docker $USER
newgrp docker

# 验证安装
docker --version        # Docker version 27.x.x
docker compose version  # Docker Compose version v2.x.x
```

### Step 3 — 克隆代码

```bash
# 创建部署目录
sudo mkdir -p /opt/tea && sudo chown $USER:$USER /opt/tea
cd /opt/tea

# 克隆仓库（主分支）
git clone https://github.com/jinheqing/shop.git .

# 或用 ssh: git clone git@github.com:jinheqing/shop.git .

# 查看当前标签（可选，切到特定版本）
git tag -l 'v*'
git checkout v1.0.0   # 推荐用 tag 而非直接 main
```

### Step 4 — 获取 GitHub PAT（GHCR 拉镜像必需）

1. 打开 https://github.com/settings/tokens/new
2. Note 填 `tea-deploy`
3. Expiration 选 `No expiration` 或 90 天（定期轮换）
4. 勾选权限：**`read:packages`**（只读容器包权限）
5. 点 **Generate token**，复制保存（只显示一次！）

```bash
# 登录 GHCR（会提示输入 username 和 PAT）
docker login ghcr.io -u jinheqing
# Username: jinheqing
# Password: ghp_xxxxxxxxxxxxxxxxxxxx（PAT）
```

### Step 5 — 生成生产密钥

```bash
# 生成安全的 JWT_SECRET（64 字符）
openssl rand -base64 64 | tr -d '\n'

# 生成 LiveKit API Secret（32 字符）
openssl rand -base64 32 | tr -d '\n'

# 生成 LiveKit API Key（32 字符）
openssl rand -hex 16
```

把三条命令的输出都复制保存好，下一步填 .env 用。

---

## 3. 配置 .env

```bash
cd /opt/tea

# 复制模板（已有完整字段）
cp .env.example .env  # 如果不存在则手动创建：vim .env
```

### 完整 .env 内容（生产版，逐项替换值）

```bash
cat > /opt/tea/.env << 'ENVEOF'
# ════════════════════════════════════════════════════════════════
# UK Tea House — 生产环境变量
# ════════════════════════════════════════════════════════════════

# ─── Docker 镜像仓库 ───
REGISTRY=ghcr.io
OWNER=jinheqing
IMAGE_TAG=latest                      # 生产建议用具体 tag，如 v1.0.0

# ─── 业务数据库 ───
DB_USER=tea_system                    # PostgreSQL 业务库用户名
DB_PASSWORD=CHANGE_ME_DB_PWD          # ⚠️ 替换为强密码（30+ 字符）
DB_NAME=tea_system

# ─── 审计数据库（独立 PostgreSQL 实例）───
AUDIT_DB_USER=tea_audit
AUDIT_DB_PASSWORD=CHANGE_ME_AUDIT_PWD # ⚠️ 替换为强密码（与业务库不同！）
AUDIT_DB_NAME=tea_audit

# ─── Redis ───
REDIS_PASSWORD=                        # 留空即无密码，生产建议设密码

# ─── RabbitMQ ───
RABBITMQ_USER=tea_mq                  # 改掉默认 guest
RABBITMQ_PASS=CHANGE_ME_MQ_PWD        # ⚠️ 替换

# ─── JWT（服务端鉴权）───
# ⚠️ 必须 ≥32 字符，上一步 openssl rand -base64 64 的输出
JWT_SECRET=CHANGE_ME_WITH_OPENSSL_JWT_SECRET

# ─── LiveKit SFU（WebRTC）───
# 上一步 openssl 生成的 key 和 secret
LIVEKIT_API_KEY=CHANGE_ME_LK_KEY
LIVEKIT_API_SECRET=CHANGE_ME_LK_SECRET

# ─── 应用域名 & 品牌 ───
APP_DOMAIN=https://ukteahouse.co.uk   # ⚠️ 你的生产域名，带 https://
APP_SERVICE_NAME="UK Tea House"
APP_CONTACT_EMAIL=tea@ukteahouse.co.uk
SERVER_VERSION=1.0.0

# ─── 默认管理员（首次启动自动创建 staff 表第一条）───
ADMIN_SEED_EMAIL=admin@ukteahouse.co.uk
ADMIN_SEED_NAME="Admin"
ADMIN_SEED_PASSWORD=CHANGE_ME_ADMIN_PWD  # ⚠️ 生产必改，否则 Validate 会报错

# ─── 邮件（魔法链接登录、通知）───
# 用 Mailgun/SendGrid/Amazon SES 任意一个
SMTP_HOST=smtp.mailgun.org
MAIL_FROM="UK Tea House <tea@ukteahouse.co.uk>"
MAILGUN_API_KEY=key-xxxxxxxxxxxxxxxxxx  # Mailgun Dashboard 里拿

# ─── 支付 ───
PAYMENT_GATEWAY=2checkout               # 或 paypal / stripe
PAYMENT_BASE_URL=https://api.2checkout.com/rest/6.0
PAYMENT_API_KEY=                        # 2Checkout/PayPal Dashboard 里拿
PAYMENT_API_SECRET=                     # 同上

ENVEOF
```

### .env 字段速查

| 变量名 | 必填 | 默认值 | 说明 |
|--------|------|--------|------|
| `DB_PASSWORD` | ✅ | — | 业务库密码，30+ 字符 |
| `AUDIT_DB_PASSWORD` | ✅ | — | 审计库密码，必须不同 |
| `JWT_SECRET` | ✅ | — | ≥32 字符，openssl 生成 |
| `LIVEKIT_API_KEY` | ✅ | — | LiveKit 鉴权 key |
| `LIVEKIT_API_SECRET` | ✅ | — | LiveKit 鉴权 secret |
| `APP_DOMAIN` | ✅ | — | 生产域名，带协议头 |
| `ADMIN_SEED_PASSWORD` | ✅ | — | 管理员密码 |
| `SMTP_HOST` | 魔法链接时必填 | — | SMTP 服务器 |
| `MAILGUN_API_KEY` | 用 Mailgun 时必填 | — | Mailgun 接口 key |
| `PAYMENT_API_KEY` | 用支付时必填 | — | 支付网关 key |

### 验证 .env 加载正确

```bash
cd /opt/tea
docker compose config > /dev/null && echo "✅ .env + compose 解析通过" || echo "❌ 检查 .env 语法"
```

---

## 4. 启动部署

### 方式 A — 生产：直接从 GHCR 拉预构建镜像（推荐）

```bash
cd /opt/tea

# 一次性拉取所有镜像（3 个自有镜像 + 7 个基础镜像）
# 首次需要下载 ~6GB（Go 150MB + Frontend 80MB + Translate 1.5GB + 模型 ~2GB）
docker compose pull

# 翻译模型预下载（可选 — 提前下载 Whisper base + NLLB-200，~2GB）
docker compose run --rm translate python3 -c "
from faster_whisper import WhisperModel
WhisperModel('base', device='cpu', compute_type='int8')
print('✅ Whisper base 已下载')
"

# 启动全部服务（不 build，用拉下来的镜像）
docker compose up -d --no-build

# 查看启动进度
docker compose ps

# 等所有 healthy（约 3-5 分钟，翻译模型首次加载慢）
watch docker compose ps
```

### 方式 B — 开发：本地 build

```bash
# 需要 Node 20+ 和 Go 1.25+ 在宿主机
cd /opt/tea
./up.sh start              # build + 启动
./up.sh rebuild            # 改代码后只重建 3 个自有服务
```

### 一键部署脚本（方式 A 的全部步骤）

```bash
# 从 Step 3 clone 完之后，这一条命令搞定 Step 4-5-启动
cd /opt/tea

# 1. 生成密钥自动填入
JWT=$(openssl rand -base64 64 | tr -d '\n')
LK_KEY=$(openssl rand -hex 16)
LK_SECRET=$(openssl rand -base64 32 | tr -d '\n')
sed -i "s/CHANGE_ME_WITH_OPENSSL_JWT_SECRET/$JWT/" .env
sed -i "s/CHANGE_ME_LK_KEY/$LK_KEY/" .env
sed -i "s/CHANGE_ME_LK_SECRET/$LK_SECRET/" .env

# 2. 拉镜像 + 启动
docker compose pull && docker compose up -d --no-build

# 3. 等健康
echo "⏳ 等待服务就绪..."
sleep 30
docker compose ps
```

---

## 5. 服务对接配置

这一节列出每个系统组件之间的连接方式。所有对接 **compose 已自动配好**，这里是理解和调优参考。

### 5.1 服务间网络拓扑

```
┌─────────────────────────────────────────────────────────────────────┐
│  Docker 内部网络 tea-net (bridge)                                    │
│                                                                     │
│  ┌──────────┐  ┌──────────────┐  ┌──────────────┐  ┌───────────┐   │
│  │  gin     │──│   postgres   │  │postgres_audit│  │  redis    │   │
│  │  :8080   │──│   :5432      │  │  :5432       │  │  :6379    │   │
│  └────┬─────┘  └──────────────┘  └──────────────┘  └───────────┘   │
│       │                                                             │
│       │   ┌──────────────┐  ┌──────────┐  ┌─────────────────────┐ │
│       ├───│  rabbitmq    │  │ translate │  │     livekit        │ │
│       │   │  :5672       │  │  :8090    │  │     :7880          │ │
│       │   └──────────────┘  └─────┬────┘  └──────────┬──────────┘ │
│       │                           │                  │             │
│       │   ┌──────────────┐       │                  │             │
│       └───│  mediamtx    │       │                  │             │
│           │  :1935 RTMP  │       │                  │             │
│           └──────────────┘       │                  │             │
│                                                                  │
│  ┌─────────┐  proxy_pass   ┌──────────┐                          │
│  │  nginx  │───────────────│ frontend │                           │
│  │  :80    │  (SPA静态)     │  :80     │                           │
│  └─────────┘               └──────────┘                           │
└─────────────────────────────────────────────────────────────────────┘
                        │
                        │ 端口映射
                        ▼
              公网可访问的端口
              80, 443, 1935, 7880, 7882
```

### 5.2 各服务对接参数

#### 🗄️ PostgreSQL 业务库
- **容器名**: `tea-postgres`
- **对接方式**: Go 后端通过 DSN 连接（compose 自动组装）
- **DSN 格式**: `postgres://tea_system:密码@postgres:5432/tea_system?sslmode=disable`
- **自动迁移**: Go 后端启动时 GORM AutoMigrate，首次启动自动建表

#### 🗄️ PostgreSQL 审计库
- **容器名**: `tea-postgres-audit`
- **对接方式**: 同上，独立 DSN + 独立实例，两张库完全隔离

#### 🔴 Redis 7
- **容器名**: `tea-redis`
- **对接方式**: `redis:6379`，go-redis 客户端
- **用途**: 会话缓存、验证码、魔法链接 token、WebSocket 在线状态
- **持久化**: `appendonly yes` 开启 AOF

#### 🐰 RabbitMQ 3
- **容器名**: `tea-rabbitmq`
- **对接方式**: `amqp://user:pass@rabbitmq:5672/`
- **管理界面**: `http://服务器IP:15672`（guest/guest 或自定义）
- **用途**: 订单状态异步通知、发票生成队列

#### 🎥 LiveKit SFU
- **容器名**: `tea-livekit`
- **对接方式**: Go 后端通过 LiveKit Server SDK 做鉴权和 token 签发
- **关键**: `LIVEKIT_API_KEY` + `LIVEKIT_API_SECRET`（前后端共用，已在 .env）
- **配置文件**: `tea-system/docker/livekit.yaml`（已通过 volume 挂载进容器）
- **需要改的地方**（如边缘节点 / 公网 IP）:
  ```yaml
  # livekit.yaml 里
  # rtc_host: 0.0.0.0                # 改成服务器公网 IP
  # rtc_addresses:                    # 如果有 WireGuard 内网
  #   - "10.10.0.1"                   # WireGuard 网关
  #   - "YOUR_PUBLIC_IP"
  ```

#### 📼 MediaMTX (RTMP → WebRTC)
- **容器名**: `tea-mediamtx`
- **对接方式**: Go 后端调 `mediamtx:8889` HTTP API 做直播控制
- **OBS 推流地址**: `rtmp://服务器IP:1935/live/房间ID`
- **配置文件**: `tea-system/docker/mediamtx.yaml`（volume 挂载）

#### 🎙️ Tea-Translate (Whisper + NLLB-200)
- **容器名**: `tea-translate`
- **对接方式**: HTTP `http://translate:8090` + WebSocket `/translate/ws_asr`
- **模型**: Whisper `base` + NLLB-200，首次启动下载 ~2GB，缓存进 `translate_models` volume
- **首次启动**: 等 health endpoint 返回 `whisper_loaded: true, nllb_loaded: true` 才算就绪

#### 🖥️ Frontend (Vue3 SPA × 2)
- **容器名**: `tea-frontend`
- **内置两个 SPA**:
  - `/` → `public-site`（消费者站点）
  - `/admin/` → `admin-dashboard`（管理后台）
- **Vite base URL**: 前端生产构建默认 `/`，管理后台需配 Vite `base: '/admin/'`（已在 Dockerfile 构建时处理）

#### 🌐 Nginx 反代
- **容器名**: `tea-nginx`
- **统一入口**: `:80`（HTTP） + `:443`（HTTPS，需配证书）
- **路由规则**:
  | 路径 | 代理目标 | 说明 |
  |------|----------|------|
  | `/` | `frontend:80` | 公共站点 SPA |
  | `/admin/` | `frontend:80` (alias) | 管理后台 SPA |
  | `/api/*` | `gin:8080` | REST API |
  | `/webhooks/*` | `gin:8080` | 第三方回调 |
  | `/ws/*` | `gin:8080` | IM WebSocket |
  | `/translate/*` | `translate:8090` | ASR + 翻译 API/WS |

### 5.3 LiveKit Webhook（可选）

如果你想让直播事件自动回调 Go 后端（直播开始/结束/录制完成）:

```yaml
# tea-system/docker/livekit.yaml
webhook:
  url: "http://gin:8080/api/webhooks/livekit"
  api_key: "任意字符串"
```

Go 后端侧 handler 在 `handlers/sgs.go` 或类似位置处理。

---

## 6. 验证清单

部署完成后，按顺序验证：

### 6.1 容器健康检查

```bash
cd /opt/tea
docker compose ps

# 期望所有容器 State 都是 Up (healthy)
# translate 的 start_period=120s，首次加载模型需要等 3-5 分钟
# 如果一直 unhealthy，看日志：
docker compose logs translate --tail=50
docker compose logs gin --tail=50
```

### 6.2 HTTP 端点冒烟测试

```bash
# 健康检查（gin）
curl http://localhost:80/health    # 应返回 {service: "tea-system", ...}

# 健康检查（translate）
curl http://localhost:8090/health  # 应返回 {whisper_loaded: true, ...}

# 公网访问（Nginx）
curl -I http://localhost/               # 200, 返回前端 index.html
curl -I http://localhost/admin/         # 200, 返回管理后台 index.html
curl http://localhost/api/health        # 200

# RabbitMQ 管理界面
curl -I http://localhost:15672          # 401 或 200（未登录）
```

### 6.3 首次管理员登录

```
浏览器打开 → http://你的服务器IP/admin/
用户名    → ADMIN_SEED_EMAIL（.env 里配置的）
密码      → ADMIN_SEED_PASSWORD（.env 里配置的）
```

> ⚠️ 首次登录后请立即在 **System Settings → Administrators** 里修改密码！
> seedAdminStaff 逻辑: 如果 staff 表为空 → 自动创建一条记录。已有管理员后不会重复创建。

### 6.4 关键业务流验证

| 功能 | 测试方法 | 预期结果 |
|------|----------|----------|
| 注册 | 公共站点点 Sign Up | 收到魔法链接邮件 → 点击激活 |
| 登录 | 魔法链接登录 | 跳转首页，显示用户名 |
| IM WebSocket | 两个登录用户互发消息 | 实时收到 |
| 直播推流 | OBS 推到 `rtmp://IP:1935/live/test123` | 直播开始 |
| 直播观看 | 另一浏览器访问直播页 | 看到画面 |
| 下单 | 管理后台手动创建订单 | 订单状态机流转 |
| 翻译 | 在聊天界面发语音消息 | 返回翻译文本 |

### 6.5 自动重启验证

```bash
# 模拟容器崩溃
docker restart tea-gin

# 观察 → 它会自动重启（restart: unless-stopped）
docker compose ps

# 模拟服务器重启（如果你有权限）
sudo reboot
# 重启后: systemd 启动 Docker → Docker 自动拉起所有容器
# 因为所有服务都是 restart: unless-stopped
# 数据卷（pg_data, redis_data, translate_models）持久化，不丢失
```

---

## 7. HTTPS 配置（生产必做）

### 方式 A — Let's Encrypt 免费证书（推荐）

```bash
# 装 certbot
sudo apt install -y certbot python3-certbot-nginx

# 申请证书（确保 DNS 已指向服务器 IP）
sudo certbot certonly --standalone \
  -d ukteahouse.co.uk \
  -d www.ukteahouse.co.uk

# 证书路径:
# /etc/letsencrypt/live/ukteahouse.co.uk/fullchain.pem
# /etc/letsencrypt/live/ukteahouse.co.uk/privkey.pem
```

然后修改 `docker-compose.yml` 里的 nginx 服务加证书挂载：

```yaml
  nginx:
    volumes:
      - ./docker/nginx.conf:/etc/nginx/conf.d/default.conf:ro
      - /etc/letsencrypt/live/ukteahouse.co.uk/fullchain.pem:/etc/nginx/ssl/cert.pem:ro
      - /etc/letsencrypt/live/ukteahouse.co.uk/privkey.pem:/etc/nginx/ssl/key.pem:ro
```

同时在 `docker/nginx.conf` 里加 443 server block：

```nginx
server {
    listen 443 ssl http2 default_server;
    server_name _;

    ssl_certificate     /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         HIGH:!aNULL:!MD5;

    # 这里复制 80 port server block 的全部 location
    # ... /, /admin/, /api/, /ws/, /translate/ ...
}

# HTTP → HTTPS 自动跳转
server {
    listen 80 default_server;
    server_name _;
    return 301 https://$host$request_uri;
}
```

然后重建 nginx 容器：

```bash
docker compose up -d --no-deps nginx
```

### 方式 B — Cloudflare / 阿里云 CDN 做 HTTPS 终止

如果用 CDN 做 HTTPS，服务器只开 80 端口就行。CDN 自动处理证书续期。

---

## 8. 运维手册

### 8.1 日常命令速查

```bash
cd /opt/tea

# 全部停止
docker compose down

# 查看状态
docker compose ps

# 实时日志（指定服务）
docker compose logs -f --tail=100 gin
docker compose logs -f --tail=100 translate
docker compose logs -f nginx

# 进入容器调试
docker compose exec gin bash
docker compose exec translate bash

# 重启单个服务（自动拉取最新镜像）
docker compose pull gin && docker compose up -d --no-deps gin

# 清理无用镜像 / 容器 / 卷
docker system prune -a        # ⚠️ 删所有未用的镜像（包括 translate 模型缓存会丢！）
docker volume ls              # 查看数据卷
```

### 8.2 数据库备份 & 恢复

```bash
# 业务库备份
docker compose exec -T postgres \
  pg_dump -U tea_system tea_system > backup_tea_$(date +%Y%m%d_%H%M%S).sql

# 审计库备份
docker compose exec -T postgres_audit \
  pg_dump -U tea_audit tea_audit > backup_audit_$(date +%Y%m%d_%H%M%S).sql

# 备份目录归档
mkdir -p /opt/tea/backups
mv backup_tea_*.sql backup_audit_*.sql /opt/tea/backups/

# 恢复（⚠️ 会清空现有数据！）
cat backup_tea_20260918_120000.sql | \
  docker compose exec -T postgres psql -U tea_system tea_system
```

### 8.3 定时自动备份（crontab）

```bash
# 编辑 root crontab
sudo crontab -e

# 每天凌晨 2 点自动备份两个数据库，保留 7 天
0 2 * * *  cd /opt/tea && \
  docker compose exec -T postgres pg_dump -U tea_system tea_system > backups/tea_$(date +\%Y\%m\%d).sql && \
  docker compose exec -T postgres_audit pg_dump -U tea_audit tea_audit > backups/audit_$(date +\%Y\%m\%d).sql && \
  find backups/ -name '*.sql' -mtime +7 -delete

# 确保 backups 目录存在且权限正确
mkdir -p /opt/tea/backups
sudo chown root:root /opt/tea/backups
```

### 8.4 升级部署

```bash
cd /opt/tea

# 方式 1 — 拉最新镜像（推荐，CI 已 build 好）
git pull origin main                   # 更新 compose / .env 模板 / nginx.conf
docker compose pull gin frontend translate
docker compose up -d gin frontend translate

# 方式 2 — 切特定版本
git checkout v1.1.0                    # 切 tag
sed -i 's/IMAGE_TAG=latest/IMAGE_TAG=v1.1.0/' .env
docker compose pull && docker compose up -d --no-build

# 方式 3 — 回滚
git log --oneline -5                   # 看历史
git checkout v1.0.0                    # 切回上一个稳定版
# 镜像如果已被覆盖，可能需要手动 pull 或本地 build
docker compose up -d --no-build
```

### 8.5 翻译服务模型更新

```bash
# 模型缓存在 translate_models volume，不会自动更新
# 手动更新 Whisper 模型：
docker compose run --rm translate python3 -c "
from faster_whisper import WhisperModel
WhisperModel('small', device='cpu', compute_type='int8')  # 换成更大的模型
"
# 然后重启 translate 容器：
docker compose restart translate
```

### 8.6 清理数据（⚠️ 危险）

```bash
# 停服务 + 删所有数据卷（重新部署）
docker compose down -v
# 这会删除：PostgreSQL 数据、Redis、RabbitMQ、翻译模型缓存
# 镜像不会删（除非加 -rmi all）

# 但代码 / .env / docker-compose.yml 都保留，重新 up 即可
docker compose up -d --no-build
```

---

## 9. 故障排查

### 服务启动失败

```bash
# 看全局日志
docker compose ps                    # 哪个 unhealthy
docker compose logs gin --tail=50    # 看具体服务报错
docker compose logs postgres
docker compose logs translate

# 常见原因：
# 1. 端口被占用
sudo lsof -i :5432                   # 或 netstat -tlnp
#    → 停掉冲突进程，或改 compose 的端口映射

# 2. 数据库首次启动太久（健康检查超时）
#    postgres healthcheck interval=10s retries=10 → 给够时间

# 3. 翻译服务卡在模型下载
docker compose exec translate ls /app/models/
#    → 看是否有模型文件，或手动触发预下载（见 Step 4）
```

### Go 后端连不上数据库

```bash
# 检查 DSN 解析
docker compose config | grep DB_BUSINESS_DSN
# 应该看到完整的 postgres://...@postgres:5432/tea_system DSN

# 测试容器内联通性
docker compose exec gin bash -c "apt update && apt install -y postgresql-client && psql -h postgres -U tea_system -d tea_system -c 'SELECT 1'"
# 输入密码（.env 里的 DB_PASSWORD），能返回 1 就通
```

### Nginx 代理 502 Bad Gateway

```bash
# 原因：被代理的后端容器没起来或端口不对
docker compose ps gin frontend
# 应该是 Up (healthy)

# 检查 nginx upstream 配置
docker compose exec nginx nginx -t    # 语法检查
docker compose logs nginx --tail=20
```

### 前端页面空白 / 404

```bash
# 检查 frontend 容器静态文件是否在正确位置
docker compose exec frontend ls /usr/share/nginx/html/
# 应该看到 index.html, assets/

docker compose exec frontend ls /usr/share/nginx/html/admin/
# 应该也看到 index.html

# SPA try_files 配置
docker compose exec frontend cat /etc/nginx/conf.d/default.conf
# 应该有 try_files $uri $uri/ /index.html;
```

### LiveKit 推流 / 观看失败

```bash
# OBS 推流测试
rtmp://服务器公网IP:1935/live/任意房间ID

# LiveKit 端口开放检查（在服务器上）
sudo ufw allow 7880/tcp
sudo ufw allow 7882/tcp
sudo ufw allow 7882/udp
sudo ufw allow 1935/tcp

# LiveKit 日志
docker compose logs livekit --tail=30

# 检查 livekit.yaml 的 rtc_addresses 是否配置了公网 IP
```

### HTTPS 证书续期失败

```bash
# certbot 自动续期命令（应该已在 crontab/systemd）
sudo certbot renew --dry-run

# 续期成功后需要 reload nginx
docker compose exec nginx nginx -s reload
```

### 磁盘空间不足

```bash
# 查看各组件占用
docker system df                      # Docker 总占用
docker system df -v                   # 详细到每个容器/卷/镜像

# 重点：翻译模型缓存
docker volume inspect tea_translate_models
du -sh $(docker volume inspect --format '{{.Mountpoint}}' tea_translate_models)

# 清理未用镜像
docker image prune -a --filter "until=24h"   # 删除 24h 前的镜像
```

---

## 附录 A — 自动重启机制详解

所有 10 个服务都设了 `restart: unless-stopped`：

| 场景 | 行为 |
|------|------|
| 容器内进程崩溃 | Docker 自动重启容器（不丢失数据，只重启进程） |
| `docker compose restart gin` | 重新启动 gin 容器 |
| 服务器物理重启 | Docker daemon 启动 → 自动拉起所有 `restart: unless-stopped` 的容器 |
| `docker compose down` | 停止所有容器，**不会**因服务器重启自动起来（除非 up） |
| `docker compose down -v` | 停 + 删数据卷（⚠️ 数据不可恢复） |
| `docker rm -f tea-gin` | 删除容器，compose 会认为服务不存在 |

**数据持久化**: 所有关键数据存在 Docker named volumes 中，不随容器删除而丢失：
- `pg_data` — PostgreSQL 业务库
- `pg_audit_data` — PostgreSQL 审计库
- `redis_data` — Redis AOF
- `rabbitmq_data` — RabbitMQ 消息队列
- `translate_models` — Whisper + NLLB-200 模型缓存（~2GB）

## 附录 B — 服务器一键重置脚本（谨慎使用）

```bash
#!/bin/bash
# reset.sh — 从裸服务器开始的最短路径
set -e
cd /opt/tea

echo "🧹 清空所有容器 + 数据卷"
docker compose down -v --rmi local

echo "📦 重新拉取"
docker compose pull

echo "🚀 重新启动"
docker compose up -d --no-build

echo "⏳ 等待健康检查..."
sleep 60
docker compose ps

echo "✅ 完成！"
```

## 附录 C — GitHub Actions 流程说明

推到 main 分支 → GitHub Actions 自动：
1. Build 三个自有服务（gin/frontend/translate）的 Docker image
2. Push 到 `ghcr.io/jinheqing/shop-{gin,frontend,translate}:latest`
3. 同时打 tag `${{ github.sha }}`（可回溯的 commit hash tag）

打 tag `v1.0.0` → 额外打 `v1.0.0` tag，生产环境推荐用固定 tag。

服务器侧只需 `docker compose pull` 就能拿到最新镜像。

---

*文档版本: v1.0 | 最后更新: 2026-09-18*
*项目地址: https://github.com/jinheqing/shop*
