package database

import (
	"errors"
	"fmt"
	"time"

	"Blog/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 定义一个初始化mysql 使用sqlx
var DB *sqlx.DB

func Init() error {

	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",

		settings.Cnf.Mysql.User,
		settings.Cnf.Mysql.Password,
		settings.Cnf.Mysql.Host,
		settings.Cnf.Mysql.Port,
		settings.Cnf.Mysql.DBName,
		settings.Cnf.Mysql.Charset,
	)
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return errors.New("初始化mysql失败")
	}

	//空闲时保存链接10分钟
	DB.SetConnMaxIdleTime(10 * time.Minute)

	//链接创建后保持多久时间
	// DB.SetConnMaxLifetime(10 )

	//最多同时打开多少条链接
	DB.SetMaxOpenConns(10)

	//链接池最多保存多少空闲链接
	DB.SetMaxIdleConns(10)

	return nil
}

func Close() {
	DB.Close()
}
