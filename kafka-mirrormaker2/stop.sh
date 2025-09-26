#!/bin/bash
# Usage: ./stop.sh [target]
# target: mm2 (only stop mirrormaker2 containers)
#         kafka (only stop kafka and kafka-ui containers)
#         all (stop all containers and remove network)
# target is required
if [ "$#" -ne 1 ] || { [ "$1" != "mm2" ] && [ "$1" != "kafka" ] && [ "$1" != "all" ]; }; then
    echo "Usage: $0 [mm2|kafka|all]"
    echo "target: mm2 (only stop mirrormaker2 containers)"
    echo "        kafka (only stop kafka and kafka-ui containers)"
    echo "        all (stop all containers and remove network)"
    exit 1
fi

target=$1

if [[ "$target" == "mm2" || "$target" == "all" ]]; then
  echo "Stopping MirrorMaker2 containers..."

  # 删除所有 mirrormaker2-xx 容器
  containers=$(podman ps -a --filter "name=mirrormaker2-" --format "{{.Names}}")
  if [ -n "$containers" ]; then
    podman stop -t 0 $containers
    podman rm -f $containers
  fi

  echo "MirrorMaker2 containers stopped."
fi

if [[ "$target" == "kafka" || "$target" == "all" ]]; then
  echo "Stopping Kafka and Kafka UI containers..."

  # 删除 kafka 相关容器
  podman stop -t 0 kafka-target kafka-source kafka-ui
  podman rm -f kafka-target kafka-source kafka-ui

  echo "Kafka and Kafka UI containers stopped."
fi

if [ "$target" == "all" ]; then
  echo "Removing network..."

  # 确保没有容器在使用网络
  remaining=$(podman network inspect kafka-net --format '{{ len .Containers }}' 2>/dev/null)
  if [ "$remaining" == "0" ]; then
    podman network rm kafka-net
    echo "Network removed."
  else
    podman network rm -f kafka-net
    echo "Network forcibly removed."
  fi
fi

echo "Environment stopped and cleaned up."
