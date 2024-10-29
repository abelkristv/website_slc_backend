package api

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchProfilePicture fetches the profile picture from the given PictureId and returns its base64 representation
func FetchProfilePicture(pictureId string) (string, error) {
	url := fmt.Sprintf("https://bluejack.binus.ac.id/lapi/api/Account/GetThumbnail?id=%s", pictureId)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch profile picture: %s", resp.Status)
	}

	// Read the image data
	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert to base64
	encoded := base64.StdEncoding.EncodeToString(imageData)
	return encoded, nil
}
