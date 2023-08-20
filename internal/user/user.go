package user

import (
	"errors"
	"gcipher/internal/db/models"
	"gcipher/internal/db/repositories"
	"gcipher/internal/util"
)

func Authenticate(username, password string) (*models.User, error) {
	user, err := repositories.GetUserRepository().FindByUsername(username)
	if err != nil {
		return nil, err
	}

	// Hash the provided password and compare it with the stored hash
	ok, err := util.ComparePasswordAndHash(user.Password, password)
	if !ok || err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
