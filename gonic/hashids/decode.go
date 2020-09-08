package main

import (
	"flag"

	"github.com/speps/go-hashids"
)

func main() {
	var id string
	flag.StringVar(&id, "id", "gY", "specify id")
	flag.Parse()

	h, _ := hashids.New()
	patientIDs, err := h.DecodeWithError(id)

	if err != nil || len(patientIDs) != 1 {
		panic("解析患者ID失败")
	}

	println(patientIDs[0])
}