package api_models

type Assistant struct {
	BinusianID string `json:"BinusianId"`
	Major      string `json:"major"`
	Name       string `json:"name"`
	PictureID  string `json:"pictureid"`
	UserID     string `json:"userID"`
	Username   string `json:"username"`
	Emails     []struct {
		Email string `json:"Email"`
	} `json:"Emails"`
}

type AssistantDataResponse struct {
	Active   []Assistant `json:"active"`
	Inactive []Assistant `json:"inactive"`
}
