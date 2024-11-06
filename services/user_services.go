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

func (s *UserService) GetCurrentUser(userID uint) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
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
