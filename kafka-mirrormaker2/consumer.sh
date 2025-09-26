#!/bin/bash

echo "Consuming mirrored data with detailed message tracking..."
podman exec -it kafka-target /opt/bitnami/kafka/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --from-beginning --property print.partition=true --property print.offset=true --property print.topic=true --formatter kafka.tools.DefaultMessageFormatter --property print.key=false --whitelist "source\.test-topic.*"
