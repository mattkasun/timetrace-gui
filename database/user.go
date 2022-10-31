package database

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mattkasun/timetrace-gui/models"
)

func SaveUser(user *models.User) error {
	user.Updated = time.Now()
	value, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return insert(user.Username, string(value), USERS_TABLE_NAME)
}

func GetUser(name string) (models.User, error) {
	var user models.User
	records, err := fetch(USERS_TABLE_NAME)
	if err != nil {
		return user, err
	}
	for key, record := range records {
		if key == name {
			if err := json.Unmarshal([]byte(record), &user); err != nil {
				return user, err
			}
			return user, nil
		}
	}
	return user, errors.New("no such user")
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	var user models.User
	records, err := fetch(USERS_TABLE_NAME)
	if err != nil {
		return users, err
	}
	for _, record := range records {
		if err := json.Unmarshal([]byte(record), &user); err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func DeleteUser(name string) error {
	return delete(name, USERS_TABLE_NAME)
}
