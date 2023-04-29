package interactor

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type userRepository interface {
	Create(ctx context.Context, u model.User) (string, error)
	//GetByName(ctx context.Context, username string) (model.User, error)
}

type balanceRepository interface {
	Сreate(ctx context.Context, userID string) error
}

type passHasher interface {
	Hash(password string) string
	IsPwdsMatched(savedHashedPwd, pwd string, salt []byte) bool
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

func (ai *authInteractor) SignUp(ctx context.Context, username, pwd string) error {
	hashedPass := ai.passHasher.Hash(pwd)

	u := model.User{
		Username:   username,
		HashedPass: hashedPass,
	}

	id, err := ai.userRepository.Create(ctx, u)
	if err != nil {
		return err
	}

	if err := ai.balanceRepository.Сreate(ctx, id); err != nil {
		return err
	}

	return nil
}

/*
func New(userRepo auth.AuthRepository) (*AuthService, error) {
	hashSalt, err := generateRandomSalt(saltSize)
	if err != nil {
		return nil, err
	}

	return &AuthService{
		userRepo: userRepo,
		hashSalt: hashSalt,
	}, nil
}
*/

/*
func (s *AuthService) SignUp(ctx context.Context, u model.User) error {
	u.Password = hashPassword(u.Password, s.hashSalt)
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	u.ID = uuid.String()

	return s.userRepo.CreateUser(ctx, u)
}

func (s *AuthService) SignIn(ctx context.Context, u model.User) (string, error) {
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
