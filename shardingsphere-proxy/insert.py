# 编写一个脚本，可以向数据库中插入任意数量的数据
# Usage: python insert.py -n 1000 -user_id_start 1

import argparse
import random
import pymysql

# CREATE TABLE t_user (
#     id BIGINT NOT NULL comment '主键',
#     /* 主键 */
#     user_id BIGINT NOT NULL unique comment '用户ID',
#     /* 用户ID, 分表 key */
#     mch_id BIGINT NOT NULL comment '商户ID',
#     /* 商户ID, 分库 key */
#     PRIMARY KEY (user_id)
# ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

# mch_ids = [100 - 199]
mch_ids = [i for i in range(100, 200)]

def insert_data(user_id_start:int, n: int):
    conn = pymysql.connect(
        host='localhost',
        user='root',
        password='root',
        database='sharding_db',
        port=3307,
        charset='utf8mb4',
        cursorclass=pymysql.cursors.DictCursor
    )

    with conn.cursor() as cursor:
        for i in range(n):
            # mch_id = random.choice(mch_ids)
            # user_id = user_id_start + i
            sql = f"INSERT INTO t_user (user_id, mch_id) VALUES ({user_id_start + i}, {random.choice(mch_ids)})"
            cursor.execute(sql)
        conn.commit()
    
    conn.close()


# CREATE TABLE t_order (
#     id BIGINT NOT NULL comment '主键',
#     /* 主键 */
#     order_id BIGINT NOT NULL unique comment '订单ID',
#     /* 订单ID */
#     status INT NOT NULL comment '订单状态',
#     /* 订单状态 */
#     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
#     /* 创建时间 */
#     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
#     /* 更新时间 */
# ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

def insert_order_data(order_id_start:int, n: int):
    conn = pymysql.connect(
        host='localhost',
        user='root',
        password='root',
        database='sharding_db',
        port=3307,
        charset='utf8mb4',
        cursorclass=pymysql.cursors.DictCursor
    )

    with conn.cursor() as cursor:
        for i in range(n):
            sql = f"INSERT INTO t_order (order_id, status) VALUES ({order_id_start + i}, {random.choice([0, 1, 2])})"
            cursor.execute(sql)
        conn.commit()

    conn.close()

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Insert data to sharding_db.t_user')
    parser.add_argument('-n', type=int, default=1000, help='Number of data to insert')
    parser.add_argument('-table', type=str, default='t_user', help='Table name')
    parser.add_argument('-user_id_start', type=int, default=1, help='Start user_id')
    parser.add_argument('-order_id_start', type=int, default=20241114, help='Start order_id') 
    args = parser.parse_args()

    user_id_start = args.user_id_start
    n = args.n

    if user_id_start < 1:
        user_id_start = 1000
    if n < 1:
        n = 1000

    print(f"Insert {args.n} data to t_user, start from user_id {args.user_id_start}...")

    if args.table == 't_user':
        insert_data(args.user_id_start, args.n)
    if args.table == 't_order':
        insert_order_data(args.order_id_start, args.n)

    print(f"Insert {args.n} data successfully!")
