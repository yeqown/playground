## nats 超级集群实战

通过 nerdctl（docker）compose 部署 nats 超级集群，总共构建 2 个集群，每个集群 3 个节点，共 6 个节点。

集群1信息：
节点数量：3
cluster-id: cluster1
subject: cluster1-subject
port: 4222 到 4224
jetstream: enabled

集群2信息：
节点数量：3
cluster-id: cluster2
subject: cluster2-subject
port: 4332 到 4334
jetstream: enabled

### 部署

```bash
$ pwd
/path/to/nats/supercluster

$ nerdctl.lima compose -f cluster1.compose.yaml up -d
$ nerdctl.lima compose -f cluster2.compose.yaml up -d

# 观察集群状态
$ nats -s nats://admin:cluster1@localhost:4222 server ls
╭────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│                                                              Server Overview                                               │
├────────────┬──────────┬──────┬─────────┬─────┬───────┬───────┬────────┬─────┬────────┬───────┬───────┬──────┬────────┬─────┤
│ Name       │ Cluster  │ Host │ Version │ JS  │ Conns │ Subs  │ Routes │ GWs │ Mem    │ CPU % │ Cores │ Slow │ Uptime │ RTT │
├────────────┼──────────┼──────┼─────────┼─────┼───────┼───────┼────────┼─────┼────────┼───────┼───────┼──────┼────────┼─────┤
│ cluster1_1 │ cluster1 │ 0    │ 2.10.12 │ yes │ 1     │ 179   │      8 │   1 │ 18 MiB │ 0     │     4 │    0 │ 9m24s  │ 3ms │
│ cluster2_3 │ cluster2 │ 0    │ 2.10.12 │ yes │ 0     │ 169   │      8 │   1 │ 15 MiB │ 0     │     4 │    0 │ 1m44s  │ 3ms │
│ cluster1_2 │ cluster1 │ 0    │ 2.10.12 │ yes │ 0     │ 179   │      8 │   1 │ 17 MiB │ 0     │     4 │    0 │ 9m23s  │ 3ms │
│ cluster2_1 │ cluster2 │ 0    │ 2.10.12 │ yes │ 0     │ 169   │      8 │   1 │ 16 MiB │ 0     │     4 │    0 │ 1m44s  │ 3ms │
│ cluster1_3 │ cluster1 │ 0    │ 2.10.12 │ yes │ 0     │ 179   │      8 │   1 │ 16 MiB │ 0     │     4 │    0 │ 9m23s  │ 3ms │
│ cluster2_2 │ cluster2 │ 0    │ 2.10.12 │ yes │ 0     │ 169   │      8 │   1 │ 16 MiB │ 0     │     4 │    0 │ 1m44s  │ 3ms │
├────────────┼──────────┼──────┼─────────┼─────┼───────┼───────┼────────┼─────┼────────┼───────┼───────┼──────┼────────┼─────┤
│            │ 2        │ 6    │         │ 6   │ 1     │ 1,044 │        │     │ 97 MIB │       │       │    0 │        │     │
╰────────────┴──────────┴──────┴─────────┴─────┴───────┴───────┴────────┴─────┴────────┴───────┴───────┴──────┴────────┴─────╯

╭─────────────────────────────────────────────────────────────────────────────╮
│                                 Cluster Overview                            │
├──────────┬────────────┬───────────────────┬───────────────────┬─────────────┤
│ Cluster  │ Node Count │ Outgoing Gateways │ Incoming Gateways │ Connections │
├──────────┼────────────┼───────────────────┼───────────────────┼─────────────┤
│ cluster2 │          3 │                 3 │                 3 │           0 │
│ cluster1 │          3 │                 3 │                 3 │           1 │
├──────────┼────────────┼───────────────────┼───────────────────┼─────────────┤
│          │          6 │                 6 │                 6 │           1 │
╰──────────┴────────────┴───────────────────┴───────────────────┴─────────────╯
```

从集群信息可以看到，集群1和集群2个，各有3个出口和入口网关。cluster1 和 cluster2 将三个节点都指定为了 gateway，而每个 gateway 都只会连接到另一个集群的一个节点。

### 测试

当部署好一个集群时，使用如下的 nats 命令进行测试。如下结果代表集群工作正常。

```bash
$ nats sub -s nats://cluster1:cluster1@localhost:4222 cluster1-subject
10:26:54 Subscribing on cluster1-subject
[#1] Received on "cluster1-subject"
hello

$ nats pub -s nats://cluster1:cluster1@localhost:4222 cluster1-subject "hello"
10:27:02 Published 5 bytes to "cluster1-subject"
```

#### 测试超级集群的订阅场景

1. 乐观模式

乐观模式是指，在 cluster1 中投递一条消息后，cluster1 会检查其是否是 cluster2 不感兴趣的主题，如果是，则不会将消息投递到 cluster2 中。否则，会将消息投递到 cluster2 中。

cluster2 对于主题的感兴趣与否是通过，判断当前集群中是否有订阅者来决定的。如果一条消息被投递到 cluster2 之后，没有订阅者订阅该主题，cluster2 会发送一个网关协议给 cluster1，告诉 cluster1 该主题不感兴趣。当 cluster2 上有订阅者订阅该主题时，cluster2 会发送一条网关协议给 cluster1 “我对这个主题感兴趣”。

```bash
+-----------+     Is the message topic     +-----------+
|           |  not cluster2's cup of tea?  |           |
|  Cluster1 | ---------------------------> | Cluster2  |
+-----------+                              +-----------+
      |                                       |
      | YES, not interesting                  | NO, interesting or not sure
      |                                       |
      v                                       v
+-----------------------------------+   +-----------------------------------+
| The message is not delivered to   |   | The message is delivered to       |
| Cluster2 because the topic is not |   | Cluster2. If there are no         |
| interesting to Cluster2.          |   | subscribers for the topic in      |
+-----------------------------------+   | Cluster2, Cluster2 sends a        |
                                        | gateway protocol to Cluster1,     |
                                        | indicating that the topic is not  |
                                        | interesting.                      |
                                        +-----------------------------------+
```


```bash
# 在 cluster2 中订阅消息
$ nats sub -s nats://admin:cluster2@localhost:4333 foo

# 在 cluster1 中发布消息
$ nats pub -s nats://admin:cluster1@localhost:4222 foo "hello"
```

OUTPUT:

```bash
```

2. 兴趣模式

TODO://


3. 队列模式

服务器优先为本地队列的订阅者投递消息，如果本地没有订阅者，服务器会将消息投递到远程集群中。选择 RTT 最低的节点进行投递。

```bash
# 在 cluster1 中订阅消息, queue 模式, queue = cross-cluster-queue
$ nats -s nats://admin:cluster1@localhost:4222 --queue=cross-cluster-queue sub foo
# 在 cluster2 中订阅消息, queue 模式, queue = cross-cluster-queue
$ nats -s nats://admin:cluster2@localhost:4333 --queue=cross-cluster-queue sub foo
```

```bash
# 在 cluster1 中发布消息
$ nats pub -s nats://admin:cluster1@localhost:4222 foo "hello"
```

OUTPUT:

```bash
# 当 cluster1 中有订阅者时, cluster2 中的订阅者不会收到消息,
# 仅当 cluster1 中没有订阅者时, cluster2 中的订阅者才会收到消息
```

