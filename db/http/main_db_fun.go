package http

import (
	"wz1025/db"
	"wz1025/utils"
	"database/sql"
)

type MainDbFun uint8

var mainDbFun = MainDbFun(1)

//登录
//result：返回登录后用户的信息在map里，否则返回nil
func (self *MainDbFun) Login(args map[string]interface{}) map[string]interface{} {
	row := db.GetDb().QueryRow("SELECT ID,PHONE_NUM,NICK_NAME,REG_TIME,AIQIYI_EXPIRE,YOUKU_EXPIRE,LETV_EXPIRE,TENTCENT_EXPIRE,ACTIVE FROM zj_member WHERE PHONE_NUM=? AND PWD=? LIMIT 1", args["phone_number"], args["pwd"])

	var (
		id                                                                                        uint64
		active                                                                                    uint8
		phone_num, nick_name, reg_time, aiqiyi_expire, youku_expire, letv_expire, tentcent_expire string
	)

	err := row.Scan(&id, &phone_num, &nick_name, &reg_time, &aiqiyi_expire, &youku_expire, &letv_expire, &tentcent_expire, &active)
	if err != nil {
		utils.ErrorLog("[error]main_db_fun.go Login method.row.Scan err is ", err)
		return nil
	}

	reply := map[string]interface{}{
		"id":              id,
		"phone_num":       phone_num,
		"nick_name":       nick_name,
		"active":          active,
		"reg_time":        utils.FormatDataTime(reg_time),
		"aiqiyi_expire":   utils.FormatDataTime(aiqiyi_expire),
		"youku_expire":    utils.FormatDataTime(youku_expire),
		"letv_expire":     utils.FormatDataTime(letv_expire),
		"tentcent_expire": utils.FormatDataTime(tentcent_expire),
	}
	return reply
}

//注册会员
//resule:成功true，否则false
func (self *MainDbFun) RegMember(args map[string]interface{}) bool {
	//判断手机号是否已经存在
	row_has := db.GetDb().QueryRow("SELECT ID from zj_member WHERE PHONE_NUM=? LIMIT 1", args["phone_number"])
	var pre_id int64
	err := row_has.Scan(&pre_id)
	if err != nil {
		if err == sql.ErrNoRows {
			pre_id = int64(0)
		} else {
			utils.ErrorLog("[error]main_db_fun.go RegMember method.db.GetDb().QueryRow err is ", err)
			return false
		}
	}

	//如果用户已经存在则直接更改
	if pre_id > 0 {
		_, err = db.GetDb().Exec("UPDATE zj_member SET PWD=?, NICK_NAME=? WHERE ID=?", args["pwd"], args["nick_name"], pre_id)
		if err != nil {
			utils.ErrorLog("[error]main_db_fun.go RegMember method.db.GetDb() err is ", err)
			return false
		}
		return true
	}

	_, err = db.GetDb().Exec("INSERT INTO zj_member (NICK_NAME,PHONE_NUM,REG_TIME,AIQIYI_EXPIRE,YOUKU_EXPIRE,LETV_EXPIRE,TENTCENT_EXPIRE,PWD,ACTIVE) VALUES(?,?,CURRENT_TIME,CURRENT_DATE,CURRENT_DATE,CURRENT_DATE,CURRENT_DATE,?,1)", args["nick_name"], args["phone_number"], args["pwd"])
	if err != nil {
		utils.ErrorLog("[error]main_db_fun.go RegMember method. db.GetDb().Exec err is ", err)
		return false
	}
	return true
}

//得到操作关于main数据的结构信息
func NewMainDbFun() *MainDbFun {
	return &mainDbFun
}
