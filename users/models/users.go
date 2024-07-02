package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email          string         `json:"email"`
	HashedPassword string         `json:"-"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      sql.NullString `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(email string, password string) error {
	var existedEmail string
	row := m.DB.QueryRow("SELECT email FROM users WHERE email=$1", email)
	err := row.Scan(&existedEmail)
	if err == nil {
		return errors.New("duplicated email")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec("INSERT INTO users(email, hashed_password) VALUES($1, $2)", email, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}
