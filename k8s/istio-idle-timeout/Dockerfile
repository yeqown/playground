# 编写一个包含 telnet 的 Dockerfile 镜像 arm64v8/alpine-telnet

FROM arm64v8/alpine:3.12

RUN apk add --no-cache busybox-extras

# pause command
CMD ["sh", "-c", "trap : TERM INT; sleep infinity & wait"]