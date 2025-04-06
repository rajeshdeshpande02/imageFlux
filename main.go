package main

import (
	"imageFlux/internal/controller"
	"imageFlux/internal/k8s"
	"log"
)

func main() {
	clientSet, err := k8s.GetKubeClient()
	if err != nil {
		log.Fatalf("Error getting kube client: %v", err)
	}
	stopCh := make(chan struct{})
	controller.Start(clientSet, stopCh)
}
