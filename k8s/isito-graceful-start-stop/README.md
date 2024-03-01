## Istio sidecar 和 应用程序的优雅启动和停止

### 问题描述

在使用 istio 时，应用程序的启动和停止过程中，会出现一些问题：

1. 应用启动时经常 RESTARTS 多次，原因是 istio sidecar 还未启动完成，应用就启动了。
2. 在应用停止时，如果应用内有优雅关闭的逻辑，enovy 在应用程序关闭之前就已经将连接断开了，导致应用程序内部逻辑执行时导致异常。


### 设计复现

1. 编写一个 python 服务，应用启动时会连接到一个外部的 redis 服务，启动后会每隔 1s 向 redis 发送一个 ping 消息。
    - 应用配置 readyz 探针，当应用启动完成后，会返回 200 状态码。
2. 在应用程序内部，有一个优雅关闭的逻辑，当应用程序接收到 SIGTERM 信号时，会等待 5s 后再关闭。
3. 部署应用到 k8s 集群中，观察应用启动和停止的过程。

### 相关脚本

#### docker 镜像打包和推送

```bash
# build image
nerdctl.lima build -t yeqown/istio-graceful-start-stop:v1 .

# test image to ensure it works
nerdctl.lima run --rm -it yeqown/istio-graceful-start-stop:v1 sh

# push to docker hub
nerdctl.lima push yeqown/istio-graceful-start-stop:v1
```

#### 部署 deployment

```bash
kubectl apply -f deployment.yaml
```

#### 修改 istio 配置

- https://www.zhaohuabing.com/istio-guide/docs/best-practice/startup-dependence/
- https://www.zhaohuabing.com/istio-guide/docs/best-practice/graceful-termination/
