package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type TeachingSchedule struct {
	ClassName    string `json:"ClassName"`
	CourseCode   string `json:"CourseCode"`
	CourseTitle  string `json:"CourseTitle"`
	DeliveryMode string `json:"DeliveryMode"`
	StartAt      string `json:"StartAt"`
	EndAt        string `json:"EndAt"`
}

func FetchTeachingHistory(binusianId, semesterId, token, assistantName, periodName string) ([]TeachingSchedule, error) {
	apiURL := fmt.Sprintf(
		"https://bluejack.binus.ac.id/lapi/api/Lecturer/GetLecturerTeachingSchedules?binusianId=%s&semesterId=%s&startDate=&endDate=",
		url.QueryEscape(binusianId), url.QueryEscape(semesterId),
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
