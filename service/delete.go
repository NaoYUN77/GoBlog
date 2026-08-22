package service

import "Blog/repository/database"

func DeletePost(id int64) error {
	return database.DeletePost(id)
}
