## DORIS

在 K8S 中部署存算分离 Doris 集群。这里记录的是在 ubuntu 22.04 上配合 minikube 搭建的 k8s 集群，搭建 doris 存算分离集群的过程。

os=Ubuntu 24.04.1 LTS x86_64 
kernel=6.8.0-52-generic
minikube=v1.35.0
kubernetes=v1.32.0
doris=3.0.3

### 1. 安装 FoundationDB

> https://doris.apache.org/zh-CN/docs/3.0/install/deploy-on-kubernetes/separating-storage-compute/install-fdb

#### 安装 CRD 部署 FoundationDB 相关资源定义

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/FoundationDB/fdb-kubernetes-operator/main/config/crd/bases/apps.foundationdb.org_foundationdbclusters.yaml
    kubectl apply -f https://raw.githubusercontent.com/FoundationDB/fdb-kubernetes-operator/main/config/crd/bases/apps.foundationdb.org_foundationdbbackups.yaml
    kubectl apply -f https://raw.githubusercontent.com/FoundationDB/fdb-kubernetes-operator/main/config/crd/bases/apps.foundationdb.org_foundationdbrestores.yaml
    ```

#### 部署 fdb-kubernetes-operator 服务

```bash
wget https://raw.githubusercontent.com/apache/doris-operator/master/config/operator/fdb-operator.yaml

kubectl apply -f fdb-operator.yaml
```

#### 部署 FDB 集群

```bash
wget https://raw.githubusercontent.com/foundationdb/fdb-kubernetes-operator/main/config/samples/cluster.yaml -O fdb-cluster.yaml

kubectl apply -f fdb-cluster.yaml

# 查看集群状态
kubectl get fdb

# 预期输出（启动需要时间，需要等待几分钟）
bint@Z590-D:~/doris$ kubectl get fdb
NAME           GENERATION   RECONCILED   AVAILABLE   FULLREPLICATION   VERSION   AGE
test-cluster   1            1            true        true              7.1.26    3m30s
```

### 部署 Doris Operator

#### 下发资源定义

```bash
kubectl create -f https://raw.githubusercontent.com/apache/doris-operator/master/config/crd/bases/crds.yaml
```

#### 部署 Doris Operator

```bash
wget https://raw.githubusercontent.com/apache/doris-operator/master/config/operator/disaggregated-operator.yaml -O disaggregated-operator.yaml

kubectl apply -f disaggregated-operator.yaml

# 查看部署状态
kubectl get pod -n doris

# 预期输出
bint@Z590-D:~/doris$ kubectl -n doris get pods
NAME                              READY   STATUS    RESTARTS   AGE
doris-operator-5fd65d8d69-rgqlk   1/1     Running   0          79s
```

### 部署存算分离集群

#### 下载示例配置文件

```bash
curl -O https://raw.githubusercontent.com/apache/doris-operator/master/doc/examples/disaggregated/cluster/ddc-sample.yaml
```

#### 配置 FDB 集群地址

> 对于 ddc-sample.yaml 配置进行调整配置

```bash
暂不用处理
```

#### 配置存算分离集群 

> 对于 ddc-sample.yaml 配置进行调整配置。这三个都需要分别配置 ConfigMap 并修改集群中的配置挂载。

    - 配置元数据 https://doris.apache.org/zh-CN/docs/3.0/install/deploy-on-kubernetes/separating-storage-compute/config-ms
    - 配置 FE 集群 https://doris.apache.org/zh-CN/docs/3.0/install/deploy-on-kubernetes/separating-storage-compute/config-fe
    - 配置计算资源组 https://doris.apache.org/zh-CN/docs/3.0/install/deploy-on-kubernetes/separating-storage-compute/config-cg

```yaml
apiVersion: v1
data:
  doris_cloud.conf: |
    # // meta_service
    brpc_listen_port = 5000
    brpc_num_threads = -1
    brpc_idle_timeout_sec = 30
    http_token = greedisgood9999

    # // doris txn config
    label_keep_max_second = 259200
    expired_txn_scan_key_nums = 1000

    # // logging
    log_dir = ./log/
    # info warn error
    log_level = info
    log_size_mb = 1024
    log_filenum_quota = 10
    log_immediate_flush = false
    log_verbose_modules = *

    # //max stage num
    max_num_stages = 40
