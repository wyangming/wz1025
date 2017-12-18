package spider

import (
	"wz1025/model"
	"wz1025/db"
	"wz1025/utils"
	"sync"
	"os"
)

type VideoRecDbFun uint8

var videoRecDbFun = VideoRecDbFun(1)

//批量保存电影资源数据
func (this *VideoRecDbFun) SpiderVideoRecsSave(videoRecs *[]*model.SpiderVideoRec) {
	if len(*videoRecs) < 1 {
		return
	}

	db := db.GetDb()
	tx, err := db.Begin()
	stmt_obj, err := tx.Prepare("INSERT INTO zj_video_rec (video_type,video_plot_type,video_place,video_spider_url,video_spider_update_date,video_name,video_title,video_down_urls) VALUES(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE video_down_urls=?,video_spider_update_date=?")
	if err != nil {
		utils.ErrorLog("[error]VideoRecDbFun.SpiderVideoRecsSave db.Prepare ", err)
	}
	stmt_wait := &sync.WaitGroup{}
	stmt_wait.Add(len(*videoRecs))
	for _, videoRec := range *videoRecs {
		go func(tmp_videoRecs *model.SpiderVideoRec) {
			_, err := stmt_obj.Exec(tmp_videoRecs.VideoType, tmp_videoRecs.VideoPlotType, tmp_videoRecs.VideoPlace, tmp_videoRecs.VideoSpiderUrl, tmp_videoRecs.VideoSpiderUpdateDate, tmp_videoRecs.VideoName, tmp_videoRecs.VideoTitle, tmp_videoRecs.VideoDownUrls, tmp_videoRecs.VideoDownUrls, tmp_videoRecs.VideoSpiderUpdateDate)
			if err != nil {
				utils.ErrorLog("[error]VideoRecDbFun.SpiderVideoRecsSave stmt_obj.Exec ", err)
				utils.ErrorLog("[error]info is ", *tmp_videoRecs)
				os.Exit(-1)
			}
			stmt_wait.Done()
		}(videoRec)

	}
	stmt_wait.Wait()

	tx.Commit()
	//关闭链接
	stmt_obj.Close()
}

func NewVideoRecDbFun() *VideoRecDbFun {
	return &videoRecDbFun
}
