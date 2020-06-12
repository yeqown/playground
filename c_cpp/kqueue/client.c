// Copyright 2020 The yeqown Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
// #include <time.h>
#include <unistd.h>

// #include <sys/epoll.h>
#include <sys/event.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <errno.h>
#include <err.h>
#include <arpa/inet.h>

#define MAXSIZE 1024
#define IPADDRESS "127.0.0.1"
#define SERV_PORT 8080
#define FDSIZE 1024
#define EPOLLEVENTS 20

int dial(char *addr, int port);
// void recv_loop(int kq);

int main()
{
    int client_fd = dial(IPADDRESS, SERV_PORT);
    if (client_fd < 0)
    {
        fprintf(stderr, "dial failed\n");
        return -1;
    }
    int kq = kqueue();
    write(client_fd, "hello server", 32);

    struct kevent evt;                                                    // 创建
    EV_SET(&evt, client_fd, EVFILT_READ, EV_ADD | EV_ENABLE, 0, 0, NULL); // 赋值
    int ret = kevent(kq, &evt, 1, NULL, 0, NULL);
    if (ret == -1)
    {
        err(EXIT_FAILURE, "kevent register");
    }
    if (evt.flags & EV_ERROR)
    {
        errx(EXIT_FAILURE, "Event error: %s", strerror(evt.data));
    }

    char buf[MAXSIZE];
    int n;
    while (1)
    {
        ret = kevent(kq, &evt, 1, NULL, 0, NULL);
        if (ret == -1)
        {
            err(EXIT_FAILURE, "kevent wait");
        }
        else if (ret > 0)
        {
            printf("Something was written");
            n = read(client_fd, buf, MAXSIZE);
            if (n <= 0)
            {
                close(client_fd);
                return 0;
            }
            printf("client recv msg is:%s\n", buf);
        }
    }

    return 0;
}

// int dial(char *addr, int port)
// ret = -1, means dial failed
//     else means fd
int dial(char *addr, int port)
{
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd == -1)
    {
        fprintf(stderr, "create socket failed, errno=%d, reason=%s\n",
                errno, strerror(errno));
        return -1;
    }

    struct sockaddr_in sockaddr;
    bzero(&sockaddr, sizeof(sockaddr));
    sockaddr.sin_family = AF_INET;
    sockaddr.sin_port = htons(8080);
    inet_pton(AF_INET, addr, &sockaddr.sin_addr);

    int ret = connect(fd, (struct sockaddr *)&sockaddr, sizeof(sockaddr));
    if (ret < 0)
    {
        fprintf(stderr, "connect failed, err=%s\n", strerror(errno));
        return -1;
    }

    return fd;
}

// void recv_loop(int kq)
// {
//     int kq;
//     while (1)
//     {
//         kq = kqueue();
//     }
// }