// Note: the example only works with the code within the same release/branch.
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

//DeploymentManager K8s deployment wrapper interface
type DeploymentManager interface {
	CreateDeployment(cfg, namespace, image, name string) (string, error)
	DeleteDeployment(cfg, namespace, name string) (string, error)
}

//DeploymentManagerImpl DeploymentManager implementation
type DeploymentManagerImpl struct {
}

//CreateDeployment creates a kubernetes deployment with the given parameters
func (dp *DeploymentManagerImpl) CreateDeployment(cfg, namespace, image, name string) (string, error) {

	namespace, client, err := configSetup(cfg, namespace)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	deployment := createDeploymentFromTemplate(namespace, image, name)

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := client.Resource(deploymentRes).Namespace(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetName())

	return "foobar", err
}

func (dp *DeploymentManagerImpl) listDeployments(err error, client dynamic.Interface, deploymentRes schema.GroupVersionResource, namespace string) {
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := client.Resource(deploymentRes).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		replicas, found, err := unstructured.NestedInt64(d.Object, "spec", "replicas")
		if err != nil || !found {
			fmt.Printf("Replicas not found for deployment %s: error=%s", d.GetName(), err)
			continue
		}
		fmt.Printf(" * %s (%d replicas)\n", d.GetName(), replicas)
	}
}

//DeleteDeployment deletes the specified deployment
func (dp *DeploymentManagerImpl) DeleteDeployment(cfg, namespace, name string) (string, error) {
	namespace, client, err := configSetup(cfg, namespace)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := client.Resource(deploymentRes).Namespace(namespace).Delete(context.TODO(), name, deleteOptions); err != nil {
		panic(err)
	}

	fmt.Println("Deleted deployment.")
	return "foobar", err
}

func createDeploymentFromTemplate(namespace, image, name string) *unstructured.Unstructured {
	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": name,
				"labels": map[string]interface{}{
					"app": name,
				},
			},
			"spec": map[string]interface{}{
				"replicas": 1,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": name,
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": name,
						},
					},

					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  name,
								"image": image,
								/*								"ports": []map[string]interface{}{
																{
																	"name":          "http",
																	"protocol":      "TCP",
																	"containerPort": 80,
																},
															},*/
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
