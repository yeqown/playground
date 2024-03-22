package main

import (
	"flag"
	"net"
	"os"
	"syscall"
	"fmt"
	"time"
	"errors"
	"os/signal"
	"sync/atomic"
	"path/filepath"
	"strconv"
)

var (
	processPrefix = flag.String("prefix", "", "-prefix to mark old or new process")
	port          = flag.Int("port", 9000, "-port to use specified listen port")
)

type gracefulTcpServer struct {
	listener     *net.TCPListener
	shutdownChan chan struct{}
	conns        map[net.Conn]struct{}

	servingConnCount atomic.Int32
	serveRunning     atomic.Bool
}

func start(port int) (*gracefulTcpServer, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &gracefulTcpServer{
		listener:         ln.(*net.TCPListener),
		shutdownChan:     make(chan struct{}, 1),
		conns:            make(map[net.Conn]struct{}, 16),
		servingConnCount: atomic.Int32{},
		serveRunning:     atomic.Bool{},
	}

	return s, nil
}

func startFromFork() (*gracefulTcpServer, error) {
	tmpdir := os.TempDir()

	var (
		nfd int = 0
		err error
	)

	if nfdStr := os.Getenv(__GRACE_ENV_NFDS); nfdStr == "" {
		panic("not nfds env")
	} else if nfd, err = strconv.Atoi(nfdStr); err != nil {
		panic(err)
	}

	// restore conn fds, 0, 1, 2 has been used by os.Stdin, os.Stdout, os.Stderr
	lfd := os.NewFile(3, filepath.Join(tmpdir, "graceful"))
	ln, err := net.FileListener(lfd)
	if err != nil {
		panic(err)
	}

	s := &gracefulTcpServer{
		listener:         ln.(*net.TCPListener),
		shutdownChan:     make(chan struct{}, 1),
		conns:            make(map[net.Conn]struct{}, 16),
		servingConnCount: atomic.Int32{},
		serveRunning:     atomic.Bool{},
	}

	for i := 0; i < nfd; i++ {
		fd := os.NewFile(uintptr(4+i), filepath.Join(tmpdir, strconv.Itoa(4+i)))
		conn, err := net.FileConn(fd)
		if err != nil {
			panic("restore conn failed: " + err.Error())
		}

		go s.handleConn(conn)
	}

	return s, nil
}

func (s *gracefulTcpServer) serve() {
	fmt.Printf(*processPrefix+"| server listening: %d\n", *port)
	s.serveRunning.Store(true)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				break
			} else {
				fmt.Printf(*processPrefix + "| server accept error")
			}
		}

		go s.handleConn(conn)
	}

	s.serveRunning.Store(false)
}

func (s *gracefulTcpServer) handleConn(conn net.Conn) {
	fmt.Printf(*processPrefix+"| handling conn: %d\n", *port)
	ticker := time.NewTicker(3 * time.Second)
	data := []byte(*processPrefix + "| sending\n")

	s.conns[conn] = struct{}{}

	s.servingConnCount.Add(1)
	defer s.servingConnCount.Add(-1)

	for {
		select {
		case <-ticker.C:
			_, _ = conn.Write(data)
		case <-s.shutdownChan:
			return
		}
	}
}

func (s *gracefulTcpServer) gracefulRestart() {
	_ = s.listener.SetDeadline(time.Now())
	lfd, err := s.listener.File()
	if err != nil {
		panic(err)
	}

	os.Setenv(__GRACE_ENV_FLAG, "true")
	os.Setenv(__GRACE_ENV_NFDS, strconv.Itoa(len(s.conns)))

	files := make([]uintptr, 4, 3+1+len(s.conns))
	copy(files[:4], []uintptr{
		os.Stdin.Fd(),
		os.Stdout.Fd(),
		os.Stderr.Fd(),
		lfd.Fd(),
	})

	for conn := range s.conns {
		connFd, _ := conn.(*net.TCPConn).File()
		files = append(files, connFd.Fd())
	}

	procAttr := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: files,
		Sys:   nil,
	}

	childPid, err := syscall.ForkExec(os.Args[0], os.Args, procAttr)
	if err != nil {
		panic(err)
	}

	fmt.Println("child process is running with pid: ", childPid)
}

func (s *gracefulTcpServer) waitForSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGHUP)

	// wait for os signals
	sig := <-sigChan
	if sig == syscall.SIGHUP {
		s.gracefulRestart()
	}

	s.shutdownChan <- struct{}{}

	// wait serve quit
	for s.serveRunning.Load() {
		time.Sleep(10 * time.Millisecond)
	}
	// wait all handleConn quit
	var runningConn = int32(1)
	for ; runningConn != 0; runningConn = s.servingConnCount.Load() {
		time.Sleep(10 * time.Millisecond)
	}

	// other signal, just quit
	fmt.Println("server quit")
}

const (
	__GRACE_ENV_FLAG       = "GRACEFUL_RESTART"
	__GRACE_ENV_NFDS       = "GRACEFUL_NFDS"
	__GRACE_ENV_NFDS_START = "GRACEFUL_NFDS_START"
)

func main() {
	flag.Parse()

	var (
		s   *gracefulTcpServer
		err error
	)

	if v := os.Getenv(__GRACE_ENV_FLAG); v != "" {
		s, err = startFromFork()
	} else {
		s, err = start(*port)
	}
	if err != nil {
		panic(err)
	}

	go s.serve()

	s.waitForSignals()
}
