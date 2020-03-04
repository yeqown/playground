# gonic

golang code for test and some note 

## http1.x与http2 对比测试

#### [HTTP1.x](http)

测试结果
```sh
➜  http2 git:(master) ✗ curl -Ik --http1.1 https://localhost:8080/http1x/std

HTTP/1.1 200 OK
Date: Mon, 11 Feb 2019 03:10:09 GMT
Content-Length: 21
Content-Type: text/plain; charset=utf-8

➜  ~ bombardier --http1 -k -l https://localhost:8080/http1x/std

Bombarding https://localhost:8080/http1x/std for 10s using 125 connection(s)
[==========================================================================] 10s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec     28540.24    5945.34   72867.04
  Latency        4.36ms     3.73ms   190.76ms
  Latency Distribution
     50%     3.72ms
     75%     4.80ms
     90%     6.67ms
     95%     8.47ms
     99%    17.16ms
  HTTP codes:
    1xx - 0, 2xx - 280844, 3xx - 0, 4xx - 0, 5xx - 0
    others - 4847
  Errors:
    Get https://localhost:8080/http1x/std: dial tcp 127.0.0.1:8080: socket: too many open files - 4838
    Get https://localhost:8080/http1x/std: dial tcp [::1]:8080: socket: too many open files - 9
  Throughput:     8.10MB/s%  
```

#### [HTTP2](http2)

准备TLS
```sh
1，生成服务端私钥
openssl genrsa -out http2.key 2048
2，生成服务端证书
openssl req -new -x509 -key http2.key -out http2.pem -days 3650
```

测试结果
```sh
➜  http2 git:(master) ✗ curl -Ik --http2 https://localhost:8080/http2/std

HTTP/2 200
content-type: text/plain; charset=utf-8
content-length: 20
date: Mon, 11 Feb 2019 03:11:10 GMT

➜  http2 git:(master) ✗ bombardier -c 125 -t 5s -l -d 10s --http2 -k https://localhost:8080/http2/std

Bombarding https://localhost:8080/http2/std for 10s using 125 connection(s)
[==========================================================================] 10s
Done!
Statistics        Avg      Stdev        Max
Reqs/sec     12883.03    3540.87   21067.48
Latency        9.74ms     5.94ms   222.08ms
Latency Distribution
    50%     8.70ms
    75%    11.30ms
    90%    15.17ms
    95%    18.00ms
    99%    26.74ms
HTTP codes:
1xx - 0, 2xx - 128244, 3xx - 0, 4xx - 0, 5xx - 0
others - 0
Throughput:     1.69MB/s%
```