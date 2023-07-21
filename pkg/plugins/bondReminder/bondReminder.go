package bondReminder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"pigeon/config"
	"pigeon/pkg/framework"
	"time"
)

const Name = "bond-reminder"

type BondReminder struct {
	Cron string `json:"cron"`
}

func New(Arguments config.Arguments) framework.Plugin {
	log.Println("weather plugin init", Arguments)
	return &BondReminder{
		Cron: Arguments["cron"].(string),
	}
}
func (b BondReminder) Name() string {
	return Name
}

// SendMessage send msg to channel
func (b BondReminder) SendMessage(bindMsg interface{}, Msg chan config.Msg) error {
	msgjson, err := json.Marshal(bindMsg)
	if err != nil {
		return err
	}
	msg := config.Msg{
		Title:       Name,
		Description: string(msgjson),
		Channel:     9,
	}

	Msg <- msg
	return nil
}

func (b BondReminder) Run(Msg chan config.Msg, config config.Config, c *cron.Cron) error {

	_, err := c.AddFunc(b.Cron, func() {
		log.Printf("%s 开始执行任务", Name)
		bonds := BondParser(BondData())
		msg := BondFilter(bonds)
		if len(msg) == 0 {
			log.Printf("今日 %v 无可转债", time.Now())
			return
		}
		if err := b.SendMessage(msg, Msg); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func BondParser(data []byte) any {
	var bond any
	json.NewDecoder(&JsonpWrapper{
		Underlying: bytes.NewBuffer(data),
		Prefix:     "_",
	}).Decode(&bond)

	return bond
}

type JsonpWrapper struct {
	Prefix     string
	Underlying io.Reader

	gotPrefix bool
}

func (jpw *JsonpWrapper) Read(b []byte) (int, error) {
	if jpw.gotPrefix {
		return jpw.Underlying.Read(b)
	}

	prefix := make([]byte, len(jpw.Prefix))
	n, err := io.ReadFull(jpw.Underlying, prefix)
	if err != nil {
		return n, err
	}

	if string(prefix) != jpw.Prefix {
		return n, fmt.Errorf("JSONP prefix mismatch: expected %q, got %q",
			jpw.Prefix, prefix)
	}

	char := make([]byte, 1)
	for char[0] != '(' {
		n, err = jpw.Underlying.Read(char)
		if n == 0 || err != nil {
			return n, err
		}
	}

	jpw.gotPrefix = true
	return jpw.Underlying.Read(b)
}

func BondData() []byte {
	client := http.Client{Timeout: 10 * time.Second}
	log.Println("Getting convertible bonds data")
	for {
		resp, err := client.Get("https://datacenter-web.eastmoney.com/api/data/v1/get?callback=_&sortColumns=PUBLIC_START_DATE&sortTypes=-1&pageNumber=1&quoteType=0&reportName=RPT_BOND_CB_LIST&columns=ALL&quoteColumns=f2~01~CONVERT_STOCK_CODE~CONVERT_STOCK_PRICE,f235~10~SECURITY_CODE~TRANSFER_PRICE,f236~10~SECURITY_CODE~TRANSFER_VALUE,f2~10~SECURITY_CODE~CURRENT_BOND_PRICE,f237~10~SECURITY_CODE~TRANSFER_PREMIUM_RATIO,f239~10~SECURITY_CODE~RESALE_TRIG_PRICE,f240~10~SECURITY_CODE~REDEEM_TRIG_PRICE,f23~01~CONVERT_STOCK_CODE~PBV_RATIO")
		if err != nil {
			log.Println("Retry due to network failure")
			continue
		}
		defer resp.Body.Close()

		b, _ := ioutil.ReadAll(resp.Body)
		return b
	}
}

func BondFilter(data any) string {
	type Bonds struct {
		Result struct {
			Data []struct {
				Name   string `json:"SECURITY_NAME_ABBR"`
				Code   string `json:"SECURITY_CODE"`
				Date   string `json:"VALUE_DATE"`
				Rating string `json:"RATING"`
			} `json:"data"`
		} `json:"result"`
	}
	var (
		message string = ""
		bonds   Bonds
	)
	json.Unmarshal(func(data any) []byte {
		b, _ := json.Marshal(data)
		return b
	}(data), &bonds)

	for _, v := range bonds.Result.Data {
		// 匹配今天
		if v.Date == time.Now().Format("2006-01-02")+" 00:00:00" {
			message += "·" + v.Name + "（" + v.Code + " / " + v.Rating + "）\n"
		}
		// 匹配明天
		if v.Date == time.Now().Add(time.Hour*24).Format("2006-01-02")+" 00:00:00" {
			message += "·" + v.Name + "（" + v.Code + " / " + v.Rating + " / 预约）\n"
		}
		// 匹配后天
		if v.Date == time.Now().Add(time.Hour*24*2).Format("2006-01-02")+" 00:00:00" {
			message += "·" + v.Name + "（" + v.Code + " / " + v.Rating + " / 预约）\n"
		}
	}

	if len(message) == 0 {
		message = "今天没有可转债供申购或预约"
	}
	return message
}
