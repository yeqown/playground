# aggr_mongo_update 用来将对 mongodb 的操作进行封装聚合
# 避免对 mongodb 多次操作，如：
# collection.update_one({'name': 'Tom'}, {'$set': {'age': 28}}, upsert=True) 
# 执行了20次, 更新内容各不相同，但是都是对同一个文档进行操作
# aggr_mongo_update 则可以将这20次操作聚合成一次，最终才发送到 mongodb.

import mongomock
from pymongo import MongoClient

class AggrMongoColl(object):

    def __init__(self, database, coll, flush_interval):
        self.client = mongomock.MongoClient()
        db = self.client[database]
        self.collection = db[coll]

        actual_db = MongoClient(host='localhost', port=27017)[database]
        self.actual_collection = actual_db[coll]

        self.flush_interval = flush_interval

    def get_collection(self):
        return self.collection

    def batch_sync(self, batch_size: int):
        # TODO: 这部分可以允许自定义
        # 将 self.collection 中的数据按照 batch_size 分批插入(upsert)到 self.actual_collection

        # 1. 获取 self.collection 中的数据
        # 2. 按照 batch_size 分批插入(upsert)到 self.actual_collection
        items = list(self.collection.find())
        self.actual_collection.insert_many(items)

    def flush(self):
        # 将 self.collection 中的数据按照 batch_size 分批插入(upsert)到 self.actual_collection
        self.batch_sync(batch_size=10)
        self.clear()

    def clear(self):
        # TODO: 这部分也允许自定义
        # 根据规则清除已经同步的数据
        pass

# 预期的使用方式如下：

# 1. 创建一个 aggr_mongo_update 对象, 指明数据库和集合
# 这样会创建两个 collection 对象
# 一个 collection 代表 mongomock 的 collection 对象，客户端操作这个对象
# 一个 collection 代表真实的 mongodb 的 collection 对象，后台操作这个对象
aggr = AggrMongoColl(database='test_database', coll='test_collection', flush_interval=10)
coll = aggr.get_collection()

# 2. 客户端操作 collection 对象
# 2.1 插入数据
coll.insert_one({'name': 'Alice', 'age': 30})

# 2.2 多次更新数据
for i in range(20):
    coll.update_one({'name': 'Alice'}, {'$set': {'age': 30+i}}, upsert=True)

# 等待 flush_interval 时间，或者手动调用 flush 方法
# time.sleep(10)
aggr.flush()