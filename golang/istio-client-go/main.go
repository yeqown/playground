package main

import (
	"log"

	istio "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	mcs := buildMyClientSet()
	_ = mcs
}

func buildMyClientSet() myClientSet {
	config, err := clientcmd.BuildConfigFromFlags("", "~/.kube/config")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new clientset for the given config.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new clientset for the given config.
	istioclientset, err := istio.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new clientset for the given config.
	myclientset := myClientSet{
		k8s:   clientset,
		istio: istioclientset,
	}

	return myclientset
}

type myClientSet struct {
	k8s   *kubernetes.Clientset
	istio *istio.Clientset
}

// reviews version 增加一个版本：其部署在 另一个 namespace(v2) 中
// 1. 添加到 istio 的 virtualservice 中 v2 的 match 配置
// - match:
//   - headers:
//   x-mesh-swimlane:
//     exact: swimlane-v1.2.3
//   route:
//   - destination:　
//     host: reviews.bookinfo-v2.svc.cluster.local
func deployVersion(mcs myClientSet) {

}

// reviews version 释放一个版本：其部署在 另一个 namespace(v2) 中
// 1. 从 istio 的 virtualservice 中删除 v2 的 match 配置
func releaseVersion(mcs myClientSet) {

}
