package main

import "github.com/nxadm/tail"

func main() {
	t, err := tail.TailFile("./foo", tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		panic(err)
	}

	for line := range t.Lines {
		println(line.Text)
	}
}
