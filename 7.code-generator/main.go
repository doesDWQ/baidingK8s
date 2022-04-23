package main

import (
	clientset "baidingK8s/7.code-generator/pkg/generated/clientset/versioned"
	"baidingK8s/7.code-generator/pkg/generated/informers/externalversions"
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// 严重自定义资源是否生效
func main() {

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	list, err := clientset.CrdV1().Foos("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, foo := range list.Items {
		fmt.Println(foo.Name)
	}

	// 监听改动事件
	factory := externalversions.NewSharedInformerFactory(clientset, 0)
	factory.Crd().V1().Foos().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("增加了一个自定义资源。。。。")
		},
	})

	stopCh := make(chan struct{})
	factory.Start(stopCh)
	<-stopCh
}
