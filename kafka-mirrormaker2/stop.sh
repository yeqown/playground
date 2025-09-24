#!/bin/bash

echo "Stopping Kafka MirrorMaker2 environment..."

podman stop mirrormaker2 kafka-target kafka-source zookeeper-target zookeeper-source 2>/dev/null
podman rm mirrormaker2 kafka-target kafka-source zookeeper-target zookeeper-source 2>/dev/null
podman network rm kafka-net 2>/dev/null

echo "Environment stopped and cleaned up."
