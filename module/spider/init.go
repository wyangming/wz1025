package spider

import (
	"github.com/jakecoffman/cron"
	"github.com/astaxie/beego"
	"wz1025/utils"
)
//爬虫任务
var spider_CronJob *cron.Cron
//初始化爬虫内容
func init() {
	//把爬取的任务添加进去
	spider_CronJob= cron.New()
	spider_CronJob.AddFunc("23 13 */1 * * ?", spider_video_film_new, "spider_video_film")
	spider_CronJob.Start()

	//启动时执行任务
	spider_video_film_new()
	spider_video_film_detail(utils.FormatDataTime(beego.AppConfig.String("spiderFilimVideoTime")))
}
