#!/bin/bash

# run.sh is a script to manage shardingsphere-proxy containers those
# are created by nerdctl.lima. It can be used to start, stop, restart
# and remove containers.

# Usage:
#   run.sh [COMMAND]
#   COMMAND: start, stop, restart, remove, ls

# Functions
start() {
    echo "Starting services..."
    nerdctl.lima compose up -d
}

stop() {
    echo "Stopping services..."
    nerdctl.lima compose down
}

remove() {
    echo "Removing services..."
    nerdctl.lima compose down -v
}

restart() {
    echo "Restarting services..."
    nerdctl.lima compose down
    nerdctl.lima compose up -d
}

ls() {
    echo "Listing services..."
    nerdctl.lima compose ps
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
    ls
    ;;
*)
    echo "Usage: $0 [start|stop|restart|remove|ls]"
    ;;
esac