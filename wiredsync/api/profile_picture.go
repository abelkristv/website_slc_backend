package api

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/abelkristv/slc_website/wiredsync/api/config"
)

func FetchProfilePicture(pictureId string) (string, error) {
	url := fmt.Sprintf("%s/Account/GetThumbnail?id=%s", config.BaseURL, pictureId)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch profile picture: %s", resp.Status)
	}

	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(imageData)
	return encoded, nil
}
