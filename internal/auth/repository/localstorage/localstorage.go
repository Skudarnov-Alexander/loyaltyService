package localstorage

import (
	"context"
	"errors"
	"fmt"

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

func (ls *Localstorage) CreateUser(ctx context.Context, u *model.User) error {
	if _, ok := ls.storage[u.Username]; ok {
		fmt.Println("login is exist")
		return errors.New("login is exist")
	}

	ls.storage[u.Username] = u.Password
	fmt.Printf("User is stored: %+v\n", ls.storage)
	return nil
}

func (ls *Localstorage) GetUser(ctx context.Context, username string) (*model.User, error) {
	pwd, ok := ls.storage[username]
	if !ok {
		return nil, errors.New("user is not exist")
	}

	return &model.User{
		Username: username,
		Password: pwd,
	}, nil
}
