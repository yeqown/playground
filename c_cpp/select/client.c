// Copyright 2020 The yeqown Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#include <stdio.h>
#include <errno.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>
#include <sys/select.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <arpa/inet.h>

#define MAXLINE 1024

#define max(a, b) (a > b) ? a : b

static void do_write(int sockfd, char *buf);
static void handle_connection(int sockfd);

int main(int argc, char *argv[])
{
    int sockfd = socket(AF_INET, SOCK_STREAM, 0);

    struct sockaddr_in servaddr;
    bzero(&servaddr, sizeof(servaddr));
    servaddr.sin_family = AF_INET;
    servaddr.sin_port = htons(8080);
    inet_pton(AF_INET, "127.0.0.1", &servaddr.sin_addr);

    int retval = connect(sockfd, (struct sockaddr *)&servaddr, sizeof(servaddr));
    if (retval < 0)
    {
        fprintf(stderr, "connect fail,error:%s\n", strerror(errno));
        return -1;
    }

    printf("client send to server .\n");
    write(sockfd, "hello server", 32);
    handle_connection(sockfd);

    return 0;
}

static void do_write(int sockfd, char *buf)
{
    printf("client recv msg is:%s\n", buf);
    sleep(5);
    write(sockfd, buf, strlen(buf) + 1);
}

static void handle_connection(int sockfd)
{
    char sendline[MAXLINE], recvline[MAXLINE];
    int maxfdp, stdineof;
    fd_set readfds;
    int n;
    struct timeval tv;
    int retval = 0;

    while (1)
    {
        FD_ZERO(&readfds);
        FD_SET(sockfd, &readfds);
        maxfdp = sockfd;

        tv.tv_sec = 5;
        tv.tv_usec = 0;

        retval = select(maxfdp + 1, &readfds, NULL, NULL, &tv);

        if (retval == -1)
        {
            return;
        }

        if (retval == 0)
        {
            printf("client timeout.\n");
            continue;
        }

        if (FD_ISSET(sockfd, &readfds))
        {
            n = read(sockfd, recvline, MAXLINE);
            if (n <= 0)
            {
                fprintf(stderr, "client: server is closed.\n");
                close(sockfd);
                FD_CLR(sockfd, &readfds);
                return;
            }
            do_write(sockfd, recvline);
        }
    }
}
