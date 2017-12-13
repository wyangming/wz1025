package utils

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

//返回正则表达式匹配的所有字符串，如果没有则返回nil
//二维数组行是匹配的次数，列是原字符串跟()里的字符串
func SpiderRegInfo(reg_str string, str *string) [][]string {
	if len(reg_str) < 1 || len(*str) < 1 {
		ErrorLog("[error]SpiderRegInfo reg_str or str is nil ", nil)
		return nil
	}

	reg, err := regexp.Compile(reg_str)
	if err != nil {
		ErrorLog("[error]regexp.Compile ", err)
		return nil
	}
	return reg.FindAllStringSubmatch(*str, -1)
}

//非常简单的抓取网页源代码
//用的是get方式抓取
func SpiderHtmlSrc(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ErrorLog("[error]http.NewRequest ", err)
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req_client := &http.Client{}
	res, err := req_client.Do(req)
	if err != nil {
		ErrorLog("[error] ", err)
		return ""
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		InfoLog("网站爬取错误，状态为", res.Status)
		return ""
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ErrorLog("[error]ioutil.ReadyAll", err)
		return ""
	}
	return string(body)
}
