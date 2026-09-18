-- ============================================================
-- 英国高端定制普洱茶系统 — 初始化 Schema
-- 17 张业务表 + audit_logs（独立 PostgreSQL 实例）
-- 所有表: id bigserial PK, created_at, updated_at, deleted_at(soft delete)
-- ============================================================

-- ---------- 1. users ----------
CREATE TABLE IF NOT EXISTS users (
    id                      bigserial PRIMARY KEY,
    name                    varchar(100) NOT NULL,
    email                   varchar(200) NOT NULL UNIQUE,
    phone                   varchar(30),
    password_hash           varchar(255),
    billing_address         jsonb,
    delivery_address        jsonb,
    preferred_language      varchar(10) NOT NULL DEFAULT 'en',
    preferred_timezone      varchar(50) NOT NULL DEFAULT 'Europe/London',
    preferred_advisor_id    bigint,
    consent_marketing       boolean NOT NULL DEFAULT false,
    consent_analytics       boolean NOT NULL DEFAULT false,
    data_delete_requested_at timestamp,
    data_delete_completed_at timestamp,
    dsar_request_count      integer NOT NULL DEFAULT 0,
    last_login_at           timestamp,
    created_at              timestamp NOT NULL,
    updated_at              timestamp NOT NULL,
    deleted_at              timestamp
);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_preferred_advisor ON users (preferred_advisor_id);

-- ---------- 2. staff ----------
CREATE TABLE IF NOT EXISTS staff (
    id                      bigserial PRIMARY KEY,
    name                    varchar(100) NOT NULL,
    email                   varchar(200) NOT NULL UNIQUE,
    password_hash           varchar(255) NOT NULL,
    role                    varchar(30) NOT NULL,
    permissions             jsonb NOT NULL DEFAULT '[]'::jsonb,
    mfa_enabled             boolean NOT NULL DEFAULT false,
    mfa_secret              varchar(100),
    work_timezone           varchar(50) NOT NULL DEFAULT 'Europe/London',
    assigned_farm_id        bigint,
    is_active               boolean NOT NULL DEFAULT true,
    last_login_at           timestamp,
    created_at              timestamp NOT NULL,
    updated_at              timestamp NOT NULL,
    deleted_at              timestamp
);
CREATE INDEX IF NOT EXISTS idx_staff_deleted_at ON staff (deleted_at);

-- ---------- 3. conversations ----------
CREATE TABLE IF NOT EXISTS conversations (
    id                  bigserial PRIMARY KEY,
    conversation_type   varchar(20) NOT NULL,
    title               varchar(200),
    creator_id          bigint NOT NULL,
    last_message_at     timestamp,
    created_at          timestamp NOT NULL,
    updated_at          timestamp NOT NULL,
    deleted_at          timestamp
);
CREATE INDEX IF NOT EXISTS idx_conversations_deleted_at ON conversations (deleted_at);
CREATE INDEX IF NOT EXISTS idx_conversations_creator ON conversations (creator_id);

-- ---------- 4. conversation_participants ----------
CREATE TABLE IF NOT EXISTS conversation_participants (
    id              bigserial PRIMARY KEY,
    conversation_id bigint NOT NULL,
    user_id         bigint REFERENCES users(id),
    staff_id        bigint REFERENCES staff(id),
    role            varchar(20) NOT NULL DEFAULT 'member',
    last_read_at    timestamp,
    created_at      timestamp NOT NULL,
    deleted_at      timestamp,
    CONSTRAINT chk_participant_owner CHECK (user_id IS NOT NULL OR staff_id IS NOT NULL)
);
CREATE INDEX IF NOT EXISTS idx_conv_part_conv ON conversation_participants (conversation_id);
CREATE INDEX IF NOT EXISTS idx_conv_part_user ON conversation_participants (user_id);
CREATE INDEX IF NOT EXISTS idx_conv_part_staff ON conversation_participants (staff_id);

