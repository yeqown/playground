#!/bin/bash

echo "=== Memberlist Test Scenarios ==="
echo

echo "1. 启动单个节点："
echo "   go run main.go 1"
echo

echo "2. 启动并加入集群："
echo "   go run main.go 2 127.0.0.1:7947"
echo

echo "3. 启动完整集群："
echo "   ./start-cluster.sh"
echo

echo "4. 模拟节点故障（在另一个终端中杀死进程）："
echo "   kill <node-pid>"
echo

echo "5. 动态添加新节点："
echo "   go run main.go 4 127.0.0.1:7947"
echo

echo "6. 查看特定端口的进程："
echo "   lsof -i :7947"
echo

echo "开始演示..."
echo "按任意键继续..."
read

echo "启动节点1（种子节点）..."
go run main.go 1 &
NODE1_PID=$!
echo "节点1 PID: $NODE1_PID"
sleep 3

echo "启动节点2并加入集群..."
go run main.go 2 127.0.0.1:7947 &
NODE2_PID=$!
echo "节点2 PID: $NODE2_PID"
sleep 3

echo "启动节点3并加入集群..."
go run main.go 3 127.0.0.1:7947 &
NODE3_PID=$!
echo "节点3 PID: $NODE3_PID"
sleep 5

echo "集群运行中... 10秒后模拟节点2故障"
sleep 10

echo "模拟节点2故障（杀死进程）..."
kill -9 $NODE2_PID
echo "节点2已停止"
sleep 10

echo "添加新节点4..."
go run main.go 4 127.0.0.1:7947 &
NODE4_PID=$!
echo "节点4 PID: $NODE4_PID"
sleep 10

echo "清理所有节点..."
kill -9 $NODE1_PID $NODE3_PID $NODE4_PID 2>/dev/null
echo "演示完成"
