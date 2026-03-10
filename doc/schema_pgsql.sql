-- ============================================================
-- 投壶股票系统 - PostgreSQL 建表脚本
-- 基于 doc 下产品说明、实现方案、接口列表及项目 model 生成
-- ============================================================

-- 1. 建库（需以 superuser 执行，或单独执行）
-- CREATE DATABASE stock_db
--   WITH ENCODING = 'UTF8'
--        LC_COLLATE = 'zh_CN.UTF-8'
--        LC_CTYPE = 'zh_CN.UTF-8'
--        TEMPLATE = template0;

-- \c stock_db;

-- ============================================================
-- 2. 用户表 user
-- ============================================================
CREATE TABLE IF NOT EXISTS "user" (
    user_id         VARCHAR(32) PRIMARY KEY,  -- 雪花ID，与前端交互用string
    user_name       VARCHAR(255) NOT NULL,
    pass_word       VARCHAR(255) NOT NULL,
    real_name       VARCHAR(255) DEFAULT '',
    name            VARCHAR(255) DEFAULT '',
    we_chat         VARCHAR(255) DEFAULT '',
    phone_number    VARCHAR(255) DEFAULT '',
    address         TEXT DEFAULT '',
    luck            BIGINT NOT NULL DEFAULT 0,  -- 幸运值（货币）
    referrer        VARCHAR(255) DEFAULT '',
    del_flg         SMALLINT NOT NULL DEFAULT 0,  -- 标记删除 0=正常 1=已注销
    deletion_reason TEXT DEFAULT '',
    ctime           BIGINT NOT NULL,  -- 注册时间戳
    CONSTRAINT uk_user_name UNIQUE (user_name)
);

COMMENT ON TABLE "user" IS '用户表';
COMMENT ON COLUMN "user".luck IS '幸运值（货币单位）';
COMMENT ON COLUMN "user".del_flg IS '0=正常 1=已注销';

-- ============================================================
-- 3. 项目（物品）表 item
-- ============================================================
CREATE TABLE IF NOT EXISTS item (
    id              VARCHAR(32) PRIMARY KEY,  -- 雪花ID
    name            VARCHAR(255) NOT NULL,
    description     TEXT DEFAULT '',
    purchase_price  BIGINT NOT NULL,  -- 进货价（幸运值）
    total_value     BIGINT NOT NULL DEFAULT 0,  -- 当前总价值
    total_quantity  BIGINT NOT NULL DEFAULT 0,  -- 当前总股数
    status          SMALLINT NOT NULL DEFAULT 1,  -- 1=立项中 2=已淘汰 3=游戏中 4=已结算
    ctime           BIGINT NOT NULL,
    mtime           BIGINT NOT NULL
);

COMMENT ON TABLE item IS '项目/物品表，每个物品对应一支股票';
COMMENT ON COLUMN item.status IS '1=立项中 2=已淘汰 3=游戏中 4=已结算';

-- ============================================================
-- 4. 用户持仓表 stock_user
-- ============================================================
CREATE TABLE IF NOT EXISTS stock_user (
    stock_user_id   VARCHAR(32) PRIMARY KEY,  -- 雪花ID
    stock_id        VARCHAR(32) NOT NULL,  -- 项目ID，对应 item.id
    stock_name      VARCHAR(255) DEFAULT '',  -- 冗余，方便查询展示
    stock_number    BIGINT NOT NULL DEFAULT 0,  -- 持有股数
    user_id         VARCHAR(32) NOT NULL,
    ctime           BIGINT,  -- 可选
    mtime           BIGINT,  -- 可选
    CONSTRAINT fk_stock_user_user FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_stock_user_item FOREIGN KEY (stock_id) REFERENCES item(id) ON DELETE CASCADE,
    CONSTRAINT uk_stock_user_user_item UNIQUE (user_id, stock_id)
);

CREATE INDEX idx_stock_user_user_id ON stock_user(user_id);
CREATE INDEX idx_stock_user_stock_id ON stock_user(stock_id);

COMMENT ON TABLE stock_user IS '用户持仓表';

