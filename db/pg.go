package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresSql(dataSourceName string) DB {
	database, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}
	return &Postgres{
		DB: database,
	}
}

type Postgres struct {
	DB *sqlx.DB
}

func (pg *Postgres) Create(queryList []string) error {
	defer pg.DB.Close()
	tx, err := pg.DB.Beginx()
	if err != nil {
		return err
	}
	for _, query := range queryList {
		_, err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
	return nil
}

func (pg *Postgres) Insert(queryList []string) error {
	defer pg.DB.Close()
	tx, err := pg.DB.Beginx()
	if err != nil {
		return err
	}
	for _, query := range queryList {
		_, err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
	return nil
}

func (pg *Postgres) Update(queryList []string) error {
	defer pg.DB.Close()
	tx, err := pg.DB.Beginx()
	if err != nil {
		return err
	}
	for _, query := range queryList {
		_, err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
	return nil
}

func (pg *Postgres) Delete(queryList []string) error {
	defer pg.DB.Close()
	tx, err := pg.DB.Beginx()
	if err != nil {
		return err
	}
	for _, query := range queryList {
		_, err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
	return nil
}

//TODO select with interface{}
//func (pg *Postgres) Select(object interface{}, query string) (interface{}, error) {
//	defer pg.DB.Close()
//
//	err := pg.DB.Get(&object, query)
//	if err != nil {
//		return nil, err
//	}
//
//	return object, nil
//}
