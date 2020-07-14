package cmd

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//UnwrapNodePort Retrieves the nodeport assined from the k8s api result
func unwrapNodePort(result *unstructured.Unstructured) (map[string]int64,error) {
	res := make(map[string]int64)

	keys := []string{"nodePort","targetPort"}
	ports, ok, err := unstructured.NestedSlice(result.UnstructuredContent(), "spec", "ports")

	if ok && err == nil && len(ports) > 0 {
		for _, ele := range ports {
			port, _ := ele.(map[string]interface{})
			for _, key := range keys {
				num ,_ := port[key].(int64)
				res[key] = num
			}
		}
	}

	return res, nil
}

