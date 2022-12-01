cmake_minimum_required(VERSION 3.22)
project(c_cpp)

set(CMAKE_CXX_STANDARD 14)

add_executable(c_cpp
        epoll/client.c
        epoll/server.c
        kqueue/client.c
        kqueue/server.c
        select/client.c
        select/server.c
        hello.c
        hello.cpp)
