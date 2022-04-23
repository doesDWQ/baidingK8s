1, 下载cotroller到本地
 git clone https://github.com/kubernetes-sigs/controller-tools.git

2, 生成type模板
type-scaffold --kind Foo > .foo_type.template.txt

3，生成deepcopy等方法
controller-gen object paths=./pkg/apis/baiding.teach/v1/types.go

4, 生成crd
controller-gen crd paths=./... output:crd:dir=config/crd

5, go run main.go