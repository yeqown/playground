## ShardingSphere Proxy

在本地容器内搭建起 ShardingSphere Proxy 集群，用于测试数据库的读写分离、分库分表等功能。

### 搭建

这里使用本地的 mysql 作为数据源，ShardingSphere Proxy 作为中间件，将请求转发到 mysql 数据源。
- proxy 启动 3 个节点，分别监听 3307、3308、3309 端口。
- 其中分库配置为 2 个库 均在本地 mysql 中创建，分表配置为 2 张表。
- 逻辑库叫做：sharding_db。分库算法使用取模算法，
- 逻辑表为：t_table。分表算法使用取模算法。

> 使用 nerdctl.lima 代替 Docker

1. 下载 ShardingSphere Proxy 的 Docker 镜像 （5.4.0）

```bash
nerdctl.lima pull apache/shardingsphere-proxy:5.0.0
```

2. 创建配置文件

```bash
mkdir -p shardingsphere-proxy
```

在 shardingsphere-proxy 目录下创建 `config-*.yaml` 文件，内容如下：


3. 启动 shardingsphere-proxy 集群

```bash
nerdctl.lima run -d --name shardingsphere-proxy -p 3307:3307 -v $(pwd)/shardingsphere-proxy:/etc/shardingsphere-proxy apache/shardingsphere-proxy:5.0.0
```

4. 连接测试

获取容器的 IP 地址

```bash
nerdctl.lima inspect --format '{{.NetworkSettings.IPAddress}}' shardingsphere-proxy
```

```bash
mysql -h YOUR_CONTAINER_IP -P 3307 -u root -proot
```

5. 数据库配置

```sql
CREATE DATABASE sharding_db;
