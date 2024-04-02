## Nuclio 安装体验说明

### 安装 Nuclio

> 在本地的 minikube 搭建的 k8s 集群上安装 Nuclio

```bash
# 添加 nuclio helm repo
helm repo add nuclio https://nuclio.github.io/nuclio/charts

# 创建 nuclio namespace
kubectl create namespace nuclio

# 安装 nuclio helm chart
helm install nuclio \
    --set registry.pushPullUrl=localhost:5000 \
    --set controller.image.tag=1.13.0-arm64 \
    --set dashboard.image.tag=1.13.0-arm64 \
    --namespace nuclio \
    nuclio/nuclio
```
