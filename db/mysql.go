package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMysql(dataSourceName string) DB {
	database, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	return &Mysql{
		DB: database,
	}
}

type Mysql struct {
	DB *sqlx.DB
}

func (mysql *Mysql) Create(queryList []string) error {
	defer mysql.DB.Close()
	for _, query := range queryList {
		_, err := mysql.DB.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mysql *Mysql) Insert(queryList []string) error {
	defer mysql.DB.Close()
	for _, query := range queryList {
		_, err := mysql.DB.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mysql *Mysql) Update(queryList []string) error {
	defer mysql.DB.Close()
	for _, query := range queryList {
		_, err := mysql.DB.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mysql *Mysql) Delete(queryList []string) error {
	defer mysql.DB.Close()
	for _, query := range queryList {
		_, err := mysql.DB.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mysql *Mysql) Select(object interface{}, query string) (interface{}, error) {
	defer mysql.DB.Close()
	err := mysql.DB.Select(&object, query)
	if err != nil {
		return nil, err
	}
	return object, nil
}
