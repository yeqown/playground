## OpenFaaS on Kubernetes 实践

OpenFaaS 是一个 Serverless 平台，可以在 Kubernetes 上运行。这里对 OpenFaaS 的整体功能进行一个使用演示。

### 搭建

> 在本地的 minikube 搭建的 Kubernetes 环境上进行搭建。

#### 安装 OpenFaaS CE

参考：https://github.com/openfaas/faas-netes/blob/master/chart/openfaas/README.md

```bash
# 创建 Namespace
kubectl apply -f https://raw.githubusercontent.com/openfaas/faas-netes/master/namespaces.yml

# 添加 OpenFaaS Helm Chart 仓库
helm repo add openfaas https://openfaas.github.io/faas-netes/

# 安装 helm chart
helm repo update \
 && helm upgrade openfaas \
  --install openfaas/openfaas \
  --namespace openfaas

# 查看安装状态
kubectl -n openfaas get deployments -l "release=openfaas, app=openfaas"

# 查看 admin 密码
PASSWORD=$(kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode) && \
echo "OpenFaaS admin password: $PASSWORD"
```

### 访问 FaaS Portal

> 本地的 minikube 是通过虚拟机启动的，所以通过 port-forward 的方式访问 OpenFaaS Portal。

```bash
export OPENFAAS_URL=http://127.0.0.1:8080
kubectl port-forward -n openfaas svc/gateway 8080:8080 &
```

### 部署 Function

#### 通过 Portal 部署

在 portal 中快速的选择现有的模板，部署一个 env function。

```bash
curl -X GET http://127.0.0.1:8080/function/env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=env-54b4ffd889-kh5r7
fprocess=env
ENV_PORT_8080_TCP=tcp://10.97.3.140:8080
ENV_PORT_8080_TCP_PORT=8080
ENV_PORT_8080_TCP_ADDR=10.97.3.140
KUBERNETES_PORT=tcp://10.96.0.1:443
ENV_PORT=tcp://10.97.3.140:8080
ENV_PORT_8080_TCP_PROTO=tcp
KUBERNETES_SERVICE_HOST=10.96.0.1
KUBERNETES_SERVICE_PORT=443
KUBERNETES_SERVICE_PORT_HTTPS=443
KUBERNETES_PORT_443_TCP_PROTO=tcp
KUBERNETES_PORT_443_TCP_PORT=443
ENV_SERVICE_HOST=10.97.3.140
KUBERNETES_PORT_443_TCP_ADDR=10.96.0.1
ENV_SERVICE_PORT_HTTP=8080
KUBERNETES_PORT_443_TCP=tcp://10.96.0.1:443
ENV_SERVICE_PORT=8080
HOME=/home/app
Http_User_Agent=curl/8.4.0
Http_Accept=*/*
Http_Accept_Encoding=gzip
Http_X_Call_Id=ddbd1391-cc5d-4be5-9015-302d66268e2c
Http_X_Forwarded_For=127.0.0.1:45810
Http_X_Forwarded_Host=127.0.0.1:8080
Http_X_Start_Time=1712047032998171663
Http_Method=GET
Http_ContentLength=0
Http_Content_Length=0
Http_Path=/
Http_Host=10.244.0.10:8080
```
