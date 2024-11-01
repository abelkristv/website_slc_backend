package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/abelkristv/slc_website/wiredsync/api/config"
)

type TeachingSchedule struct {
	Assistant       string  `json:"Assistant"`
	Campus          string  `json:"Campus"`
	Class           string  `json:"Class"`
	CourseOutlineId string  `json:"CourseOutlineId"`
	Date            *string `json:"Date"`
	Day             int     `json:"Day"`
	Id              string  `json:"Id"`
	LecturerCode    *string `json:"LecturerCode"`
	LecturerName    *string `json:"LecturerName"`
	Note            *string `json:"Note"`
	Number          int     `json:"Number"`
	Realization     string  `json:"Realization"`
	Room            string  `json:"Room"`
	SemesterId      string  `json:"SemesterId"`
	Session         *string `json:"Session"`
	Shift           string  `json:"Shift"`
	Subject         string  `json:"Subject"`
	SubjectId       string  `json:"SubjectId"`
	TheoryClass     *string `json:"TheoryClass"`
	TotalStudent    int     `json:"TotalStudent"`
}

func FetchTeachingHistory(username, semesterId, token, assistantName, periodName string) ([]TeachingSchedule, error) {
	apiURL := fmt.Sprintf(
		"%sAssistant/GetClassTransactionByAssistantUsername?username=%s&semesterId=%s&startDate=&endDate=",
		config.BaseURL, url.QueryEscape(username), url.QueryEscape(semesterId),
	)

	log.Printf("Fetching teaching history data for %s- %s", assistantName, periodName)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch teaching history, status code: %d", resp.StatusCode)
	}

	var schedules []TeachingSchedule
	if err := json.NewDecoder(resp.Body).Decode(&schedules); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Print(schedules)

	return schedules, nil
}
