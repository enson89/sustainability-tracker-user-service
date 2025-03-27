package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/enson89/sustainability-tracker-user-service/internal/model"
	"github.com/enson89/sustainability-tracker-user-service/internal/repository"
)

type AuthService interface {
	Register(username, email, password string) (*model.User, error)
	Login(email, password string) (string, error)
	UpdateProfile(id int64, username, email, password string) (*model.User, error)
}

type authService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(username, email, password string) (*model.User, error) {
	// Check if user already exists
	_, err := s.repo.GetUserByEmail(email)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	token, err := s.generateJWT(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) UpdateProfile(id int64, username, email, password string) (*model.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	user.Username = username
	user.Email = email
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) generateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
