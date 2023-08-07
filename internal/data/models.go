package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflit")
)

type Model struct {
	Permissions PermissionModel
	Token       TokenModel
	Track       TrackModel
	User        UserModel
}

func New(db *pgxpool.Pool) Model {
	return Model{
		Permissions: PermissionModel{db},
		Token:       TokenModel{db},
		Track:       TrackModel{db},
		User:        UserModel{db},
	}
}
