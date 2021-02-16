package postgres

import "github.com/orbis-challenge/src/models"

// GetUserByEmail gets  user model by email if exists.
func (q DBQuery) GetUserByEmail(email string) (user models.User, err error) {
	err = q.Model(&user).
		Where(`"user".email = ?`, email).
		First()

	return user, err
}

// CreateUser creates a new user
func (q DBQuery) CreateUser(user *models.User) (models.User, error) {
	_, err := q.Model(user).
		Returning("*").
		Insert()

	if err != nil {
		return models.User{}, err
	}

	return *user, nil
}
