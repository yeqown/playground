#!/bin/bash

# Usage: ./start-mm2.sh [mode, 1|2|3]
# mode: 1 (two instances with same config mm2.properties)
#       2 (two instances with different config mm2-01.properties and mm2-02.properties)
#       3 (kafka connect cluster mode)
# mode is required
if [ "$#" -ne 1 ] || { [ "$1" != "1" ] && [ "$1" != "2" ]; }; then
    echo "Usage: $0 [mode, 1|2]"
    echo "mode: 1 (two instances with same config mm2.properties)"
    echo "      2 (two instances with different config mm2-01.properties and mm2-02.properties)"
    echo "      3 (kafka connect cluster mode) - NOT supported now"
    exit 1
fi

echo "Starting Kafka MirrorMaker2 environment with Bitnami Kafka..."


MODE=$1

if [ $MODE -eq "1" ]; then
  echo "Mode 1: Starting two MirrorMaker2 instances with the same configuration (mm2.properties)"

  podman run -d --name mirrormaker2-01 --network kafka-net --platform linux/arm64 \
    -v $(pwd)/mm2.properties:/opt/bitnami/kafka/config/mm2.properties \
    bitnami/kafka:3.3.2-debian-12-r36 \
    /opt/bitnami/kafka/bin/connect-mirror-maker.sh /opt/bitnami/kafka/config/mm2.properties

  # podman run -d --name mirrormaker2-02 --network kafka-net --platform linux/arm64 \
  #   -v $(pwd)/mm2.properties:/opt/bitnami/kafka/config/mm2.properties \
  #   bitnami/kafka:3.3.2-debian-12-r36 \
  #   /opt/bitnami/kafka/bin/connect-mirror-maker.sh /opt/bitnami/kafka/config/mm2.properties
else
  echo "Mode 2: Starting two MirrorMaker2 instances with different configurations (mm2-01.properties and mm2-02.properties)"
  echo "Dedicated Mode DO not support multiple instances with same config, quit now."
  exit 1

  podman run -d --name mirrormaker2-01 --network kafka-net --platform linux/arm64 \
    -v $(pwd)/mm2-01.properties:/opt/bitnami/kafka/config/mm2.properties \
    bitnami/kafka:3.3.2-debian-12-r36 \
    /opt/bitnami/kafka/bin/connect-mirror-maker.sh /opt/bitnami/kafka/config/mm2.properties

  podman run -d --name mirrormaker2-02 --network kafka-net --platform linux/arm64 \
    -v $(pwd)/mm2-02.properties:/opt/bitnami/kafka/config/mm2.properties \
    bitnami/kafka:3.3.2-debian-12-r36 \
    /opt/bitnami/kafka/bin/connect-mirror-maker.sh /opt/bitnami/kafka/config/mm2.properties
fi

echo "Environment ready!"
echo "Run ./producer.sh and ./consumer.sh to produce and consume data"
