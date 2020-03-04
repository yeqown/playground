# The guess API is already defined for you.
# @param num, your guess
# @return -1 if my number is lower, 1 if my number is higher, otherwise return 0
# def guess(num):

import time
import math


def guess(num):
    if num == pick:
        return 0
    elif num > pick:
        return -1
    else:
        return 1


class Solution(object):
    def guessNumber(self, n):
        """
        :type n: int
        :rtype: int
        """
        bottom = 0
        top = n
        mid = 0
        while True:
            mid = math.floor((bottom + top)/2)
            result = guess(mid)
            if result > 0:
                bottom = mid + 1
            elif result < 0:
                top = mid - 1
        return mid


if __name__ == "__main__":
    global pick
    pick = 12817
    print(Solution().guessNumber(20000))
