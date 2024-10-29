package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	BinusianID string `json:"BinusianId"`
	Major      string `json:"major"`
	Name       string `json:"name"`
	PictureID  string `json:"pictureid"`
	UserID     string `json:"userID"`
	Username   string `json:"username"`
}

type UserDataResponse struct {
	Active   []User `json:"active"`
	Inactive []User `json:"inactive"`
}

func FetchDataFromAPI() (UserDataResponse, error) {
	resp, err := http.Get("https://bluejack.binus.ac.id/lapi/api/Assistant/All")
	if err != nil {
		return UserDataResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return UserDataResponse{}, err
	}

	var data UserDataResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return UserDataResponse{}, err
	}

	return data, nil
}

func FetchAssistantRoles(username string) ([]string, error) {
	url := fmt.Sprintf("https://bluejack.binus.ac.id/lapi/api/Assistant/GetAssistantRoles?username=%s", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch roles: %s", resp.Status)
	}

	var roles []string
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}

	return roles, nil
}
