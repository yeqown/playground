-- 通过 sharding-sphere proxy 初始化表 t_user
CREATE TABLE t_user (
    id BIGINT NOT NULL comment '主键',
    /* 主键 */
    user_id BIGINT NOT NULL unique comment '用户ID',
    /* 用户ID, 分表 key */
    mch_id BIGINT NOT NULL comment '商户ID',
    /* 商户ID, 分库 key */
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
    /* 创建时间 */
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    /* 更新时间 */
    PRIMARY KEY (user_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE t_order (
    id BIGINT NOT NULL comment '主键',
    order_id BIGINT NOT NULL unique comment '订单ID',
    `status` INT NOT NULL comment '订单状态',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;