# 🔧 普洱茶私域系统 — 问题修复追踪清单

> 审计时间：2026-09-17  
> 修复状态 legend：`[ ]` 待修 · `[🔄]` 进行中 · `[x]` 已修 · `[n/a]` 合理降级/不在 scope

---

## 问题分类

- **P0 阻断** — 系统不可用 / 用户核心旅程断
- **P1 功能缺失** — 功能有路由有页面但不真存 DB 或不真调 API
- **P2 体验问题** — 页面是空壳 / 硬编码 demo 数据
- **P3 合理降级** — 外部依赖没跑但代码正确（不阻断）

---

## 一、后端 Handler 不写 DB（最严重，路由通但数据不落库）

| # | 文件 | 方法 | 问题 | 优先级 | 状态 |
|---|------|------|------|--------|------|
| B-01 | dsar.go | CreateRequest | 不写 DB，内存 struct 返回，List 永远空数组 | **P0** | [x] ✅ |
| B-02 | dsar.go | List | 硬返回 `items: []dsarTicket{}` | P1 | [x] ✅ |
| B-03 | dsar.go | Export / Delete | 只返回 message，不真操作 | P1 | [x] ✅ |
| B-04 | site_content.go | List | 硬返回 4 条 gin.H{}，不查 DB | **P0** | [x] ✅ |
| B-05 | site_content.go | Update | 只返回 message，不 UPDATE DB | P1 | [x] ✅ |
| B-06 | cookie_consent.go | Submit | 不写 DB | P1 | [x] ✅ |
| B-07 | slow_preset.go | 全部 8 方法 | repo=0，所有 CRUD 不写 DB | **P0** | [x] ✅（通过 SlowPresetService → LiveRoomRepo） |
| B-08 | live_room.go | 全部 11 方法 | repo=0，所有 CRUD 不写 DB | **P0** | [x] ✅（通过 LiveRoomService → LiveRoomRepo；CustomerRequest 仍是 stub） |
| B-09 | translate.go | TranslateText / Status / ASR | 不写 DB（合理）但 FastAPI 外部依赖没跑 | P3 | [n/a] |
| B-10 | im_ws.go | Serve | WebSocket Hub 真运行，但 Message 不存 DB | P2 | [ ] |
| B-11 | health.go | Health / Ready | 合理：Health 不写 DB | P3 | [n/a] |

---

## 二、router.go 匿名 func（绕过 handler 层，没有 DB 查询）

| # | 路由 | 问题 | 状态 |
|---|------|------|------|
| R-01 | GET `/auth/staff/audit-logs` | 匿名 func，返回 gin.H{} 不查 audit 库 | [x] ✅ → StaffAuthHandler.AuditLogs（已注入 auditDB） |
| R-02 | GET `/auth/orders/:id/timeline` | 匿名 func | [x] ✅ → OrderHandler.Timeline（聚合 audit 库事件） |
| R-03 | POST `/auth/declarations/:id/void` | 匿名 func | [x] ✅ → DeclarationHandler.Void（调 SoftDeleteDeclaration repo） |
| R-04 | POST `/auth/live-rooms/customer-requests/:id/approve` | 匿名 func | [x] ✅ → LiveRoomHandler.ApproveRequest（查 room 后 patch status） |
| R-05 | POST `/auth/live-rooms/system-create` | 匿名 func | [x] ✅ → LiveRoomHandler.SystemCreate（真调 svc.Create） |
| R-06 | GET `/auth/audit-logs` | 匿名 func，只是提示用另一个路径 | [x] ✅ → 别名路由复用 StaffAuth.AuditLogs |
| R-07 | GET `/auth/test` | 测试路由 | [x] ✅ 已删除 |

---

## 三、Admin Dashboard 页面 api=0（空壳不调后端）

| # | 页面 | 应调的 API | 状态 |
|---|------|------------|------|
| A-01 | SiteContent.vue | GET/PUT /site-contents | [ ] |
| A-02 | LiveCalendar.vue | GET /live-rooms/calendar | [ ] |
| A-03 | TranslateStatus.vue | GET /translate/status | [ ] |
| A-04 | WebhookLogs.vue | 需新增 webhook_logs model + API | [ ] |
| A-05 | DeliveryInspection.vue | 需新增发货检验 model + API | [ ] |
| A-06 | Login.vue | 已通过 Pinia store → api/staff/login，非问题 | [n/a] |
| A-07 | Layout.vue | 布局组件，非问题 | [n/a] |

