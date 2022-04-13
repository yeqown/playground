FROM alpine

RUN apk update && apk add curl nmap-ncat
EXPOSE 3306

# ncat --sh-exec "ncat REMOTE_HOST REMOTE_PORT" -l LOCAL_PORT  --keep-open