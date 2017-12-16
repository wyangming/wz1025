package http

import (
	"database/sql"
	"fmt"
	"strings"
	"wz1025/db"
	"wz1025/utils"
	"bytes"
)

type AdminDbFun uint8

var adminDbFun = AdminDbFun(1)

//更新解析地址状态
//参数 action_type 1 恢复，0作废，默认0
func (self *AdminDbFun) UpdateVideoExplainActive(ids, action_type string) bool {
	if ids == "" || action_type == "" {
		return false
	}

	action_val := 0
	if action_type == "1" {
		action_val = 1
	}
	_, err := db.GetDb().Exec(fmt.Sprintf("UPDATE ZJ_EXPLAIN_URL SET ACTIVE=? WHERE ID IN (%s)", ids), action_val)
	if err != nil {
		utils.ErrorLog("admin_db_fun.go UpdateVideoExplainActive method.", err)
		return false
	}
	return true
}

//更新会员状态
//参数 action_type 1 启用，0禁用，默认0
func (self *AdminDbFun) UpdateMemberActive(ids, action_type string) bool {
	if ids == "" || action_type == "" {
		return false
	}

	action_val := 0
	if action_type == "1" {
		action_val = 1
	}
	_, err := db.GetDb().Exec(fmt.Sprintf("UPDATE ZJ_MEMBER SET ACTIVE=? WHERE ID IN (%s)", ids), action_val)
	if err != nil {
		utils.ErrorLog("admin_db_fun.go UpdateMemberActive method.", err)
		return false
	}
	return true
}

//更改会员过期时间
func (self *AdminDbFun) UpdateMemberExpire(args map[string]interface{}) bool {
	id_val, id_has := args["id"]
	field_val, field_has := args["field"]
	date_val, date_has := args["value"]
	if !id_has || !field_has || !date_has {
		return false
	}

	sql_strs := bytes.NewBufferString("UPDATE ZJ_MEMBER SET ")
	field_str, field_str_ok := field_val.(string)
	if field_str_ok {
		sql_strs.Write([]byte(field_str))
	}
	sql_strs.Write([]byte("=? WHERE ID=?"))

	_, err := db.GetDb().Exec(sql_strs.String(), date_val, id_val)
	if err != nil {
		utils.ErrorLog("admin_db_fun.go AddVideoExplainUrl method.", err)
		return false
	}
	return true
}

//添加视频解析地址
func (self *AdminDbFun) AddVideoExplainUrl(args map[string]interface{}) bool {
	_, err := db.GetDb().Exec("INSERT INTO ZJ_EXPLAIN_URL (URL_ADDR,CREATE_TIME,ACTIVE,TYPE) VALUES(?, CURRENT_TIME, 1,?)", args["url_addr"], args["type"])
	if err != nil {
		utils.ErrorLog("admin_db_fun.go AddVideoExplainUrl method.", err)
		return false
	}
	return true
}

//更新全部的url地址
func (self *AdminDbFun) UpdateAllExplainUrl(explains []map[string]interface{}) bool {
	tx, err := db.GetDb().Begin()
	if err != nil {
		utils.ErrorLog("admin_db_fun.go UpdateAllExplainUrl method. init tx error.", err)
		return false
	}

	//删除原来的数据
	tx.Exec("TRUNCATE TABLE ZJ_EXPLAIN_URL")

	//批量添加
	stmt, err := tx.Prepare("INSERT INTO ZJ_EXPLAIN_URL (URL_ADDR,CREATE_TIME,ACTIVE,TYPE) VALUES(?, CURRENT_TIME, 1,?)")
	if err != nil {
		utils.ErrorLog("admin_db_fun.go UpdateAllExplainUrl method. init stmt error.", err)
		return false
	}
	defer stmt.Close()
	for _, explain := range explains {
		_, err = stmt.Exec(explain["url"], explain["type"])
		if err != nil {
			return false
			utils.ErrorLog("admin_db_fun.go UpdateAllExplainUrl method. init stmt Exec error.", err)
		}
	}

	tx.Commit()

	return true
}

