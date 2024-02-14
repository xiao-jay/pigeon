package taopiaopiao

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
	"pigeon/config"
	"pigeon/pkg/framework"
	"regexp"
	"strings"
)

const Name = "taopiaopiao"

type Taopiaopiao struct {
	Cron string `json:"cron"`
}

func (t Taopiaopiao) Name() string {
	return Name
}

func (t Taopiaopiao) SendMessage(msg interface{}, Msg chan config.Msg) error {
	movie_infos := msg.([]Movie)
	var total_msg string
	for _, movie := range movie_infos {
		total_msg += movie.Name + "," + movie.Score + "," + movie.LeadingActor + "\n\n"
	}

	msgs := config.Msg{
		Title:       Name,
		Description: total_msg,
		Channel:     9,
	}

	Msg <- msgs
	return nil
}

func (t Taopiaopiao) Run(Msg chan config.Msg, config config.Config, c *cron.Cron) error {
	_, err := c.AddFunc(t.Cron, func() {
		log.Printf("%s 开始执行任务", Name)
		msg := GetHotMovie()

		if err := t.SendMessage(msg, Msg); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func New(Arguments config.Arguments) framework.Plugin {
	log.Println("taopiaopiao plugin init", Arguments)
	return &Taopiaopiao{
		Cron: Arguments["cron"].(string),
	}
}

type Movie struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Score        string `json:"score"`
	LeadingActor string `json:"leadingActor"`
}

func GetHotMovie() []Movie {
	movies := make([]Movie, 0)
	url := "https://dianying.taobao.com/showList.htm?spm=a1z21.6646273.city.3.56e44890Op1oMB&city=330100"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		return nil
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP request failed:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return nil
	}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))

	hot_movies := doc.Find(".tab-content").Find(".tab-movie-list").Eq(0).Find(".movie-card-wrap")

	// soon_movies = doc('.tab-content').find('.tab-movie-list').eq(1).find('.movie-card-wrap')
	for i := 0; i < hot_movies.Length(); i++ {
		url, _ = hot_movies.Eq(i).Find(".movie-card").Attr("href")
		pattern := regexp.MustCompile(`showId=(\d+)`)
		//log.Println(hot_movies.Eq(i).Find(".movie-card-list").Text())
		re := regexp.MustCompile(`主演：(.*)`)
		movie_info := Movie{
			ID:           pattern.FindStringSubmatch(url)[1],
			Name:         hot_movies.Eq(i).Find(".movie-card-name").Find(".bt-l").Text(),
			Score:        hot_movies.Eq(i).Find(".movie-card-name").Find(".bt-r").Text(),
			LeadingActor: re.FindStringSubmatch(hot_movies.Eq(i).Find(".movie-card-list").Text())[1],
		}
		movies = append(movies, movie_info)
	}
	return movies
}
