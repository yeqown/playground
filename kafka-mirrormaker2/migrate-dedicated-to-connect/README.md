## 从 Kafka Dedicated 迁移到 Kafka Connect 的 MirrorMaker2

‼️ 此方案未测试通过

> 注意：提取 offset 要事先停止 MirrorMaker2 任务，避免 offset 变化

### 运行 offset 提取脚本

```bash
podman run -it --rm --network kafka-net -v $(pwd):/workspace -w /workspace python:3.9 \
bash -c "pip install kafka-python && python3 migrate.py extract --brokers kafka-target:9092 --dedicated-offset-topic mm2-offsets.source.internal --output-file /workspace/offsets"
```

> 注意：使用 offsets 文件前，请检查内容，确保没有包含内部 topic 的 offset; 保证 connect 已经启动，connector 任务没有创建！！！！

### 使用提取的 offsets 文件

```bash
podman run -it --rm --network kafka-net -v $(pwd):/workspace -w /workspace python:3.9 bash -c "pip install kafka-python && python3 migrate.py publish --brokers kafka-target:9092 --connector-name mm2-source-connector --input-file /workspace/offsets --connect-offset-topic connect-offsets"
```