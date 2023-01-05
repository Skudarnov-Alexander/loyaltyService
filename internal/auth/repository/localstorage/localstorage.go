package localstorage

import (
	"context"
	"errors"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type Localstorage struct {
	storage map[string]string
}

func New() *Localstorage {
	return &Localstorage{
		storage: make(map[string]string),
	}
}

func (ls *Localstorage) CreateUser(ctx context.Context, User *model.User) error {
	if _, ok := ls.storage[User.Username]; !ok {
		return errors.New("login is exist")
	}

	ls.storage[User.Username] = User.Password
	return nil
}

func (ls *Localstorage) GetUser(ctx context.Context, username, password string) (*model.User, error) {
	return nil, nil
}
