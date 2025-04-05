package k8s

import (
	"context"
	"fmt"
	"imageFlux/internal/registry"
	"log"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeployment() {
	clientSet, err := GetKubeClient()
	if err != nil {
		log.Printf("Error getting kube client: %v", err)

	}

	deployments, err := clientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error getting deployments: %v", err)
	}

	for _, deployment := range deployments.Items {
		annotations := deployment.GetAnnotations()
		if val := annotations["ImageFlux"]; val == "Enabled" {

			image := deployment.Spec.Template.Spec.Containers[0].Image
			imageParts := strings.Split(image, ":")
			if len(imageParts) == 2 {
				imageName := imageParts[0]
				imageTag := imageParts[1]
				fmt.Println("Image Name:", imageName)
				fmt.Println("Image Tag:", imageTag)
				currentTag, err := registry.GetImageTags(imageName)
				if err != nil {
					log.Printf("Failed to fetach image tag for image %s: %v", imageName, err)
					continue
				}
				if currentTag != imageTag {
					fmt.Println("Image tag is not the latest:")
					deployment.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s", imageName, currentTag)
					_, err := clientSet.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), &deployment, metav1.UpdateOptions{})
					if err != nil {
						log.Printf("Error updating deployment image: %v", err)
					} else {
						log.Printf("Deployment image updated to the latest tag: %v", currentTag)
					}
				} else {
					log.Printf("Image tag is already the latest: %v", imageTag)
				}

			} else {
				log.Printf("Image does not have a tag: %v", image)
			}
		}

	}
}
