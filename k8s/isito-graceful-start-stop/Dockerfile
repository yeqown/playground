# 使用官方的Python镜像作为基础镜像
FROM python:3.8-slim-buster

# 设置工作目录
WORKDIR /app

# 将当前目录的内容复制到工作目录中
COPY . /app

# 安装依赖
RUN pip install --no-cache-dir -r requirements.txt

# 运行app.py当Docker启动时
CMD ["python", "app.py"]