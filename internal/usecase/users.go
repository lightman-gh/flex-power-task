package usecase

import (
	"context"
	"flex/internal/dbschema/model"
	"flex/internal/dbschema/model/user"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const BcryptCost = 14

type UsersService struct {
	client *model.Client
}

func NewUsersService(client *model.Client) *UsersService {
	return &UsersService{
		client: client,
	}
}

func (u *UsersService) Create(ctx context.Context, usr *model.User) error {
	// we need to hash password to save it in the database
	hash, err := hashPassword(usr.Password)
	if err != nil {
		return fmt.Errorf("can not hash password, err: %v", err)
	}

	// add database transaction
	_, err = u.client.User.Create().
		SetPassword(hash).
		SetUsername(usr.Username).
		Save(ctx)

	return err
}

func (u *UsersService) GetByPassword(ctx context.Context, usr *model.User) (*model.User, error) {
	dbUser, err := u.client.User.Query().
		Where(user.Username(usr.Username)).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("can not select user, err: %v", err)
	}

	if err := verifyPassword(usr.Password, dbUser.Password); err != nil {
		return nil, fmt.Errorf("incorrect passwordr, err: %v", err)
	}

	return dbUser, nil
}

func hashPassword(clearTextPw string) (string, error) {
	hashedPwBytes, err := bcrypt.GenerateFromPassword([]byte(clearTextPw), BcryptCost)

	return string(hashedPwBytes), err
}

func verifyPassword(clearTextPw, hashedPw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(clearTextPw))
}
