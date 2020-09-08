package main

import (
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"fmt"
)

func testMD() {
	ctx0 := context.Background()

	ctx1 := metadata.AppendToOutgoingContext(ctx0, "key", "value")
	val1 := valueOfGRPCContext(ctx1, "key")
	fmt.Printf("val1=%s\n", val1)

	// ctx2 := metadata.AppendToOutgoingContext(ctx1, "key", "value2")
	// val2 := valueOfGRPCContext(ctx2, "key")
	// fmt.Printf("val2=%s\n", val2)
}

// 从 GRPC 上下文中获取 key 对应的值
func valueOfGRPCContext(ctx context.Context, key string) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	delete(md, "key")

	if !ok {
		return ""
	}

	for _, data := range md.Get(key) {
		data = strings.TrimSpace(data)
		if data == "" {
			continue
		}
		return data
	}

	return ""
}

func main() {
	testMD()
}
