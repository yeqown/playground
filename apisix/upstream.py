# upstream 服务
# 运行在 3001 端口上的 http server, handler 打印所有的请求头信息

import http.server
import socketserver

PORT = 3001

class MyHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        print(self.headers)
        self.send_response(200)
        self.end_headers()
        # 把 headers 信息返回给客户端
        self.wfile.write(str(self.headers).encode())

with socketserver.TCPServer(("", PORT), MyHandler) as httpd:
    print("serving at port", PORT)
    httpd.serve_forever()