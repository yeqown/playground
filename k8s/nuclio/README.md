## Nuclio 安装体验说明

### 安装 Nuclio

> 在本地的 minikube 搭建的 k8s 集群上安装 Nuclio


### 添加 dockerhub 密码

```bash
# 创建 nuclio-registry-secret
read -s mypassword

kubectl -n nuclio create secret docker-registry nuclio-registry-credentials \
    --docker-username yeqown@gmail.com \
    --docker-password $mypassword \
    --docker-server registry.hub.docker.com \
    --docker-email yeqown@gmail.com

unset mypassword
```

```bash
# 添加 nuclio helm repo
helm repo add nuclio https://nuclio.github.io/nuclio/charts

# 创建 nuclio namespace
kubectl create namespace nuclio

# 安装 nuclio helm chart
helm install nuclio \
    --set registry.pushPullUrl=docker.io/yeqown \
    --set controller.image.tag=1.13.0-arm64 \
    --set dashboard.image.tag=1.13.0-arm64 \
    --set registry.secretName=nuclio-registry-credentials \
    --namespace nuclio \
    nuclio/nuclio
```

### 访问 Nuclio

```bash
# 获取 nuclio dashboard 地址
kubectl port-forward -n nuclio $(kubectl get pod -n nuclio -l nuclio.io/app=dashboard -o jsonpath='{.items[0].metadata.name}') 8070:8070
```