package db

import (
	"fmt"
)

type DB interface {
	Create(queryList []string) error
	Insert(queryList []string) error
	Update(queryList []string) error
	Delete(queryList []string) error
	//Select(object interface{}, query string) (interface{}, error)
}

func NewDB(config Config) DB {
	switch config.Type {
	case "mysql":
		return NewMysql(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.UserName, config.Password, config.Host, config.Port, config.DataBase))
	case "postgres":
		return NewPostgresSql(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", config.Host, config.Port, config.UserName, config.DataBase, config.Password))
	}
	return nil
}

type Config struct {
	Type     string
	UserName string
	Password string
	Host     string
	Port     string
	DataBase string
}
