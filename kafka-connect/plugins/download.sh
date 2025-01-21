#!/bin/bash

# 创建插件目录
echo "创建插件目录..."
mkdir -p elasticsearch
mkdir -p debezium

# Download Elasticsearch Connector and dependencies
echo "下载 Elasticsearch Connector 及其依赖..."
wget https://d2p6pa21dvn84.cloudfront.net/api/plugins/confluentinc/kafka-connect-elasticsearch/versions/14.0.8/confluentinc-kafka-connect-elasticsearch-14.0.8.zip -P elasticsearch/
cd elasticsearch && unzip confluentinc-kafka-connect-elasticsearch-14.0.8.zip && rm confluentinc-kafka-connect-elasticsearch-14.0.8.zip && cd ..

# Check if download was successful
if [ $? -eq 0 ]; then
    echo "Elasticsearch Connector 及依赖下载成功"
else
    echo "Elasticsearch Connector 及依赖下载失败"
    exit 1
fi

# 下载 MySQL Connector
echo "下载 MySQL Connector..."
wget https://repo1.maven.org/maven2/io/debezium/debezium-connector-mysql/2.2.1.Final/debezium-connector-mysql-2.2.1.Final-plugin.tar.gz -P debezium/
cd debezium && tar -xzf debezium-connector-mysql-2.2.1.Final-plugin.tar.gz && rm debezium-connector-mysql-2.2.1.Final-plugin.tar.gz

# 检查下载是否成功
if [ $? -eq 0 ]; then
    echo "MySQL Connector 下载成功"
else
    echo "MySQL Connector 下载失败"
    exit 1
fi

# 设置权限
echo "设置权限..."
chmod -R 755 elasticsearch
chmod -R 755 debezium

echo "插件下载和配置完成"
