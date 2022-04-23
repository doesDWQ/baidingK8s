package main

import (
	"fmt"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// discoveryClient 获取到 api-version 相当于 kubectl api-versions
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}

	// 获取gvr
	a, al, err := dc.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}

	fmt.Println("g=======================================")
	// 打印出所有的api group
	for _, v := range a {
		fmt.Printf("group:%s\n", v.Name)
	}
	fmt.Println("gv=======================================")
	for _, v := range al {
		fmt.Printf("gv:%s=======================================\n\n", v.GroupVersion)

		for _, item := range v.APIResources {
			fmt.Printf("resourceName:%s \n", item.Name)
		}

	}
}
