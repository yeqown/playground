# 把 python3 upstream.py 打包成一个镜像
# nerdctl.lima build -t apisix-debug-upstream:1.0.0 .
FROM python:3.7

WORKDIR /app

EXPOSE 3001

# upstream.py 为当前目录下的文件
COPY upstream.py .

CMD ["python3", "upstream.py"]
