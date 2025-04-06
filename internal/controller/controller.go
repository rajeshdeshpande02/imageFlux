package controller

import (
	"context"
	"fmt"
	"imageFlux/internal/registry"
	"log"
	"strings"
	"sync"
	"time"

	appsv1 "k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	imageFluxAnnotation = "ImageFlux"
	checkInterval       = 60 * time.Second
)

func Start(clientSet *kubernetes.Clientset, stopCh <-chan struct{}) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			log.Println("Shutting down controller...")
			return
		case <-ticker.C:
			processDeployments(clientSet)
		}
	}
}

func processDeployments(clientSet *kubernetes.Clientset) {
	deployments, err := clientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error getting deployments: %v", err)
		return
	}

	var wg sync.WaitGroup
	for _, deployment := range deployments.Items {
		annotations := deployment.GetAnnotations()
		if annotations[imageFluxAnnotation] != "Enabled" {
			continue
		}

		wg.Add(1)
		go func(deploymentCopy appsv1.Deployment) {
			defer wg.Done()
			updateDeploymentIfNeeded(clientSet, &deploymentCopy)
		}(deployment)
	}
	wg.Wait()
}

func updateDeploymentIfNeeded(clientSet *kubernetes.Clientset, deployment *appsv1.Deployment) {
	containers := deployment.Spec.Template.Spec.Containers
	for i, container := range containers {
		imageParts := strings.Split(container.Image, ":")
		if len(imageParts) != 2 {
			log.Printf("Image does not have a tag: %v", container.Image)
			continue
		}

		imageName, imageTag := imageParts[0], imageParts[1]
		currentTag, err := registry.GetImageTags(imageName)
		if err != nil {
			log.Printf("Failed to fetch image tag for image %s: %v", imageName, err)
			continue
		}

		if currentTag != imageTag {
			log.Printf("Updating image for deployment %s: %s -> %s", deployment.Name, imageTag, currentTag)
			containers[i].Image = fmt.Sprintf("%s:%s", imageName, currentTag)
			deployment.Spec.Template.Spec.Containers = containers
			_, err := clientSet.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
			if err != nil {
				log.Printf("Error updating deployment image: %v", err)
			} else {
				log.Printf("Deployment %s image updated to the latest tag: %v", deployment.Name, currentTag)
			}
		} else {
			log.Printf("Image tag is already the latest for deployment %s: %v", deployment.Name, imageTag)
		}
	}
}
