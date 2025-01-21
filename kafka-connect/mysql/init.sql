-- 创建测试表
CREATE TABLE IF NOT EXISTS test.users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 设置权限
GRANT ALL PRIVILEGES ON test.* TO 'root' @'%';

GRANT REPLICATION SLAVE,
REPLICATION CLIENT ON *.* TO 'root' @'%';

FLUSH PRIVILEGES;

-- 插入测试数据
INSERT INTO
    test.users (name, email, created_at)
VALUES
    (
        'John',
        'zhangsan@example.com',
        CURRENT_TIMESTAMP()
    ),
    ('Teri', 'lisi@example.com', CURRENT_TIMESTAMP());