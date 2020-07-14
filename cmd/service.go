package cmd

import (
	"context"
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

type Port struct {
	Port int64
	NodePort int64
}

type ServiceResponse struct {
	Ports map[string]int64
}

//ServiceManager K8s service wrapper interface
type ServiceManager interface {
	CreateService(cfg, namespace, name string, port uint16) (ServiceResponse, error)
	DeleteService(cfg, namespace, name string) error
}

//ServiceManagerImpl ServiceManager implementation
type ServiceManagerImpl struct {
}

//CreateService asdsad
func (sm *ServiceManagerImpl) CreateService(cfg, namespace, name string, port uint16) (ServiceResponse, error) {

	res := ServiceResponse{}

	namespace1, client, err := configSetup(cfg, namespace)

	if err != nil {
		return res, err
	}

	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}

	service := sm.createServiceFromTemplate(namespace1, name, port)

	// Create service
	var result *unstructured.Unstructured

	result, err = client.Resource(serviceRes).Namespace(namespace1).Create(context.TODO(), service, metav1.CreateOptions{})

	if err != nil {
		return res, err
	}

	res.Ports, err = unwrapNodePort(result);

	return res,err
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