package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"pigeon/config"
)

type FangTang struct {
}

func (f FangTang) SendMessage(msgs []config.Msg, sendKeys any) error {
	sendkeysStringList := sendKeys.([]string)
	for _, sendkey := range sendkeysStringList {
		client := &http.Client{}
		url := fmt.Sprintf("https://sctapi.ftqq.com/%s.send", sendkey)
		for _, msg := range msgs {
			jsondata, err := json.Marshal(msg)
			if err != nil {
				return err
			}
			req, err := http.NewRequest("POST", url, bytes.NewReader(jsondata))
			if err != nil {
				return err
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			if resp.StatusCode != 200 {
				return fmt.Errorf("发送消息失败，HTTP 状态码: %d", resp.StatusCode)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println(err)
				}
			}(resp.Body)
		}
	}
	return nil
}
