package database

import (
	"database/sql"
	"errors"
)

func GetAdmin(username, password string) error {
	var dbPassword string
	sqlstr := "select password from tb_user where username = ?"
	err := DB.QueryRow(sqlstr, username).Scan(&dbPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("用户不存在")
		}
		return errors.New("查询用户失败")
	}
	if dbPassword != password {
		return errors.New("密码错误")
	}
	return nil
}
