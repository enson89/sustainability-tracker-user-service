package repository_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/enson89/sustainability-tracker-user-service/internal/model"
	"github.com/enson89/sustainability-tracker-user-service/internal/repository"
	"github.com/jmoiron/sqlx"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewUserRepository(sqlxDB)

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, now()) RETURNING id, created_at")).
		WithArgs(user.Username, user.Email, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	err = repo.CreateUser(user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if user.ID != 1 {
		t.Errorf("expected user ID to be 1, got %v", user.ID)
	}
}
