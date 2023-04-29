package dto

import (
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/google/uuid"
)

type User struct {
	UserID   uuid.UUID `db:"user_id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
}

func UserToModel(u User) (model.User, error) {
	uuid, err := u.UserID.Value()
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:       uuid.(string),
		Username: u.Username,
		//Password: u.Password,
	}, nil
}