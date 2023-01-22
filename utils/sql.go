package utils

import (
	"database/sql"
	"psyWeb/configuration"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

/* 数据库初始化操作*/
type psyWebDataBase struct {
	Db *sql.DB
}

var psyWebDataBaseInstance *psyWebDataBase
var psyWebDataBaseOnce sync.Once

func GetPsyWebDataBaseInstance() *psyWebDataBase {
	psyWebDataBaseOnce.Do(func() {
		psyWebDataBaseInstance = &psyWebDataBase{}
	})
	return psyWebDataBaseInstance
}

func (me *psyWebDataBase) ConnectToSQL() (err error) {
	admin_name := &configuration.GetConfigInstance().DB.ID
	admin_password := &configuration.GetConfigInstance().DB.Password
	db_name := &configuration.GetConfigInstance().DB.Name
	// 打开数据库
	me.Db, err = sql.Open("mysql", (*admin_name)+":"+(*admin_password)+"@/"+(*db_name))
	if err == nil {
		err = me.Db.Ping() // 尝试与数据库建立连接
	}
	return err
}
