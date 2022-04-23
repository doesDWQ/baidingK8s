package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

func main() {
	// 直接使用rest client的方式获取到 pod 信息
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	config.APIPath = "/api"

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	podList := &v1.PodList{}
	err = restClient.Get().
		Namespace("default").
		Resource("pods").
		Do(context.TODO()).Into(podList)
	if err != nil {
		panic(err)
	}

	for _, v := range podList.Items {
		fmt.Println(v.Name)
	}
}