kind: ConfigMap
metadata:
  name: doris-metaservice
  namespace: default
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: fe-configmap
  namespace: default
  labels:
    app.kubernetes.io/component: fe
data:
  fe.conf: |
    CUR_DATE=`date +%Y%m%d-%H%M%S`
    # Log dir
    LOG_DIR = ${DORIS_HOME}/log
    # For jdk 17, this JAVA_OPTS will be used as default JVM options
    JAVA_OPTS_FOR_JDK_17="-Djavax.security.auth.useSubjectCredsOnly=false -Xmx8192m -Xms8192m -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=$LOG_DIR -Xlog:gc*:$LOG_DIR/fe.gc.log.$CUR_DATE:time,uptime:filecount=10,filesize=50M --add-opens=java.base/java.nio=ALL-UNNAMED --add-opens java.base/jdk.internal.ref=ALL-UNNAMED"
    # INFO, WARN, ERROR, FATAL
    sys_log_level = INFO
    # NORMAL, BRIEF, ASYNC
    sys_log_mode = NORMAL
    # Default dirs to put jdbc drivers,default value is ${DORIS_HOME}/jdbc_drivers
    # jdbc_drivers_dir = ${DORIS_HOME}/jdbc_drivers
    http_port = 8030
    rpc_port = 9020
    query_port = 9030
    edit_log_port = 9010
    enable_fqdn_mode=true
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: be-configmap
  labels:
    app.kubernetes.io/component: be
data:
  be.conf: |
    # For jdk 17, this JAVA_OPTS will be used as default JVM options
    JAVA_OPTS_FOR_JDK_17="-Xmx1024m -DlogPath=$LOG_DIR/jni.log -Xlog:gc*:$LOG_DIR/be.gc.log.$CUR_DATE:time,uptime:filecount=10,filesize=50M -Djavax.security.auth.useSubjectCredsOnly=false -Dsun.security.krb5.debug=true -Dsun.java.command=DorisBE -XX:-CriticalJNINatives -XX:+IgnoreUnrecognizedVMOptions --add-opens=java.base/java.lang=ALL-UNNAMED --add-opens=java.base/java.lang.invoke=ALL-UNNAMED --add-opens=java.base/java.lang.reflect=ALL-UNNAMED --add-opens=java.base/java.io=ALL-UNNAMED --add-opens=java.base/java.net=ALL-UNNAMED --add-opens=java.base/java.nio=ALL-UNNAMED --add-opens=java.base/java.util=ALL-UNNAMED --add-opens=java.base/java.util.concurrent=ALL-UNNAMED --add-opens=java.base/java.util.concurrent.atomic=ALL-UNNAMED --add-opens=java.base/sun.nio.ch=ALL-UNNAMED --add-opens=java.base/sun.nio.cs=ALL-UNNAMED --add-opens=java.base/sun.security.action=ALL-UNNAMED --add-opens=java.base/sun.util.calendar=ALL-UNNAMED --add-opens=java.security.jgss/sun.security.krb5=ALL-UNNAMED --add-opens=java.management/sun.management=ALL-UNNAMED"
    file_cache_path = [{"path":"/mnt/disk1/doris_cloud/file_cache","total_size":107374182400,"query_limit":107374182400}]
    deploy_mode = cloud
```

#### 启动

```bash
# 部署集群
kubectl apply -f ddc-sample.yaml

# 查看集群状态
kubectl get ddc

# 预期输出
bint@Z590-D:~/doris$ kubectl get ddc
NAME                         CLUSTERHEALTH   MSPHASE   FEPHASE   CGCOUNT   CGAVAILABLECOUNT   CGFULLAVAILABLECOUNT
test-disaggregated-cluster   green           Ready     Ready     2         2                  2
```

#### 创建远程存储后端

```bash
kubectl get svc

# test-disaggregated-cluster-fe            ClusterIP   10.104.79.145    <none>        8030/TCP,9020/TCP,9030/TCP,9010/TCP   35m


# MySQL 客户端启动
kubectl run mysql-client --image=mysql:5.7 -it --rm --restart=Never -- /bin/bash

# 连接 Doris FE
mysql -uroot -P9030 -h test-disaggregated-cluster-fe

