package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abelkristv/slc_website/wiredsync/api/config"
)

type EmailResponse struct {
	Emails []struct {
		Email  string `json:"Email"`
		Prefer string `json:"Prefer"`
		Type   string `json:"Type"`
	} `json:"Emails"`
}

func FetchEmail(binusianID string) (string, error) {
	url := fmt.Sprintf("%sAssistant/GetBinusianByBinusianId?binusianId=%s", config.BaseURL, binusianID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch email: %s", resp.Status)
	}

	var emailData EmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&emailData); err != nil {
		return "", err
	}

	for _, email := range emailData.Emails {
		return email.Email, nil
	}

	return "", nil

}
