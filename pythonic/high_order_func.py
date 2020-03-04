"""
    # 高阶函数学习练手
    # Author: yeqiang
    # Date：170908
"""

import sys
import functools

####################### wraps 的效果展示 ###############################
# 不带functools.wraps
def trace(func):
    if sys.stderr:
        def call_func(*args, **kwargs):
            """a wrapper functoins"""
            # sys.stderr.write("Calling functions； {}\n".format(func.__name__))
            ret_val = func(*args, **kwargs)
            sys.stderr.write("The Func of {}'s retrun value is {}\n".format(
                func.__name__, ret_val))
            return ret_val
        return call_func
    else:
        return func


def trace_with_wraps(func):
    if sys.stderr:
        @functools.wraps(func)
        def call_func(*args, **kwargs):
            """a wrapper functoins"""
            # sys.stderr.write("Calling functions； {}\n".format(func.__name__))
            ret_val = func(*args, **kwargs)
            sys.stderr.write("The Func of {}'s retrun value is {}\n".format(
                func.__name__, ret_val))
            return ret_val
        return call_func
    else:
        return func

@trace
def square(val):
    """this is a normal function"""
    return val * val

@trace_with_wraps
def square2(val):
    """this is also normal function"""

class DemoPlusHandler(object):

    self._val = 1

    @property
    def val(self):
        return self._val

    @staticmethod
    def plus(one_val, other_val):
        return sone_val * other_val

    @staticmethod
    def plus_self(other_val):
        return functools.partial(DemoPlusHandler.plus, DemoPlusHandler.val)

################### 演示prtial的用法 #########################


if __name__ == "__main__":
    # print("val: 3, square: {}".format(square(3)))

    print("没有带有functools.wraps: ", square.__name__)
    print("带了functools.wraps: ", square2.__name__)



