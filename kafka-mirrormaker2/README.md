# Kafka MirrorMaker2 实验环境

## 环境组成
- Source Kafka 集群: localhost:9092 (Bitnami Kafka 3.3.2)
- Target Kafka 集群: localhost:9094 (Bitnami Kafka 3.3.2)
- MirrorMaker2: 复制 test-topic 从 source 到 target

## 使用方法

1. 启动环境:
```bash
./start.sh
```

2. 开始生产数据:
```bash
./producer.sh
```

3. 验证数据复制 (新终端):
```bash
podman exec kafka-target /opt/bitnami/kafka/bin/kafka-console-consumer.sh --topic source.test-topic --bootstrap-server localhost:9092 --from-beginning
```

## 验证数据复制

### 1. 消费 target 集群数据
```bash
podman exec kafka-target /opt/bitnami/kafka/bin/kafka-console-consumer.sh --topic source.test-topic --bootstrap-server localhost:9092 --from-beginning
```

### 2. 检查 target 集群 topic 列表
```bash
podman exec kafka-target /opt/bitnami/kafka/bin/kafka-topics.sh --list --bootstrap-server localhost:9092
```

### 3. 查看 topic 详细信息
```bash
podman exec kafka-target /opt/bitnami/kafka/bin/kafka-topics.sh --describe --topic source.test-topic --bootstrap-server localhost:9092
```

### 4. 检查消息数量
```bash
# Source 集群消息数
podman exec kafka-source /opt/bitnami/kafka/bin/kafka-run-class.sh kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic test-topic

# Target 集群消息数  
podman exec kafka-target /opt/bitnami/kafka/bin/kafka-run-class.sh kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic source.test-topic
```

## 停止环境
```bash
./stop.sh
```
