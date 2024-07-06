package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

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
	name := strings.Split(input.Email, "@")[0]
	_, err = m.DB.Exec("INSERT INTO users_settings(user_id, name, theme, avatar, phone_no) VALUES($1, $2, $3, $4, $5)",
		id, name, "LIGHT", "", "")
	if err != nil {
		log.Print("Insert user settings error. Error: " + err.Error())
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
	return m.GetUserById(ctx, &pb.GetUserByIdRequest{UserId: int32(user.Id)})
}

func (m *UserModel) GetUserById(ctx context.Context, input *pb.GetUserByIdRequest) (*pb.User, error) {
	user := &pb.User{}
	row := m.DB.QueryRow("SELECT id, email FROM users WHERE id=$1", input.UserId)
	err := row.Scan(&user.Id, &user.Email)
	if err != nil {
		return nil, errors.New("user id not found")
	}

	userSettings := &pb.UserSettings{}

	err = m.DB.QueryRow("SELECT user_id, name, avatar, theme, phone_no FROM users_settings WHERE user_id=$1", user.Id).
		Scan(&userSettings.Id, &userSettings.Name, &userSettings.Avatar, &userSettings.Theme, &userSettings.PhoneNo)

	if err == nil {
		user.Settings = userSettings
	} else {
		log.Print(err)
	}

	return user, nil
}

func (m *UserModel) UpdateUserSettings(ctx context.Context, in *pb.UserSettings) (*pb.UserSettings, error) {
	_, err := m.DB.Exec(`UPDATE users_settings 
				SET name = $1, avatar = $2, theme = $3, phone_no = $4 
				WHERE user_id = $5`, in.Name, in.Avatar, in.Theme, in.PhoneNo, in.Id,
	)
	if err != nil {
		return nil, err
	}
	return in, nil
}