//根据条件查询会员信息
func (self *AdminDbFun) FindMember(args, result map[string]interface{}) {
	after_strs := bytes.NewBufferString("")
	params := make([]interface{}, 0)

	if active_val, active_has := args["active"]; active_has {
		after_strs.Write([]byte(" AND ACTIVE=?"))
		params = append(params, active_val)
	}
	if phone_num, phone_has := args["phone_num"]; phone_has {
		after_strs.Write([]byte(" AND PHONE_NUM=?"))
		params = append(params, phone_num)
	}
	if nick_name, nick_has := args["nick_name"]; nick_has {
		after_strs.Write([]byte(" AND NICK_NAME=?"))
		params = append(params, nick_name)
	}

	after_strs.Write([]byte(" ORDER BY REG_TIME DESC"))

	//查询总条数
	count_sql := bytes.NewBufferString("")
	count_sql.Write([]byte("SELECT COUNT(ID) FROM ZJ_MEMBER WHERE 1=1"))
	count_sql.Write(after_strs.Bytes())
	row := db.GetDb().QueryRow(count_sql.String(), params...)
	var rows_count uint64
	err := row.Scan(&rows_count)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.ErrorLog("admin_db_fun.go FindMember method.", err)
		}
		return
	}
	result["count"] = rows_count

	//判断分页参数
	if startLine_val, startLine_has := args["startLine"]; startLine_has {
		after_strs.Write([]byte(" LIMIT ?"))
		params = append(params, startLine_val)
	}
	if endLine_val, endLine_has := args["endLine"]; endLine_has {
		after_strs.Write([]byte( " ,?"))
		params = append(params, endLine_val)
	}

	//查询信息
	sql_strs := bytes.NewBufferString("")
	sql_strs.Write([]byte("SELECT ID, PHONE_NUM, NICK_NAME, REG_TIME, AIQIYI_EXPIRE, YOUKU_EXPIRE, LETV_EXPIRE, TENTCENT_EXPIRE, ACTIVE FROM ZJ_MEMBER WHERE 1=1"))
	sql_strs.Write(after_strs.Bytes())
	rows, err := db.GetDb().Query(sql_strs.String(), params...)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.ErrorLog("admin_db_fun.go FindMember method.", err)
		}
		return
	}

	defer rows.Close()
	result_rows := make([]map[string]interface{}, 0)
	for rows.Next() {
		var (
			id                                                                                        uint64
			active                                                                                    uint8
			phone_num, nick_name, reg_time, aiqiyi_expire, youku_expire, letv_expire, tentcent_expire string
		)
		rows.Scan(&id, &phone_num, &nick_name, &reg_time, &aiqiyi_expire, &youku_expire, &letv_expire, &tentcent_expire, &active)
		result_row := map[string]interface{}{
			"id":              id,
			"phone_num":       phone_num,
			"nick_name":       nick_name,
			"reg_time":        reg_time[:19],
			"aiqiyi_expire":   aiqiyi_expire[:10],
			"youku_expire":    youku_expire[:10],
			"letv_expire":     letv_expire[:10],
			"tentcent_expire": tentcent_expire[:10],
			"active_str":      self.member_active2str(active),
			"active":          active,
		}
		result_rows = append(result_rows, result_row)
	}

	result["data"] = result_rows
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
			"id":              id,
			"url":             url_addr,
			"create_time":     utils.FormatDataTime(create_time[:10]),
			"create_time_str": create_time[:19],
			"active":          active,
			"active_str":      self.active2str(active),
			"type_str":        self.type2str(row_type),
			"type":            row_type,
		}
		result_rows = append(result_rows, result_row)
	}

	result["data"] = result_rows
}

//将会员地址状态转化为字符串
func (self *AdminDbFun) member_active2str(row_active uint8) string {
	str := ""
	switch row_active {
	case 1:
		str = "正常"
	default:
		str = "禁用"
	}
	return str
}

//将解析地址状态转化为字符串
func (self *AdminDbFun) active2str(row_active uint8) string {
	str := ""
	switch row_active {
	case 1:
		str = "激活"
	default:
		str = "作废"
	}
	return str
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
	return &adminDbFun
}
