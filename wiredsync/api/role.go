package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abelkristv/slc_website/wiredsync/config"
)

type Role string

type RoleDataResponse []Role

func FetchAssistantRoles(username string) (RoleDataResponse, error) {
	url := fmt.Sprintf("%s/Assistant/GetAssistantRoles?username=%s", config.BaseURL, username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch assistant roles: %s", resp.Status)
	}

	var roles RoleDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}

	return roles, nil
}
