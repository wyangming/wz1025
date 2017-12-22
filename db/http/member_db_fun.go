package http

import (
	"wz1025/db"
	"wz1025/utils"
	"database/sql"
	"encoding/json"
	"crypto/md5"
	"fmt"
	"strconv"
)

type MemberDbFun uint8

var memberDbFun = MemberDbFun(1)

//根据一个视频类型查询视频的解析地址
func (self *MemberDbFun) VideoExplainUrlByType(v_type uint8) string {
	row := db.GetDb().QueryRow("select URL_ADDR from zj_explain_url where ACTIVE=1 AND TYPE IN (0,?) ORDER BY TYPE DESC LIMIT 1", v_type)
	var url_addr string
	err := row.Scan(&url_addr)
	if err != nil {
		utils.ErrorLog("[error]member_db_fun.go VideoExplainUrlByType method. row.Scan err is", err)
		return ""
	}
	return url_addr
}

//根据用户输入的视频名称查询视频资源
func (self *MemberDbFun) FindSpiderVideoRecs(video_name string) (*[]map[string]interface{}) {
	ret := make([]map[string]interface{}, 0)
	rows, err := db.GetDb().Query("select video_name,video_title,video_down_urls,video_type from zj_video_rec where video_name like ? limit 20", "%"+video_name+"%")
	if err != nil {
		if err != sql.ErrNoRows {
			utils.ErrorLog("[error]member_db_fun.go FindSpiderVideoRecs method. db.GetDb().Query err is ", err)
		}
		return &ret
	}

	for rows.Next() {
		var (
			video_name, video_title, video_down_urls, video_type_name string
			video_type                                                uint8
		)
		down_infos := make([][]string, 0)
		rows.Scan(&video_name, &video_title, &video_down_urls, &video_type)

		//把下载链接转换为数组
		err := json.Unmarshal([]byte(video_down_urls), &down_infos)
		if err != nil {
			utils.ErrorLog("[error]member_db_fun.go FindSpiderVideoRecs method. json.Unmarshal err is ", err)
		}

		//处理下载链接的信息计算id
		md5_h := md5.New()
		for i := 0; i < len(down_infos); i++ {
			//根据下载链接计算
			md5_h.Write([]byte(down_infos[i][1]))
			md5_str := fmt.Sprintf("%x", md5_h.Sum(nil))
			down_infos[i] = append(down_infos[i], md5_str)
			md5_h.Reset()

			//如果是电视剧类型处理每集的名称
			if video_type > 0 {
				index_down_url, _ := strconv.Atoi(down_infos[i][0])
				down_infos[i][0] = fmt.Sprintf("第%d集", index_down_url)
			} else {
				down_infos[i][0] = video_title
			}
		}

		if video_type > 0 {
			video_type_name = "电视剧"
		} else {
			video_type_name = "电影"
		}

		//添加到返回信息里
		ret = append(ret, map[string]interface{}{
			"video_name":      video_name,
			"down_infos":      down_infos,
			"video_type_name": video_type_name,
		})
	}
	return &ret
}

//得到操作关于member数据的结构信息
func NewMemberDbFun() *MemberDbFun {
	return &memberDbFun
}
