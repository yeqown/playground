### Conntector 配置说明

```json
{
    "name": "elasticsearch-sink",
    "config": {
        "connector.class": "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
        "tasks.max": "1",                   // 设置 worker 任务数
        "topics.regex": "mysql-cdc\\..*",   // 设置监听的 topic 正则
        "connection.url": "http://elasticsearch:9200", // ES 地址

        // 配置 SMT transform, 这里的顺序就是执行顺序 unwrap -> RenameField -> key -> TimestampRouter
        "transforms": "unwrap,RenameField,key,TimestampRouter",
        // unwrap 用来执行删除动作
        "transforms.unwrap.type": "io.debezium.transforms.ExtractNewRecordState",
        "transforms.unwrap.drop.tombstones": "true",
        "transforms.unwrap.delete.handling.mode": "none",
        // RenameField 在 predicate 条件满足的情况下执行，将 user_id 字段改为 id
        "transforms.RenameField.type": "org.apache.kafka.connect.transforms.ReplaceField$Key",
        "transforms.RenameField.renames": "user_id:id",
        "transforms.RenameField.predicate": "userIdKey",
        // key 用来执行 key 的转换，将 key.field 指定的字段替换整个 message 的 key
        // 最终影响到 ES 的 docId
        "transforms.key.type": "org.apache.kafka.connect.transforms.ExtractField$Key",
        "transforms.key.field": "id",
        // TimestampRouter 用来执行 topic 的转换，将 topic 名称改为 topic-timestamp 的格式
        // ！！！注意：TimestampRouter 使用的是 kafka message 中的 timestamp
        "transforms.TimestampRouter.type": "org.apache.kafka.connect.transforms.TimestampRouter",
        "transforms.TimestampRouter.topic.format": "${topic}-${timestamp}",
        "transforms.TimestampRouter.timestamp.format": "yyyyMMdd",
        // predicates 配置了一条规则：如果 topic 名称匹配 user_profile
        // 提供给 RenameField 这个 transform 有条件的执行
        "predicates": "userIdKey",
        "predicates.userIdKey.type": "org.apache.kafka.connect.transforms.predicates.TopicNameMatches",
        "predicates.userIdKey.pattern": ".*\\.user_profile",
        "key.ignore": "false",
        "behavior.on.null.values": "DELETE",
        "write.method": "UPSERT",
        "flush.synchronously": "true"
    }
}
```