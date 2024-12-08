import mongomock
import unittest
from pymongo import MongoClient

class TestMongoMock(unittest.TestCase):

    def setUp(self):
        self.client = mongomock.MongoClient()
        self.db = self.client['test_database']
        self.collection = self.db['test_collection']

    def test_insert_and_find(self):
        """
        测试插入和查找
        """
        self.collection.insert_one({'name': 'Alice', 'age': 30})
        result = self.collection.find_one({'name': 'Alice'})

        print(f'case: test_insert_and_find result: {result}')

        self.assertIsNotNone(result)
        self.assertIsNotNone(result['_id'])
        self.assertEqual(result['age'], 30)
        self.assertEqual(result['name'], 'Alice')

    def test_upsert(self):
        """
        测试更新 upsert=True
        """
        self.collection.update_one({'name': 'Tom'}, {'$set': {'age': 28}}, upsert=True)
        result1 = self.collection.find_one({'name': 'Tom'})

        self.assertIsNotNone(result1)
        self.assertIsNotNone(result1['_id'])
        self.assertEqual(result1['age'], 28)
        self.assertEqual(result1['name'], 'Tom')

        self.collection.update_one({'name': 'Tom'}, {'$set': {'age': 30}})
        result2 = self.collection.find_one({'name': 'Tom'})

        print(f'case: test_upsert result1: {result1}, result2: {result2}')

        self.assertIsNotNone(result2)
        self.assertIsNotNone(result2['_id'])
        self.assertEqual(result2['age'], 30)
        self.assertEqual(result2['name'], 'Tom')
        self.assertEqual(result1['_id'], result2['_id'])

    def test_update_push(self):
        """
        测试更新数组
        """
        self.collection.insert_one({'name': 'Alice', 'age': 30, 'books': ['Python']})
        self.collection.update_one({'name': 'Alice'}, {'$push': {'books': 'Java'}})
        result = self.collection.find_one({'name': 'Alice'})

        print(f'case: test_update_push result: {result}')
        self.assertIsNotNone(result)
        self.assertIsNotNone(result['_id'])
        self.assertEqual(result['age'], 30)
        self.assertEqual(result['name'], 'Alice')
        self.assertEqual(result['books'], ['Python', 'Java'])

    def test_insert_delete(self):
        """
        测试删除
        """
        self.collection.insert_one({'name': 'Alice', 'age': 30})
        result = self.collection.find_one({'name': 'Alice'})
        self.assertIsNotNone(result)

        self.collection.delete_one({'name': 'Alice'})
        result = self.collection.find_one({'name': 'Alice'})

        print(f'case: test_insert_delete result: {result}')
        self.assertIsNone(result)

if __name__ == '__main__':
    unittest.main()