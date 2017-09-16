package db

import (
	"database/sql"
	"fmt"
	"wz1025/utils"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

//常量定义
var (
	db *sql.DB
)

//初始化函数
func init() {
	//初始化数据库连接
	db_url := beego.AppConfig.String("db_str")
	var err error
	db, err = sql.Open("mysql", db_url)
	if err != nil {
		utils.ErrorLog("Please check your net or database connection info.", err)
		return
	}
	err = db.Ping()
	if err != nil {
		utils.ErrorLog("Please check your net or database connection info.", err)
		return
	}

	fmt.Println("db start over.")
}
