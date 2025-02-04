package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pableeee/k8s-simple-wrapper/cmd"
)

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

func main() {
	deployment := cmd.DeploymentManagerImpl{}

	res, err := deployment.CreateDeployment("", "default", "nginx", "nginx")

	if err != nil {
		fmt.Println("Hubo un error")
		os.Exit(1)
	}

	fmt.Println("Deployment Created")
	fmt.Println(res)

	service := cmd.ServiceManagerImpl{}
	var port cmd.ServiceResponse

	port, err = service.CreateService("", "default", "nginx", 80)
	fmt.Println("Service Created")
	fmt.Println(port)

	prompt()

	err = service.DeleteService("","default","nginx")

	if err != nil {
		fmt.Println("Hubo un error")
		os.Exit(1)
	}
	fmt.Println("Service Deleted")

	err = deployment.DeleteDeployment("", "default", "nginx")

	if err != nil {
		fmt.Println("Hubo un error")
		os.Exit(1)
	}
	fmt.Println("Deployment Deleted")
}
