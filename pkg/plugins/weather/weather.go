package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	"log"
	"pigeon/config"
	"pigeon/pkg/framework"
)

const Name = "weather"
const (
	REQUESTUrl = "https://restapi.amap.com/v3/weather/weatherInfo"
)

var _ framework.Plugin = &Weather{}

type Weather struct {
	Key  string `json:"key"`
	City string `json:"city"`
	Cron string `json:"cron"`
}

type WeatherResponse struct {
	Status    string      `json:"status"`
	Count     string      `json:"count"`
	Info      string      `json:"info"`
	Infocode  string      `json:"infocode"`
	Forecasts []Forecasts `json:"forecasts"`
}

type SimepleWeather struct {
	City        string       `json:"city"`
	SimpleCasts []SimpleCast `json:"simplecasts"`
}

type SimpleCast struct {
	Dayweather   string `json:"dayweather"`
	Nightweather string `json:"nightweather"`
	Daytemp      string `json:"daytemp"`
	Nighttemp    string `json:"nighttemp"`
}
type Forecasts struct {
	City       string `json:"city"`
	Adcode     string `json:"adcode"`
	Province   string `json:"province"`
	Reporttime string `json:"reporttime"`
	Casts      []Cast `json:"casts"`
}

type Cast struct {
	Date            string `json:"date"`
	Week            string `json:"week"`
	Dayweather      string `json:"dayweather"`
	Nightweather    string `json:"nightweather"`
	Daytemp         string `json:"daytemp"`
	Nighttemp       string `json:"nighttemp"`
	Daywind         string `json:"daywind"`
	Nightwind       string `json:"nightwind"`
	Daypower        string `json:"daypower"`
	Nightpower      string `json:"nightpower"`
	Daytemp_float   string `json:"daytemp_float"`
	Nighttemp_float string `json:"nighttemp_float"`
}

func (w Weather) Name() string {
	return Name
}

func (w Weather) Run(Msg chan config.Msg, config config.Config, c *cron.Cron) error {
	_, err := c.AddFunc(w.Cron, func() {
		log.Printf("%s 开始执行任务", Name)
		wehtherInfo, err := w.GetWeatherInfo()
		if err != nil {
			log.Println(err)
			return
		}
		simpleWehtherInfo := w.GetSimpleWeatherInfo(*wehtherInfo)
		if err := w.SendMessage(simpleWehtherInfo, Msg); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func New(Arguments config.Arguments) framework.Plugin {
	log.Println("weather plugin init", Arguments)
	return &Weather{
		Key:  Arguments["key"].(string),
		City: Arguments["city"].(string),
		Cron: Arguments["cron"].(string),
	}
}

func (w Weather) GetSimpleWeatherInfo(weather WeatherResponse) string {
	var simpleWeathers string

	// 获取城市名称
	city := weather.Forecasts[0].City
	simpleWeathers += city + "： "
	// 遍历天气预报列表，生成简单天气信息
	for _, forecast := range weather.Forecasts[0].Casts {
		date := forecast.Date
		dayWeather := forecast.Dayweather
		nightWeather := forecast.Nightweather
		dayTemp := forecast.Daytemp
		nightTemp := forecast.Nighttemp

		// 将简单天气信息添加到列表中
		simpleWeathers += fmt.Sprintf("日期：%s，白天天气：%s，夜晚天气：%s，温度：%s-%s；\n\n",
			date, dayWeather, nightWeather, nightTemp, dayTemp)

	}

	return simpleWeathers
}

func (w Weather) GetWeatherInfo() (*WeatherResponse, error) {
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"key":        w.Key,
			"city":       w.City,
			"extensions": "all",
			"output":     "JSON",
		}).
		Get(REQUESTUrl)

	if err != nil {
		fmt.Printf("Failed to get weather information: %v", err)
		return nil, err
	}

	// 解析响应
	var weatherResponse WeatherResponse
	err = json.Unmarshal(resp.Body(), &weatherResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if weatherResponse.Status == "1" {
		// 打印每天的天气信息
		for _, forecast := range weatherResponse.Forecasts {
			fmt.Printf("%+v", forecast)
		}
		return &weatherResponse, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Failed to get weather information: %s", weatherResponse.Info))
	}
}

func (w Weather) SendMessage(weatherMsg interface{}, Msg chan config.Msg) error {
	msgjson, err := json.Marshal(weatherMsg)
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
