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

func (s *UserService) CreateUser(username, password, role string) (*models.User, error) {
	if username == "" || password == "" || role == "" {
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
		Username: username,
		Password: hashedPassword,
		Role:     role,
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
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func (s *UserService) Login(username, password string) (string, error) {
	user, err := s.userRepo.LoginByUserInitial(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
