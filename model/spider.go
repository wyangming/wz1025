package model

//爬取视频资源的结构体
type SpiderVideoRec struct {
	VideoType             uint8  //视频类型，0电影，1电视剧类
	VideoPlotType         string //剧情的类型，爱情动画之类的
	VideoPlace            string //产地
	VideoSpiderUrl        string //爬虫的来源地址
	VideoSpiderUpdateDate string //最近一次爬虫爬取的日期
	VideoName             string //视频名称
	VideoTitle            string //视频标题
	VideoDownUrls         string //下载链接,json数组格式:[["1","one"],["2","two"]];列的每一位是每几集，第二位是下载地址
	VideoReleaseDate     string //发布日期
}
