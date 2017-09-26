package db

import (
	"database/sql"
	"wz1025/utils"
	"fmt"
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
	url := beego.AppConfig.String("db_str")
	var err error
	db, err = sql.Open("mysql", url)
	if err != nil {
		utils.ErrorLog("Please check your net or database connection info.", err)
		return
	}
	err = db.Ping()
	if err != nil {
		utils.ErrorLog("Please check your net or database connection info.", err)
		return
	}
	fmt.Println("http init")
}

//得到数据源
func GetDb() *sql.DB {
	return db
}
