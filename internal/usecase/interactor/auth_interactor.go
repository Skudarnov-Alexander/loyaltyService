package interactor

import (
	"context"
	"errors"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/pkg/jwt"
)

var SampleSecretKey = []byte("SecretYouShouldHide") //TODO спрятать в конфиги / секретницу

type userRepository interface {
	Create(ctx context.Context, u model.User) (string, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
}

type balanceRepository interface {
	Сreate(ctx context.Context, userID string) error
}

type passHasher interface {
	Hash(password string) string
	IsPwdsMatched(savedHashedPwd, pwd string) bool
}

type authInteractor struct {
	userRepository
	balanceRepository
	passHasher
}

func NewAuthInteractor(userRepository userRepository, balanceRepository balanceRepository, passHasher passHasher) *authInteractor {
	return &authInteractor{
		userRepository:    userRepository,
		balanceRepository: balanceRepository,
		passHasher:        passHasher,
	}
}

func (ai *authInteractor) SignUp(ctx context.Context, username, pwd string) (string, error) {
	hashedPass := ai.passHasher.Hash(pwd)

	u := model.User{
		Username:   username,
		HashedPass: hashedPass,
	}

	id, err := ai.userRepository.Create(ctx, u)
	if err != nil {
		return "", err
	}

	if err := ai.balanceRepository.Сreate(ctx, id); err != nil {
		return "", err
	}

	return jwt.GenerateJWT(SampleSecretKey, u)

}

func (ai *authInteractor) LogIn(ctx context.Context, username, pwd string) (string, error) {
	u, err := ai.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if ok := ai.passHasher.IsPwdsMatched(u.HashedPass, pwd); !ok {
		return "", errors.New("pass is not match")
	}

	return jwt.GenerateJWT(SampleSecretKey, u)

}

/*

func (ai *authInteractor) SignIn(ctx context.Context, u model.User) (string, error) {
	user, err := s.userRepo.GetUser(ctx, u.Username)
	if err != nil {
		return "", err
	}

	if ok := doPasswordsMatch(user.Password, u.Password, s.hashSalt); !ok {
		return "", errors.New("pass is not match")
	}

	return generateJWT(user)

}
*/
