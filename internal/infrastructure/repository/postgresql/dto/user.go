package dto

import (
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
}

func (u *User) ToModel() (*model.User, error) {
	uuid, err := u.ID.Value()
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:         uuid.(string),
		Username:   u.Username,
		HashedPass: u.Password,
	}, nil
}
