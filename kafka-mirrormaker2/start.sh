#!/bin/bash

echo "Starting Kafka MirrorMaker2 environment with Bitnami Kafka..."

# Create network
podman network create kafka-net 2>/dev/null || true

# Start source cluster
podman run -d --name kafka-source --network kafka-net -p 9092:9092 --platform linux/arm64 \
  -e KAFKA_CFG_NODE_ID=1 \
  -e KAFKA_CFG_PROCESS_ROLES=controller,broker \
  -e KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 \
  -e KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT \
  -e KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@kafka-source:9093 \
  -e KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  bitnami/kafka:3.3.2-debian-12-r36

# Start target cluster  
podman run -d --name kafka-target --network kafka-net -p 9094:9092 --platform linux/arm64 \
  -e KAFKA_CFG_NODE_ID=2 \
  -e KAFKA_CFG_PROCESS_ROLES=controller,broker \
  -e KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 \
  -e KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT \
  -e KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=2@kafka-target:9093 \
  -e KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9094 \
  bitnami/kafka:3.3.2-debian-12-r36

echo "Waiting for Kafka clusters to be ready..."
sleep 30

# Start MirrorMaker2
podman run -d --name mirrormaker2 --network kafka-net --platform linux/arm64 \
  -v $(pwd)/mm2.properties:/opt/bitnami/kafka/config/mm2.properties \
  bitnami/kafka:3.3.2-debian-12-r36 \
  /opt/bitnami/kafka/bin/connect-mirror-maker.sh /opt/bitnami/kafka/config/mm2.properties

echo "Environment ready!"
echo "Run ./producer.sh to start producing data"
