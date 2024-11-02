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

type Steam struct {
	UserAgent string
	Type      string
	Language  string
}

type sgame struct {
	Title       string
	ReleaseDate string
	Link        string
	Platform    string
	Category    string
}

func (s Steam) get() (string, error) {
	url := "https://store.steampowered.com/explore/upcoming/"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	if len(s.UserAgent) > 0 {
		req.Header.Set("User-Agent", s.UserAgent)
	}

	var l string

	if s.Language == "zh_CN" {
		l = "zh-CN,zh;q=0.9"
	} else {
		l = "en-US,en;q=0.9"
	}

	req.Header.Set("Accept-Language", l)

	resp, err := http.DefaultClient.Do(req)

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (s Steam) parse() ([]sgame, error) {
	var games []sgame

	content, err := s.get()

	if err != nil {
		return games, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		return games, err
	}

	id := fmt.Sprintf("#tab_%s_comingsoon_content", s.Type)

	doc.Find(id).Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			link := s.AttrOr("href", "")

			title := s.Find(".tab_item_name").Text()

			if len(title) == 0 {
				return
			}

			date := s.Find(".release_date").Text()

			var tags []string

			s.Find(".tab_item_top_tags").Find("span").Each(func(i int, s *goquery.Selection) {
				tags = append(tags, s.Text())
			})

			var platforms []string
			s.Find(".tab_item_details").Find(".platform_img").Each(func(i int, s *goquery.Selection) {
				if s.HasClass("platform_img win") {
					platforms = append(platforms, "Windows")
				}

				if s.HasClass("platform_img mac") {
					platforms = append(platforms, "macOS")
				}

				if s.HasClass("platform_img linux") {
					platforms = append(platforms, "Linux")
				}
			})

			game := sgame{
				Title:       title,
				ReleaseDate: date,
				Link:        link,
				Platform:    strings.Join(platforms, ","),
				Category:    strings.Join(tags, ""),
			}

			games = append(games, game)
		})
	})

	return games, nil
}

func (s Steam) getCalendar() (string, error) {
	games, err := s.parse()

	if err != nil {
		return "", err
	}

	cal := ics.NewCalendar()

	cal.SetMethod("PUBLISH")
	cal.SetRefreshInterval("P1D")
	cal.SetTimezoneId("UTC")
	cal.SetName("Upcoming Releases - Steam")
	cal.SetVersion("2.0")
	cal.SetCalscale("GREGORIAN")

	for _, game := range games {
		id, _ := uuid.NewUUID()

		var start time.Time

		if s.Language == "zh_CN" {
			d := game.ReleaseDate

			d = strings.ReplaceAll(d, " ", "")
			d = strings.ReplaceAll(d, "年", "-")
			d = strings.ReplaceAll(d, "月", "-")
			d = strings.ReplaceAll(d, "日", "")

			digits := strings.Split(d, "-")

			for i, digit := range digits {
				if len(digit) == 1 {
					digits[i] = "0" + digit
				}
			}

			d = strings.Join(digits, "-")

			s, err := time.Parse("2006-01-02", d)

			if err != nil {
				continue
			}

			start = s
		} else {
			t, err := time.Parse("2 Jan, 2006", game.ReleaseDate)
			if err != nil {
				continue
			}

			start = t
		}

		event := cal.AddEvent(id.String())
		event.SetSummary(game.Title)
		event.SetCreatedTime(time.Now())
		event.SetAllDayStartAt(start)
		event.SetURL(game.Link)

		if s.Language == "zh_CN" {
			event.SetDescription(fmt.Sprintf("分类: %s\n\n平台: %s", game.Category, game.Platform))
		} else {
			event.SetDescription(fmt.Sprintf("Categories: %s\n\nPlatforms: %s", game.Category, game.Platform))
		}
	}

	return cal.Serialize(), nil
}

func (s Steam) Run() (string, error) {
	return s.getCalendar()
}
