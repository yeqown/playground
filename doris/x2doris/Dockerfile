FROM openjdk:8-jdk-slim

# 安装必要的工具
RUN apt-get update && apt-get install -y wget

# 工作目录
WORKDIR /opt/selectdb

COPY selectdb-x2doris-1.0.5_2.12-bin.tar.gz .

# 下载并解压 x2doris
RUN tar -xzvf selectdb-x2doris-1.0.5_2.12-bin.tar.gz
RUN rm selectdb-x2doris-1.0.5_2.12-bin.tar.gz
RUN ls -d */ && mv selectdb-x2doris-1.0.5-bin x2doris

# 暴露端口
EXPOSE 9091

# 启动命令
CMD ["./x2doris/bin/startup.sh"]

# 使用说明：
# 前提在 Dockerfile 目录已经有了 selectdb-x2doris-1.0.5_2.12-bin.tar.gz 文件
# docker build -t yeqown/x2doris:1.0.5 .