# MYSQL 命令执行：S3 Storage Vault 
CREATE STORAGE VAULT IF NOT EXISTS s3_vault
    PROPERTIES (
        "type"="S3",
        "s3.endpoint" = "oss-cn-beijing.aliyuncs.com", 
        "s3.region" = "bj",
        "s3.bucket" = "bucket",
        "s3.root.path" = "big/data/prefix", 
        "s3.access_key" = "your-ak",
        "s3.secret_key" = "your-sk",
        "provider" = "OSS" 
    );

# MYSQL 命令执行：设置默认数据后端
SET s3_vault AS DEFAULT STORAGE VAULT;
```

### 测试集群

```bash
# 创建一个临时的 MySQL 客户端 Pod
kubectl run mysql-client --image=mysql:5.7 -it --rm --restart=Never -- mysql -h test-disaggregated-cluster-fe -P9030 -uroot
```

#### 测试基本操作

```sql
# 执行一些查询
-- 查看 FE 节点状态
SHOW FRONTENDS;

-- 查看 BE 节点状态
SHOW BACKENDS;

-- 创建测试数据库
CREATE DATABASE test_db;
USE test_db;


-- 创建测试表
CREATE TABLE test_tbl
(
    id INT,
    name VARCHAR(50),
    score DECIMAL(10,2)
)
UNIQUE KEY(id)
DISTRIBUTED BY HASH(id) BUCKETS 3;

-- 插入测试数据
INSERT INTO test_tbl VALUES 
(1, 'Tom', 89.5),
(2, 'Jerry', 92.0),
(3, 'Jack', 85.5);

-- 查询数据
SELECT * FROM test_tbl;
```

#### 验证存算分离

```sql
-- 查看存储后端状态
SHOW STORAGE VAULTS;

