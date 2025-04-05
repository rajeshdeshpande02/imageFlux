package k8s

import (
	"context"
	"fmt"
	"imageFlux/internal/registry"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeployment() {
	clientSet, err := GetKubeClient()
	if err != nil {
		fmt.Errorf("Error getting kube client: %v", err)

	}
	deployments, err := clientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Errorf("Error getting deployments: %v", err)
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
				currentTag, _ := registry.GetImageTags(imageName)
				if currentTag != imageTag {
					fmt.Println("Image tag is not the latest:")
					deployment.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s", imageName, currentTag)
					_, err := clientSet.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), &deployment, metav1.UpdateOptions{})
					if err != nil {
						fmt.Errorf("Error updating deployment image: %v", err)
					} else {
						fmt.Println("Deployment image updated to the latest tag:", currentTag)
					}
				} else {
					fmt.Println("Image tag is already the latest:", imageTag)
				}

			} else {
				fmt.Println("Image does not have a tag:", image)
			}
		}

	}
}
