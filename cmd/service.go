package cmd

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

//

//
// Uncomment to load all auth plugins
// _ "k8s.io/client-go/plugin/pkg/client/auth"
//
// Or uncomment to load specific auth plugins
// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"

//ServiceManager K8s service wrapper interface
type ServiceManager interface {
	CreateService(cfg, namespace, name string, port uint16) (uint16, error)
	DeleteService(cfg, namespace, name string) error
}

//ServiceManagerImpl ServiceManager implementation
type ServiceManagerImpl struct {
}

//CreateService asdsad
func (sm *ServiceManagerImpl) CreateService(cfg, namespace, name string, port uint16) (uint16, error) {

	namespace1, client, err := configSetup(cfg, namespace)

	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}

	service := sm.createServiceFromTemplate(namespace1, name, port)

	// Create service
	fmt.Println("Creating service...")
	result, err := client.Resource(serviceRes).Namespace(namespace1).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	value, ok, errs := unstructured.NestedSlice(result.UnstructuredContent(), "spec", "ports")

	if ok && errs == nil && len(value) == 1 {
		v, k := value[0].(map[string]interface{})
		if k {
			fmt.Println(v["nodePort"])
		}
	}

	//fmt.Printf("Created service %q.\n", result.GetName())

	return 1, err
}

func (sm *ServiceManagerImpl) createServiceFromTemplate(namespace, name string, port uint16) *unstructured.Unstructured {
	service := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
				"labels": map[string]interface{}{
					"app": name,
				},
			},
			"spec": map[string]interface{}{
				"ports": []map[string]interface{}{
					{
						"protocol":   "TCP",
						"port":       port,
						"targetPort": port,
					},
				},
				"selector": map[string]interface{}{
					"app": name,
				},
				"type": "NodePort",
			},
			"status": map[string]interface{}{
				"loadBalancer": map[string]interface{}{},
			},
		},
	}
	return service
}

//UnwrapNodePort aasdasd
func (sm *ServiceManagerImpl) UnwrapNodePort(value unstructured.Unstructured) int64 {
	value, ok, errs := unstructured.NestedSlice(value.UnstructuredContent(), "spec", "ports")

	if ok && errs == nil && len(value) == 1 {
		v, k := value[0].(map[string]interface{})
		if k {
			fmt.Println(v["nodePort"])
		}
	}
}
