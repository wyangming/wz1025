package spider

import (
	"wz1025/utils"
	"regexp"
	"fmt"
	"strings"
	db "wz1025/db/http"
	"github.com/jakecoffman/cron"
)

//解析地址爬虫
func spider_explain_url() {
	spec := "2 13 1 * * ?" //每天的1点13分2秒更新解析地址
	cronJob := cron.New()
	//cronJob.RemoveJob()//要删除任务使用这个方法
	cronJob.AddFunc(spec, con_spider_fun, "spider_explain_url")
	cronJob.Start()
}

//定时任务的方法
func con_spider_fun() {
	//k要抓取的url，v过滤的正则表达式
	url_map := map[string]string{
		"http://www.5ifxw.com/vip":              "<option value=\"([a-zA-z]+://[^\\s]*)\">(.*)</option>",
		"http://www.qmaile.com/?qqdrsign=05efd": "<option value=\"([a-zA-z]+://[^\\s]*)\" selected=\"\">(.*)</option>",
	}

	//要更新所有的解析地址
	explain_urls := make(map[string]map[string]interface{})

	for k, v := range url_map {
		explain_url(k, v, explain_urls)
	}

	explains := make([]map[string]interface{}, 0)
	for _, v := range explain_urls {
		explains = append(explains, v)
	}

	if (db.NewAdminDbFun().UpdateAllExplainUrl(explains)) {
		utils.InfoLog("update All explain url success")
	} else {
		utils.InfoLog("update All explain url faile")
	}
}

//爬取视频解析地址
func explain_url(url, reg_str string, res map[string]map[string]interface{}) {
	html_src := utils.SpiderHtmlSrc(url)
	if html_src == "" {
		utils.InfoLog(fmt.Sprintf("client %s is error", url))
		return
	}
	reg, _ := regexp.Compile(reg_str)
	infos := reg.FindAllStringSubmatch(html_src, -1)
	if infos == nil {
		utils.InfoLog("no find explain url")
	}

	//处理解析地址
	for _, info := range infos {
		if len(info) < 3 {
			continue
		}

		info_name := info[2]
		if strings.Contains(info_name, "&nbsp;") {
			continue
		}
		explain := make(map[string]interface{})
		explain["url"] = info[1]
		type_str := ""

		if strings.Contains(info_name, "通用") || strings.Contains(info_name, "万能") {
			explain["type"] = 0
			type_str = "0"
		} else if strings.Contains(info_name, "爱奇艺") {
			explain["type"] = 1
			type_str = "1"
		} else if strings.Contains(info_name, "优酷") {
			explain["type"] = 2
			type_str = "2"
		} else if strings.Contains(info_name, "腾讯") {
			explain["type"] = 3
			type_str = "3"
		} else if strings.Contains(info_name, "乐视") {
			explain["type"] = 4
			type_str = "4"
		} else if strings.Contains(info_name, "搜狐") {
			explain["type"] = 5
			type_str = "5"
		} else if strings.Contains(info_name, "土豆") {
			explain["type"] = 6
			type_str = "6"
		} else if strings.Contains(info_name, "芒果") {
			explain["type"] = 7
			type_str = "7"
		} else if strings.Contains(info_name, "PPTV") {
			explain["type"] = 7
			type_str = "7"
		}

		//如果没有直接添加到map里
		_, has := res[info[1]+type_str]
		if !has {
			res[info[1]+type_str] = explain
		}
	}
}
