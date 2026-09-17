# 🔧 普洱茶私域系统 — 功能齐全性验收 + 修复追踪

> 审计时间：2026-09-17（第二轮：前端页面完整性 + 配置项合理性）
> 修复状态 legend：`[ ]` 待修 · `[🔄]` 进行中 · `[x]` 已修 · `[n/a]` 合理降级

---

## 🔴 关键发现摘要（第二轮）

### A. Admin Dashboard 前端空壳（9/31 菜单项是硬编码 demo 数据）

| # | 页面 | 现状 | 后端路由是否存在 | 严重程度 |
|---|------|------|------------------|----------|
| 1 | SiteContent.vue | `ref([...demo])` 硬编码 5 条 JSON | ✅ `GET/PUT /site-contents` 真实现成 | **P0 完全脱钩** |
| 2 | DSAR.vue | list 硬编码 4 条 demo；**仅** export/erase 绑了 API | ✅ 四个端点全部真实现 | **P0 列表假数据** |
| 3 | LiveCalendar.vue | `ref([...demo])` 硬编码 | ✅ `GET /live-rooms/schedule/calendar` | **P0** |
| 4 | CustomerRequests.vue | `ref([...demo])` 硬编码 | ✅ `GET /live-rooms/customer-requests` | **P0** |
| 5 | DeliveryInspection.vue | `ref([...demo])` 硬编码 | ✅ `POST /live-rooms/system-create` | **P1** |
| 6 | QRCodes.vue | `ref([...demo])` 硬编码 | ✅ `POST /qrcodes/generate` | **P1** |
| 7 | Staff.vue | `ref([...demo])` 硬编码 + `POST /staff` **后端不存在** | ❌ 无 staff CRUD 路由 | **P1 缺后端** |
| 8 | TranslateStatus.vue | 100% 硬编码（FastAPI 外部依赖合理 mock） | ✅ `GET /translate/status` | P3 合理 |
| 9 | WebhookLogs.vue | `ref([...demo])` 硬编码 | ❌ **后端根本没有 webhook_logs model/API** | **P1 需新增** |

### B. Admin Dashboard 配置页面"假保存"

| # | 页面 | 问题 | 后端路由 | 严重程度 |
|---|------|------|----------|----------|
| 10 | PaymentConfig.vue | `api.post('/system/payment-config', cfg)` — **后端不存在** | ❌ 无任何 system config 路由 | **P0 假保存** |

### C. Public Site 空壳（6/21 页面完全硬编码）

| # | 页面 | 现状 | 后端路由是否存在 |
|---|------|------|------------------|
| 11 | TeaGardens.vue | 6 条茶园 demo 数据，完全不调 API | ✅ `GET /public/slow-presets` 现成 |
| 12 | Quality.vue | 全硬编码 SGS 报告 demo | ✅ `GET /public/sgs-reports` 现成 |
| 13 | Contact.vue | 表单 submit 只弹 ElMessage，**不调任何 API** | ❌ 后端无 contact_form |
| 14 | CookieConsent.vue | 只存 localStorage，**不 POST /api/v1/cookie-consent** | ✅ `POST /cookie-consent` 现成 |
| 15 | Live.vue | API 绑了但模板 fallback 有硬编码 | ✅ 部分真调 |
| 16 | Home.vue | API 绑了但 `ref` 初始值有硬编码 fallback | ✅ 部分真调 |

### D. Admin Dashboard 完全缺失的配置页面

| 缺失模块 | 说明 |
|----------|------|
| **系统设置页（Settings / System Config）** | PaymentConfig 单独在菜单里是一个页，但没有更通用的系统配置——**没有**：GDPR SLA 天数配置、邮件 SMTP 配置、LiveKit Webhook URL 配置、站点品牌/logo/favicon 配置、运营公告配置 |
| **配送/运费规则页** | 无 `shipping_zones` / `shipping_rates` model 或页面 |
| **退款管理页** | Declarations 有 void，但无 `refunds` model / 页面 |
| **Marketing / 优惠券** | 无 `coupons` / `marketing_campaigns` model / 页面 |

### E. 已有配置页面的设计/数据完整性问题

