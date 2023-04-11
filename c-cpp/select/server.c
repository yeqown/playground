// Copyright 2020 The yeqown Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <sys/select.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <assert.h>

// int select(int maxfdp1,fd_set *readset,fd_set *writeset,fd_set *exceptset,const struct timeval *timeout) // 返回值：就绪描述符的数目，超时返回0，出错返回-1
// void FD_ZERO(fd_set *fdset);          // 清空集合
// void FD_SET(int fd, fd_set *fdset);   // 将一个给定的文件描述符加入集合之中
// void FD_CLR(int fd, fd_set *fdset);   // 将一个给定的文件描述符从集合中删除
// int FD_ISSET(int fd, fd_set *fdset);  // 检查集合中指定的文件描述符是否可以读写

#define MAXLINE 1024 //
#define LISTENQ 5    //
#define SIZE 10      //

typedef struct ServerContext
{
    int client_cnt;       // 客户端个数
    int client_fds[SIZE]; // 客户端套接字
    fd_set fds;           // 套接字集合
    int maxfd;            // 套接字最大值, select 需要的参数
} ServerContext;
static ServerContext *g_srv_ctx = NULL;

// TODO: 简化下列函数声明和调用
static int bind_socket(const char *ip, int port);

static int accept_conn(int srvfd);
static int do_write(int fd, char *buf);
static void do_read(fd_set *readfds);
static void serve(int srvfd);

static void sever_close(ServerContext *ctx);
static int server_init();

int main(int argc, char *argv[])
{
    // 初始化 g_srv_ctx
    if (server_init() < 0)
    {
        fprintf(stderr, "init server context failed.\n");
        return -1;
    }
    printf("done here");
    // 创建服务套接字，开始监听特定端口8080
    int srvfd = bind_socket("127.0.0.1", 8080);
    if (srvfd < 0)
    {
        fprintf(stderr, "socket create or bind fail.\n");
        goto err;
    }
    // 开始接收并处理客户端请求
    serve(srvfd);
    sever_close(g_srv_ctx);
    return 0;
err:
    sever_close(g_srv_ctx);
    return -1;
}

// 套接字编程
static int bind_socket(const char *ip, int port)
{
    // 创建套接字
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd == -1)
    {
        fprintf(stderr, "create socket failed, errno=%d, reason=%s\n",
                errno, strerror(errno));
        return -1;
    }

    // 一个端口释放后会等待两分钟之后才能再被使用，
    // SO_REUSEADDR是让端口释放后立即就可以被再次使用。
    int reuse = 1;
    if (setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &reuse, sizeof(reuse)) == -1)
    {
        return -1;
    }

    // 解析服务端地址
    struct sockaddr_in sockaddr;
    bzero(&sockaddr, sizeof(sockaddr));
    inet_pton(AF_INET, ip, &sockaddr.sin_addr);
    sockaddr.sin_family = AF_INET;
    sockaddr.sin_port = htons(port);

    // 绑定套接字(套接字和端口)
    if (bind(fd, (struct sockaddr *)&sockaddr, sizeof(sockaddr)) == -1)
    {
        perror("bind error: ");
        return -1;
    }

    // 监听套接字
    int code = listen(fd, LISTENQ);
    printf("bind_socket called listen got code =%d", code);
    return fd;
}

static int accept_conn(int srvfd)
{
    struct sockaddr_in cliaddr;
    socklen_t cliaddrlen;
    cliaddrlen = sizeof(cliaddr);
    int clifd = -1;

    printf("accpet clint proc is called.\n");

ACCEPT:
    clifd = accept(srvfd, (struct sockaddr *)&cliaddr, &cliaddrlen);

    if (clifd == -1)
    {
        if (errno == EINTR)
        {
            goto ACCEPT;
        }
        else
        {
            fprintf(stderr, "accept fail,error:%s\n", strerror(errno));
            return -1;
        }
    }

    fprintf(stdout, "accept a new client: %s:%d\n",
            inet_ntoa(cliaddr.sin_addr), cliaddr.sin_port);

    //将新的连接描述符添加到数组中
    int i = 0;
    for (i = 0; i < SIZE; i++)
    {
        if (g_srv_ctx->client_fds[i] < 0)
        {
            g_srv_ctx->client_fds[i] = clifd;
            g_srv_ctx->client_cnt++;
            break;
        }
    }

    if (i == SIZE)
    {
        fprintf(stderr, "too many clients.\n");
        return -1;
    }
    return 101;
}

