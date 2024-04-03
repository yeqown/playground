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
    --set controller.image.tag=1.13.0-arm64 \
    --set dashboard.image.tag=1.13.0-arm64 \
    --set registry.pushPullUrl=docker.io/yeqown \
    --set registry.secretName=nuclio-registry-credentials \
    --namespace nuclio \
    nuclio/nuclio
```

### 访问 Nuclio

```bash
# 获取 nuclio dashboard 地址
kubectl port-forward -n nuclio $(kubectl get pod -n nuclio -l nuclio.io/app=dashboard -o jsonpath='{.items[0].metadata.name}') 8070:8070
```


### 部署一个简单的函数

```python
# 创建一个简单的函数 simple.py
import os

def my_entry_point(context, event):

	# use the logger, outputting the event body
	context.logger.info_with('Got invoked',
		trigger_kind=event.trigger.kind,
		event_body=event.body,
		some_env=os.environ.get('MY_ENV_VALUE'))

	# check if the event came from cron
	if event.trigger.kind == 'cron':

		# log something
		context.logger.info('Invoked from cron')

	else:

		# return a response
		return 'A string response'
```

```bash
# 部署 simple.py 函数
nuctl deploy my-function \
	--namespace nuclio \
	--path ./simple.py \
	--runtime python \
	--handler simple:my_entry_point \
	--http-trigger-service-type nodePort \
	--registry docker.io/yeqown \
	--run-registry docker.io/yeqown

# 查看部署的函数
$ ./nuctl-arm get functions
 NAMESPACE | NAME        | PROJECT | STATE | REPLICAS | NODE PORT
 nuclio    | my-function | default | ready | 1/1      | 32768


# 调用函数
$ curl -v  localhost:32768
> GET / HTTP/1.1
> Host: localhost:32768
> User-Agent: curl/7.79.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Server: nuclio
< Date: Wed, 03 Apr 2024 01:57:13 GMT
< Content-Type: text/plain
< Content-Length: 17
<
A string response$
```

