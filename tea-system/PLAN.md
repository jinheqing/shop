# UK Tea House — 会员体系 + 直播可见性 + 短链 + 社交分享 实施计划

> 每完成一个模块，在开头加 `[x] 已完成 · YYYY-MM-DD`。
> 依赖顺序严格自上而下。

---

## 模块 0：基础设施（先做，因为后面所有模块依赖它）

### [x] 0.1 Tailwind 奢侈色板 & 公用 Class（已完成在之前）
- `ink-900 / ink-800 / ivory-100 / gold / gold-soft / sand`
- `tracking-lux / duration-lux`

### [x] 0.2 admin SlowPresets/LiveRooms/CustomProducts 上传（已完成在之前）

---

## 模块 1：后端 — 新 Models（GORM AutoMigrate 自动建表）

### [x] 1.1 `user_groups` 用户组模型
**文件**: `tea-system/internal/models/user_group.go`
```go
type UserGroup struct {
    ID          uint64    `gorm:"primaryKey"`
    Name        string    `gorm:"size:100;not null"`       // "VIP 2026" / "勐海古树 2026 批次买家"
    Description string    `gorm:"type:text"`
    AutoRule    JSONMap   `gorm:"type:jsonb"`              // {spend_threshold: 5000, registered_before: "2025-01-01"}
    Privileges  JSONArray `gorm:"type:jsonb"`              // ["preorder_priority", "private_masterclass", "offline_garden_tour:batch_2026_menghai", "vip_concierge", "deep_traceability"]
    Reciprocal  JSONArray `gorm:"type:jsonb"`              // ["annabels", "soho_house"]
    CreatedByStaffID *uint64
    CreatedAt, UpdatedAt time.Time
}
// 关联表
type UserGroupMember struct {
    ID        uint64    `gorm:"primaryKey"`
    GroupID   uint64    `gorm:"uniqueIndex:grp_user"`
    UserID    uint64    `gorm:"uniqueIndex:grp_user"`
    AddedAt   time.Time
    AddedByStaffID *uint64
}
```

### [x] 1.2 `short_links` 短链模型
**文件**: `tea-system/internal/models/short_link.go`
```go
type ShortLink struct {
    ID          uint64    `gorm:"primaryKey"`
    Code        string    `gorm:"size:8;uniqueIndex;not null"`  // auto-gen 6-8 chars
    TargetURL   string    `gorm:"size:500;not null"`
    OwnerUserID *uint64
    ClickCount  int       `gorm:"default:0"`
    ExpiresAt   *time.Time
    CreatedAt   time.Time
}
```

### [x] 1.3 `recordings` 回放模型（独立实体）
**文件**: `tea-system/internal/models/recording.go`
```go
type Recording struct {
    ID          uint64    `gorm:"primaryKey"`
    LiveRoomID  *uint64   `gorm:"index"`
    Title       string    `gorm:"size:200"`
    FileURL     string    `gorm:"size:500;not null"`  // /uploads/recordings/xxx.mp4
    DurationSec int
    FileSize    int64
    Visibility  string    `gorm:"size:20;default:'registered'"`  // public / registered / restricted
    VisibleUserIDs  JSONArray `gorm:"type:jsonb"`
    VisibleGroupIDs JSONArray `gorm:"type:jsonb"`
    StartedAt   *time.Time
    EndedAt     *time.Time
    CreatedAt   time.Time
}
```

### [x] 1.4 LiveRoom 加新字段
**文件**: 修改 `tea-system/internal/models/live_room.go`，追加：
```go
Type             string    `gorm:"column:type;size:30;default:'scheduled'"`  // slow_live / scheduled / advisor / admin
Visibility       string    `gorm:"column:visibility;size:20;default:'registered'"`  // public / registered / restricted
VisibleUserIDs   JSONArray `gorm:"column:visible_user_ids;type:jsonb"`
VisibleGroupIDs  JSONArray `gorm:"column:visible_group_ids;type:jsonb"`
EnableRecording  bool      `gorm:"column:enable_recording;default:true"`
RecordingID      *uint64   `gorm:"column:recording_id"`
```

**验证**: 启动后 GORM AutoMigrate 自动建 3 张表 + live_rooms 加 6 个字段。

---

## 模块 2：后端 — Handlers & Routes

### [x] 2.1 User Groups Handler
**文件**: `tea-system/internal/api/handlers/user_group.go`
- `Create` — Admin 创建组
- `List` — Admin 列出所有组
- `Update` — 改名 / 改 Privileges / 改 AutoRule
- `Delete` — 删除组
- `AddMember(userID)` — 手动加用户
- `RemoveMember(userID)` — 手动移除
- `ListMembers` — 列出组成员
- `BulkAddMembers` — 批量加（CSV 导入或按条件）
- `AutoSync` — 按 AutoRule 自动归类

