package repository

import (
	"errors"

	"github.com/enson89/sustainability-tracker-user-service/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user *model.User) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	query := `INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, now()) RETURNING id, created_at`
	return r.db.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password, created_at FROM users WHERE email=$1`
	err := r.db.Get(user, query, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	query := `UPDATE users SET username=$1, email=$2, password=$3 WHERE id=$4`
	res, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user not found")
	}
	return nil
}

func NewPostgresDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
