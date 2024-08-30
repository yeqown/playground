## ShardingSphere Proxy

在本地容器内搭建起 ShardingSphere Proxy 集群，用于测试数据库的读写分离、分库分表等功能。
这里集群不启用注册中心，也就是说每个节点都是独立的，不会相互感知。

### 搭建

这里使用本地的 mysql 作为数据源，ShardingSphere Proxy 作为中间件，将请求转发到 mysql 数据源。
- proxy 启动 3 个节点，分别监听 3307、3308、3309 端口。
- 其中分库配置为 2 个库 均在本地 mysql 中创建，分表配置为 2 张表。
- 逻辑库叫做：shardingdb。分库算法使用取模算法，
- 逻辑表为：t_user。分表算法使用取模算法。

> 使用 nerdctl.lima 代替 Docker

1. 下载 ShardingSphere Proxy 的 Docker 镜像 （5.4.0）

```bash
nerdctl.lima pull apache/shardingsphere-proxy:5.0.0
```

2. 创建配置文件

```bash
mkdir -p shardingsphere-proxy
mkdir -p shardingsphere-proxy/node1
mkdir -p shardingsphere-proxy/node2
mkdir -p shardingsphere-proxy/node3
```

每个节点都包含两个配置文件：`server.yaml` 和 `config-sharding.yaml`。这两个文件分别包含：

- `server.yaml`：配置 proxy 的基本信息：
    - 认证/鉴权
    - 系统配置属性
- `config-sharding.yaml`
    - 配置数据源
    - 规则（这里仅包含 数据分片 规则）

> 配置参见 [数据分片配置](https://shardingsphere.apache.org/document/5.4.0/cn/user-manual/shardingsphere-jdbc/yaml-config/rules/sharding/)

server.yaml 配置说明：

```yaml

```

config-xxx 中数据分片规则说明：

```yaml
rules:
- !SHARDING
  tables: # 数据分片规则配置
    <logic_table_name> (+): # 逻辑表名称
      actualDataNodes (?): # 由数据源名 + 表名组成（参考 Inline 语法规则）
      databaseStrategy (?): # 分库策略，缺省表示使用默认分库策略，以下的分片策略只能选其一
        standard: # 用于单分片键的标准分片场景
          shardingColumn: # 分片列名称
          shardingAlgorithmName: # 分片算法名称
        complex: # 用于多分片键的复合分片场景
          shardingColumns: # 分片列名称，多个列以逗号分隔
          shardingAlgorithmName: # 分片算法名称
        hint: # Hint 分片策略
          shardingAlgorithmName: # 分片算法名称
        none: # 不分片
      tableStrategy: # 分表策略，同分库策略
      keyGenerateStrategy: # 分布式序列策略
        column: # 自增列名称，缺省表示不使用自增主键生成器
        keyGeneratorName: # 分布式序列算法名称
      auditStrategy: # 分片审计策略
        auditorNames: # 分片审计算法名称
          - <auditor_name>
          - <auditor_name>
        allowHintDisable: true # 是否禁用分片审计hint
  autoTables: # 自动分片表规则配置
    t_order_auto: # 逻辑表名称
      actualDataSources (?): # 数据源名称
      shardingStrategy: # 切分策略
        standard: # 用于单分片键的标准分片场景
          shardingColumn: # 分片列名称
          shardingAlgorithmName: # 自动分片算法名称
  bindingTables (+): # 绑定表规则列表
    - <logic_table_name_1, logic_table_name_2, ...> 
    - <logic_table_name_1, logic_table_name_2, ...> 
  defaultDatabaseStrategy: # 默认数据库分片策略
  defaultTableStrategy: # 默认表分片策略
  defaultKeyGenerateStrategy: # 默认的分布式序列策略
  defaultShardingColumn: # 默认分片列名称
  
  # 分片算法配置
  shardingAlgorithms:
    <sharding_algorithm_name> (+): # 分片算法名称
      type: # 分片算法类型
      props: # 分片算法属性配置
      # ...
  
  # 分布式序列算法配置
  keyGenerators:
    <key_generate_algorithm_name> (+): # 分布式序列算法名称
      type: # 分布式序列算法类型
      props: # 分布式序列算法属性配置
      # ...
  # 分片审计算法配置
  auditors:
    <sharding_audit_algorithm_name> (+): # 分片审计算法名称
      type: # 分片审计算法类型
      props: # 分片审计算法属性配置
      # ...

- !BROADCAST
  tables: # 广播表规则列表
    - <table_name>
    - <table_name>
```

3. 启动 shardingsphere-proxy 集群

> 使用 bash run.sh 脚本来管理集群的启动和停止、重启、查看日志等操作。

```bash
# nerdctl.lima run -d --name shardingsphere-proxy1 -p 3307:3307 -v $(pwd)/node1:/etc/shardingsphere-proxy apache/shardingsphere-proxy:5.4.0
# nerdctl.lima run -d --name shardingsphere-proxy2 -p 3308:3307 -v $(pwd)/node2:/etc/shardingsphere-proxy apache/shardingsphere-proxy:5.4.0
# nerdctl.lima run -d --name shardingsphere-proxy3 -p 3309:3307 -v $(pwd)/node3:/etc/shardingsphere-proxy apache/shardingsphere-proxy:5.4.0

bash run.sh start
```

4. 测试连接

```bash
# 这里使用 mycli 连接数据库
mycli -h localhost -P 3307 -u root -p root -D sharding_db
mycli -h localhost -P 3308 -u root -p root -D sharding_db
mycli -h localhost -P 3309 -u root -p root -D sharding_db

# mysql
mysql -h 127.0.0.1 -P 3307 -u root -p root -D sharding_db
mysql -h 127.0.0.1 -P 3308 -u root -p root -D sharding_db
mysql -h 127.0.0.1 -P 3309 -u root -p root -D sharding_db
```

5. 随机插入数据

```bash
python insert.py -n 1000
```

### 验证问题

1. 这种部署情况下（分布式单点），多个 shardingsphere-proxy 之间的元信息是否互通？怎么互通？
  - REFRESH METADATA 能否刷新所有节点的元信息？
2. 分片数量如果不是 2 的幂，会怎么样？
  - 迁移时，需要迁移的数据量是否会增加？
3. 