**Router** (需要 Admin 权限):
```
POST   /api/v1/user-groups
GET    /api/v1/user-groups
GET    /api/v1/user-groups/:id
PUT    /api/v1/user-groups/:id
DELETE /api/v1/user-groups/:id
POST   /api/v1/user-groups/:id/members
DELETE /api/v1/user-groups/:id/members/:userId
GET    /api/v1/user-groups/:id/members
POST   /api/v1/user-groups/auto-sync
GET    /api/v1/users/me/groups       // 用户自己能看在哪些组里（前端 /account 用）
```

### [x] 2.2 Short Links Handler
**文件**: `tea-system/internal/api/handlers/short_link.go`
- `Create(target, owner)` — 生成短链
- `Get(code)` — 查
- `Redirect(c)` — **关键**: `GET /s/:code` → 302 到 target（公开路由）
- `Stats` — Admin 统计
- `List` — Admin 列出所有

**Router**:
```
GET    /s/:code                        // 公开: 302 重定向
POST   /api/v1/short-links             // 登录后可创建
GET    /api/v1/short-links             // Admin 列表
GET    /api/v1/short-links/:code       // 查
```

### [x] 2.3 Recordings Handler
**文件**: `tea-system/internal/api/handlers/recording.go`
- `List` — Admin 全部
- `Get` — 详情 + 可见性检查
- `UpdateVisibility` — 改可见性（user_ids / group_ids）
- `Delete` — 删除
- `ListMine` — 用户自己的（/account 用）

**Router**:
```
GET    /api/v1/recordings              // Admin 全部（含可见性 filter）
GET    /api/v1/recordings/:id          // 带鉴权
PUT    /api/v1/recordings/:id/visibility
DELETE /api/v1/recordings/:id
GET    /api/v1/my/recordings           // 登录用户看自己能看的
```

### [x] 2.4 LiveRoom Handler 增强
**文件**: 修改 `tea-system/internal/api/handlers/live_room.go`
- Create 时强制选 `type` / `visibility`
- `RestrictAccess(c)` 中间件 — 检查当前用户是否在 `visible_user_ids` 或 `visible_group_ids` 里（public 跳过）
- `UpdateVisibility(user_ids, group_ids)` — Admin 随时重配
- `Pause / Resume / Stop / Delete` — 管理后台操作

### [x] 2.5 Router 注册
**文件**: 修改 `tea-system/internal/api/router.go`
- 加 Handlers 结构体新字段：`UserGroup *handlers.UserGroupHandler` 等
- 注册全部上面的路由
- `GET /s/:code` 必须放在 `/api/v1` 分组之外（公开，不需要 JWT）

### [x] 2.6 直播可见性中间件（独立）
**文件**: `tea-system/internal/middleware/live_access.go`
- 输入：liveRoom model + 当前 user (可能 nil)
- 输出：true/false
- 逻辑：
  - `visibility == "public"` → true
  - `visibility == "registered"` → 检查 JWT 是否存在
  - `visibility == "restricted"` → 查 user.ID 是否在 `visible_user_ids` 里；不在则查 user 在哪些组，组是否在 `visible_group_ids` 里

---

## 模块 3：后端 — 匿名请求合并

### [ ] 3.1 登录时触发合并
**文件**: 修改 `tea-system/internal/service/magic_link_service.go` 和 `user_auth.go`
- 用户通过 Magic Link 登录成功后，返回 `{"merged_anonymous_count": N, "concierge_inquiries": [...], "custom_products": [...]}`
- 前端据此弹出确认对话框

---

## 模块 4：前端 — public-site

### [x] 4.1 Router 改成 Hash 路由
**文件**: 修改 `packages/public-site/src/router.ts`
```ts
// createWebHistory → createWebHashHistory
```

### [ ] 4.2 登录守卫改造
**文件**: 修改 `packages/public-site/src/router.ts`
- **Checkout**: 未登录 → 跳 `/login?redirect=checkout`
- **Add to Cart**: 未登录 → 弹 Magic Link 登录（不跳走当前页，登录成功后自动加购物车）
- **/account /chat**: 未登录 → 跳 `/login?redirect=...`
- **直播受限页**: 未登录 → 显示 "This broadcast is not publicly accessible" 提示 + Magic Link 入口