static int do_write(int fd, char *buf)
{
    assert(buf);
    printf("recv buf is :%s\n", buf);
    write(fd, buf, strlen(buf) + 1);
    return 0;
}

static void do_read(fd_set *readfds)
{
    int i = 0, n = 0;
    int fd;
    char buf[MAXLINE] = {0};

    // 遍历套接字，判断套接字是否就绪
    for (i = 0; i <= g_srv_ctx->client_cnt; i++)
    {
        fd = g_srv_ctx->client_fds[i];
        if (fd < 0)
        {
            continue;
        }

        // 判断客户端套接字是否有数据
        if (FD_ISSET(fd, readfds))
        {
            // 接收客户端发送的信息
            n = read(fd, buf, MAXLINE);
            if (n <= 0)
            {
                // n == 0 表示读取完成，客户都关闭套接字
                FD_CLR(fd, &g_srv_ctx->fds);
                close(fd);
                g_srv_ctx->client_fds[i] = -1;
                continue;
            }
            do_write(fd, buf);
        }
    }
}
static void serve(int srvfd)
{
    printf("server entry");
    int clifd = -1;
    int retval = 0;
    fd_set *readfds = &g_srv_ctx->fds;
    struct timeval tv;
    int i = 0;

    while (1)
    {
        // 每次调用select前都要重新设置文件描述符和时间
        // 因为事件发生后，文件描述符和时间都被内核修改了
        FD_ZERO(readfds);

        // 添加监听套接字
        FD_SET(srvfd, readfds);
        g_srv_ctx->maxfd = srvfd;

        tv.tv_sec = 30;
        tv.tv_usec = 0;
        // 添加客户端套接字
        for (i = 0; i < g_srv_ctx->client_cnt; i++)
        {
            clifd = g_srv_ctx->client_fds[i];
            // 去除无效的客户端句柄
            if (clifd != -1)
            {
                FD_SET(clifd, readfds);
            }
            g_srv_ctx->maxfd = (clifd > g_srv_ctx->maxfd ? clifd : g_srv_ctx->maxfd);
        }

        // README：select 阻塞
        // 开始轮询接收处理服务端和客户端套接字
        retval = select(g_srv_ctx->maxfd + 1, readfds, NULL, NULL, &tv);
        if (retval == -1)
        {
            fprintf(stderr, "select error:%s.\n", strerror(errno));
            return;
        }
        if (retval == 0)
        {
            fprintf(stdout, "select is timeout.\n");
            continue;
        }

        // FIXME：?
        if (FD_ISSET(srvfd, readfds))
        {
            // 监听客户端请求
            accept_conn(srvfd);
        }
        else
        {
            // 接受处理客户端消息
            do_read(readfds);
        }
    }
}

// 释放 ServerContext g_srv_ctx 这一变量
static void sever_close(ServerContext *srvCtx)
{
    if (srvCtx != NULL)
    {
        free(srvCtx);
        srvCtx = NULL;
    }
}

// 初始化 ServerContext *g_srv_ctx 这一变量
static int server_init()
{
    // 创建一个ServerContext类型，申请空间
    g_srv_ctx = (ServerContext *)malloc(sizeof(ServerContext));
    if (g_srv_ctx == NULL)
    {
        return -1;
    }

    // 初始化
    memset(g_srv_ctx, 0, sizeof(ServerContext));

    int i;
    for (i = 0; i < SIZE; i++)
    {
        g_srv_ctx->client_fds[i] = -1;
    }
    printf("server init called");
    return 0;
}