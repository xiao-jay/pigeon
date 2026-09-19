package weibo_hot

import (
	"fmt"
	"strings"
	"testing"
)

const validPayload = `{
	"ok": 1,
	"data": {
		"realtime": [
			{"note": "普通热搜", "word": "普通热搜", "num": 1234567},
			{"note": "爆热搜", "word": "爆热搜", "num": 9999999, "is_boom": 1},
			{"note": "广告", "word": "广告", "num": 1, "is_ad": 1},
			{"note": "电影热搜", "word": "电影热搜", "num": 555, "flag_desc": "电影"},
			{"note": "新热搜", "word": "新热搜", "num": 888, "is_new": 1}
		]
	}
}`

// TestParseHotSearch 覆盖正常返回：广告被过滤，各种标签被识别
func TestParseHotSearch(t *testing.T) {
	items := parseHotSearch(strings.NewReader(validPayload))

	if len(items) != 4 {
		t.Fatalf("解析出 %d 条，期望 4 条（广告应被过滤）: %+v", len(items), items)
	}

	byTitle := make(map[string]DataItem, len(items))
	for _, item := range items {
		byTitle[item.Title] = item
	}

	if item := byTitle["普通热搜"]; item.Num != 1234567 || item.Hot != "" {
		t.Errorf("普通热搜 = %+v, 期望 num=1234567 hot=\"\"", item)
	}
	if item := byTitle["爆热搜"]; item.Hot != "爆" {
		t.Errorf("爆热搜的 Hot = %q, 期望 %q", item.Hot, "爆")
	}
	if item := byTitle["电影热搜"]; item.Hot != "影" {
		t.Errorf("电影热搜的 Hot = %q, 期望 %q", item.Hot, "影")
	}
	if item := byTitle["新热搜"]; item.Hot != "新" {
		t.Errorf("新热搜的 Hot = %q, 期望 %q", item.Hot, "新")
	}

	if _, ok := byTitle["广告"]; ok {
		t.Error("广告条目应该被过滤掉")
	}

	if item := byTitle["普通热搜"]; item.URL != "https://s.weibo.com/weibo?q=%23普通热搜%23" {
		t.Errorf("URL = %q, 不符合预期", item.URL)
	}
}

// TestParseHotSearchBadPayload 覆盖接口返回异常数据：不能 panic，返回空列表
func TestParseHotSearchBadPayload(t *testing.T) {
	cases := []struct {
		name    string
		payload string
	}{
		{name: "data 为 null", payload: `{"ok":0,"data":null}`},
		{name: "缺少 data 字段", payload: `{"ok":0,"msg":"need login"}`},
		{name: "realtime 为 null", payload: `{"ok":1,"data":{"realtime":null}}`},
		{name: "data 是字符串", payload: `{"ok":1,"data":"unexpected"}`},
		{name: "非法 JSON", payload: `<html>503 Service Unavailable</html>`},
		{name: "空响应", payload: ``},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			items := parseHotSearch(strings.NewReader(c.payload))
			if len(items) != 0 {
				t.Errorf("期望返回空列表，实际返回 %+v", items)
			}
		})
	}
}

// TestParseHotSearchMessyItems 覆盖单条数据缺字段的情况：不能 panic
func TestParseHotSearchMessyItems(t *testing.T) {
	payload := `{"ok":1,"data":{"realtime":[
		{"word":"只有word没有note和num"},
		{"note":"只有note"},
		null,
		"不是对象"
	]}}`

	items := parseHotSearch(strings.NewReader(payload))
	if len(items) != 2 {
		t.Fatalf("解析出 %d 条，期望 2 条: %+v", len(items), items)
	}
	if items[0].Title != "" || items[0].Num != 0 {
		t.Errorf("缺字段的条目应降级为空值，实际 %+v", items[0])
	}
	if items[1].Title != "只有note" {
		t.Errorf("Title = %q, 期望 %q", items[1].Title, "只有note")
	}
}

// TestGetData 真实调用微博接口的冒烟测试，只验证不 panic
func TestGetData(t *testing.T) {
	fmt.Println(GetData())
}
