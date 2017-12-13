package spider

import (
	"fmt"
	"encoding/json"
	"wz1025/utils"
)

//初始化爬虫内容
func init() {
	//spider_explain_url()

	strs := [][]string{[]string{"1", "one"}, []string{"2", "two"}}
	bytes, err := json.Marshal(&strs)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(string(bytes))
	result := make([][]string, 0)
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	spider_video_film(utils.FormatDataTime("2015-12-11"))
}
