package api_repositories

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/abelkristv/slc_website/wiredsync/api/config"
	api_models "github.com/abelkristv/slc_website/wiredsync/api/models"
)

type AssistantRepository interface {
	FetchDataFromAPI() (api_models.AssistantDataResponse, error)
	FetchAssistantRoles(username string) ([]string, error)
}

type assistantRepository struct {
}

func NewAssistantRepository() AssistantRepository {
	return &assistantRepository{}
}

func (a *assistantRepository) FetchDataFromAPI() (api_models.AssistantDataResponse, error) {
	url := fmt.Sprintf("%s/Assistant/All", config.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return api_models.AssistantDataResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return api_models.AssistantDataResponse{}, err
	}

	var data api_models.AssistantDataResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return api_models.AssistantDataResponse{}, err
	}

	return data, nil
}

func (a *assistantRepository) FetchAssistantRoles(username string) ([]string, error) {
	url := fmt.Sprintf("%s/Assistant/GetAssistantRoles?username=%s", config.BaseURL, username)
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

func (a *assistantRepository) FetchAssistantEmail(binusianID string) (string, error) {
	url := fmt.Sprintf("%sAssistant/GetBinusianByBinusianId?binusianId=%s", config.BaseURL, binusianID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch email: %s", resp.Status)
	}

	var emailData api_models.Assistant
	if err := json.NewDecoder(resp.Body).Decode(&emailData); err != nil {
		return "", err
	}

	for _, email := range emailData.Emails {
		return email.Email, nil
	}

	return "", nil
}
