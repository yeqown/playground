### Readme

#### docker 方式启动

启动：nerdctl.lima compose -f cluster1.compose.yaml up -d

停止：nerjson.lima compose -f cluster1.compose.yaml down

查看日志：nerdctl.lima compose -f cluster1.compose.yaml logs -f


#### 发送消息

```shell
# 发送消息
nats -s nats://localhost:4222 pub demo-subject "Hello, Nats!"
```
