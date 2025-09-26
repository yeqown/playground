#!/bin/bash

echo "Creating test-topic in source cluster..."
podman exec kafka-source /opt/bitnami/kafka/bin/kafka-topics.sh --create --topic test-topic --bootstrap-server localhost:9092 --partitions 9 --replication-factor 1

echo "Creating test-topic-second in source cluster..."
podman exec kafka-source /opt/bitnami/kafka/bin/kafka-topics.sh --create --topic test-topic-second --bootstrap-server localhost:9092 --partitions 9 --replication-factor 1

echo "Starting continuous data production with unique message IDs..."

counter1=10000000
counter2=20000000

while true; do
    partition1=$((counter1 % 9))
    msg_id1=$(printf "%08d" $counter1)
    # 发送到 test-topic，key=msg=msg_id1
    echo "${msg_id1}\t${msg_id1}" | podman exec -i kafka-source /opt/bitnami/kafka/bin/kafka-console-producer.sh \
        --topic test-topic \
        --bootstrap-server localhost:9092 \
        --property "parse.key=true" \
        --property "key.separator=\t"
    echo "Sent to test-topic: key=$msg_id1, msg=$msg_id1, partition=$partition1"

    partition2=$((counter2 % 9))
    msg_id2=$(printf "%08d" $counter2)
    # 发送到 test-topic-second，key=msg=msg_id2
    echo "${msg_id2}\t${msg_id2}" | podman exec -i kafka-source /opt/bitnami/kafka/bin/kafka-console-producer.sh \
        --topic test-topic-second \
        --bootstrap-server localhost:9092 \
        --property "parse.key=true" \
        --property "key.separator=\t"
    echo "Sent to test-topic-second: key=$msg_id2, msg=$msg_id2, partition=$partition2"

    counter1=$((counter1 + 1))
    counter2=$((counter2 + 1))
    sleep 2
done
