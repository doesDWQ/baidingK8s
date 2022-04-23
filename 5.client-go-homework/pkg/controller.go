package pkg

import (
	"context"
	"fmt"
	"reflect"
	"time"

	v14 "k8s.io/api/core/v1"
	v12 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informer "k8s.io/client-go/informers/core/v1"
	netInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	coreLister "k8s.io/client-go/listers/core/v1"
	v1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	workNum  = 5
	maxRetry = 10
)

// 控制器定义
type controller struct {
	client        kubernetes.Interface
	ingressLister v1.IngressLister
	serviceLister coreLister.ServiceLister
	queue         workqueue.RateLimitingInterface
}

// 更新service
func (c *controller) updateService(oldObj interface{}, newObj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err != nil {
		runtime.HandleError(err)
	}

	namespaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(err)
	}

	fmt.Printf("update service:%s:%s\n", namespaceKey, name)

	// 比较 annotation
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}

	c.enqueue(newObj)
}

// 添加service处理
func (c *controller) addService(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}

	namespaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(err)
	}

	fmt.Printf("add service:%s:%s\n", namespaceKey, name)
	c.enqueue(obj)
}

func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}

	c.queue.Add(key)
}

func (c *controller) deleteIngress(obj interface{}) {
	ingress := obj.(*v12.Ingress)
	ownerReference := v13.GetControllerOf(ingress)

	if ownerReference == nil {
		return
	}

	if ownerReference.Kind != "service" {
		return
	}

	c.queue.Add(fmt.Sprintf("%s/%s", ingress.Namespace, ingress.Name))
}

func (c *controller) Run(stopCh chan struct{}) {
	for i := 0; i < workNum; i++ {
		go wait.Until(c.worker, time.Minute, stopCh)
	}
	<-stopCh
}

func (c *controller) worker() {
	for c.processNextItem() {

	}
}

func (c *controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	defer c.queue.Done(item)

	key := item.(string)

	err := c.syncService(key)

	if err != nil {
		c.handlerError(key, err)
	}

	return true
}

func (c *controller) handlerError(key string, err error) {
	if c.queue.NumRequeues(key) <= maxRetry {
		c.queue.AddRateLimited(key)
		return
	}

	runtime.HandleError(err)

	c.queue.Forget(key)
}

func (c *controller) syncService(key string) error {
	namespaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Println("此处报错了5")
		return err
	}

	// 删除
	service, err := c.serviceLister.Services(namespaceKey).Get(name)
	if err != nil {
		fmt.Println("此处报错了4")
		return err
	}

	// 新增和删除
	_, ok := service.GetAnnotations()["ingress/http"]
	ingress, err := c.ingressLister.Ingresses(namespaceKey).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		fmt.Println("此处报错了1")
		return err
	}

	if ok && errors.IsNotFound(err) {
		// create ingress
		ig := c.constructIngress(service)
		_, err := c.client.NetworkingV1().Ingresses(namespaceKey).Create(context.TODO(), ig, v13.CreateOptions{})
		if err != nil {
			fmt.Println("此处报错了2")
			return err
		}
	} else if !ok && ingress != nil {
		// delete ingress
		err := c.client.NetworkingV1().Ingresses(namespaceKey).Delete(context.TODO(), name, *&v13.DeleteOptions{})
		if err != nil {
			fmt.Println("此处报错了3")
			return err
		}
	}

	return nil
}

func (c *controller) constructIngress(service *v14.Service) *v12.Ingress {
	ingress := &v12.Ingress{}

	ingress.ObjectMeta.OwnerReferences = []v13.OwnerReference{*v13.NewControllerRef(service, v14.SchemeGroupVersion.WithKind(("service")))}
	ingress.Name = service.Name
	ingress.Namespace = service.Namespace
	pathType := v12.PathTypePrefix

	icn := "nginx"

	ingress.Spec = v12.IngressSpec{
		IngressClassName: &icn,
		Rules: []v12.IngressRule{
			{
				Host: "example.com",
				IngressRuleValue: v12.IngressRuleValue{
					HTTP: &v12.HTTPIngressRuleValue{
						Paths: []v12.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: v12.IngressBackend{
									Service: &v12.IngressServiceBackend{
										Name: service.Name,
										Port: v12.ServiceBackendPort{
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return ingress
}

// 自制controller
func NewController(client kubernetes.Interface, serviceInformer informer.ServiceInformer, ingressInformer netInformer.IngressInformer) controller {
	// 获取控制器
	c := controller{
		client:        client,
		ingressLister: ingressInformer.Lister(),
		serviceLister: serviceInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}

	// service回调设置
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	// ingress回调设置
	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})

	return c
}
