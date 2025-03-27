package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/enson89/sustainability-tracker-user-service/internal/model"
	"github.com/enson89/sustainability-tracker-user-service/internal/service"
)

// mockUserRepo implements the UserRepository interface for testing.
type mockUserRepo struct {
	user *model.User
	err  error
}

func (m *mockUserRepo) CreateUser(user *model.User) error {
	if m.err != nil {
		return m.err
	}
	user.ID = 1
	user.CreatedAt = time.Now()
	m.user = user
	return nil
}

func (m *mockUserRepo) GetUserByEmail(email string) (*model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.user != nil && m.user.Email == email {
		return m.user, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) UpdateUser(user *model.User) error {
	if m.err != nil {
		return m.err
	}
	m.user = user
	return nil
}

func TestRegister_Success(t *testing.T) {
	repo := &mockUserRepo{}
	authSvc := service.NewAuthService(repo, "testsecret")

	user, err := authSvc.Register("testuser", "test@example.com", "password")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if user.ID == 0 {
		t.Errorf("expected user ID to be set, got 0")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	repo := &mockUserRepo{err: errors.New("user not found")}
	authSvc := service.NewAuthService(repo, "testsecret")

	_, err := authSvc.Login("nonexistent@example.com", "password")
	if err == nil {
		t.Fatal("expected error for invalid credentials, got nil")
	}
}
