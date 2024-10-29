package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EmailResponse struct {
	Emails []struct {
		Email  string `json:"Email"`
		Prefer string `json:"Prefer"`
		Type   string `json:"Type"`
	} `json:"Emails"`
}

func FetchEmail(binusianID string) (string, error) {
	url := fmt.Sprintf("https://bluejack.binus.ac.id/lapi/api/Assistant/GetBinusianByBinusianId?binusianId=%s", binusianID)
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
		// if strings.HasSuffix(email.Email, ".ac.id") || strings.HasSuffix(email.Email, ".edu") {
		return email.Email, nil
		// }
	}

	// return "", fmt.Errorf("no valid email ending with .ac.id found")
	return "", nil

}
