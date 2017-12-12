package spider

import (
	"fmt"
	"wz1025/utils"
	"regexp"
	"os"
)

const (
	VIDEO_BASE_URL = "http://www.dy2018.com"
)

var VIDEO_FILM_CHANEL = []string{"/0", "/1", "/2", "/3", "/4", "/5", "/7", "/8", "/14", "/15"}
var VIDEO_TV_CHANEL = []string{"/html/tv/hytv", "/html/tv/hytv", "/html/tv/hytv"}

func spider_video_film() {
	url_str := fmt.Sprintf("%s%s", VIDEO_BASE_URL, VIDEO_FILM_CHANEL[0])
	index_chanel_str := utils.SpiderHtmlSrc(url_str)
	fmt.Println(index_chanel_str)
	if len(index_chanel_str) < 0 {
		return
	}
	page_reg, err := regexp.Compile(fmt.Sprintf("<option value='%s/index_([1-9]+?).html'>",VIDEO_FILM_CHANEL[0]))
	if err != nil {
		os.Exit(-1)
	}
	pages := page_reg.FindAllStringSubmatch(index_chanel_str, -1)
	for _, page := range pages {
		fmt.Println(page)
	}
}
