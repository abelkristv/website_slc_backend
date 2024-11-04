package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abelkristv/slc_website/wiredsync/api/course"
	"github.com/abelkristv/slc_website/wiredsync/config"
)

type CourseDescription struct {
	CourseDescription string `json:"CourseDescription"`
}

func FetchCourseOutlines(token string) ([]course.GetCourseOutlineResponse, error) {
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

	var courseOutlines []course.GetCourseOutlineResponse
	if err := json.NewDecoder(resp.Body).Decode(&courseOutlines); err != nil {
		return nil, err
	}

	return courseOutlines, nil
}

func FetchCourseDescription(courseId, token string) (CourseDescription, error) {
	url := fmt.Sprintf("%s/Course/GetCourseOutlineDetail?courseOutlineId=%s", config.BaseURL, courseId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CourseDescription{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CourseDescription{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CourseDescription{}, fmt.Errorf("failed to fetch course outlines: %s", resp.Status)
	}

	var courseDescription CourseDescription
	if err := json.NewDecoder(resp.Body).Decode(&courseDescription); err != nil {
		return CourseDescription{}, err
	}

	return courseDescription, nil
}
