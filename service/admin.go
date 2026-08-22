package service

import (
	"Blog/repository/database"
)

func GetAdminUser(username, password string) error {
	return database.GetAdmin(username, password)
}
