## 关于 istio 使用过程中遇到的 sidecar 连接偶尔断开的问题

### 问题描述

应用在未接入 istio 时，正常运行，接入 istio 后，偶尔出现连接断开的情况。日志中常出现 “invalid connection” 错误，常见于集群外使用的中间件连接。
如： memcached、redis、mysql、mongodb 等。

### 问题原因

猜测: istio sidecar 与集群外的中间件连接，由于连接空闲时间过长，导致连接断开。

- https://github.com/istio/istio/issues/24387
- https://www.envoyproxy.io/docs/envoy/latest/faq/configuration/timeouts#tcp
- https://stackoverflow.com/questions/63843610/istio-proxy-closing-long-running-tcp-connection-after-1-hour

### 设计复现

1. 集群外启动一个中间件，如 redis
2. 在集群内启动两个POD，一个接入 istio，一个未接入 istio, 两个POD 都连接到集群外的 redis。
3. 配置调整和观察
    - 通过 envoy filter 设置连接超时时间为 30s, 观察连接是否在 30s 后断开。
    - 观察连接是否在一段时间（1h）后断开, 1h 是 istio 默认的连接空闲时间。
    - 通过 envoy filter 设置连接超时时间为 24h, 观察连接是否在 24h 后断开。

可以预先设置 envoy 的连接超时时间为 30s，观察连接是否在 30s 后断开。


### 相关脚本

#### docker 镜像打包和推送

```bas
# build image
nerdctl.lima build -t yeqown/istio-idle-timeout:v1 .

# test image to ensure it works
nerdctl.lima run --rm -it yeqown/istio-idle-timeout:v1 sh

# push to docker hub
nerdctl.lima push yeqown/istio-idle-timeout:v1
```

#### 部署 deployment

1. envoyFilter 应用 envoyfilter.yaml

```bash
kubectl apply -f envoyfilter-$TIMEOUT.yaml
```

2. 部署应用

```bash
kubectl create ns istio-idle-timeout && kubectl label ns istio-idle-timeout istio-injection=enabled

# apply deployment without istio
kubectl apply -f deployment.yaml -n default
# apply deployment with istio
kubectl apply -f deployment.yaml -n istio-idle-timeout
```

#### 连接外部端口

集群外部已经启动了一个 TCP 服务, 在容器内启动 telnet 连接外部服务

> macOS 使用 minikube 在本地上启动了一个 redis 服务，运行在 6379 端口
> 可以通过 host.minikube.internal:6379 连接到 redis 服务

```bash
# 连接外部服务
telnet host.minikube.internal 6379
```

## 复现结果

1. 没有 istio 的情况下，连接不会断开。
2. 有 istio 的情况下，连接会在 10s 后断开, 更新配置为 24h 后，并不会在 10s 内断开。如下输出是在 istio 中执行，且 ilde-timeout 为 10s 的情况下。

```bash
/ # time telnet mysql 6379
Connected to mysql
Connection closed by foreign host
Command exited with non-zero status 1
real	0m 10.00s
user	0m 0.00s
sys	0m 0.00s
```

3. 在没有设置 envoy filter 的情况下，连接会在 1h 后断开。

```bash
/ # time telnet mysql 6379
Connected to mysql
Connection closed by foreign host
Command exited with non-zero status 1
real	1h 0m 00s
user	0m 0.00s
sys	0m 0.00s
```

### 其他

1. 或许相关的日志输出

```bash
# 调整日志级别
curl -X POST http://127.0.0.1:15000/logging?level=debug
```

通过日志发现，连接在 10s 后断开时伴随着 `invoking idle callbacks` 的日志输出。如下：(前后一共测试了4次)

```bash
➜  istio-1.19.3 klf istio-idle-timeout-demo-5b4894b67d-tdpb5 -n istio-idle-timeout -c istio-proxy | grep "invoking idle callbacks"
2024-03-01T07:22:21.481180Z	debug	envoy pool external/envoy/source/common/conn_pool/conn_pool_base.cc:454	invoking idle callbacks - is_draining_for_deletion_=false	thread=22
2024-03-01T07:22:21.481310Z	debug	envoy pool external/envoy/source/common/conn_pool/conn_pool_base.cc:454	invoking idle callbacks - is_draining_for_deletion_=false	thread=22
2024-03-01T07:24:03.947896Z	debug	envoy pool external/envoy/source/common/conn_pool/conn_pool_base.cc:454	invoking idle callbacks - is_draining_for_deletion_=false	thread=23
2024-03-01T07:24:03.947903Z	debug	envoy pool external/envoy/source/common/conn_pool/conn_pool_base.cc:454	invoking idle callbacks - is_draining_for_deletion_=false	thread=23


2024-03-01T07:25:04.884340Z	debug	envoy pool external/envoy/source/common/conn_pool/conn_pool_base.cc:454	invoking idle callbacks - is_draining_for_deletion_=false	thread=22
2024-03-01T07:25:04.884394Z	debug	envoy pool external/envoy/source/common/conn_pool/conn_pool_base.cc:454	invoking idle callbacks - is_draining_for_deletion_=false	thread=22
2024-03-01T07:25:17.629907Z	debug	envoy pool external/envoy/source/common/conn_pool/conn_pool_base.cc:454	invoking idle callbacks - is_draining_for_deletion_=false	thread=22
2024-03-01T07:25:17.629912Z	debug	envoy pool external/envoy/source/common/conn_pool/conn_pool_base.cc:454	invoking idle callbacks - is_draining_for_deletion_=false	thread=22
```

2. 应用 envoy filter 中的超时之后，已经建立的连接并不会应用新的超时时间，只有新的连接才会应用新的超时时间？

3. envoyFilter 的删除方法

```bash
kubectl apply -f envoyfilter-remove.yaml
```