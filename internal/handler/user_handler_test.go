package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/enson89/sustainability-tracker-user-service/internal/handler"
	"github.com/enson89/sustainability-tracker-user-service/internal/model"
	"github.com/gin-gonic/gin"
)

// stubAuthService implements the service.AuthService interface.
type stubAuthService struct{}

func (s *stubAuthService) Register(username, email, password string) (*model.User, error) {
	return &model.User{
		ID:       1,
		Username: username,
		Email:    email,
	}, nil
}

func (s *stubAuthService) Login(email, password string) (string, error) {
	return "dummy-token", nil
}

func (s *stubAuthService) UpdateProfile(id int64, username, email, password string) (*model.User, error) {
	return &model.User{
		ID:       id,
		Username: username,
		Email:    email,
	}, nil
}

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stubSvc := &stubAuthService{}
	userHandler := handler.NewUserHandler(stubSvc)
	router := gin.New()
	router.POST("/register", userHandler.Register)

	payload := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}
