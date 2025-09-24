#!/bin/bash

echo "Creating test-topic in source cluster..."
podman exec kafka-source /opt/bitnami/kafka/bin/kafka-topics.sh --create --topic test-topic --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1

echo "Starting continuous data production to test-topic..."
counter=1
while true; do
    message="Message $counter - $(date)"
    echo "$message" | podman exec -i kafka-source /opt/bitnami/kafka/bin/kafka-console-producer.sh --topic test-topic --bootstrap-server localhost:9092
    echo "Sent: $message"
    counter=$((counter + 1))
    sleep 2
done