-- 查看表的存储位置
SHOW CREATE TABLE test_db.test_tbl;
```

### 部署期间遇到的问题和解决


🚀这里在部署的时候遇到了一些情况：

1. Be 无法启动，输入日志为：

```bash
Defaulted container "compute" out of: compute, default-init (init)
[Fri Feb 28 02:56:24 UTC 2025] [info] Process conf file be.conf ...
/opt/apache-doris/be_disaggregated_entrypoint.sh: line 73: /opt/apache-doris/be/conf/: Is a directory
[Fri Feb 28 02:56:24 UTC 2025] [info] use root no password show backends result 10221	test-disaggregated-cluster-cg1-1.test-disaggregated-cluster-cg1.default.svc.cluster.local	9050	-1	-1	-1	-1	NULL	NULL	false	false	0	0.000 	0.000 	1.000 B	0.000 	0.00 %	0.00 %	0.000 	{"cloud_unique_id" : "1:1751150972:t1Ws6Mrv", "compute_group_status" : "NORMAL", "private_endpoint" : "", "compute_group_name" : "cg1", "location" : "default", "public_endpoint" : "", "compute_group_id" : "ZNES_zRC"}	java.net.UnknownHostException: test-disaggregated-cluster-cg1-1.test-disaggregated-cluster-cg1.default.svc.cluster.local		{"lastStreamLoadTime":-1,"isQueryDisabled":false,"isLoadDisabled":false,"isActive":true,"currentFragmentNum":0,"lastFragmentUpdateTime":0}	287		1	0.00
10222	test-disaggregated-cluster-cg1-0.test-disaggregated-cluster-cg1.default.svc.cluster.local	9050	-1	-1	-1	-1	NULL	NULL	false	false	0	0.000 	0.000 	1.000 B	0.000 	0.00 %	0.00 %	0.000 {"cloud_unique_id" : "1:1751150972:YA6LXvvg", "compute_group_status" : "NORMAL", "private_endpoint" : "", "compute_group_name" : "cg1", "location" : "default", "public_endpoint" : "", "compute_group_id" : "ZNES_zRC"}	java.net.UnknownHostException: test-disaggregated-cluster-cg1-0.test-disaggregated-cluster-cg1.default.svc.cluster.local		{"lastStreamLoadTime":-1,"isQueryDisabled":false,"isLoadDisabled":false,"isActive":true,"currentFragmentNum":0,"lastFragmentUpdateTime":0}	287		1	0.00
10223	test-disaggregated-cluster-cg1-2.test-disaggregated-cluster-cg1.default.svc.cluster.local	9050	-1	-1	-1	-1	NULL	NULL	false	false	0	0.000 	0.000 	1.000 B	0.000 	0.00 %	0.00 %	0.000 {"cloud_unique_id" : "1:1751150972:Ox_RJuee", "compute_group_status" : "NORMAL", "private_endpoint" : "", "compute_group_name" : "cg1", "location" : "default", "public_endpoint" : "", "compute_group_id" : "ZNES_zRC"}	java.net.UnknownHostException: test-disaggregated-cluster-cg1-2.test-disaggregated-cluster-cg1.default.svc.cluster.local		{"lastStreamLoadTime":-1,"isQueryDisabled":false,"isLoadDisabled":false,"isActive":true,"currentFragmentNum":0,"lastFragmentUpdateTime":0}	287		1	0.00
10224	test-disaggregated-cluster-cg2-0.test-disaggregated-cluster-cg2.default.svc.cluster.local	9050	-1	-1	-1	-1	NULL	NULL	false	false	0	0.000 	0.000 	1.000 B	0.000 	0.00 %	0.00 %	0.000 {"cloud_unique_id" : "1:1751150972:E_SJoMU8", "compute_group_status" : "NORMAL", "private_endpoint" : "", "compute_group_name" : "cg2", "location" : "default", "public_endpoint" : "", "compute_group_id" : "oZ2gH5Ml"}	java.net.UnknownHostException: test-disaggregated-cluster-cg2-0.test-disaggregated-cluster-cg2.default.svc.cluster.local		{"lastStreamLoadTime":-1,"isQueryDisabled":false,"isLoadDisabled":false,"isActive":true,"currentFragmentNum":0,"lastFragmentUpdateTime":0}	287		1	0.00
10251	test-disaggregated-cluster-cg2-1.test-disaggregated-cluster-cg2.default.svc.cluster.local	9050	-1	-1	-1	-1	NULL	NULL	false	false	0	0.000 	0.000 	1.000 B	0.000 	0.00 %	0.00 %	0.000 {"cloud_unique_id" : "1:1751150972:B_h0m9vp", "compute_group_status" : "NORMAL", "private_endpoint" : "", "compute_group_name" : "cg2", "location" : "default", "public_endpoint" : "", "compute_group_id" : "oZ2gH5Ml"}	java.net.UnknownHostException: test-disaggregated-cluster-cg2-1.test-disaggregated-cluster-cg2.default.svc.cluster.local		{"lastStreamLoadTime":-1,"isQueryDisabled":false,"isLoadDisabled":false,"isActive":true,"currentFragmentNum":0,"lastFragmentUpdateTime":0}	287		1	0.00
10252	test-disaggregated-cluster-cg2-2.test-disaggregated-cluster-cg2.default.svc.cluster.local	9050	-1	-1	-1	-1	NULL	NULL	false	false	0	0.000 	0.000 	1.000 B	0.000 	0.00 %	0.00 %	0.000 {"cloud_unique_id" : "1:1751150972:nmt5aHJC", "compute_group_status" : "NORMAL", "private_endpoint" : "", "compute_group_name" : "cg2", "location" : "default", "public_endpoint" : "", "compute_group_id" : "oZ2gH5Ml"}	java.net.UnknownHostException: test-disaggregated-cluster-cg2-2.test-disaggregated-cluster-cg2.default.svc.cluster.local		{"lastStreamLoadTime":-1,"isQueryDisabled":false,"isLoadDisabled":false,"isActive":true,"currentFragmentNum":0,"lastFragmentUpdateTime":0}	287		1	0.00  .
[Fri Feb 28 02:56:24 UTC 2025] [info] Check myself (test-disaggregated-cluster-cg1-0.test-disaggregated-cluster-cg1.default.svc.cluster.local:9050) exist in FE, start be directly ...
/etc/podinfo/annotationsis not exists.
[Fri Feb 28 02:56:24 UTC 2025] run start_be.sh
Disable swap memory before starting be
```

问题详情参见：https://github.com/apache/doris/issues/48460。一开始关注了 `/opt/apache-doris/be_disaggregated_entrypoint.sh: line 73: /opt/apache-doris/be/conf/: Is a directory` 以为是脚本 be_disaggregated_entrypoint.sh 有问题, 经过自己修改重新打包后发现问题依旧。重新检查日志和脚本发现在 start_be.sh 中, 检查了交换内存的命令，如果发现交换内存启用，就会退出。

```bash
# line 193
        if [[ "$(swapon -s | wc -l)" -gt 1 ]]; then
            echo "Disable swap memory before starting be"
            exit 1
        fi
