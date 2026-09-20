-- ============================================================
-- 00002: 用户资料字段 + 订单物流字段
-- ============================================================

-- users: 引荐人姓名
ALTER TABLE users ADD COLUMN IF NOT EXISTS referred_by_name varchar(100);

-- users: 社交账号 (WhatsApp / WeChat / Instagram 等)
ALTER TABLE users ADD COLUMN IF NOT EXISTS social_accounts jsonb;

-- orders: 物流单号
ALTER TABLE orders ADD COLUMN IF NOT EXISTS tracking_number varchar(60);

-- orders: 承运商
ALTER TABLE orders ADD COLUMN IF NOT EXISTS shipping_carrier varchar(50);

-- orders: 发货时间
ALTER TABLE orders ADD COLUMN IF NOT EXISTS shipped_at timestamp;
