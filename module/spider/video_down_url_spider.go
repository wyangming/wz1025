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
	"wz1025/db/spider"
	"sync"
)

const (
	VIDEO_BASE_URL = "http://www.dytt8.net"
)

func spider_video_film_new() {
	base_html_src := utils.SpiderHtmlSrc(VIDEO_BASE_URL)
	if len(base_html_src) < 1 {
		return
	}
	base_html_src = utils.StrGbk2Utf8(base_html_src)
	is_spider := false
	videoRecs := make([]*model.SpiderVideoRec, 0)
	pre_date := time.Now()
	date_str := pre_date.Format("2006-01-02")

	for i, out_reg := range VIDEO_INDEX_PAGE {

		new_film_infos := utils.SpiderRegInfo(out_reg, &base_html_src)
		if new_film_infos == nil {
			return
		}
		new_film_info := new_film_infos[0][1]
		if i == 2 || i == 5 {
			new_film_info = new_film_infos[1][1]
		}
		last_film_infos := utils.SpiderRegInfo(`<a href='(.*?)'>.*?</a><br/>`, &new_film_info)
		if last_film_infos == nil {
			return
		}
		video_type := uint8(0)
		if i > 0 && i != 5 {
			video_type = uint8(1)
		}
		spider_video_film_list(&last_film_infos, &is_spider, &videoRecs, &pre_date, &date_str, video_type)
	}
	spider.NewVideoRecDbFun().SpiderVideoRecsSave(&videoRecs)
}

var VIDEO_INDEX_PAGE = []string{`<!--{start:最新影视下载-->([\s\S]*?)<!--}end:最新下载--->`, `<!--{start:最新TV下载-->([\s\S]*?)<!--}end:最新TV下载--->`, `<!--{start:最新TV下载-->([\s\S]*?)<!--}end:最新TV下载--->`, `<!--{start:最新欧美剧集下载-->([\s\S]*?)<!--}end:最新欧美剧集下载--->`, `<!--{start:日韩电视推荐-->([\s\S]*?)<!--}end:日韩电视推荐-->`, `<!--{start:最新影视下载-->([\s\S]*?)<!--}end:最新下载--->`}
var VIDEO_FILM_CHANELS = []string{"/html/gndy/dyzz", "/html/gndy/rihan", "/html/gndy/oumei", "/html/gndy/china", "/html/gndy/jddy", "/html/gndy/rihan"}
var VIDEO_TV_CHANELS = []string{"/html/tv/gangtai", "/html/tv/hepai", "/html/tv/hytv", "/html/tv/rihantv", "/html/tv/oumeitv"}

func spider_video_convert_down_url(url_str *string) {
	*url_str = fmt.Sprintf("AA%sZZ", url.PathEscape(*url_str))
	*url_str = strings.Replace(*url_str, "%2F", "/", -1)
	*url_str = strings.Replace(*url_str, "%5B", "[", -1)
	*url_str = strings.Replace(*url_str, "%5D", "]", -1)
	*url_str = fmt.Sprintf("thunder://%s", base64.StdEncoding.EncodeToString([]byte(*url_str)))
}

func spider_video_film() {
	fmt.Println("开始爬取内容")
	h, _ := time.ParseDuration("-1h")
	pre_date := time.Now().Add(1 * h)
	spider_video_film_detail(pre_date)
}

//爬取时间点之后的电影
func spider_video_film_detail(pre_date time.Time) {
	spider_video_by_chanels(&VIDEO_FILM_CHANELS, 0, pre_date)
	spider_video_by_chanels(&VIDEO_TV_CHANELS, 1, pre_date)
}
func spider_video_by_chanels(chanels *[]string, video_type int, pre_date time.Time) {
	for _, video_film_chanel := range *chanels {
		now_date_str := time.Now().Format("2006-01-02")
		//一个栏目保存了次
		var videoRecs []*model.SpiderVideoRec
		is_spider := true
		videoRecs = make([]*model.SpiderVideoRec, 0)
		url_str := fmt.Sprintf("%s%s", VIDEO_BASE_URL, video_film_chanel)

		//解析一个列表里所有的电影地址
		index_chanel_str := funcListVideoFilms(&url_str, &now_date_str, &is_spider, &videoRecs, &pre_date, uint8(video_type))

		//判断是否还需爬取下一页
		if !is_spider {
			spider.NewVideoRecDbFun().SpiderVideoRecsSave(&videoRecs)
			continue
		}

		pagesInfos := utils.SpiderRegInfo("<option value='(list_[0-9]+_[0-9]+?.html)'>", index_chanel_str)
		if pagesInfos == nil {
			spider.NewVideoRecDbFun().SpiderVideoRecsSave(&videoRecs)
			continue
		}
		//多线程
		pagesInfos_count := len(pagesInfos)
		for i := 0; i < pagesInfos_count; {
			page_waite := &sync.WaitGroup{}
			waite_count := 5
			page_waite.Add(waite_count)
			for waite_count > 0 {

				tmp_url_str := fmt.Sprintf("%s%s/%s", VIDEO_BASE_URL, video_film_chanel, pagesInfos[i][1])
				go func() {
					funcListVideoFilms(&tmp_url_str, &now_date_str, &is_spider, &videoRecs, &pre_date, uint8(video_type))
					page_waite.Done()
				}()

				waite_count -= 1
				i += 1
				if i >= pagesInfos_count {
					break
				}
			}
			//处理剩余的线程
			for waite_count > 0 {
				page_waite.Done()
				waite_count -= 1
			}
			page_waite.Wait()
		}
		spider.NewVideoRecDbFun().SpiderVideoRecsSave(&videoRecs)
	}
}

