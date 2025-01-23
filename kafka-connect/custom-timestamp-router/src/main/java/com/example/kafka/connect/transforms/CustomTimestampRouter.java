package com.example.kafka.connect.transforms;

import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.connect.connector.ConnectRecord;
import org.apache.kafka.connect.components.Versioned;
import org.apache.kafka.connect.transforms.Transformation;
import org.apache.kafka.connect.data.Struct;
import org.apache.kafka.connect.data.Field;
import java.text.SimpleDateFormat;
import org.apache.kafka.connect.transforms.util.SimpleConfig;
import java.util.Map;
import java.util.Date;

import java.util.TimeZone;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * 自定义时间戳路由器：根据 record 中的 timestamp 字段值，将 record 调整为新的主题
 * 如：
 * {"created_at": "2023-02-01T12:00:00Z", "topic": "my-topic"} 
 * 调整为：
 * {"created_at": "2023-02-01T12:00:00Z", "topic": "my-topic-20230201"}
 */
public class CustomTimestampRouter<R extends ConnectRecord<R>> implements Transformation<R>, AutoCloseable, Versioned {
    private String topicFormat;
    private String messageTimestampField;
    private ThreadLocal<SimpleDateFormat> timestampFormat;

    private static final Logger log = LoggerFactory.getLogger(CustomTimestampRouter.class);

    @Override
    public String version() {
        return "v0.0.1";
    }

    private interface ConfigName {
        String MESSAGE_TIMESTAMP_FIELD = "message.timestamp.field";
        String TOPIC_FORMAT = "topic.format";
        String TIMESTAMP_FORMAT = "timestamp.format";
    }

    /**
     * 配置参数
     * message.timestamp.field: 时间戳字段名
     * topic.format: 主题格式
     * timestamp.format: 时间戳处理格式化
     */
    public static final ConfigDef CONFIG_DEF = new ConfigDef()
        .define(ConfigName.MESSAGE_TIMESTAMP_FIELD, ConfigDef.Type.STRING, "created_at", ConfigDef.Importance.HIGH, "时间戳字段名")
        .define(ConfigName.TOPIC_FORMAT, ConfigDef.Type.STRING, "${topic}-${timestamp}", ConfigDef.Importance.HIGH, "主题格式")
            .define(ConfigName.TIMESTAMP_FORMAT, ConfigDef.Type.STRING, "yyyyMMdd", ConfigDef.Importance.HIGH, "时间戳处理格式");

    @Override
    public R apply(R record) {
        Long timestamp = null;
        
        // 尝试从记录值中获取时间戳
        if (record.value() instanceof Struct) {
            Struct value = (Struct) record.value();
            try {
                Field field = value.schema().field(messageTimestampField);
                if (field != null) {
                    timestamp = value.getInt64(messageTimestampField);
                }
            } catch (Exception e) {
                log.warn("Failed to extract timestamp from field {}: {}", messageTimestampField, e.getMessage());
            }
        }

        // 如果无法从字段获取时间戳，则使用记录的时间戳
        if (timestamp == null) {
            log.warn("No timestamp found in record.value {}, trying record.timestamp", messageTimestampField);
            timestamp = record.timestamp();
            if (timestamp == null) {
                log.warn("No timestamp found in record, using current time");
                timestamp = System.currentTimeMillis();
            }
        }

        // 格式化时间戳
        String formattedTimestamp = timestampFormat.get().format(new Date(timestamp));
        
        // 替换主题格式中的变量
        String updatedTopic = topicFormat
            .replace("${topic}", record.topic())
            .replace("${timestamp}", formattedTimestamp);

        // log.info("Updated topic: {}", updatedTopic);

        // 创建新记录
        return record.newRecord(
            updatedTopic,
            record.kafkaPartition(),
            record.keySchema(),
            record.key(),
            record.valueSchema(),
            record.value(),
            record.timestamp()
        );
    }

    @Override
    public ConfigDef config() {
        return CONFIG_DEF;
    }

    @Override
    public void configure(Map<String, ?> props) {
        final SimpleConfig config = new SimpleConfig(CONFIG_DEF, props);

        topicFormat = config.getString(ConfigName.TOPIC_FORMAT);
        messageTimestampField = config.getString(ConfigName.MESSAGE_TIMESTAMP_FIELD);

        final String timestampFormatStr = config.getString(ConfigName.TIMESTAMP_FORMAT);
        timestampFormat = ThreadLocal.withInitial(() -> {
            final SimpleDateFormat fmt = new SimpleDateFormat(timestampFormatStr);
            fmt.setTimeZone(TimeZone.getTimeZone("UTC"));
            return fmt;
        });

        log.info("CustomTimestampRouter configured with topic format: {}, timestamp field: {}, timestamp format: {}",
            topicFormat, messageTimestampField, timestampFormatStr);
    }

    @Override
    public void close() {
        timestampFormat.remove();
    }
}