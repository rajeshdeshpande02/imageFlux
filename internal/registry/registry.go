package registry

import (
	"encoding/json"
	"fmt"
	"log"
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
		return "", fmt.Errorf("error getting image tags")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		return "", fmt.Errorf("error: received non-200 response code")
	}
	var tagsResponse TagsResponse

	if err := json.NewDecoder(resp.Body).Decode(&tagsResponse); err != nil {
		return "", fmt.Errorf("error decoding respo")
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
	log.Printf("Tag for upgrade: %v", goldenTag)
	return goldenTag, nil

}
