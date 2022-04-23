package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 使用dynamic client 可以操作内部client，也可以操作自定资源
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	unstructObj, err := client.Resource(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}).Namespace("default").
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	podList := &v1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), podList)
	if err != nil {
		panic(err)
	}

	for _, v := range podList.Items {
		fmt.Println(v.Name)
	}
}