---

## 四、Public Site 页面 api=0（全静态）

| # | 页面 | 应调的 API | 状态 |
|---|------|------------|------|
| P-01 | Home.vue | 从 site-contents + custom-products 拉 CMS + 报价预览 | [ ] |
| P-02 | Live.vue | POST /livekit/token-for-obs + LiveKit Web Player | [ ] |
| P-03 | TeaGardens.vue | 从 slow-presets 拉茶园列表 | [ ] |
| P-04 | Quality.vue | 从 sgs-reports 拉公开数据 | [ ] |
| P-05 | Contact.vue | 需新增 contact_form API + 邮件通知 | [ ] |
| P-06 | About / FAQ / Privacy | 纯 CMS，合理静态 | [n/a] |

---

## 五、Flutter 页面 api=0（UI 空壳）

| # | 页面 | 应调的 API | 状态 |
|---|------|------------|------|
| F-01 | advisor/quote_create_page.dart | POST /custom-products | [ ] |
| F-02 | chat/chat_page.dart | `/ws/im` WebSocket | [ ] |
| F-03 | farmer/slow_presets.dart | GET /slow-presets | [ ] |
| F-04 | farmer/sgs_reports.dart | GET /public/sgs-reports | [ ] |
| F-05 | farmer/order_detail.dart | GET /orders/:id/timeline | [ ] |
| F-06 | home/home_page.dart | GET /public/slow-presets + /public/sgs-reports | [ ] |
| F-07 | user_login.dart | 已真调 `/user/login`，非问题 | [n/a] |

---

## 六、运维/基础设施缺失

| # | 项目 | 状态 |
|---|------|------|
| I-01 | RabbitMQ 未运行 → InvoiceService sync 降级 | 正确行为（有日志），非阻断 | [n/a] |
| I-02 | MailService mock → MagicLink 邮件发不出 | 需配置 SMTP | [ ] |
| I-03 | WireGuard wg binary 缺失 → mock deploy | 需装 wireguard-tools | [ ] |
| I-04 | FastAPI 翻译引擎未运行 → translate 503 | 正确降级，非阻断 | [n/a] |
| I-05 | 2Checkout/PayPal 无测试账号 → Webhook handler 虽真但不能真打 | 非阻断（handler 就绪） | [n/a] |
| I-06 | 雪花 ID 未启用 → PostgreSQL bigserial 自增 | 多节点可能冲突，建议引入 Sonyflake | [ ] |

---

## 修复后回归测试 checklist

- [x] ✅ Go build 零 error（go build ./... exit=0）
- [x] ✅ go vet 零 warning（exit=0）
- [ ] 后端启动 102 路由全注册无 panic
- [ ] Staff 登录 (admin@ukteahouse.co.uk / Admin!Tea2026) 真取 JWT
- [ ] Public API 真测（slow-presets / sgs-reports / qrcodes / cookie-consent）
- [ ] 28 字段 POST custom-products 全链路
- [ ] 28 字段 PUT 真持久化
- [ ] DSAR Create/List/Export/Delete 四端点（已真写 DB 代码就绪）
- [ ] SiteContent CMS List/PUT（已真查/真写 DB 代码就绪）
- [ ] CookieConsent Submit 真存（已真写 DB + 错误日志代码就绪）
- [ ] SlowPreset CRUD + enable/disable 真存（已通过 service→repo）
- [ ] LiveRoom CRUD + start/end 真存（已通过 service→repo）
- [x] ✅ AuditLogs 真查 audit 库（StaffAuthHandler.AuditLogs + auditDB 注入）
- [x] ✅ Order timeline 真返回（OrderHandler.Timeline + auditDB 聚合）
- [x] ✅ Declarations void 真操作（DeclarationHandler.Void → SoftDeleteDeclaration repo）
- [x] ✅ CustomerRequests approve 真写（LiveRoomHandler.ApproveRequest → svc.Update）
- [x] ✅ LiveRooms system-create 真写（LiveRoomHandler.SystemCreate → svc.Create）
- [ ] Flutter quote_create_page 真 POST（待前端接入）
- [ ] Live.vue 集成 HLS/LiveKit Web Player（待前端接入）
- [x] ✅ 全项目 mountain / GPS / 山头名 零残留（grep 扫描仅迁移脚本 rename SQL）
- [ ] 无 JWT 访问 protected → 401
- [ ] 非法 JWT → 401 invalid
