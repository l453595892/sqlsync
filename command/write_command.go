package command

import (
	"github.com/l453595892/raft"
	"github.com/l453595892/sqlsync/db"
	"strings"
)

type WriteCommand struct {
	SqlList []string `json:"sql"`
}

func NewWriteCommand(sqlList []string) *WriteCommand {
	return &WriteCommand{
		SqlList: sqlList,
	}
}

func (c *WriteCommand) CommandName() string {
	return "Write"
}

func (c *WriteCommand) Apply(server raft.Server) (interface{}, error) {
	config := server.Context().(db.Config)
	createList := make([]string, 0)
	insertList := make([]string, 0)
	updateList := make([]string, 0)
	deleteList := make([]string, 0)

	for _, sql := range c.SqlList {
		if strings.Contains(sql, "create") || strings.Contains(sql, "CREATE") {
			createList = append(createList, sql)
		}
		if strings.Contains(sql, "insert") || strings.Contains(sql, "INSERT") {
			insertList = append(insertList, sql)
		}
		if strings.Contains(sql, "update") || strings.Contains(sql, "UPDATE") {
			updateList = append(updateList, sql)
		}
		if strings.Contains(sql, "delete") || strings.Contains(sql, "DELETE") {
			deleteList = append(deleteList, sql)
		}
	}

	driver := db.NewDB(config)
	if len(createList) != 0 {
		err := driver.Create(createList)
		if err != nil {
			return nil, err
		}
	}
	if len(insertList) != 0 {
		err := driver.Insert(insertList)
		if err != nil {
			return nil, err
		}
	}
	if len(deleteList) != 0 {
		err := driver.Delete(deleteList)
		if err != nil {
			return nil, err
		}
	}
	if len(updateList) != 0 {
		err := driver.Update(updateList)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
