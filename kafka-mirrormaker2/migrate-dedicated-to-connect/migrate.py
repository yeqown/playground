#!/usr/bin/env python3

# 这个脚本用于将 Kafka MirrorMaker2 (MM2) 的 offset 从专用模式迁移到 Connect 模式
# 专用模式使用一个专用的 topic 来存储 offset，而 Connect 模式则使用 Kafka Connect 的内置 offset 存储机制。
#
# 专用模式 topic 的格式通常是 "mm2-offsets-<source-cluster-alias>.internal"。
# 其中的消息 key 和 value 都是 JSON 格式，如下：
# Key: ["MirrorSourceConnector", {"cluster": "<source-cluster-alias>", "partition": <partition>, "topic": "<topic-name>"}]
# Value: {"offset": <offset>}
#
# Connect topic 的格式是 "connect-offsets"（默认）。
# 其中的消息 key 和 value 也是 JSON 格式，如下：
# Key: ["<connector-name>", {"cluster": "<source-cluster-alias>", "partition": <partition>, "topic": "<topic-name>"}]
# Value: {"offset": <offset>}
#
# 脚本的主要步骤：
# 1. 从专用模式的 offset topic 中读取所有消息，解析出 topic、partition 和 offset 信息，只保留每个 topic-partition 的最新 offset
# 2. 将这些信息转换为 Connect 模式所需的格式
# 3. 将转换后的消息发布到 Connect 的 offset topic 中
# 
#
# Usage:
# 提取专用模式的 offsets 并保存到文件
# python migrate.py extract --source-brokers <source-brokers> --dedicated-offset-topic <dedicated-offset-topic> --output-file <output-file>
# 转换为 Connect 格式并发布到 Connect offset topic
# python migrate.py publish --target-brokers <target-brokers> --connect-offset-topic <connect-offset-topic> --connector-name <connector-name> --input-file <output-file>
#
# 注意：确保在运行此脚本前，Kafka 集群和 Kafka Connect 集群均已启动，并且 Kafka Connect 集群已配置好 Connect offset topic。

import json
import sys
import argparse
from kafka import KafkaConsumer, KafkaProducer, TopicPartition
from kafka.errors import KafkaError

def extract_dedicated_offsets(brokers, dedicated_offset_topic, output_file):
    """从专用模式的 offset topic 中提取 offset 信息并保存到文件"""
    
    print(f"从 {dedicated_offset_topic} 提取 offsets, 创建 consumer 连接到 {brokers}")
    
    consumer = KafkaConsumer(
        dedicated_offset_topic,
        bootstrap_servers=brokers,
        auto_offset_reset='earliest',
        enable_auto_commit=False,
        consumer_timeout_ms=10000,  # 10秒内没有新消息就退出
        key_deserializer=lambda x: json.loads(x.decode('utf-8')) if x else None,
        value_deserializer=lambda x: json.loads(x.decode('utf-8')) if x else None
    )
    
    print("开始消费消息...")
    
    # 只保留每个 topic-partition 的最新 offset
    offsets = {}
    
    for message in consumer:
        if message.key and message.value:
            if message.key[0] != "MirrorSourceConnector":
                # 只处理 MirrorSourceConnector 的消息
                continue
            
            # 专用模式格式: Key[1] = {"cluster": "...", "partition": ..., "topic": "..."}
            partition_info = message.key[1]
            cluster = partition_info["cluster"]
            topic = partition_info["topic"]
            partition = partition_info["partition"]
            offset = message.value["offset"]
            
            # 过滤掉内部 topic (heartbeats, checkpoints 等)
            if topic.endswith('heartbeats') or 'checkpoints' in topic or 'internal' in topic:
                print(f"跳过内部 topic: {topic}")
                continue
            
            # 使用 topic-partition 作为唯一标识，保留最新的 offset
            key = f"{cluster}-{topic}-{partition}"
            offsets[key] = {
                "cluster": cluster,
                "topic": topic,
                "partition": partition,
                "offset": offset
            }
            print(f"提取 offset: {key} -> {offset}")
        else:
            print(f"跳过其他消息: {message.key[0] if message.key else 'None'}")
    
    consumer.close()

    # 保存到文件
    with open(output_file, 'w') as f:
        json.dump(offsets, f, indent=2)
    
    print(f"提取了 {len(offsets)} 个 offset 记录到 {output_file}")

def publish_connect_offsets(brokers, connect_offset_topic, connector_name, input_file):
    """从文件读取 offset 信息，转换为 Connect 格式并发布到 Connect offset topic"""
    # 从文件读取 offset 信息
    with open(input_file, 'r') as f:
        offsets = json.load(f)
    
    producer = KafkaProducer(
        bootstrap_servers=brokers,
        key_serializer=lambda x: json.dumps(x).encode('utf-8'),
        value_serializer=lambda x: json.dumps(x).encode('utf-8')
    )
    
    success_count = 0
    
    for offset_data in offsets.values():
        # 转换为 Connect 格式
        connect_key = [
            connector_name,
            {
                "cluster": offset_data["cluster"],
                "partition": offset_data["partition"],
                "topic": offset_data["topic"]
            }
        ]
        
        connect_value = {
            "offset": offset_data["offset"]
        }
        
        try:
            # print(f"发布 offset: {connect_key} -> {connect_value} 到 {connect_offset_topic}")
            future = producer.send(connect_offset_topic, key=connect_key, value=connect_value)
            future.get(timeout=10)
            success_count += 1
        except KafkaError as e:
            print(f"发送 offset 失败: {e}", file=sys.stderr)
    
    producer.flush()
    producer.close()
    
    print(f"成功发布了 {success_count}/{len(offsets)} 个 offset 记录")

def main():
    parser = argparse.ArgumentParser(description='将 MM2 offset 从专用模式迁移到 Connect 模式')
    subparsers = parser.add_subparsers(dest='command', help='可用命令')
    
    # extract 子命令
    extract_parser = subparsers.add_parser('extract', help='提取专用模式的 offsets')
    extract_parser.add_argument('--brokers', required=True, help='Kafka brokers')
    extract_parser.add_argument('--dedicated-offset-topic', required=True, help='专用模式 offset topic')
    extract_parser.add_argument('--output-file', required=True, help='输出文件')
    
    # publish 子命令
    publish_parser = subparsers.add_parser('publish', help='发布 Connect 格式的 offsets')
    publish_parser.add_argument('--brokers', required=True, help='Kafka brokers')
    publish_parser.add_argument('--connect-offset-topic', default='connect-offsets', help='Connect offset topic')
    publish_parser.add_argument('--connector-name', required=True, help='Connect connector 名称')
    publish_parser.add_argument('--input-file', required=True, help='输入文件')
    
    args = parser.parse_args()
    
    if args.command == 'extract':
        extract_dedicated_offsets(
            args.brokers.split(','),
            args.dedicated_offset_topic,
            args.output_file
        )
    elif args.command == 'publish':
        publish_connect_offsets(
            args.brokers.split(','),
            args.connect_offset_topic,
            args.connector_name,
            args.input_file
        )
    else:
        parser.print_help()

if __name__ == "__main__":
    main()
