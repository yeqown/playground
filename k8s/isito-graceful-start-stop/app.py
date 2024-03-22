# 1. 编写一个 python 服务，应用启动时会连接到一个外部的 redis 服务，启动后会每隔 1s 向 redis 发送一个 ping 消息。
#     - 应用配置 readyz 探针，当应用启动完成后，会返回 200 状态码。

import time
import os
import redis
import signal
import threading
from flask import Flask

# get redis host from env var REDIS_HOST
REDIS_HOST = os.getenv('REDIS_HOST', 'localhost')
app = Flask(__name__)
cache = redis.Redis(host=REDIS_HOST, port=6379, password="123456")

@app.route('/')
def hello():
    return 'Hello World!'

@app.route('/readyz')
def readyz():
    try:
        cache.ping()
        return 'ok', 200
    except:
        return 'failed', 500
    
@app.route('/liveness')
def liveness():
    return 'alive', 200

ping_loop_enabled = True
def ping_loop():
    global ping_loop_enabled
    while ping_loop_enabled:
        cache.ping()
        print('ping redis {}'.format(time.time()))
        time.sleep(2)

ping_loop_thread = threading.Thread(target=ping_loop)

def sigterm_handler(_signo, _stack_frame):
    print('sigterm_handler called')

    # mock 5s to release resources and execute graceful shutdown
    time.sleep(5)

    # stop ping loop
    print('stopping ping loop')
    ping_loop_enabled = False

    # close redis connection
    print('closing redis connection')
    cache.close()

    print('exiting web')

    os._exit(0)

def sigkill_handler(_signo, _stack_frame):
    os._exit(0)

if __name__ == "__main__":
    # register sigterm handler
    signal.signal(signal.SIGTERM, sigterm_handler)
    signal.signal(signal.SIGINT, sigterm_handler)
    signal.signal(signal.SIGKILL, sigkill_handler)

    # start ping loop
    ping_loop_thread.start()

    # start a timer to send ping message to redis every 1s in another thread
    app.run(host="0.0.0.0", debug=True)