-- ---------- 5. messages ----------
CREATE TABLE IF NOT EXISTS messages (
    id                  bigserial PRIMARY KEY,
    conversation_id     bigint NOT NULL REFERENCES conversations(id),
    sender_type         varchar(10) NOT NULL,
    sender_id           bigint NOT NULL,
    message_type        varchar(20) NOT NULL,
    content             text NOT NULL,
    translation_zh      text,
    translation_en      text,
    translation_status  varchar(20) NOT NULL DEFAULT 'pending',
    reply_to_id         bigint,
    attachment_urls     jsonb NOT NULL DEFAULT '[]'::jsonb,
    created_at          timestamp NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_messages_conv ON messages (conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_created ON messages (created_at);

-- ---------- 6. custom_products ----------
CREATE TABLE IF NOT EXISTS custom_products (
    id                      bigserial PRIMARY KEY,
    product_token           varchar(64) NOT NULL UNIQUE,
    version                 integer NOT NULL DEFAULT 1,
    is_bespoke              boolean NOT NULL DEFAULT true,
    non_refundable          boolean NOT NULL DEFAULT true,
    title                   varchar(200) NOT NULL,
    raw_tea_source          text NOT NULL,
    custom_requirement      text NOT NULL,
    tea_type                varchar(20) NOT NULL,
    tea_shape               varchar(20) NOT NULL,
    tea_shape_weight        integer,
    smoked_with_flower      boolean NOT NULL DEFAULT false,
    flower_type             varchar(30),
    inner_packaging         varchar(100) NOT NULL,
    outer_packaging         varchar(100) NOT NULL,
    product_card_text       text,
    product_card_format     varchar(20) NOT NULL DEFAULT 'vertical',
    qr_code_position        varchar(20) NOT NULL,
    sku                     varchar(50) UNIQUE,
    unit_price              numeric(10,2) NOT NULL,
    quantity                integer NOT NULL,
    shipping_cost           numeric(10,2) NOT NULL,
    total_amount            numeric(10,2) NOT NULL,
    lead_time               varchar(100) NOT NULL,
    status                  varchar(20) NOT NULL,
    created_by_staff_id     bigint REFERENCES staff(id),
    reviewed_by_staff_id    bigint REFERENCES staff(id),
    reviewed_at             timestamp,
    -- 溯源字段
    harvest_date            date NOT NULL,
    roasting_date           date NOT NULL,
    mountain_location       varchar(200) NOT NULL,
    master_name             varchar(100) NOT NULL,
    storage_location        varchar(200) NOT NULL,
    sgs_report_id           bigint,
    created_at              timestamp NOT NULL,
    updated_at              timestamp NOT NULL,
    deleted_at              timestamp
);
CREATE INDEX IF NOT EXISTS idx_custom_products_deleted_at ON custom_products (deleted_at);

-- ---------- 7. orders ----------
CREATE TABLE IF NOT EXISTS orders (
    id                          bigserial PRIMARY KEY,
    order_no                    varchar(30) NOT NULL UNIQUE,
    user_id                     bigint NOT NULL REFERENCES users(id),
    staff_id                    bigint NOT NULL REFERENCES staff(id),
    custom_product_id           bigint NOT NULL REFERENCES custom_products(id),
    custom_product_snapshot     jsonb NOT NULL,
    state                       varchar(30) NOT NULL,
    unit_price                  numeric(10,2) NOT NULL,
    quantity                    integer NOT NULL,
    shipping_cost               numeric(10,2) NOT NULL,
    total_amount                numeric(10,2) NOT NULL,
    hs_code                     varchar(20) NOT NULL,
    country_of_origin           varchar(50) NOT NULL,
    billing_address_snapshot    jsonb NOT NULL,
    delivery_address_snapshot  jsonb NOT NULL,
    live_room_id                bigint,
    created_at                  timestamp NOT NULL,
    updated_at                  timestamp NOT NULL,
    deleted_at                  timestamp
);
CREATE INDEX IF NOT EXISTS idx_orders_user ON orders (user_id);
CREATE INDEX IF NOT EXISTS idx_orders_state ON orders (state);
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders (deleted_at);

-- ---------- 8. payment_transactions ----------
CREATE TABLE IF NOT EXISTS payment_transactions (
    id                      bigserial PRIMARY KEY,
    order_id                bigint NOT NULL REFERENCES orders(id),
    payment_gateway         varchar(20) NOT NULL,
    gateway_transaction_id  varchar(100) NOT NULL UNIQUE,  -- 🔑 数据库层幂等关键约束
    amount                  numeric(10,2) NOT NULL,
    currency                varchar(3) NOT NULL DEFAULT 'GBP',
    status                  varchar(20) NOT NULL,
    raw_callback            jsonb,
    paid_at                 timestamp,
    created_at              timestamp NOT NULL,
    updated_at              timestamp NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_payment_tx_order ON payment_transactions (order_id);

-- ---------- 9. invoices ----------
CREATE TABLE IF NOT EXISTS invoices (
    id              bigserial PRIMARY KEY,
    order_id        bigint NOT NULL REFERENCES orders(id) UNIQUE,
    invoice_no      varchar(30) NOT NULL UNIQUE,
    issue_date      date NOT NULL,
    pdf_url         varchar(500) NOT NULL,
    locked          boolean NOT NULL DEFAULT true,
    hs_code         varchar(20) NOT NULL,
    country_of_origin varchar(50) NOT NULL,
    total_amount    numeric(10,2) NOT NULL,
    created_at      timestamp NOT NULL,
    updated_at      timestamp NOT NULL
);

-- ---------- 10. declaration_ledger ----------
CREATE TABLE IF NOT EXISTS declaration_ledger (
    id                      bigserial PRIMARY KEY,
    order_id                bigint NOT NULL REFERENCES orders(id) UNIQUE,
    customs_declaration_no  varchar(50),
    hs_code                 varchar(20) NOT NULL,
    commodity_desc          text NOT NULL,
    gross_weight            numeric(10,3),
    net_weight              numeric(10,3),
    declared_value          numeric(10,2),
    customs_status          varchar(30),
    declaration_date        date,
    cleared_date            date,
    remarks                 text,
    voided_at               timestamp,
    voided_reason           text,
    created_at              timestamp NOT NULL,
    updated_at              timestamp NOT NULL
);

-- ---------- 11. foreign_exchange_ledger ----------
CREATE TABLE IF NOT EXISTS foreign_exchange_ledger (
    id                      bigserial PRIMARY KEY,
    order_id                bigint NOT NULL REFERENCES orders(id) UNIQUE,
    payment_gateway         varchar(20) NOT NULL,
    gateway_transaction_id  varchar(100) NOT NULL,
    amount_gbp              numeric(10,2) NOT NULL,
    exchange_rate           numeric(10,6),
    amount_cny              numeric(12,2),
    settlement_date         date,
    bank_account            varchar(100),
    remarks                 text,
    created_at              timestamp NOT NULL,
    updated_at              timestamp NOT NULL
);

-- ---------- 12. live_rooms ----------
CREATE TABLE IF NOT EXISTS live_rooms (
    id                      bigserial PRIMARY KEY,
    room_id                 varchar(64) NOT NULL UNIQUE,
    room_name               varchar(200) NOT NULL,
    room_type               varchar(30) NOT NULL,
    location                varchar(200),
    cover_image             varchar(500),
    description             text,
    push_source             varchar(20) NOT NULL,
    camera_rtmp_url         varchar(500),
    obs_rtmp_url            varchar(500),
    obs_rtmp_key            varchar(100),
    order_id                bigint REFERENCES orders(id),
    host_staff_id           bigint REFERENCES staff(id),
    livekit_token_for_host  text,
    livekit_token_for_viewers text,
    status                  varchar(20) NOT NULL,
    scheduled_start         timestamp,
    started_at              timestamp,
    ended_at                timestamp,
    peak_viewers            integer NOT NULL DEFAULT 0,
    recording_url           varchar(500),
    translation_session_id  bigint,
    created_by_staff_id     bigint REFERENCES staff(id),
    created_at              timestamp NOT NULL,
    updated_at              timestamp NOT NULL,
    deleted_at              timestamp
);
CREATE INDEX IF NOT EXISTS idx_live_rooms_type ON live_rooms (room_type);
CREATE INDEX IF NOT EXISTS idx_live_rooms_status ON live_rooms (status);
CREATE INDEX IF NOT EXISTS idx_live_rooms_deleted_at ON live_rooms (deleted_at);

-- ---------- 13. sgs_reports ----------
CREATE TABLE IF NOT EXISTS sgs_reports (
    id          bigserial PRIMARY KEY,
    report_no   varchar(100) NOT NULL UNIQUE,
    batch_no    varchar(100) NOT NULL,
    tea_type    varchar(30) NOT NULL,
    test_date   date NOT NULL,
    issue_date  date NOT NULL,
    pdf_url     varchar(500) NOT NULL,
    test_items  jsonb,
    created_at  timestamp NOT NULL,
    updated_at  timestamp NOT NULL
);

-- ---------- 14. nodes ----------
-- ⚠️ SSH 密码**不存入此表**（部署期间仅在 Go 进程内存存在）
CREATE TABLE IF NOT EXISTS nodes (
    id                      bigserial PRIMARY KEY,
    node_name               varchar(100) NOT NULL,
    node_type               varchar(30) NOT NULL,
    public_ip               varchar(50) NOT NULL,
    wireguard_ip            varchar(50) NOT NULL,
    wireguard_public_key    varchar(200) NOT NULL,
    status                  varchar(20) NOT NULL,
    last_health_check       timestamp,
    last_cpu_percent        integer,
    last_memory_percent     integer,
    last_bandwidth_mbps     integer,
    last_viewers_count      integer,
    deployed_by_staff_id    bigint REFERENCES staff(id),
    deployed_at             timestamp,
    deleted_at              timestamp,
    created_at              timestamp NOT NULL,
    updated_at              timestamp NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes (status);
CREATE INDEX IF NOT EXISTS idx_nodes_deleted_at ON nodes (deleted_at);

-- ---------- 15. translation_sessions ----------
CREATE TABLE IF NOT EXISTS translation_sessions (
    id          bigserial PRIMARY KEY,
    source_type varchar(20) NOT NULL,
    source_id   bigint,
    direction   varchar(10) NOT NULL,
    status      varchar(20) NOT NULL,
    created_at  timestamp NOT NULL,
    updated_at  timestamp NOT NULL
);

-- ---------- 16. audit_logs（独立 PostgreSQL 实例，无软删除）----------
-- 注意：此表在独立 PostgreSQL 实例上创建，DDL 放在此处保持完整参考
CREATE TABLE IF NOT EXISTS audit_logs (
    id          bigserial PRIMARY KEY,
    staff_id    bigint,
    action      varchar(100) NOT NULL,
    target_type varchar(50),
    target_id   bigint,
    detail      jsonb,
    ip_address  varchar(50),
    user_agent  text,
    created_at  timestamp NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_staff ON audit_logs (staff_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs (action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs (created_at);

-- ---------- 17. site_contents ----------
CREATE TABLE IF NOT EXISTS site_contents (
    id                  bigserial PRIMARY KEY,
    page_key            varchar(50) NOT NULL UNIQUE,
    section_key         varchar(50) NOT NULL,
    content             jsonb NOT NULL,
    updated_by_staff_id bigint REFERENCES staff(id),
    updated_at          timestamp NOT NULL
);
