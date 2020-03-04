"""
python 实现栈push,pop O(1)时间复杂度取min,max值
"""

from random import randint
from time import sleep


class Stack(object):

    def __init__(self):
        self._data = []
        self._min = None
        self._second_min = None
        self._max = None
        self._second_max = None

    def pop(self):
        if not len(self._data):
            return None
        item = self._data.pop()
        if item == self._min:
            self._min = self._second_min
        if item == self._max:
            self._max = self._second_max
        # 怎么持续更新最小值呢？

    def push(self, item):
        self._data.append(item)
        if not self._min:
            self._min = item
        if not self._second_min:
            self._second_min = item
        if not self._max:
            self._max = item
        if not self._second_max:
            self._second_max = item

        # 设置最小值
        if item < self._min:
            self._min, self._second_min = item, self._min
            return
        # 设置最大值
        if item > self._max:
            self._max, self._second_max = item, self._max
            return
        # 设置第二小值
        if self._second_min > item > self._min:
            self._second_min = item
            return
        # 设置第二大值
        if self._max > item > self._second_max:
            self._second_max = item
            return

    @property
    def min(self):
        return self._min

    @property
    def max(self):
        return self._max


def test_push():

    s = Stack()

    for i in range(0, 100):
        item = randint(1, 67)
        print("Push item: {}".format(item))
        s.push(item)
        print("Current Max: {max}, SecondMax: {}, Min: {min} SecondMin".format(
            max=s.max, min=s.min))
        sleep(1)


def test_pop():

    s = Stack()
    action = lambda e: s.push(randint(2, 999))
    # push data into Stack
    _ = [action for i in range(100)]
    # test Stack.pop() method, how to keep min and max
    for _ in range(100):
        item = s.pop()
        print(item)


test_push()
