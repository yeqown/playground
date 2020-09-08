#!/bin/bash

docker run --rm \
    -p 8003:80 \
    -v /Users/med/projects/opensource/playground/nginx/conf.d/:/etc/nginx/conf.d/ \
    nginx:1.19.2