package com.example.kafka.connect.transforms;

import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.data.SchemaBuilder;
import org.apache.kafka.connect.data.Struct;
import org.apache.kafka.connect.sink.SinkRecord;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.util.HashMap;
import java.util.Map;

import static org.junit.Assert.assertTrue;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class CustomTimestampRouterTest {

    private CustomTimestampRouter<SinkRecord> router;
    private Map<String, String> configs;

    @Before
    public void setUp() {
        router = new CustomTimestampRouter<>();
        configs = new HashMap<>();
        configs.put("message.timestamp.field", "created_at");
        configs.put("topic.format", "${topic}-${timestamp}");
        configs.put("timestamp.format", "yyyyMMdd");
        router.configure(configs);
    }

    @After
    public void tearDown() {
        router.close();
    }

    @Test
    public void testVersion() {
        assertEquals("v0.0.1", router.version());
    }

    @Test
    public void testApplyWithValidTimestamp() {
        // 创建带有时间戳的记录
        Schema schema = SchemaBuilder.struct()
            .field("id", Schema.INT64_SCHEMA)
            .field("created_at", Schema.INT64_SCHEMA)
            .build();

        Struct value = new Struct(schema)
            .put("id", 1L)
            .put("created_at", 1677628800000L); // 2023-03-01 00:00:00

        SinkRecord record = new SinkRecord(
            "test-topic",
            0,
            null,
            null,
            schema,
            value,
            0
        );

        SinkRecord transformedRecord = router.apply(record);
        assertEquals("test-topic-20230301", transformedRecord.topic());
    }

    @Test
    public void testApplyWithMissingTimestampField() {
        // 创建没有时间戳字段的记录
        Schema schema = SchemaBuilder.struct()
            .field("id", Schema.INT64_SCHEMA)
            .build();

        Struct value = new Struct(schema)
            .put("id", 1L);

        long currentTime = System.currentTimeMillis();
        SinkRecord record = new SinkRecord(
            "test-topic",
            0,
            null,
            null,
            schema,
            value,
            currentTime
        );

        SinkRecord transformedRecord = router.apply(record);
        assertNotNull(transformedRecord.topic());
        // 验证使用了记录的时间戳
        assertEquals("test-topic-" + new java.text.SimpleDateFormat("yyyyMMdd")
            .format(new java.util.Date(currentTime)), transformedRecord.topic());
    }

    @Test
    public void testApplyWithNullValue() {
        long currentTime = System.currentTimeMillis();
        SinkRecord record = new SinkRecord(
            "test-topic",
            0,
            null,
            null,
            null,
            null,
            currentTime
        );

        SinkRecord transformedRecord = router.apply(record);
        assertNotNull(transformedRecord.topic());
        // 验证使用了记录的时间戳
        assertEquals("test-topic-" + new java.text.SimpleDateFormat("yyyyMMdd")
            .format(new java.util.Date(currentTime)), transformedRecord.topic());
    }

    @Test
    public void testConfig() {
        ConfigDef config = router.config();
        assertNotNull(config);
        assertEquals(3, config.names().size());
        assertTrue(config.names().contains("message.timestamp.field"));
        assertTrue(config.names().contains("topic.format"));
        assertTrue(config.names().contains("timestamp.format"));
    }
}