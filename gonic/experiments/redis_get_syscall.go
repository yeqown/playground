package main

import "github.com/go-redis/redis"

func main() {
	c := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	val, err := c.Get("nonexist").Result()

	_ = val
	_ = err
}

// dtruss -c ./main
// count syscall of the once redis get.
//
//
//CALL                                        COUNT
//access                                          1
//bsdthread_register                              1
//connect                                         1
//csops                                           1
//exit                                            1
//fstat64                                         1
//getpeername                                     1
//getpid                                          1
//getsockname                                     1
//getsockopt                                      1
//getuid                                          1
//issetugid                                       1
//kqueue                                          1
//pipe                                            1
//proc_info                                       1
//psynch_mutexdrop                                1
//psynch_mutexwait                                1
//shm_open                                        1
//bind                                            2
//csops_audittoken                                2
//getegid                                         2
//getentropy                                      2
//gettid                                          2
//ioctl                                           2
//geteuid                                         3
//open                                            3
//sysctl                                          3
//write                                           3
//__pthread_kill                                  4
//sigreturn                                       4
//socket                                          4
//stat64                                          4
//bsdthread_create                                5
//thread_selfid                                   5
//read                                            6
//setsockopt                                      6
//close                                           7
//mprotect                                        8
//kevent                                         10
//sigaltstack                                    12
//psynch_cvsignal                                13
//__pthread_sigmask                              17
//psynch_cvwait                                  17
//mmap                                           21
//fcntl                                          22
//sigaction                                      50
//__semwait_signal                               62
//madvise                                        86
