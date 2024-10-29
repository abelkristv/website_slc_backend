package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CourseOutline struct {
	CourseOutlineId string `json:"CourseOutlineId"`
	Name            string `json:"Name"`
	Subjects        any    `json:"Subjects"` // If Subjects is expected to contain more specific types, update this type accordingly
}

func FetchCourseOutlines(token string) ([]CourseOutline, error) {
	url := "https://bluejack.binus.ac.id/lapi/api/Course/GetCourseOutlines"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add the Authorization header with the Bearer token
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch course outlines: %s", resp.Status)
	}

	var courseOutlines []CourseOutline
	if err := json.NewDecoder(resp.Body).Decode(&courseOutlines); err != nil {
		return nil, err
	}

	return courseOutlines, nil
}
