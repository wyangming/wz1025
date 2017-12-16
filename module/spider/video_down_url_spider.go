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

var VIDEO_FILM_CHANELS = []string{"/html/gndy/dyzz", "/html/gndy/rihan", "/html/gndy/oumei", "/html/gndy/china", "/html/gndy/jddy", "/html/gndy/rihan"}
var VIDEO_TV_CHANELS = []string{"/html/tv/gangtai", "/html/tv/hepai", "/html/tv/hytv", "/html/tv/rihantv", "/html/tv/oumeitv"}

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
	now_date_str := time.Now().Format("2006-01-02")

	for _, video_film_chanel := range VIDEO_TV_CHANELS {
		//一个栏目保存了次
		var videoRecs []*model.SpiderVideoRec
		is_spider := true
		videoRecs = make([]*model.SpiderVideoRec, 0)
		url_str := fmt.Sprintf("%s%s", VIDEO_BASE_URL, video_film_chanel)

		//解析一个列表里所有的电影地址
		index_chanel_str := funcListVideoFilms(&url_str, &now_date_str, &is_spider, &videoRecs, &pre_date)

		//判断是否还需爬取下一页
		if !is_spider {
			fmt.Println(fmt.Sprintf("url %s is over", video_film_chanel))
			spider.NewVideoRecDbFun().SpiderVideoRecsSave(&videoRecs)
			continue
		}

		pagesInfos := utils.SpiderRegInfo("<option value='(list_[0-9]+_[0-9]+?.html)'>", index_chanel_str)
		if pagesInfos == nil {
			fmt.Println(fmt.Sprintf("url %s is over", video_film_chanel))
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
					funcListVideoFilms(&tmp_url_str, &now_date_str, &is_spider, &videoRecs, &pre_date)
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
		fmt.Println(fmt.Sprintf("url %s is over", video_film_chanel))
	}
}

//解析一个列表里所有的电影地址
func funcListVideoFilms(tmp_url_str, tmp_now_date_str *string, tmp_is_spider *bool, tmp_videoRecs *[]*model.SpiderVideoRec, tmp_pre_date *time.Time) *string {
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

	spider_video_film_list(&list_urls, tmp_is_spider, tmp_videoRecs, tmp_pre_date, tmp_now_date_str)
	return &index_chanel_str
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

		//标题与名字
		title_names := utils.SpiderRegInfo(`<title>(.*?)</title>`, &film_html_src)
		if title_names == nil {
			continue
		}
		title := title_names[0][1]
		name := strings.Replace(title,"迅雷下载_电影天堂","",-1)
		if strings.Index(title, "《") > 0 && strings.Index(title, "》") > 0 {
			name = title[(strings.Index(title, "《") + len("《")):strings.LastIndex(title, "》")]
		}

		//产地
		place_infos := utils.SpiderRegInfo(`<br />◎产[\s]{0,}地(.*?)<br />`, &film_html_src)
		place := ""
		if place_infos != nil {
			place = strings.TrimSpace(place_infos[0][1])
		}
		place=""

		//剧情
		plot_infos := utils.SpiderRegInfo(`<br />◎类[\s]{0,}别(.*?)<br />`, &film_html_src)
		plot := ""
		if plot_infos != nil {
			plot = strings.TrimSpace(plot_infos[0][1])
		}
		plot=""

		//组装模型信息
		videoRec := &model.SpiderVideoRec{
			VideoType:             1,
			VideoPlotType:         plot,
			VideoPlace:            place,
			VideoSpiderUrl:        signle_film_url,
			VideoSpiderUpdateDate: *now_date_str,
			VideoName:             name,
			VideoTitle:            title,
		}

		//转换utl
		urls := make([][]string,0)
		if videoRec.VideoType<0{
			down_url := down_urls[0][1]
			spider_video_convert_down_url(&down_url)
			urls = append(urls,[]string{"1", down_url})
		}else {
			for i,down_url:=range down_urls{
				spider_video_convert_down_url(&down_url[1])
				urls = append(urls,[]string{fmt.Sprintf("%d",(i+1)), down_url[1]})
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