```

通过在宿主机上执行 `sudo swapoff -a` 关闭交换内存后，Be 就可以正常启动了。

2. 重启 minikube 后 FE 跟 meta service 交互提示错误

```bash
RuntimeLogger 2025-02-28 06:38:51,460 WARN (main|1) [CloudEnv.getLocalTypeFromMetaService():165] failed to get cloud cluster due to incomplete response, cloud_unique_id=1:1751150972:fe, clusterId=RESERVED_CLUSTER_ID_FOR_SQL_SERVER, response=status {
code: INVALID_ARGUMENT
msg: "empty instance_id"
}
```

这个问题没有找到原因，在 issue 中有相同问题：https://github.com/apache/doris/issues/47678

最终解决，通过重新搭建 fdb + doris 集群解决。

3. doris 集群和 Storage Vault 都部署成功后

通过 mysql 连接到 FE 创建 test_db, test_tbl 并插入数据。插入数据时卡顿，且最终失败：

```bash
> INSERT INTO test_tbl VALUES 
(1, 'Tom', 89.5),
(2, 'Jerry', 92.0),
(3, 'Jack', 85.5);
> 
ERROR 1105 (HY000): errCode = 2, detailMessage = Backend Backend [id=10262, host=test-disaggregated-cluster-cg1-2.test-disaggregated-cluster-cg1.default.svc.cluster.local, heartbeatPort=9050, alive=false, lastStartTime=2025-02-28 08:46:45, process epoch=1740732405712, isDecommissioned=false, tags: {cloud_cluster_id=Ppt5BJ3g, cloud_unique_id=1:1751150972:JxbgZEBP, cloud_cluster_status=NORMAL, cloud_cluster_name=cg1, location=default, cloud_cluster_private_endpoint=, cloud_cluster_public_endpoint=}], backendStatus: [lastSuccess
```

与此同时 cg1 整体都触发了重启，查看日志发现：

```bash
RuntimeLogger F20250228 09:32:36.952682   878 storage_engine.cpp:120] Check failed: _type == Type::CLOUD ( vs. )

*** Check failure stack trace: ***

RuntimeLogger I20250228 09:32:36.955466  1376 workload_group_manager.cpp:200]

Process Memory Summary: process memory used 856.95 MB(= 935.89 MB[vm/rss] - 78.94 MB[tc/jemalloc_cache] + 0[reserved] + 0B[waiting_refresh]), sys available memory 37.77 GB(= 37.77 GB[proc/available] - 0[reserved] - 0B[waiting_refresh]), all workload groups memory usage: 384.00 B, weighted_memory_limit_ratio: 0.9851607626358986

@     0x5f7414439556  google::LogMessage::SendToLog()

@     0x5f7414435fa0  google::LogMessage::Flush()

@     0x5f7414439d99  google::LogMessageFatal::~LogMessageFatal()

@     0x5f7409adbc3a  doris::BaseStorageEngine::to_cloud()

@     0x5f7409d600a7  doris::LoadChannel::open()

@     0x5f7409d5ab0a  doris::LoadChannelMgr::open()

@     0x5f7409eaa68d  std::_Function_handler<>::_M_invoke()

@     0x5f7409ec4ebb  doris::WorkThreadPool<>::work_thread()

@     0x5f74173d1b20  execute_native_thread_routine

@     0x7ce578024ac3  (unknown)

@     0x7ce5780b6850  (unknown)

@              (nil)  (unknown)

*** Query id: 66dcfbe890a240d8-b45c17c9fa19c124 ***

*** is nereids: 0 ***

*** tablet id: 0 ***

*** Aborted at 1740735157 (unix time) try "date -d @1740735157" if you are using GNU date ***

*** Current BE git commitID: 62a58bff4c ***

*** SIGABRT unknown detail explain (@0x86) received by PID 134 (TID 878 OR 0x7ce39d600640) from PID 134; stack trace: ***

RuntimeLogger I20250228 09:32:37.674818   463 wal_manager.cpp:485] Scheduled(every 10s) WAL info: [/opt/apache-doris/be/storage/wal: limit 3993782272 Bytes, used 0 Bytes, estimated wal bytes 0 Bytes, available 3993782272 Bytes.];

