package main

import (
	v1 "baidingK8s/8.controller-tools/pkg/apis/baiding.tech/v1"
	"context"
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	config.APIPath = "/apis"
	config.NegotiatedSerializer = v1.Codecs.WithoutConversion()
	config.GroupVersion = &v1.GroupVersion

	client, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	foo := &v1.Foo{}
	err = client.Get().Namespace("default").Resource("foos").Name("crd-test").Do(context.TODO()).Into(foo)
	if err != nil {
		panic(err)
	}

	newObj := foo.DeepCopy()
	newObj.Spec.Name = "test2"

	fmt.Println(foo.Spec.Name)
	fmt.Println(foo.Spec.Replicas)

	fmt.Println(newObj.Spec.Name)
}
