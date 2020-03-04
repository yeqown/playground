package main

import (
	"crypto/rc4"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	cipher, err := rc4.NewCipher([]byte("thisiskey"))
	if err != nil {
		log.Fatalf("wrong with NewCipher: %v", err)
	}

	c := map[string]bool{}
	for i := 0; i < 10000000; i++ {
		src := []byte(fmt.Sprintf("%7d", i))
		dst := make([]byte, len(src))
		cipher.XORKeyStream(dst, src)
		s := toString(dst) // 密文是不可读的字节流，这里采用hex编码
		// println(s)         // 形如：4ae90ca9
		c[s] = true
	}

	println(len(c))
}

func toString(src []byte) string {
	return hex.EncodeToString(src)
}
