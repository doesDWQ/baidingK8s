// +groupName=baiding.tech
package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	Scheme = runtime.NewScheme()
	// 设置gv
	GroupVersion = schema.GroupVersion{Group: "baiding.tech", Version: "v1"}
	// 解码器
	Codecs = serializer.NewCodecFactory(Scheme)
)
