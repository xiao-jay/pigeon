package weibo_hot

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const Name = "weibo"

type DataItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Num   int    `json:"num"`
	Hot   string `json:"hot"`
}

func GetData() []DataItem {
	var data []DataItem

	response, err := http.Get("https://weibo.com/ajax/side/hotSearch")
	if err != nil {
		fmt.Println(err)
		return data
	}
	defer response.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return data
	}

	dataJSON := result["data"].(map[string]interface{})["realtime"].([]interface{})
	jyzy := map[string]string{
		"电影": "影",
		"剧集": "剧",
		"综艺": "综",
		"音乐": "音",
	}

	for _, dataItem := range dataJSON {
		hot := ""
		dataMap := dataItem.(map[string]interface{})

		// 如果是广告，则不添加
		if _, ok := dataMap["is_ad"]; ok {
			continue
		}

		if flagDesc, ok := dataMap["flag_desc"]; ok {
			if hotValue, ok := jyzy[flagDesc.(string)]; ok {
				hot = hotValue
			}
		}
		if _, ok := dataMap["is_boom"]; ok {
			hot = "爆"
		}
		if _, ok := dataMap["is_hot"]; ok {
			hot = "热"
		}
		if _, ok := dataMap["is_fei"]; ok {
			hot = "沸"
		}
		if _, ok := dataMap["is_new"]; ok {
			hot = "新"
		}

		dic := DataItem{
			Title: dataMap["note"].(string),
			URL:   "https://s.weibo.com/weibo?q=%23" + dataMap["word"].(string) + "%23",
			Num:   int(dataMap["num"].(float64)),
			Hot:   hot,
		}
		data = append(data, dic)
	}

	return data
}
