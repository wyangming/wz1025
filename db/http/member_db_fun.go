package http

import (
	"wz1025/db"
	"wz1025/utils"
)

type MemberDbFun struct {
}

var memberDbFun = &MemberDbFun{}

//根据一个视频类型查询视频的解析地址
func (self *MemberDbFun) VideoExplainUrlByType(v_type uint8) string {
	row := db.GetDb().QueryRow("select URL_ADDR from zj_explain_url where TYPE=? LIMIT 1", v_type)
	var url_addr string
	err := row.Scan(&url_addr)
	if err != nil {
		utils.ErrorLog("member_db_fun.go VideoExplainUrlByType method.", err)
		return ""
	}
	return url_addr
}

//得到操作关于member数据的结构信息
func NewMemberDbFun() *MemberDbFun {
	return memberDbFun
}
