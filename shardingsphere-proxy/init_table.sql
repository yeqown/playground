-- 通过 sharding-sphere proxy 初始化表 t_user
CREATE TABLE t_user (
    id BIGINT NOT NULL comment '主键',
    /* 主键 */
    user_id BIGINT NOT NULL unique comment '用户ID',
    /* 用户ID, 分表 key */
    mch_id BIGINT NOT NULL comment '商户ID',
    /* 商户ID, 分库 key */
    PRIMARY KEY (user_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;