### [ ] 4.3 社交分享组件
**文件**: 新建 `packages/public-site/src/components/SharePanel.vue`
- 接收 `url: string` 和 `title: string`
- 按钮：Copy Link / WhatsApp / X (Twitter) / Facebook / LinkedIn
- 只生成各平台专属 URL 格式（WhatsApp 用 `https://wa.me/?text={url}`）
- 所有展示的都是**短链**（先调 `/api/v1/short-links` 创建）

### [ ] 4.4 统一 OG 图
**文件**: `index.html` 加默认 OG meta + 运行时动态注入
```html
<meta property="og:image" content="/og-default.jpg">
<meta property="og:site_name" content="UK Tea House">
<meta property="og:locale" content="en_GB">
```

### [ ] 4.5 /account 页面改造
**文件**: 修改 `packages/public-site/src/views/Account.vue`
加 4 个 Tab：
1. Orders & Bespoke Quotes
2. Chat History
3. My Private Broadcasts & Recordings（调 `/api/v1/my/recordings`）
4. Garden Visits & Masterclass Invitations（Privileges 里的 offline_garden_tour / private_masterclass）

### [ ] 4.6 LiveRoom 可见性守卫
**文件**: 修改 `packages/public-site/src/views/LiveRoom.vue`
- 页面加载时先调 `GET /api/v1/live-rooms/:id`
- 如果返回 403 + `visibility: "restricted"` → 显示 "This broadcast is private — your advisor has not granted access" + Magic Link 入口

### [ ] 4.7 Concierge 匿名合并弹窗
**文件**: 登录成功后如果 `merged_anonymous_count > 0` → 弹确认框

---

## 模块 5：前端 — admin-dashboard

### [ ] 5.1 UserGroups 管理页
**文件**: 新建 `packages/admin-dashboard/src/views/UserGroups.vue`
- 列表 + 新建 + 编辑 + 删除
- 每组显示：Privileges、AutoRule 预览、成员数
- 成员管理：单个加 / 批量加（CSV）/ 移除
- **权限**: Admin + Supervisor

### [ ] 5.2 LiveRooms 管理页增强
**文件**: 修改 `packages/admin-dashboard/src/views/LiveRooms.vue`
- 创建表单加：`type` 下拉、`visibility` 下拉（public / registered / restricted）、`visible_user_ids` 多选、`visible_group_ids` 多选、`enable_recording` 开关（默认 ON）
- 操作按钮：**Pause / Resume / Stop / Delete / Reconfigure Visibility**
- 合并表顶部 Filter：按 created_by_role / visibility / type / status 筛选

### [ ] 5.3 Recordings 独立管理页
**文件**: 新建 `packages/admin-dashboard/src/views/Recordings.vue`
- 列表：从哪个 LiveRoom 来、标题、时长、文件大小、可见性
- 操作：改可见性（同 LiveRooms）、下载、删除
- 关联到用户/组

### [ ] 5.4 ShortLinks 管理页
**文件**: 新建 `packages/admin-dashboard/src/views/ShortLinks.vue`
- 列表：code → target、点击次数、过期时间、创建者
- 手动创建短链
- 统计 Top 10 热门分享

### [ ] 5.5 Admin Router 注册
**文件**: 修改 `packages/admin-dashboard/src/router.ts`
- 加 4 个新 Route，放在 `/` 布局下

---

## 模块 6：端到端测试

### [ ] 6.1 后端编译 + 全部 Go 测试
```bash
cd tea-system && go build ./... && go test ./...
```

### [ ] 6.2 前端双包构建
```bash
cd tea-frontend && pnpm run build
```

### [ ] 6.3 关键 API curl 测试（需要服务启动后）
- 创建 user_group → 加成员 → 查
- 创建 live_room (restricted) → 给用户加 visible → 验证可见性
- 创建 short_link → GET /s/code 验证 302
- 登录用户查 /my/recordings

### [ ] 6.4 回归测试
- 旧功能（SlowPresets/LiveRooms/CustomProducts CRUD）全部跑一遍
- Magic Link 登录流程
- 旧 ConciergeWidget 不受影响

---

## 用户组 Privileges 枚举（前端显示用）
```
preorder_priority           → "优先锁货权" (稀缺茶优先购买)
private_masterclass         → "闭门品鉴 & Masterclass" (顾问专属活动)
offline_garden_tour:*       → "茶园线下参观" (精准到批次)
vip_concierge               → "专属顾问直连" (聊天自动路由)
deep_traceability           → "深度溯源" (完整茶园/茶农/SGS)
offline_invite              → "年度私宴邀请" (后期)
reciprocal_access           → "全球 reciprocal clubs" (后期, 需配置 Reciprocal)
```