| 页面 | 问题 |
|------|------|
| PaymentConfig.vue | 1) `cfg.twocheckout_seller_id` vs `cfg.twocheckout_account_id` — 字段命名不一致；2) 没有 3DS / Strong Customer Authentication 配置开关；3) 没有 sandbox/live URL 切换；4) **"Save All"按钮调的路由根本不存在** |
| SiteContent.vue | 1) 5 个 demo section 覆盖不全（缺 Live Calendar section、Checkout hero、FAQ 条目）；2) 编辑弹窗里 `draft.content` 是 JSON 字符串文本框，但 `el-input` 绑定用 `.content` 不对；3) **Save 按钮不调后端** |
| CustomProducts.vue | 比较完整（215 行），CRUD + publish/review + 28 字段快照 — 但缺"生成 QR Code"联动按钮 |
| LiveRooms.vue | 1) presets 里 `app_scheduled`、`slow_247`、`customer_requested`、`custom_private` 与后端常量不一致（后端是 `open_calendar`、`slow_preset`、`customer_request`、`custom_private`）；2) **只有 Delete，缺 Start/End/Status 切换按钮** |
| Orders.vue | 1) Transition 下拉里 `ready_for_delivery` 与后端状态机不一致（后端是 `ready_for_production`）；2) **缺 Timeline 查看**（后端有 `GET /orders/:id/timeline`）；3) 缺 Declarations 快速跳转 |
| SlowPresets.vue | 1) 表单里有 `camera_rtmp_url` 字段但后端会自动生成，不该让用户填；2) 缺 Edit 按钮 |
| Staff.vue | 1) 全 demo 数据；2) **后端没有 Staff CRUD 路由**——需要新增；3) "Create & Send Invite"按钮但 MailService 是 mock |
| LiveCalendar.vue | 纯 demo 硬编码，不调 `GET /live-rooms/schedule/calendar` |

---

## 📊 完整验收矩阵：后端路由 vs Admin 菜单 vs 页面 API 绑定

### ✅ 后端路由存在 + 菜单有 + 页面真调后端（已打通）

| 模块 | 后端路由 | Admin 菜单 | 页面 API |
|------|----------|------------|----------|
| 订单 | CRUD + state + cancel + timeline | ✅ | ✅ Orders.vue / OrderDetail.vue |
| 发票 | GET /orders/:id/invoice | ✅ Invoices | ⚠️ Invoices.vue 未绑定（62行但没 api.get） |
| 报关 | CRUD + void | ✅ | ✅ Declarations.vue |
| 外汇台账 | CRUD | ✅ | ✅ Ledgers.vue |
| SGS 报告 | CRUD | ✅ | ✅ SGS.vue |
| 直播间 | CRUD + start + end + calendar | ✅ | ⚠️ LiveRooms.vue 有 API 但缺 Start/End 按钮 |
| 慢直播 | CRUD + enable/disable | ✅ | ✅ SlowPresets.vue |
| 节点 | CRUD + health | ✅ | ✅ Nodes.vue |
| IM 会话 | CRUD + messages | ✅ | ✅ Conversations.vue |
| 定制产品 | CRUD + publish + review | ✅ | ✅ CustomProducts.vue |
| 审计日志 | GET /staff/audit-logs | ✅ | ✅ AuditLogs.vue |
| Dashboard | — | ✅ | ✅ 聚合 4 个 API |

### 🔴 后端路由存在 + 菜单有 + 页面硬编码 demo（完全没打通）

| 模块 | 后端路由 | Admin 菜单 | 页面 | 空壳程度 |
|------|----------|------------|------|----------|
| CMS | GET/PUT /site-contents | ✅ | SiteContent.vue | 100% 硬编码 5 条 JSON |
| DSAR | CRUD + export + delete | ✅ | DSAR.vue | list 硬编码；export/erase 绑了 API |
| 直播日历 | GET /live-rooms/schedule/calendar | ✅ | LiveCalendar.vue | 100% 硬编码 |
| 客户申请 | GET /live-rooms/customer-requests + approve | ✅ | CustomerRequests.vue | 100% 硬编码 |
| 验货间 | POST /live-rooms/system-create | ✅ | DeliveryInspection.vue | 100% 硬编码 |
| QR Code | POST /qrcodes/generate | ✅ | QRCodes.vue | 100% 硬编码 |
| LiveKit 房间 | GET/POST/DELETE /livekit/rooms | ✅ | LiveKitRooms.vue | 基本空壳 |

### 🔴 后端根本没有对应路由（需新增后端 + 前端）

