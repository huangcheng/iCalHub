package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

import (
	"github.com/PuerkitoBio/goquery"
	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

type Douban struct {
	UserAgent string
}

type movie struct {
	Title       string
	ReleaseDate string
	Link        string
	Region      string
	Category    string
}

func (d Douban) get() (string, error) {
	url := "https://movie.douban.com/coming"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	if len(d.UserAgent) > 0 {
		req.Header.Set("User-Agent", d.UserAgent)
	}

	resp, err := http.DefaultClient.Do(req)

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (d Douban) parse() ([]movie, error) {
	var movies []movie

	content, err := d.get()

	if err != nil {
		return movies, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		return movies, err
	}

	doc.Find(".coming_list tbody").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			tds := s.Find("td")

			date := strings.TrimSpace(tds.Eq(0).Text())

			if strings.Index(date, "日") == -1 {
				return
			}

			if strings.Index(date, "年") == -1 {
				now := time.Now()
				year := now.Year()

				date = fmt.Sprintf("%d年%s", year, date)
			}

			releaseDate := strings.Replace(date, "年", "-", -1)
			releaseDate = strings.Replace(releaseDate, "月", "-", -1)
			releaseDate = strings.Replace(releaseDate, "日", "", -1)

			releaseDate = fmt.Sprintf("%s 00:00:00", releaseDate)

			anchor := tds.Eq(1).Find("a")
			title := anchor.Text()
			link, _ := anchor.Attr("href")

			category := tds.Eq(2).Text()
			region := tds.Eq(3).Text()

			movie := movie{
				Title:       strings.TrimSpace(title),
				ReleaseDate: strings.TrimSpace(releaseDate),
				Link:        strings.TrimSpace(link),
				Region:      strings.TrimSpace(region),
				Category:    strings.TrimSpace(category),
			}

			movies = append(movies, movie)
		})
	})

	return movies, nil
}

func (d Douban) getCalendar() (string, error) {
	movies, err := d.parse()

	if err != nil {
		return "", err
	}

	cal := ics.NewCalendar()

	cal.SetMethod("PUBLISH")
	cal.SetRefreshInterval("P1D")
	cal.SetTimezoneId("Asia/Shanghai")
	cal.SetName("即将上映电影")
	cal.SetVersion("2.0")
	cal.SetCalscale("GREGORIAN")

	for _, movie := range movies {
		id, _ := uuid.NewUUID()

		start, err := time.Parse(time.DateTime, movie.ReleaseDate)

		if err != nil {
			continue
		}

		event := cal.AddEvent(id.String())
		event.SetSummary(movie.Title)
		event.SetCreatedTime(time.Now())
		event.SetAllDayStartAt(start)
		event.SetURL(movie.Link)
		event.SetDescription(fmt.Sprintf("类型：%s\n地区：%s", movie.Category, movie.Region))
	}

	return cal.Serialize(), nil
}

func (d Douban) Run() (string, error) {
	return d.getCalendar()
}
