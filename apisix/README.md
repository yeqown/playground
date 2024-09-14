### APISIX

在本地启动 APISIX 服务，然后编写简单的 lua 插件进行测试。

#### 1. 启动 APISIX 服务

```shell
nerdctl.lima compose up -d
```

验证 APISIX 服务是否启动成功：

```shell
curl http://127.0.0.1:9180/apisix/admin/routes
```


#### 新增一个路由

```shell
curl -i -X POST http://127.0.0.1:9180/apisix/admin/routes -d '
{
    "uri": "/hello",
    "methods": ["GET"],
    "plugins": {
        "ip2location": {}
    },
    "upstream": {
        "type": "roundrobin",
        "nodes": {
            "upstream:3001": 1
        }
    }
}'

测试路由是否生效：

```shell
curl http://127.0.0.1:9080/hello
```

#### 参考

- https://apisix.apache.org/zh/blog/2022/02/16/file-logger-api-gateway/
- https://apisix.apache.org/zh/docs/apisix/build-apisix-dev-environment-on-mac/