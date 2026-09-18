# ==========================================================================
# 🍵 UK Tea House — 一键启动/停止脚本
#
#   🔧 开发模式 (本地 build):
#     ./up.sh start          全部 build + 启动
#     ./up.sh rebuild        只 rebuild 三个自有服务 (gin/frontend/translate)
#
#   🚀 生产模式 (从 GHCR 拉预构建镜像):
#     ./up.sh deploy         pull 镜像 + 启动（不 build）
#     ./up.sh pull           只拉取最新镜像
#
#   其他:
#     ./up.sh stop           全部停止
#     ./up.sh logs           实时日志（gin + translate）
#     ./up.sh reset          清空所有 volume（⚠️ 会删数据库！）
#     ./up.sh status         服务状态
# ==========================================================================
set -euo pipefail
cd "$(dirname "$0")"

ACTION=${1:-start}

# 颜色
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'

case "$ACTION" in
  start)
    echo "🔧  开发模式: docker compose up -d --build"
    docker compose up -d --build
    echo ""
    echo "✅ 启动完成！等待健康检查..."
    sleep 5
    docker compose ps
    echo ""
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "  🌐  公共站点:      http://localhost/"
    echo -e "  🔧  管理后台:      http://localhost/admin/"
    echo -e "  🎙️  Translate:     http://localhost:8090/health"
    echo -e "  🐰  RabbitMQ:      http://localhost:15672  (guest/guest)"
    echo -e "  📼  OBS RTMP:      rtmp://localhost:1935/live/{roomId}"
    echo -e "  🎥  LiveKit SFU:   ws://localhost:7880"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    ;;

  deploy)
    echo "🚀  生产模式: pull GHCR 镜像 + 启动"
    echo ""
    echo "📦  登录 GHCR (如未登录)..."
    # 尝试 pull 触发认证提示，已登录会跳过
    docker compose pull || {
      echo -e "${YELLOW}⚠️  需要先登录 GHCR:${NC}"
      echo "     docker login ghcr.io -u <your-github-user> -p \$GITHUB_PAT"
      echo "     然后再运行: ./up.sh deploy"
      exit 1
    }
    echo ""
    echo "▶️  启动服务 (--no-build)..."
    docker compose up -d --no-build
    echo ""
    echo "✅ 部署完成！"
    sleep 5
    docker compose ps
    ;;

  pull)
    echo "📦  拉取最新镜像..."
    docker compose pull
    ;;

  stop)
    echo "🛑  docker compose down"
    docker compose down
    ;;

  logs)
    docker compose logs -f --tail=100 gin translate frontend
    ;;

  rebuild)
    echo "🔨  只 rebuild 三个自有服务 (gin/frontend/translate)"
    docker compose build --no-cache gin frontend translate
    docker compose up -d gin frontend translate
    ;;

  reset)
    echo -e "${RED}⚠️  这会删除所有数据 (postgres, redis, rabbitmq, translate models)${NC}"
    read -p "确认? [y/N] " yes
    if [ "$yes" = "y" ]; then
      docker compose down -v
      echo "✅ 已清空"
    fi
    ;;

  status)
    docker compose ps
    ;;

  *)
    echo "用法: $0 {start|deploy|pull|stop|logs|rebuild|reset|status}"
    echo ""
    echo "  start    — 开发模式: 本地 build + 启动"
    echo "  deploy   — 生产模式: 从 GHCR pull 镜像 + 启动"
    echo "  pull     — 只拉取最新镜像"
    echo "  stop     — 全部停止"
    echo "  logs     — 实时日志"
    echo "  rebuild  — rebuild 三个自有服务"
    echo "  reset    — 清空所有数据卷"
    echo "  status   — 服务状态"
    ;;
esac
