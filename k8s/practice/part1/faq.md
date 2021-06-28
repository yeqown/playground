## FAQ

- 学习 kubectl proxy 命令及其含义。回答如何通过 proxy 访问 kubernetes 集群？
    ```sh
    kubectl apply -f ./deployment_nginx.yaml
    # curl http://localhost:8088/version
    kubectl proxy --port=8080 --api-prefix=/k8s 
    # minikue Dashboard 的实现使用了 proxy
    ```
- 学习 kubectl port-forward 命令及其含义。回答如何通过 port-forward 访问应用？
    ```sh
    # 内部服务
    kubectl expose deployment nginx --port=8080 --target-port=80
    # 将本地的8080和service的端口转发
    # localPort:remotePort
    kubectl port-forward service/nginx 8080:8080
    # curl localhost:8080
    ```
- 修改 Pod label 使其与 Deployment 不相符，集群有什么变化？
    ```sh
    # 部署一个 deployment app=nginx
    kubectl apply -f ./deployment_nginx.yaml
    # 修改 label --overwirte 覆盖已经有的label
    kubectl label pod nginx-845d4d9dff-p6tbn app=nginx-2 --overwrite
    # 获取 pod，deployment 新启动了一个 pod, 老的POD 保持运行
    $ kgp
    NAME                     READY   STATUS    RESTARTS   AGE
    nginx-845d4d9dff-hqwp9   1/1     Running   0          11s
    nginx-845d4d9dff-p6tbn   1/1     Running   0          37m
    ```

    同时：
    1. 如果恢复被修改的pod的标签，那么新创建的pod会被销毁； 
    2. 如果删除被修改的pod，deployment并不会新启动一个POD；
    3. service 中被修改的POD会被替换为新启动的POD_IP。

- 进一步学习 kubectl rollout。回答如何通过 kubectl rollout 将应用回滚到指定版本？

    ```sh
    # 更新 image 版本
    kubectl set image deployment nginx nginx=nginx:1.21.0
    # 回滚
    kubectl rollout undo deployment/nginx
    # 回滚到指定版本，可以使用 更新镜像 的方式
    ```
- Pod LivenessProbe 实验中，检查方式采用的是 http 模式。回答如何使用 exec 进行健康检查？请写出 yaml 文件。
    
    查看资源的帮助文档
    ```sh
    kubectl explain deployment.spec.template.spec.containers.livenessProbe.exec
    ```

    ```yaml
    livenessProbe:
        # httpGet:
        #     path: /
        #     port: 80
        exec:
            command:
                - "curl"
                - "http://localhost:80"
        initialDelaySeconds: 5
        periodSeconds: 5
    ```

- 进一步学习 Pod Lifecycle。回答如何使用 PostStart Hook？请写出 yaml 文件。

    A: TODO

- 登录宿主机，使用 docker ps 查看 Pod，如何理解 docker ps 输出？

    ```sh
    minikube ssh
    # 列出所有的运行中的容器
    docker ps -a
    ```

    k8s是在虚拟机使用docker引擎部署的，因此k8s本身的功能组件在输出列表中；同时k8s CRI 默认使用 docker 因此部署的POD也作为在输出列表中。

- 学习使用 Secret，然后创建一个 Secret 并在 Pod 内访问。请写出 secret 和 pod 的 yaml 文件。 ConfigMap 实验中，我们采用文件加载的方式使用 ConfigMap。请写出利用环境变量加载 configmap 的例子。