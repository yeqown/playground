# 在 python 中使用 memcached 客户端
# 验证 memcached 集群路由算法

from pymemcache.client.hash import HashClient

client = HashClient([
  "127.0.0.1:11211",
  "127.0.0.1:11212",
])

# 写入 100 30s 内过期的数据, 来找到两条会被写入到不同节点的数据
# for i in range(100):
#   key = f"key-{i}"
#   value = f"value-{i}"
#   client.set(key, value, expire=30)
  
# # 读取数据, 验证数据是否被写入到不同节点
# for i in range(100):
#   key = f"key-{i}"
#   value = client.get(key)
#   print(key, value)

client.set("key-0", "I'll be sent to node1(11211)", expire=30)
client.set("key-2", "I'll be sent to node2(11212)", expire=30)