RuntimeLogger I20250228 09:32:37.726341  1376 daemon.cpp:239] os physical memory 62.66 GB. process memory used 1.12 GB(= 1.20 GB[vm/rss] - 79.15 MB[tc/jemalloc_cache] + 0[reserved] + 0B[waiting_refresh]), limit 56.40 GB, soft limit 50.76 GB. sys available memory 37.55 GB(= 37.55 GB[proc/available] - 0[reserved] - 0B[waiting_refresh]), low water mark 3.13 GB, warning water mark 6.27 GB.

0# doris::signal::(anonymous namespace)::FailureSignalHandler(int, siginfo_t*, void*) at /home/zcp/repo_center/doris_release/doris/be/src/common/signal_handler.h:421

1# 0x00007CE577FD2520 in /lib/x86_64-linux-gnu/libc.so.6

2# pthread_kill in /lib/x86_64-linux-gnu/libc.so.6

3# raise in /lib/x86_64-linux-gnu/libc.so.6

4# abort in /lib/x86_64-linux-gnu/libc.so.6

5# 0x00005F7414443E2D in /opt/apache-doris/be/lib/doris_be

6# 0x00005F741443646A in /opt/apache-doris/be/lib/doris_be

7# google::LogMessage::SendToLog() in /opt/apache-doris/be/lib/doris_be

8# google::LogMessage::Flush() in /opt/apache-doris/be/lib/doris_be

9# google::LogMessageFatal::~LogMessageFatal() in /opt/apache-doris/be/lib/doris_be

10# doris::BaseStorageEngine::to_cloud() in /opt/apache-doris/be/lib/doris_be

11# doris::LoadChannel::open(doris::PTabletWriterOpenRequest const&) at /home/zcp/repo_center/doris_release/doris/be/src/runtime/load_channel.cpp:131

12# doris::LoadChannelMgr::open(doris::PTabletWriterOpenRequest const&) at /home/zcp/repo_center/doris_release/doris/be/src/runtime/load_channel_mgr.cpp:108

13# std::_Function_handler<void (), doris::PInternalService::tablet_writer_open(google::protobuf::RpcController*, doris::PTabletWriterOpenRequest const*, doris::PTabletWriterOpenResult*, google::protobuf::Closure*)::$_0>::_M_invoke(std::_Any_data const&) at /var/local/ldb-toolchain/bin/../lib/gcc/x86_64-linux-gnu/11/../../../../include/c++/11/bits/std_function.h:291

14# doris::WorkThreadPool<false>::work_thread(int) at /home/zcp/repo_center/doris_release/doris/be/src/util/work_thread_pool.hpp:159

15# execute_native_thread_routine at ../../../../../libstdc++-v3/src/c++11/thread.cc:84

16# 0x00007CE578024AC3 in /lib/x86_64-linux-gnu/libc.so.6

17# 0x00007CE5780B6850 in /lib/x86_64-linux-gnu/libc.so.6

/opt/apache-doris/be/bin/start_be.sh: line 433:   134 Aborted                 (core dumped) ${LIMIT:+${LIMIT}} "${DORIS_HOME}/lib/doris_be" "$@" 2>&1 < /dev/null
```

根据堆栈找到 storage_engine.cpp 中的 120 行：

```cpp
StorageEngine& BaseStorageEngine::to_local() {
    CHECK_EQ(_type, Type::LOCAL);
    return *static_cast<StorageEngine*>(this);
}

CloudStorageEngine& BaseStorageEngine::to_cloud() {
    CHECK_EQ(_type, Type::CLOUD);
    return *static_cast<CloudStorageEngine*>(this);
}

// CHECK_EQ 是 glog 的一个宏，用于检查表达式是否为真。如果表达式为假，则会打印错误消息并终止程序。
// https://github.com/google/glog/blob/master/src/glog/logging.h#L688
```

从这里可以得知 `CHECK_EQ(_type, Type::CLOUD);` 这个检查失败，导致 BE 重启。

解决手段：将 BE 的存储引擎模式设置为 `cloud`。参见 https://doris.apache.org/zh-CN/docs/3.0/compute-storage-decoupled/compilation-and-deployment#541-%E9%85%8D%E7%BD%AE-beconf
