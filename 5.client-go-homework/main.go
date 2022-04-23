package main

import (
	"baidingK8s/5.client-go-homework/pkg"
	"fmt"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1，获取到config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	// 2，获取k8s client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 3，获取informer factory
	factory := informers.NewSharedInformerFactory(clientset, 0)

	// 获取service的informer
	serviceInformer := factory.Core().V1().Services()
	// 获取ingress的informer
	ingressInformer := factory.Networking().V1().Ingresses()

	// 获取控制器
	controller := pkg.NewController(clientset, serviceInformer, ingressInformer)

	// 5，启动informer
	stopCh := make(chan struct{})
	fmt.Println("启动infromer中。。。。")
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)
	controller.Run(stopCh)
}
