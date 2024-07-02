package models

import "database/sql"

type Models struct {
	UserModel *UserModel
}

func NewModels(DB *sql.DB) *Models {
	return &Models{
		UserModel: &UserModel{
			DB: DB,
		},
	}
}