-- ============================================================
-- 5. 交易大厅挂单表 "order"
-- order 为 PostgreSQL 保留字，需加双引号
-- ============================================================
CREATE TABLE IF NOT EXISTS "order" (
    id              VARCHAR(32) PRIMARY KEY,  -- 雪花ID
    order_type      SMALLINT NOT NULL,  -- 1=卖单 2=买单
    initiator_id    VARCHAR(32) NOT NULL,  -- 发起人 user_id
    item_id         VARCHAR(32) NOT NULL,
    price           BIGINT NOT NULL,  -- 单价（幸运值）
    quantity        BIGINT NOT NULL,  -- 挂单数量
    remaining       BIGINT NOT NULL,  -- 剩余未成交数量
    status          SMALLINT NOT NULL DEFAULT 1,  -- 1=挂单中 2=部分成交 3=完全成交 4=已撤销
    ctime           BIGINT NOT NULL,
    mtime           BIGINT NOT NULL,
    CONSTRAINT fk_order_user FOREIGN KEY (initiator_id) REFERENCES "user"(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_order_item FOREIGN KEY (item_id) REFERENCES item(id) ON DELETE CASCADE
);

CREATE INDEX idx_order_item_id ON "order"(item_id);
CREATE INDEX idx_order_status ON "order"(status);
CREATE INDEX idx_order_initiator ON "order"(initiator_id);

COMMENT ON TABLE "order" IS '交易大厅挂单表';
COMMENT ON COLUMN "order".order_type IS '1=卖单 2=买单';
COMMENT ON COLUMN "order".status IS '1=挂单中 2=部分成交 3=完全成交 4=已撤销';

-- ============================================================
-- 6. 成交记录表 trade_record
-- ============================================================
CREATE TABLE IF NOT EXISTS trade_record (
    id              VARCHAR(32) PRIMARY KEY,  -- 雪花ID
    order_id        VARCHAR(32) NOT NULL,
    buyer_id        VARCHAR(32) NOT NULL,
    seller_id       VARCHAR(32) NOT NULL,
    item_id         VARCHAR(32) NOT NULL,
    quantity        BIGINT NOT NULL,
    price           BIGINT NOT NULL,
    total_amount    BIGINT NOT NULL,  -- 总金额（幸运值）
    ctime           BIGINT NOT NULL,
    CONSTRAINT fk_trade_record_order FOREIGN KEY (order_id) REFERENCES "order"(id) ON DELETE CASCADE,
    CONSTRAINT fk_trade_record_buyer FOREIGN KEY (buyer_id) REFERENCES "user"(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_trade_record_seller FOREIGN KEY (seller_id) REFERENCES "user"(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_trade_record_item FOREIGN KEY (item_id) REFERENCES item(id) ON DELETE CASCADE
);

CREATE INDEX idx_trade_record_order_id ON trade_record(order_id);
CREATE INDEX idx_trade_record_buyer ON trade_record(buyer_id);
CREATE INDEX idx_trade_record_seller ON trade_record(seller_id);

COMMENT ON TABLE trade_record IS '成交记录表';

-- ============================================================
-- 7. 投壶记录表 throw_record
-- ============================================================
CREATE TABLE IF NOT EXISTS throw_record (
    id              VARCHAR(32) PRIMARY KEY,  -- 雪花ID
    player_pay      BIGINT NOT NULL,  -- 玩家付费（幸运值）
    total_items     INT NOT NULL,  -- 总项目数，用于算单项目收益
    ctime           BIGINT NOT NULL
);

COMMENT ON TABLE throw_record IS '投壶记录表';

-- ============================================================
-- 8. 投壶命中明细表 throw_hit
-- ============================================================
CREATE TABLE IF NOT EXISTS throw_hit (
    id              VARCHAR(32) PRIMARY KEY,  -- 雪花ID
    throw_record_id VARCHAR(32) NOT NULL,
    item_id         VARCHAR(32) NOT NULL,
    hit_count       INT NOT NULL DEFAULT 0,  -- 该物品被投中次数
    ctime           BIGINT NOT NULL,
    CONSTRAINT fk_throw_hit_record FOREIGN KEY (throw_record_id) REFERENCES throw_record(id) ON DELETE CASCADE,
    CONSTRAINT fk_throw_hit_item FOREIGN KEY (item_id) REFERENCES item(id) ON DELETE CASCADE
);

CREATE INDEX idx_throw_hit_record ON throw_hit(throw_record_id);

COMMENT ON TABLE throw_hit IS '投壶命中明细表';

-- ============================================================
-- 9. 结算批次表 settlement_batch
-- ============================================================
CREATE TABLE IF NOT EXISTS settlement_batch (
    id              VARCHAR(32) PRIMARY KEY,  -- 雪花ID
    date            BIGINT NOT NULL,  -- 结算日期时间戳
    status          SMALLINT NOT NULL DEFAULT 1,  -- 1=处理中 2=已完成
    ctime           BIGINT NOT NULL,
    mtime           BIGINT NOT NULL
);

COMMENT ON TABLE settlement_batch IS '每日结算批次表';
COMMENT ON COLUMN settlement_batch.status IS '1=处理中 2=已完成';

-- ============================================================
-- 10. 结算明细表 settlement_detail
-- ============================================================
CREATE TABLE IF NOT EXISTS settlement_detail (
    id              VARCHAR(32) PRIMARY KEY,  -- 雪花ID
    batch_id        VARCHAR(32) NOT NULL,
    user_id         VARCHAR(32) NOT NULL,
    item_id         VARCHAR(32) NOT NULL,
    quantity        BIGINT NOT NULL,  -- 持有股数
    price_per_unit  BIGINT NOT NULL,  -- 每股价值
    refund_amount   BIGINT NOT NULL,  -- 退还幸运值
    ctime           BIGINT NOT NULL,
    CONSTRAINT fk_settlement_batch FOREIGN KEY (batch_id) REFERENCES settlement_batch(id) ON DELETE CASCADE,
    CONSTRAINT fk_settlement_user FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_settlement_item FOREIGN KEY (item_id) REFERENCES item(id) ON DELETE CASCADE
);

CREATE INDEX idx_settlement_detail_batch ON settlement_detail(batch_id);

COMMENT ON TABLE settlement_detail IS '结算明细表';

-- ============================================================
-- 表创建顺序说明（如有外键依赖，需按序创建）：
-- 1. user
-- 2. item
-- 3. stock_user (依赖 user, item)
-- 4. order (依赖 user, item)
-- 5. trade_record (依赖 order, user, item)
-- 6. throw_record
-- 7. throw_hit (依赖 throw_record, item)
-- 8. settlement_batch
-- 9. settlement_detail (依赖 settlement_batch, user, item)
-- ============================================================
