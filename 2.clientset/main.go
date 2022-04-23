package main

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// clientset
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploys, err := clientset.AppsV1().
		Deployments("default").
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, item := range deploys.Items {
		fmt.Println(item.Name)
	}
}