| 模块 | Admin 菜单 | 页面 | 需新增后端 |
|------|------------|------|------------|
| Staff CRUD | ✅ | Staff.vue | **必须加**：GET/POST/PUT/DELETE /staff |
| Webhook Logs | ✅ | WebhookLogs.vue | **必须加**：webhook_logs model + CRUD |
| Payment Config 保存 | ✅ | PaymentConfig.vue | **必须加**：system_config model + GET/PUT /system/payment-config |
| Customers (Users) | ✅ | Users.vue | **必须加**：GET /users |
| Payment Transactions | ✅ | Transactions.vue | **必须加**：GET /payment/transactions（或从 payment_transactions table 查） |

---

## 🔧 第一轮修复记录（后端 P0/P1）

| # | 文件 | 方法 | 优先级 | 状态 |
|---|------|------|--------|------|
| B-01 | dsar.go | CreateRequest | **P0** | [x] ✅ |
| B-02 | dsar.go | List | P1 | [x] ✅ |
| B-03 | dsar.go | Export / Delete | P1 | [x] ✅ |
| B-04 | site_content.go | List | **P0** | [x] ✅ |
| B-05 | site_content.go | Update | P1 | [x] ✅ |
| B-06 | cookie_consent.go | Submit | P1 | [x] ✅ |
| B-07 | slow_preset.go | 全部 8 方法 | **P0** | [x] ✅ |
| B-08 | live_room.go | 全部 11 方法 | **P0** | [x] ✅ |
| R-01 | router.go | staff/audit-logs | P0 | [x] ✅ → StaffAuthHandler.AuditLogs |
| R-02 | router.go | orders/:id/timeline | P0 | [x] ✅ → OrderHandler.Timeline |
| R-03 | router.go | declarations/:id/void | P0 | [x] ✅ → DeclarationHandler.Void |
| R-04 | router.go | customer-requests/:id/approve | P0 | [x] ✅ → LiveRoomHandler.ApproveRequest |
| R-05 | router.go | live-rooms/system-create | P0 | [x] ✅ → LiveRoomHandler.SystemCreate |
| R-06 | router.go | audit-logs 别名 | P1 | [x] ✅ |
| R-07 | router.go | /test 测试路由 | P1 | [x] ✅ 已删除 |

---

## 🔧 第二轮修复清单（前端空壳 + 后端缺失路由）

### 优先级 P0 — 有后端现成路由，前端硬编码导致数据断裂

| # | 页面 | 修复方式 | 状态 |
|---|------|----------|------|
| F-01 | SiteContent.vue | 改 `ref([...demo])` → `api.get('/site-contents')`；save 调 `api.put('/site-contents/' + id, draft)`；编辑用 JSON 文本框而非 `.content` 错误绑定 | [ ] |
| F-02 | DSAR.vue | list 改硬编码 → `api.get('/dsar/requests')`；export/erase 已绑 | [ ] |
| F-03 | LiveCalendar.vue | 改硬编码 → `api.get('/live-rooms/schedule/calendar')` | [ ] |
| F-04 | CustomerRequests.vue | 改硬编码 → `api.get('/live-rooms/customer-requests')`；Approve 调 `api.post('/live-rooms/customer-requests/' + id + '/approve')` | [ ] |
| F-05 | PaymentConfig.vue | 先**新增后端** `GET/PUT /system/payment-config` + system_config model → 前端再接 | [ ] |

### 优先级 P1 — 有后端现成路由，前端空壳

| # | 页面 | 修复方式 | 状态 |
|---|------|----------|------|
| F-06 | QRCodes.vue | 改硬编码 → 先 `api.get('/custom-products')` 选产品生成；展示调 `api.post('/qrcodes/generate')` | [ ] |
| F-07 | DeliveryInspection.vue | 改硬编码 → 从 `api.get('/live-rooms', { room_type: 'delivery_inspection' })` 过滤；Auto-Create 按钮调 `api.post('/live-rooms/system-create', { order_id })` | [ ] |
| F-08 | LiveRooms.vue | 1) presets 常量对齐后端（`slow_preset`、`obs_tasting`、`delivery_inspection`）；2) 加 Start/End 按钮 | [ ] |
| F-09 | SlowPresets.vue | 1) 表单去掉 `camera_rtmp_url`（后端自动生成）；2) 加 Edit 按钮 | [ ] |
| F-10 | Orders.vue | 1) state 对齐后端（`ready_for_production`）；2) 加 Timeline 跳转按钮 | [ ] |
| F-11 | Invoices.vue | 目前 62 行空壳 → 绑 `api.get('/orders/:id/invoice')` + PDF 下载 | [ ] |
| F-12 | CookieConsent.vue | accept/decline 时调 `api.post('/cookie-consent', {...})` | [ ] |

