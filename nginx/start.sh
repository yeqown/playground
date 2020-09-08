#!/bin/bash

docker run --rm \
    -p 8003:80 \
    -v /Users/med/projects/opensource/playground/nginx/conf.d:/etc/nginx/conf.d \
    -v /Users/med/projects/opensource/playground/nginx/nginx.conf:/etc/nginx/nginx.conf \
    --add-host=host.docker.internal.test0:192.168.65.2 \
    --add-host=host.docker.internal.test1:192.168.65.2 \
    nginx:1.19.2
