package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgres_Create(t *testing.T) {
	db := NewPostgresSql()

	company := `CREATE TABLE COMPANY(
   ID INT PRIMARY KEY     NOT NULL,
   NAME           TEXT    NOT NULL,
   AGE            INT     NOT NULL,
   ADDRESS        CHAR(50),
   SALARY         REAL
);`

	err := db.Create([]string{company})
	assert := assert.New(t)
	assert.Nil(err)
}

func TestPostgres_Insert(t *testing.T) {
	db := NewPostgresSql()

	insert := `insert into company (id,name,age) VALUES (1,'李伟杰',27)`

	err := db.Insert([]string{insert})
	assert := assert.New(t)
	assert.Nil(err)
}

func TestPostgres_Update(t *testing.T) {
	db := NewPostgresSql()

	update := `update company set age=18 where id=1`

	err := db.Update([]string{update})
	assert := assert.New(t)
	assert.Nil(err)
}

func TestPostgres_Delete(t *testing.T) {
	db := NewPostgresSql()

	delete := `delete from company where id=1`

	err := db.Delete([]string{delete})
	assert := assert.New(t)
	assert.Nil(err)
}

//func TestPostgres_Select(t *testing.T) {
//	db := NewPostgresSql()
//
//	company := Company{}
//
//	selectQuery := `select * from company`
//
//	b, _ := json.Marshal(company)
//	fmt.Println(db.Select(b, selectQuery))
//}
//
//type Company struct {
//	Id      int    `json:"id"`
//	Name    string `json:"name"`
//	Age     int    `json:"age"`
//	Address string `json:"address"`
//	Salary  string `json:"salary"`
//}
