package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TagsResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func GetImageTags(imageName string) (string, error) {
	url := "https://registry.hub.docker.com/v2/repositories/library/" + imageName + "/tags"

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error getting image tags")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		return "", fmt.Errorf("Error: received non-200 response code")
	}
	var tagsResponse TagsResponse

	if err := json.NewDecoder(resp.Body).Decode(&tagsResponse); err != nil {
		return "", fmt.Errorf("Error decoding respo")
	}

	if len(tagsResponse.Results) == 0 {
		return "", fmt.Errorf("no tags found")
	}

	var goldenTag string

	if tagsResponse.Results[0].Name == "latest" {
		goldenTag = tagsResponse.Results[1].Name
	} else {
		goldenTag = tagsResponse.Results[0].Name
	}
	fmt.Println("Tag for upgrade:", goldenTag)
	return goldenTag, nil

}
