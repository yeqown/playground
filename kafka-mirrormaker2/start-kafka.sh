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
  -e KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka-source:9092 \
  bitnami/kafka:3.3.2-debian-12-r36

# Start target cluster  
podman run -d --name kafka-target --network kafka-net -p 9094:9092 --platform linux/arm64 \
  -e KAFKA_CFG_NODE_ID=2 \
  -e KAFKA_CFG_PROCESS_ROLES=controller,broker \
  -e KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 \
  -e KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT \
  -e KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=2@kafka-target:9093 \
  -e KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka-target:9092 \
  bitnami/kafka:3.3.2-debian-12-r36

echo "Waiting for Kafka clusters to be ready..."
sleep 15

# Start Kafka UI
podman run -d --name kafka-ui --network kafka-net -p 8080:8080 --platform linux/arm64 \
  -e KAFKA_CLUSTERS_0_NAME=source \
  -e KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS=kafka-source:9092 \
  -e KAFKA_CLUSTERS_1_NAME=target \
  -e KAFKA_CLUSTERS_1_BOOTSTRAPSERVERS=kafka-target:9092 \
  provectuslabs/kafka-ui:latest

echo "Environment ready!"
echo "Run ./start-mm2.sh --mode [mode] to start MirrorMaker2"
echo "Access Kafka UI at http://localhost:8080"
