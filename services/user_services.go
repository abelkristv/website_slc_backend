package services

import (
	"errors"
	"time"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) GetCurrentUser(userID uint) (map[string]interface{}, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	userResponse := make(map[string]interface{})
	userResponse["Username"] = user.Username

	groupedHistory := make(map[string]interface{})
	groupedHistory["ID"] = user.Assistant.ID
	groupedHistory["Email"] = user.Assistant.Email
	groupedHistory["Bio"] = user.Assistant.Bio
	groupedHistory["FullName"] = user.Assistant.FullName
	groupedHistory["ProfilePicture"] = user.Assistant.ProfilePicture
	groupedHistory["Initial"] = user.Assistant.Initial
	groupedHistory["Generation"] = user.Assistant.Generation
	groupedHistory["Status"] = user.Assistant.Status

	type SocialMediaResponse struct {
		AssistantId         int
		GithubLink          string
		InstagramLink       string
		LinkedInLink        string
		WhatsappLink        string
		PersonalWebsiteLink string
	}

	socialMediaResponse := &SocialMediaResponse{
		GithubLink:          user.Assistant.AssistantSocialMedia.GithubLink,
		InstagramLink:       user.Assistant.AssistantSocialMedia.InstagramLink,
		LinkedInLink:        user.Assistant.AssistantSocialMedia.LinkedInLink,
		WhatsappLink:        user.Assistant.AssistantSocialMedia.WhatsappLink,
		PersonalWebsiteLink: user.Assistant.AssistantSocialMedia.PersonalWebsiteLink,
	}

	groupedHistory["SocialMedia"] = socialMediaResponse

	var teachingHistoryEntries []TeachingHistoryEntry
	var assistantAwardEntries []AssistantAwardEntry

	for _, award := range user.Assistant.AssistantAward {
		awardTitle := award.Award.AwardTitle
		AwardDescription := award.Award.AwardDescription

		found := false

		if !found {
			assistantAwardEntries = append(assistantAwardEntries, AssistantAwardEntry{
				AwardTitle:       awardTitle,
				AwardDescription: AwardDescription,
			})
		}
	}

	for _, history := range user.Assistant.TeachingHistory {
		periodTitle := history.Period.PeriodTitle
		courseData := map[string]interface{}{
			"CourseTitle":       history.Course.CourseTitle,
			"CourseCode":        history.Course.CourseCode,
			"CourseDescription": history.Course.CourseDescription,
		}

		found := false
		for i := range teachingHistoryEntries {
			if teachingHistoryEntries[i].PeriodTitle == periodTitle {
				teachingHistoryEntries[i].Courses = append(teachingHistoryEntries[i].Courses, courseData)
				found = true
				break
			}
		}
		if !found {
			teachingHistoryEntries = append(teachingHistoryEntries, TeachingHistoryEntry{
				PeriodTitle: periodTitle,
				Courses:     []map[string]interface{}{courseData},
			})
		}
	}

	sortedTeachingHistory := make([]map[string]interface{}, len(teachingHistoryEntries))
	for i, entry := range teachingHistoryEntries {
		sortedTeachingHistory[i] = map[string]interface{}{
			"PeriodTitle": entry.PeriodTitle,
			"Courses":     entry.Courses,
		}
	}

	groupedHistory["TeachingHistories"] = sortedTeachingHistory
	groupedHistory["Awards"] = assistantAwardEntries

	var positionEntries []map[string]interface{}
	for _, position := range user.Assistant.AssistantPosition {
		var startDate string
		if position.StartDate.Format("2006-01-02 15:04:05-07") != "0001-01-01 00:00:00+00" {
			startDate = position.StartDate.Format("2006-01-02 15:04:05-07")
		} else {
			startDate = ""
		}
		var endDate string
		if position.EndDate.Format("2006-01-02 15:04:05-07") != "0001-01-01 00:00:00+00" {
			endDate = position.EndDate.Format("2006-01-02 15:04:05-07")
		} else {
			endDate = ""
		}

		positionData := map[string]interface{}{
			"PositionName":        position.Position.Name,
			"PositionDescription": position.Description,
			"StartDate":           startDate,
			"EndDate":             endDate,
		}
		positionEntries = append(positionEntries, positionData)
	}
	groupedHistory["Positions"] = positionEntries

	userResponse["Assistant"] = groupedHistory

	return userResponse, nil
}

func (s *UserService) CreateUser(username, password, role string, assistantId int) (*models.User, error) {
	if username == "" || password == "" || role == "" || assistantId < 0 {
		return nil, errors.New("all fields are required")
	}

	// existingUser, err := s.userRepo.Get(email)
	// if err != nil {
	// 	return nil, err
	// }

	// if existingUser != nil {
	// 	return nil, errors.New("email already in use")
	// }

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		Username:    username,
		Password:    hashedPassword,
		Role:        role,
		AssistantId: assistantId,
	}

	err = s.userRepo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	existingUser, err := s.userRepo.GetUserByID(user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	existingUser.Username = user.Username
	existingUser.Password, err = HashPassword(user.Password)
	if err != nil {
		return err
	}
	existingUser.Role = user.Role

	return s.userRepo.UpdateUser(existingUser)
}

func (s *UserService) DeleteUser(id uint) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	return s.userRepo.DeleteUser(user)
}

var jwtKey = []byte("hehe")

type Claims struct {
	Id int `json:"id"`
	// Role     string `json:"role"`
	jwt.StandardClaims
}
type UserResponse struct {
	ID        uint             `json:"id"`
	Username  string           `json:"username"`
	Assistant models.Assistant `json:"assistant"`
}

func (s *UserService) Login(username, password string) (string, *UserResponse, error) {
	user, err := s.userRepo.LoginByUserInitial(username)
	if err != nil {
		return "", &UserResponse{}, err
	}

	userToReturn, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", &UserResponse{}, err
	}

	userResponse := &UserResponse{
		ID:        userToReturn.ID,
		Username:  userToReturn.Username,
		Assistant: userToReturn.Assistant,
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", &UserResponse{}, errors.New("invalid credentials")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Id: int(user.ID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", &UserResponse{}, err
	}

	return tokenString, userResponse, nil
}

func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	if err == nil {
		return errors.New("new password must not be same with old password")
	}

	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters long")
	}

	// if !isValidPassword(newPassword) {
	// 	return errors.New("new password must contain at least one number and one special character")
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
