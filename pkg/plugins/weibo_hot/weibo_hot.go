package weibo_hot

import (
	"encoding/json"
	"fmt"
	"io"
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
	response, err := http.Get("https://weibo.com/ajax/side/hotSearch")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer response.Body.Close()

	return parseHotSearch(response.Body)
}

// parseHotSearch 解析微博热搜接口返回的 JSON，数据格式异常时返回空列表而不是 panic
func parseHotSearch(r io.Reader) []DataItem {
	var data []DataItem

	var result map[string]interface{}
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		fmt.Println(err)
		return data
	}

	dataField, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Printf("微博热搜返回数据异常，缺少 data 字段: %v\n", result)
		return data
	}

	dataJSON, ok := dataField["realtime"].([]interface{})
	if !ok {
		fmt.Printf("微博热搜返回数据异常，缺少 realtime 字段: %v\n", dataField)
		return data
	}

	jyzy := map[string]string{
		"电影": "影",
		"剧集": "剧",
		"综艺": "综",
		"音乐": "音",
	}

	for _, dataItem := range dataJSON {
		hot := ""
		dataMap, ok := dataItem.(map[string]interface{})
		if !ok {
			continue
		}

		// 如果是广告，则不添加
		if _, ok := dataMap["is_ad"]; ok {
			continue
		}

		if flagDesc, ok := dataMap["flag_desc"]; ok {
			if flagDescStr, ok := flagDesc.(string); ok {
				if hotValue, ok := jyzy[flagDescStr]; ok {
					hot = hotValue
				}
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

		note, _ := dataMap["note"].(string)
		word, _ := dataMap["word"].(string)
		num, _ := dataMap["num"].(float64)

		dic := DataItem{
			Title: note,
			URL:   "https://s.weibo.com/weibo?q=%23" + word + "%23",
			Num:   int(num),
			Hot:   hot,
		}
		data = append(data, dic)
	}

	return data
}
