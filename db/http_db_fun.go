package db

import (
	//"database/sql"
	"wz1025/utils"
)

type HttpDbFun struct {
}

var httpDbFun = &HttpDbFun{}

//登录
//result：返回登录后用户的信息在map里，否则返回nil
func (this *HttpDbFun) Login(args map[string]interface{}) map[string]interface{} {
	row := db.QueryRow("SELECT ID,PHONE_NUM,NICK_NAME,REG_TIME,AIQIYI_EXPIRE,YOUKU_EXPIRE,LETV_EXPIRE,TENTCENT_EXPIRE,ACTIVE FROM zj_member WHERE PHONE_NUM=? AND PWD=? LIMIT 1", args["phone_number"], args["pwd"])

	var (
		id                                                                                        uint64
		active                                                                                    uint8
		phone_num, nick_name, reg_time, aiqiyi_expire, youku_expire, letv_expire, tentcent_expire string
	)

	err := row.Scan(&id, &phone_num, &nick_name, &reg_time, &aiqiyi_expire, &youku_expire, &letv_expire, &tentcent_expire, &active)
	if err != nil {
		utils.ErrorLog("http_db_fun.go Login method.", err)
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
func (this *HttpDbFun) RegMember(args map[string]interface{}) bool {
	_, err := db.Exec("INSERT INTO zj_member (NICK_NAME,PHONE_NUM,REG_TIME,AIQIYI_EXPIRE,YOUKU_EXPIRE,LETV_EXPIRE,TENTCENT_EXPIRE,PWD,ACTIVE) VALUES(?,?,CURRENT_TIME,CURRENT_DATE,CURRENT_DATE,CURRENT_DATE,CURRENT_DATE,?,1)", args["nick_name"], args["phone_number"], args["pwd"])
	if err != nil {
		utils.ErrorLog("http_db_fun.go RegMember method.", err)
		return false
	}
	return true
}

//得到操作关于http数据的结构信息
func NewHttpDbFun() *HttpDbFun {
	return httpDbFun
}
