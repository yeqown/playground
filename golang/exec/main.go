package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
	"syscall"
)

var (
	addrFlag = flag.String("addr", ":8080", "the address to listen on")
)

func main() {
	flag.Parse()

	if err := serve(*addrFlag); err != nil {
		panic(err)
	}
}

// start a new http server with a handler to exec the command:
// "./main"
// and return the output
func serve(addr string) error {
	http.HandleFunc("/serve", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()

		addr := r.Form.Get("addr")
		if addr == "" {
			_, _ = w.Write([]byte("no addr"))
			return
		}

		cmd := exec.Command("/bin/sh", "-c", "./main", "-addr", addr)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			//Setsid:  true,
			Setpgid: true,
		}
		out, err := cmd.CombinedOutput()
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		// write the output to the http response
		_, _ = w.Write(out)
	})

	log.Printf("Listening on: %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