### 优先级 P1 — 需新增后端路由 + 前端

| # | 模块 | 需新增后端 | 状态 |
|---|------|------------|------|
| B-12 | Staff CRUD | StaffHandler: List / Create / ToggleActive / ResetMFA（staff_repo 已有 CRUD） | [ ] |
| B-13 | Users (Customers) | `GET /users`（查 users table） | [ ] |
| B-14 | Payment Transactions | `GET /payment/transactions`（查 payment_transactions table） | [ ] |
| B-15 | System Payment Config | system_config model + `GET/PUT /system/payment-config` | [ ] |
| B-16 | Webhook Logs | webhook_logs model + `GET /webhook-logs`（或作为审计日志的 webhook target_type 过滤） | [ ] |

### 优先级 P2 — Public Site 空壳

| # | 页面 | 修复方式 | 状态 |
|---|------|----------|------|
| P-01 | TeaGardens.vue | 改硬编码 6 条 demo → `api.get('/public/slow-presets')` | [ ] |
| P-02 | Quality.vue | 改硬编码 → `api.get('/public/sgs-reports')` | [ ] |
| P-03 | Contact.vue | submit 弹 ElMessage → 先加 `POST /contact-form` 后端（或对接 MagicLink service） | [ ] |
| P-04 | Home.vue | 去掉 fallback 硬编码 → 纯 API 驱动 | [ ] |
| P-05 | Live.vue | 去掉 fallback 硬编码 → 纯 API 驱动 | [ ] |

### 优先级 P3 — 设计/配置项完善

| # | 项目 | 说明 | 状态 |
|---|------|------|------|
| D-01 | PaymentConfig.vue | 加 3DS 开关、sandbox/live URL、Webhook Test Send | [ ] |
| D-02 | SiteContent.vue | 覆盖更多 section（Live Calendar、Checkout hero、FAQ 条目、GDPR SLA） | [ ] |
| D-03 | 系统设置总页 | 合并 Staff、Payment、Mail、LiveKit、站点品牌到统一 Settings 页 | [ ] |
| D-04 | 配送规则 | 新增 shipping_zones model + Admin 配置页 | [ ] |
| D-05 | 退款管理 | 新增 refunds model + Admin 页（关联 Declarations void） | [ ] |

---

## 🎯 下一步优先级建议

### 立即可做（不新增后端）
1. **F-01 SiteContent.vue** — 改 API 绑定（55 行 → 需扩展到 ~120 行，加编辑对话框）
2. **F-02 DSAR.vue** — list 绑 API + 状态刷新
3. **F-03 LiveCalendar.vue** — 绑 API + 改日历组件
4. **F-04 CustomerRequests.vue** — 绑 API + approve 按钮
5. **F-08 LiveRooms.vue** — 修常量 + 加 Start/End 按钮
6. **F-09 SlowPresets.vue** — 去 RTMP 字段 + 加 Edit
7. **F-10 Orders.vue** — 修状态名 + 加 Timeline
8. **F-12 CookieConsent.vue** — 绑后端
9. **P-01 TeaGardens.vue** — 改硬编码 → API
10. **P-02 Quality.vue** — 改硬编码 → API

### 需要先加后端
11. **B-12 Staff CRUD**（4 个简单 handler 方法）
12. **B-15 System Payment Config**（system_config model + handler）
13. **F-05 PaymentConfig.vue**（依赖 B-15）
14. **F-07 DeliveryInspection.vue**（已基本可用，调 `/live-rooms/system-create`）
15. **F-11 Invoices.vue**（绑 `/orders/:id/invoice`）
16. **B-13/B-14 Users + Transactions**（查现有 table）
17. **B-16 Webhook Logs**（复用 audit_logs + target_type 过滤）

---

## 回归测试 checklist

- [x] ✅ Go build 零 error
- [x] ✅ go vet 零 warning
- [x] ✅ 全项目 mountain/GPS 零残留
- [ ] 后端启动全路由注册无 panic
- [ ] Staff 登录真取 JWT
- [ ] Public API 真测
- [ ] 前端所有 admin 页面 API 可通
- [ ] 前端所有 public 页面 API 可通
