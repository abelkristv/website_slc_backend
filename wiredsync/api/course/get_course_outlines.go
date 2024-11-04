package course

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abelkristv/slc_website/wiredsync/config"
)

type GetCourseOutlineResponse struct {
	CourseOutlineId string `json:"CourseOutlineId"`
	Name            string `json:"Name"`
	Subjects        any    `json:"Subjects"`
}

func FetchCourseOutlines(token string) ([]GetCourseOutlineResponse, error) {
	url := fmt.Sprintf("%s/Course/GetCourseOutlines", config.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

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

	var courseOutlines []GetCourseOutlineResponse
	if err := json.NewDecoder(resp.Body).Decode(&courseOutlines); err != nil {
		return nil, err
	}

	return courseOutlines, nil
}
