package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Period struct {
	Description string `json:"Description"`
	End         string `json:"End"`
	SemesterID  string `json:"SemesterID"`
	Start       string `json:"Start"`
}

type PeriodDataResponse []Period

func FetchPeriods() ([]Period, error) {
	url := "https://bluejack.binus.ac.id/lapi/api/Semester/GetSemestersWithActiveDate"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch periods: %s", resp.Status)
	}

	var periods []Period
	if err := json.NewDecoder(resp.Body).Decode(&periods); err != nil {
		return nil, err
	}

	return periods, nil
}
