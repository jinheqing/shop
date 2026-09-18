# ==========================================================================
# tea-frontend — 多阶段构建
# Stage 1: pnpm install + vite build (public-site + admin-dashboard)
# Stage 2: nginx:alpine 托管两个 SPA 的 dist/
# ==========================================================================
FROM node:20-alpine AS builder
RUN apk add --no-cache git python3 make g++ \
    && npm install -g pnpm@9
WORKDIR /workspace
# Layer cache: 先拷依赖文件
COPY pnpm-workspace.yaml package.json pnpm-lock.yaml* ./
COPY packages/public-site/package.json packages/public-site/
COPY packages/admin-dashboard/package.json packages/admin-dashboard/
RUN pnpm install --frozen-lockfile --reporter=append-only 2>&1 | tail -5
# 源码
COPY . .
# 构建
RUN pnpm -r build 2>&1 | tail -15

# -------------------------- Stage 2: nginx --------------------------
FROM nginx:1.27-alpine
# public-site → /
COPY --from=builder /workspace/packages/public-site/dist/ /usr/share/nginx/html/
# admin-dashboard → /admin/
RUN mkdir -p /usr/share/nginx/html/admin
COPY --from=builder /workspace/packages/admin-dashboard/dist/ /usr/share/nginx/html/admin/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
