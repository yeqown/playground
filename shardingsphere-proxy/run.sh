#!/bin/bash

# run.sh is a script to manage shardingsphere-proxy containers those
# are crreated by nerdctl.lima. It can be used to start, stop, restart
# and remove containers.

# Usage:
#   run.sh [COMMAND]
#   COMMAND: start, stop, restart, remove, ls

# Start containers command reference:
# nerdctl.lima run -d --name shardingsphere-proxy1 -p 3307:3307 -v $(pwd)/node1:/etc/shardingsphere-proxy apache/shardingsphere-proxy:5.4.0
# nerdctl.lima run -d --name shardingsphere-proxy2 -p 3308:3307 -v $(pwd)/node2:/etc/shardingsphere-proxy apache/shardingsphere-proxy:5.4.0
# nerdctl.lima run -d --name shardingsphere-proxy3 -p 3309:3307 -v $(pwd)/node3:/etc/shardingsphere-proxy apache/shardingsphere-proxy:5.4.0

# Variables
PROXY1=shardingsphere-proxy1
PROXY2=shardingsphere-proxy2
PROXY3=shardingsphere-proxy3

# Functions
start() {
    echo "Starting containers..."
    nerdctl.lima run -d --name $PROXY1 -p 3307:3307 -v $(pwd)/node1:/etc/shardingsphere-proxy apache/shardingsphere-proxy:5.4.0
    nerdctl.lima run -d --name $PROXY2 -p 3308:3307 -v $(pwd)/node2:/etc/shardingsphere-proxy apache/shardingsphere-proxy:5.4.0
    nerdctl.lima run -d --name $PROXY3 -p 3309:3307 -v $(pwd)/node3:/etc/shardingsphere-proxy apache/shardingsphere-proxy:5.4.0
}

stop() {
    echo "Stopping containers..."
    nerdctl.lima stop $PROXY1 $PROXY2 $PROXY3
}

restart() {
    nerdctl.lima restart $PROXY1 $PROXY2 $PROXY3
}

remove() {
    stop

    echo "Removing containers..."
    nerdctl.lima rm $PROXY1 $PROXY2 $PROXY3
}

list() {
    echo "Listing containers..."
    nerdctl.lima ps -a | grep -E "$PROXY1|$PROXY2|$PROXY3"
}

# Main
case $1 in
start)
    start
    ;;
stop)
    stop
    ;;
restart)
    restart
    ;;
remove)
    remove
    ;;
ls)
    list
    ;;
*)
    echo "Usage: $0 [start|stop|restart|remove|ls]"
    ;;
esac
