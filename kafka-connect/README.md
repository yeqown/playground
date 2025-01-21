## kafka-connect 

搭建一套 mysql -> kafka -> ES 的数据同步演示环境。

mysql 采用容器部署映射到本地的 3306 端口
kafka 使用容器部署映射到本地 9092 端口
ES 使用容器部署映射到本地 9200 端口
kafka-connect 采用容器部署部署映射到本地 8083 端口

### Get Started

1. 启动 compose 服务组

```bash
nerdctl.lima compose up -d
```

2. 检查服务是否正常

检查 mysql 服务是否正常
```bash
nerdctl.lima exec mysql mysql -uroot -p123456 -e "show databases;"
```

检查 kafka 服务是否正常
```bash
nerdctl.lima exec kafka kafka-topics --list --zookeeper zookeeper:2181
```

检查 ES 服务是否正常
```bash
nerdctl.lima exec es curl nerdctl.lima exec es curl http://localhost:9200/_cluster/health
```

检查 kafka-connect 服务是否正常
```bash
nerdctl.lima exec kafka-connect curl http://localhost:8083/connectors

3. 注册 kafka-connecor tasks

注册 mysql source connector

```bash
curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @/Users/apple/projects/opensource/playground/kafka-connect/connector-tasks/mysql-source.test.users.json
```

注册 es sink connector

```bash
curl -X POST curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @/Users/apple/projects/opensource/playground/kafka-connect/connector-tasks/es-sink.test.users.json
```

检查任务是否正常
```bash
curl http://localhost:8083/connectors/mysql-source/status

curl http://localhost:8083/connectors/elasticsearch-sink/status
```


4. 插入数据

```bash
nerdctl.lima exec -it kafka-connect-mysql-1 mysql -uroot -proot -e "INSERT INTO test.users (name, email) VALUES ('测试用户', 'test@example.com');"
```

5. 验证数据同步

查看所有索引：
```bash
curl -X GET "curl -X GET "http://localhost:9200/_cat/indices?v"
```

