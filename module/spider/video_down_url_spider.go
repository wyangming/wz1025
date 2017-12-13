package spider

import (
	"fmt"
	"wz1025/utils"
	"time"
	"strings"
	"net/url"
	"encoding/base64"
	"wz1025/model"
	"encoding/json"
)

const (
	VIDEO_BASE_URL = "http://www.dy2018.com"
)

var VIDEO_FILM_CHANELS = []string{"/0", "/1", "/2", "/3", "/4", "/5", "/7", "/8", "/14", "/15"}
var VIDEO_TV_CHANELS = []string{"/html/tv/hytv", "/html/tv/hytv", "/html/tv/hytv"}

func spider_video_convert_down_url(url_str *string) {
	*url_str = fmt.Sprintf("AA%sZZ", url.PathEscape(*url_str))
	*url_str = strings.Replace(*url_str, "%2F", "/", -1)
	*url_str = strings.Replace(*url_str, "%5B", "[", -1)
	*url_str = strings.Replace(*url_str, "%5D", "]", -1)
	*url_str = fmt.Sprintf("thunder://%s", base64.StdEncoding.EncodeToString([]byte(*url_str)))
}

//爬取那日之前的电影
func spider_video_film(pre_date time.Time) {
	//表示是否还需要爬取与是否是第一页
	is_spider := true
	now_date_str := time.Now().Format("2006-01-02")

	var videoRecs []*model.SpiderVideoRec

	for _, video_film_chanel := range VIDEO_FILM_CHANELS {
		//一个栏目保存了次
		videoRecs = make([]*model.SpiderVideoRec, 0)
		url_str := fmt.Sprintf("%s%s", VIDEO_BASE_URL, video_film_chanel)
		index_chanel_str := utils.SpiderHtmlSrc(url_str)
		if len(index_chanel_str) < 0 {
			return
		}
		index_chanel_str = utils.StrGbk2Utf8(index_chanel_str)

		//过滤列表电影的url地址
		list_urls := utils.SpiderRegInfo(`<a href="(.*?)" class="ulink" title=".*?">.*?</a>[\s\S]*?<font color="#8F8C89">日期：(.*?)  </font>`, &index_chanel_str)
		if list_urls == nil {
			return
		}

		spider_video_film_list(&list_urls, &is_spider, &videoRecs, &pre_date, &now_date_str)

		//判断是否还需爬取下一页
		if !is_spider {
			return
		}

		pagesInfos := utils.SpiderRegInfo(fmt.Sprintf("<option value='(%s/index_([1-9]+?).html)'>", video_film_chanel), &index_chanel_str)
		if pagesInfos == nil {
			return
		}
		//pages := make([]string, len(pagesInfos))
		for _, page := range pagesInfos {
			fmt.Println(page[1])
		}
		break
	}
}

func spider_video_film_list(list_urls *[][]string, is_spider *bool, videoRecs *[]*model.SpiderVideoRec, pre_date *time.Time, now_date_str *string) {
	for _, url_str := range *list_urls {
		*is_spider = utils.FormatDataTime(url_str[2]).After(*pre_date)
		if !*is_spider {
			break
		}

		//爬取单个电影页面
		signle_film_url := fmt.Sprintf("%s%s", VIDEO_BASE_URL, url_str[1])
		film_html_src := utils.SpiderHtmlSrc(signle_film_url)
		if len(film_html_src) < 0 {
			continue
		}
		film_html_src = utils.StrGbk2Utf8(film_html_src)

		//下载链接
		down_urls := utils.SpiderRegInfo(`<td style="WORD-WRAP: break-word" bgcolor="#fdfddf"><a href="(ftp://.*?)">.*?</a></td>`, &film_html_src)
		if down_urls == nil {
			continue
		}
		down_url := down_urls[0][1]
		spider_video_convert_down_url(&down_url)

		//标题与名字
		title_names := utils.SpiderRegInfo(`<title>(.*?)</title>`, &film_html_src)
		if title_names == nil {
			continue
		}
		title := title_names[0][1]
		name := title[(strings.Index(title, "《") + len("《")):strings.LastIndex(title, "》")]

		//产地
		place_infos := utils.SpiderRegInfo(`<p>◎产　　地　(.*?)</p>`, &film_html_src)
		if place_infos == nil {
			continue
		}
		place := place_infos[0][1]

		//剧情
		plot_infos := utils.SpiderRegInfo(`<p>◎类　　别　(.*?)</p>`, &film_html_src)
		if plot_infos == nil {
			continue
		}
		plot := plot_infos[0][1]

		//组装模型信息
		videoRec := &model.SpiderVideoRec{
			VideoType:             0,
			VideoPlotType:         plot,
			VideoPlace:            place,
			VideoSpiderUrl:        signle_film_url,
			VideoSpiderUpdateDate: *now_date_str,
			VideoName:             name,
			VideoTitle:            title,
		}
		urls := [][]string{[]string{"1", down_url}}
		bytes, err := json.Marshal(&urls)
		if err != nil {
			utils.ErrorLog("[error]json.Marshal ", err)
			continue
		}
		videoRec.VideoDownUrls = string(bytes)
		*videoRecs = append(*videoRecs, videoRec)
	}
}
