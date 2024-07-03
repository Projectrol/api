package models

import (
	"context"
	"database/sql"
	"errors"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             int            `json:"id"`
	Email          string         `json:"email"`
	HashedPassword string         `json:"-"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      sql.NullString `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(DB *sql.DB) *UserModel {
	return &UserModel{DB}
}

func (m *UserModel) Insert(ctx context.Context, input *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	row := m.DB.QueryRow("SELECT email FROM users WHERE email=$1", input.Email)
	var existedEmail string
	err := row.Scan(&existedEmail)
	if err == nil {
		return &pb.CreateUserResponse{Id: -1}, errors.New("duplicated email")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return &pb.CreateUserResponse{Id: -1}, errors.New("somthing error")
	}
	_, err = m.DB.Exec("INSERT INTO users(email, hashed_password) VALUES($1, $2)", input.Email, hashedPassword)
	if err != nil {
		return &pb.CreateUserResponse{Id: -1}, errors.New("somthing error")
	}
	var id int32
	row = m.DB.QueryRow("SELECT id from users WHERE email=$1", input.Email)
	err = row.Scan(&id)
	if err != nil {
		return &pb.CreateUserResponse{Id: -1}, errors.New("somthing error")
	}
	return &pb.CreateUserResponse{Id: id}, nil
}

func (m *UserModel) Login(ctx context.Context, input *pb.LoginRequest) (*pb.User, error) {
	user := &User{}
	row := m.DB.QueryRow("SELECT id, email, hashed_password FROM users WHERE email=$1", input.Email)
	err := row.Scan(&user.Id, &user.Email, &user.HashedPassword)
	if err != nil {
		return nil, errors.New("email not existed")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password))
	if err != nil {
		return nil, errors.New("password incorrect")
	}
	return &pb.User{
		Id:    int32(user.Id),
		Email: user.Email,
	}, nil
}

func (m *UserModel) GetUserById(ctx context.Context, input *pb.GetUserByIdRequest) (*pb.User, error) {
	user := &User{}
	row := m.DB.QueryRow("SELECT id, email FROM users WHERE id=$1", input.UserId)
	err := row.Scan(&user.Id, &user.Email)
	if err != nil {
		return nil, errors.New("user id not found")
	}
	return &pb.User{
		Id:    int32(user.Id),
		Email: user.Email,
	}, nil
}
