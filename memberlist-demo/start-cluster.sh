#!/bin/bash

echo "Starting 3-node memberlist cluster..."

# 启动第一个节点（种子节点）
echo "Starting node-1 (seed node)..."
go run main.go 1 &
NODE1_PID=$!
echo "node-1 PID: $NODE1_PID"
sleep 2

# 启动第二个节点并加入集群
echo "Starting node-2..."
go run main.go 2 127.0.0.1:7947 &
NODE2_PID=$!
echo "node-2 PID: $NODE2_PID"
sleep 2

# 启动第三个节点并加入集群
echo "Starting node-3..."
go run main.go 3 127.0.0.1:7947 &
NODE3_PID=$!
echo "node-3 PID: $NODE3_PID"
sleep 2

echo "All nodes started. PIDs: $NODE1_PID, $NODE2_PID, $NODE3_PID"
echo "Press Ctrl+C to stop all nodes"

# 等待用户中断
trap "kill $NODE1_PID $NODE2_PID $NODE3_PID; exit" INT
wait
