package http

import (
	"wz1025/db"
	"wz1025/utils"
	"strings"
	"database/sql"
)

type AdminDbFun struct {
}

var adminDbFun = &AdminDbFun{}

//添加视频解析地址
func (self *AdminDbFun) AddVideoExplainUrl(args map[string]interface{}) bool {
	_, err := db.GetDb().Exec("INSERT INTO ZJ_EXPLAIN_URL() VALUES(URL_ADDR,CREATE_TIME,ACTIVE,TYPE) VALUES(?, CURRENT_TIME, 1,?)", args["url_addr"], args["type"])
	if err != nil {
		utils.ErrorLog("admin_db_fun.go AddVideoExplainUrl method.", err)
		return false
	}
	return true
}

//根据查询条件查询视频解释地址
func (self *AdminDbFun) FindVideoExplain(args, result map[string]interface{}) {
	after_strs := make([]string, 0)
	params := make([]interface{}, 0)

	if type_val, type_has := args["type"]; type_has {
		after_strs = append(after_strs, " AND TYPE=?")
		params = append(params, type_val)
	}
	if active_val, active_has := args["active"]; active_has {
		after_strs = append(after_strs, " AND ACTIVE=?")
		params = append(params, active_val)
	}

	//排序
	after_strs = append(after_strs, " ORDER BY CREATE_TIME DESC")

	//查询总条数
	count_sql := make([]string, 0, len(after_strs)+1)
	count_sql = append(count_sql, "SELECT COUNT(ID) FROM ZJ_EXPLAIN_URL WHERE 1=1")
	count_sql = append(count_sql, after_strs...)
	row := db.GetDb().QueryRow(strings.Join(count_sql, ""), params...)
	var rows_count uint64
	err := row.Scan(&rows_count)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.ErrorLog("admin_db_fun.go FindVideoExplain method.", err)
		}
		return
	}
	result["count"] = rows_count

	//判断分页参数
	if startLine_val, startLine_has := args["startLine"]; startLine_has {
		after_strs = append(after_strs, " LIMIT ?")
		params = append(params, startLine_val)
	}
	if endLine_val, endLine_has := args["endLine"]; endLine_has {
		after_strs = append(after_strs, " ,?")
		params = append(params, endLine_val)
	}

	//查询信息
	sql_strs := make([]string, 0, len(after_strs)+1)
	sql_strs = append(sql_strs, "SELECT ID, URL_ADDR,CREATE_TIME,ACTIVE,TYPE FROM ZJ_EXPLAIN_URL WHERE 1=1")
	sql_strs = append(sql_strs, after_strs...)
	rows, err := db.GetDb().Query(strings.Join(sql_strs, ""), params...)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.ErrorLog("admin_db_fun.go FindVideoExplain method.", err)
		}
		return
	}
	defer rows.Close()
	result_rows := make([]map[string]interface{}, 0)
	for rows.Next() {
		var (
			id                    uint64
			active, row_type      uint8
			url_addr, create_time string
		)
		rows.Scan(&id, &url_addr, &create_time, &active, &row_type)
		result_row := map[string]interface{}{
			"id":          id,
			"url":         url_addr,
			"create_time": utils.FormatDataTime(create_time[:10]),
			"active":      active,
			"type_str":    self.type2str(row_type),
			"type":        row_type,
		}
		result_rows = append(result_rows, result_row)
	}

	result["data"] = result_rows
}

//类型转换为字符
func (self *AdminDbFun) type2str(row_type uint8) string {
	str := ""
	switch row_type {
	case 1:
		str = "爱奇艺"
	case 2:
		str = "优酷"
	case 3:
		str = "腾讯"
	case 4:
		str = "乐视"
	case 5:
		str = "搜狐"
	case 6:
		str = "土豆"
	case 7:
		str = "芒果TV"
	case 8:
		str = "PPTV"
	default:
		str = "万能"
	}
	return str
}

//得到操作关于member数据的结构信息
func NewAdminDbFun() *AdminDbFun {
	return adminDbFun
}
