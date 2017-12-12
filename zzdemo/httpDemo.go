package zzdemo

//问题：1、在爬http://studygolang.com/pkgdoc这个网站后会被过滤为400的错误
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"encoding/base64"
)

//爬虫测试
func HttpDemo() {
	//thunder://
	//AA
	//ftp://d:d@dygodj8.com:12311/[电影天堂www.dy2018.com]王牌特工2：黄金圈BD中英双字.mp4
	//ZZ
	url_str := "AA" + url.PathEscape("ftp://d:d@dygodj8.com:12311/[电影天堂www.dy2018.com]王牌特工2：黄金圈BD中英双字.mp4") + "ZZ"
	url_str = strings.Replace(url_str, "%2F", "/", -1)
	url_str = strings.Replace(url_str, "%5B", "[", -1)
	url_str = strings.Replace(url_str, "%5D", "]", -1)
	fmt.Println("thunder://" + base64.StdEncoding.EncodeToString([]byte(url_str)))
	var result = true
	if result {
		return
	}
	req, err := http.NewRequest("GET", "http://wzshipin.com/", nil)
	if err != nil {
		fmt.Println("[error]http.NewRequest ", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req_client := &http.Client{}
	res, err := req_client.Do(req)
	if err != nil {
		fmt.Println("[error]req_client.Do ", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("网站爬取错误，状态为", res.Status)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("[error]ioutil.ReadyAll", err)
		return
	}
	fmt.Println(res.Cookies())
	html_str := string(body)
	fmt.Println(html_str)
	fmt.Println("cookies信息为：", res.Cookies())
}

//电影列表页面得到单个电影的链接
//<a href="(.*?)" class="ulink" title=".*?">.*?</a>
//电影列表匹配更新日期
//<font color="#8F8C89">日期：(.*?)  </font>
//电影列表里得到所有的页数
//<option value='/3/index_([1-9]\d).html'>
//单个电影页面的下载转换信息
//<td style="WORD-WRAP: break-word" bgcolor="#fdfddf"><a href="(ftp://.*?)">.*?</a></td>

//thunder://QUFmdHA6Ly9kOmRAZHlnb2RqOC5jb206MTIzMTEvWyVFNyU5NCVCNSVFNSVCRCVCMSVFNSVBNCVBOSVFNSVBMCU4Mnd3dy5keTIwMTguY29tXSVFNyU4RSU4QiVFNyU4OSU4QyVFNyU4OSVCOSVFNSVCNyVBNTIlRUYlQkMlOUElRTklQkIlODQlRTklODclOTElRTUlOUMlODhCRCVFNCVCOCVBRCVFOCU4QiVCMSVFNSU4RiU4QyVFNSVBRCU5Ny5tcDRaWg==