//解析一个列表里所有的电影地址
func funcListVideoFilms(tmp_url_str, tmp_now_date_str *string, tmp_is_spider *bool, tmp_videoRecs *[]*model.SpiderVideoRec, tmp_pre_date *time.Time, video_type uint8) *string {
	index_chanel_str := utils.SpiderHtmlSrc(*tmp_url_str)
	if len(index_chanel_str) < 0 {
		return &index_chanel_str
	}
	index_chanel_str = utils.StrGbk2Utf8(index_chanel_str)

	//过滤列表电影的url地址
	list_urls := utils.SpiderRegInfo(`<a href="(.*?)" class="ulink">.*?</a>[\s\S]*?<font color="#8F8C89">日期：(.*? .*?) [\s\S]*?点击：[0-9]+ </font>`, &index_chanel_str)
	if list_urls == nil {
		list_urls = utils.SpiderRegInfo(`<a href="(.*?)" class="ulink" title=".*?">.*?</a>[\s\S]*?<font color="#8F8C89">日期：(.*?)[\s\S]点击：[0-9]+ </font>`, &index_chanel_str)
		if list_urls == nil {
			return &index_chanel_str
		}
	}

	spider_video_film_list(&list_urls, tmp_is_spider, tmp_videoRecs, tmp_pre_date, tmp_now_date_str, video_type)
	return &index_chanel_str
}
func spider_video_film_list(list_urls *[][]string, is_spider *bool, videoRecs *[]*model.SpiderVideoRec, pre_date *time.Time, now_date_str *string, video_type uint8) {
	for _, url_str := range *list_urls {
		if len(url_str) > 2 {
			*is_spider = utils.FormatDataTime(url_str[2]).After(*pre_date)
			if !*is_spider {
				break
			}
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
			utils.ErrorLog("[error]find down load url error, no find down_urls from ", url_str)
			continue
		}

		//标题与名字
		title_names := utils.SpiderRegInfo(`<title>(.*?)</title>`, &film_html_src)
		if title_names == nil {
			continue
		}
		title := title_names[0][1]
		name := strings.Replace(title, "迅雷下载_电影天堂", "", -1)
		if strings.Index(title, "《") > 0 && strings.Index(title, "》") > 0 {
			name = title[(strings.Index(title, "《") + len("《")):strings.LastIndex(title, "》")]
		}

		//产地
		place_infos := utils.SpiderRegInfo(`<br />◎产[\s]{0,}地(.*?)<br />`, &film_html_src)
		place := ""
		if place_infos != nil {
			place = strings.TrimSpace(place_infos[0][1])
		}
		place = ""

		//剧情
		plot_infos := utils.SpiderRegInfo(`<br />◎类[\s]{0,}别(.*?)<br />`, &film_html_src)
		plot := ""
		if plot_infos != nil {
			plot = strings.TrimSpace(plot_infos[0][1])
		}
		plot = ""
		//发布时间
		res_time_infos := utils.SpiderRegInfo(`发布时间：(([0-9]{3}[1-9]|[0-9]{2}[1-9][0-9]{1}|[0-9]{1}[1-9][0-9]{2}|[1-9][0-9]{3})-(((0[13578]|1[02])-(0[1-9]|[12][0-9]|3[01]))|((0[469]|11)-(0[1-9]|[12][0-9]|30))|(02-(0[1-9]|[1][0-9]|2[0-8]))))`, &film_html_src)
		res_time := (*now_date_str)[:10]
		if res_time_infos != nil {
			res_time = strings.TrimSpace(res_time_infos[0][1])
		}

		//组装模型信息
		videoRec := &model.SpiderVideoRec{
			VideoType:             video_type,
			VideoPlotType:         plot,
			VideoPlace:            place,
			VideoSpiderUrl:        signle_film_url,
			VideoSpiderUpdateDate: *now_date_str,
			VideoName:             name,
			VideoTitle:            title,
			VideoReleaseDate:      res_time,
		}

		//转换utl
		urls := make([][]string, 0)
		if videoRec.VideoType < 0 {
			down_url := down_urls[0][1]
			spider_video_convert_down_url(&down_url)
			urls = append(urls, []string{"1", down_url})
		} else {
			for i, down_url := range down_urls {
				spider_video_convert_down_url(&down_url[1])
				urls = append(urls, []string{fmt.Sprintf("%d", (i + 1)), down_url[1]})
			}
		}

		bytes, err := json.Marshal(&urls)
		if err != nil {
			utils.ErrorLog("[error]json.Marshal ", err)
			continue
		}

		videoRec.VideoDownUrls = string(bytes)

		*videoRecs = append(*videoRecs, videoRec)
	}
}
