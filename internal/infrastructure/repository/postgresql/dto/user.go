package dto

import (
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID       pgtype.UUID `db:"id"`
	Username string      `db:"username"`
	Password string      `db:"password"`
}

func (u *User) ToModel() (model.User, error) {
	uuid, err := u.ID.Value()
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:         uuid.(string),
		Username:   u.Username,
		HashedPass: u.Password,
	}, nil
}
