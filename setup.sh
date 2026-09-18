#!/usr/bin/env bash
# =================================================================
# 🍵 UK Tea House — 一键配置脚本
# 作用: 自动生成所有密钥 + 写好 .env + 等待替换服务器 IP
# 用法: bash setup.sh
# =================================================================
set -e

cd "$(dirname "$0")"

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║  🍵 UK Tea House — 一键配置                                 ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""

# ── 自动获取服务器 IP ──
SERVER_IP=$(curl -s ifconfig.me || echo "REPLACE_WITH_YOUR_IP")
echo "🌐 检测到服务器 IP: ${SERVER_IP}"
echo ""

# ── 生成全部密钥 ──
echo "🔑 生成安全密钥..."
JWT_SECRET=$(openssl rand -base64 64 | tr -d '\n')
LK_KEY=$(openssl rand -hex 16)
LK_SECRET=$(openssl rand -base64 32 | tr -d '\n')
DB_PWD=$(openssl rand -base64 24 | tr -d '\n')
AUDIT_PWD=$(openssl rand -base64 24 | tr -d '\n')
MQ_PWD=$(openssl rand -base64 16 | tr -d '\n')
ADMIN_PWD=$(openssl rand -base64 16 | tr -d '\n')
echo "✅ 密钥生成完成"
echo ""

# ── 写 .env ──
echo "📝 写入 .env 文件..."
cat > .env << EOF
# ════════════════════════════════════════════════════════════
# UK Tea House — 自动生成的配置 (setup.sh at $(date))
# ════════════════════════════════════════════════════════════

# 镜像仓库
REGISTRY=ghcr.io
OWNER=jinheqing
IMAGE_TAG=latest

# 业务数据库
DB_USER=tea_system
DB_PASSWORD=${DB_PWD}
DB_NAME=tea_system

# 审计数据库
AUDIT_DB_USER=tea_audit
AUDIT_DB_PASSWORD=${AUDIT_PWD}
AUDIT_DB_NAME=tea_audit

# RabbitMQ
RABBITMQ_USER=tea_mq
RABBITMQ_PASS=${MQ_PWD}

# JWT
JWT_SECRET=${JWT_SECRET}

# LiveKit
LIVEKIT_API_KEY=${LK_KEY}
LIVEKIT_API_SECRET=${LK_SECRET}

# 应用域名
APP_DOMAIN=http://${SERVER_IP}
APP_SERVICE_NAME="UK Tea House"
APP_CONTACT_EMAIL=tea@ukteahouse.co.uk
SERVER_VERSION=1.0.0

# 默认管理员 (首次启动自动创建)
ADMIN_SEED_EMAIL=admin@ukteahouse.co.uk
ADMIN_SEED_NAME="Admin"
ADMIN_SEED_PASSWORD=${ADMIN_PWD}

# 邮件 (暂时留空, 管理后台 System Settings 里可以配)
SMTP_HOST=
MAIL_FROM="UK Tea House <tea@ukteahouse.co.uk>"
MAILGUN_API_KEY=

# 支付 (暂时留空)
PAYMENT_GATEWAY=2checkout
PAYMENT_BASE_URL=
PAYMENT_API_KEY=
PAYMENT_API_SECRET=
EOF

echo "✅ .env 写入完成: $(pwd)/.env"
echo ""

# ── 打印全部密钥 (必须保存!) ──
echo "╔═══════════════════════════════════════════════════════════╗"
echo "║  ⚠️  重要! 下面这张清单请全部复制保存到你的记事本!          ║"
echo "║     只显示这一次, 关掉终端就没了!                           ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""
echo "────────────────────────────────────────────────────────────"
echo "  服务器 IP           :  ${SERVER_IP}"
echo "  管理后台            :  http://${SERVER_IP}/admin/"
echo ""
echo "  管理员 Email        :  admin@ukteahouse.co.uk"
echo "  管理员密码          :  ${ADMIN_PWD}"
echo ""
echo "  数据库密码          :  ${DB_PWD}"
echo "  审计库密码          :  ${AUDIT_PWD}"
echo "  RabbitMQ 密码       :  ${MQ_PWD}"
echo "  JWT Secret          :  ${JWT_SECRET}"
echo "  LiveKit API Key     :  ${LK_KEY}"
echo "  LiveKit API Secret  :  ${LK_SECRET}"
echo "────────────────────────────────────────────────────────────"
echo ""
echo "下一步: 运行 'docker compose pull && docker compose up -d --no-build'"
