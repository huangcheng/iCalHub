package handlers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

import (
	"github.com/PuerkitoBio/goquery"
	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

type IMDb struct {
	UserAgent string
	URL       string
}

type imovie struct {
	Title        string
	ReleaseDate  string
	Category     string
	Presentation string
	Link         string
}

func (i IMDb) get() (string, error) {
	req, err := http.NewRequest("GET", i.URL, nil)

	if err != nil {
		return "", err
	}

	if len(i.UserAgent) > 0 {
		req.Header.Set("User-Agent", i.UserAgent)
	}

	resp, err := http.DefaultClient.Do(req)

	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (i IMDb) parse() ([]imovie, error) {
	content, err := i.get()

	var movies []imovie

	if err != nil {
		return movies, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	doc.Find(".ipc-page-section.ipc-page-section--base article").Each(func(i int, s *goquery.Selection) {
		date := s.Find("h3.ipc-title__text").Text()

		s.Find(".ipc-metadata-list.ipc-metadata-list--base li").Each(func(i int, s *goquery.Selection) {
			anchor := s.Find("a")

			title := anchor.Text()

			regex, err := regexp.Compile(`\s+\(.*\)$`)

			if err != nil {
				return
			}

			title = regex.ReplaceAllString(title, "")
			title = strings.TrimSpace(title)

			link, _ := anchor.Attr("href")

			link = fmt.Sprintf("https://www.imdb.com%s", link)

			var categories []string
			var presentations []string

			s.Find("ul").Each(func(index int, s *goquery.Selection) {
				s.Find("li span").Each(func(i int, s *goquery.Selection) {
					text := s.Text()

					if index == 0 {
						categories = append(categories, text)
					} else {
						presentations = append(presentations, text)
					}
				})
			})

			if len(title) == 0 {
				return
			}

			movies = append(movies, imovie{
				Title:        title,
				ReleaseDate:  date,
				Category:     strings.Join(categories, " · "),
				Presentation: strings.Join(presentations, " · "),
				Link:         link,
			})
		})
	})

	return movies, nil
}

func (i IMDb) getCalendar() (string, error) {
	movies, err := i.parse()

	if err != nil {
		return "", err
	}

	println(movies)

	cal := ics.NewCalendar()

	cal.SetMethod("PUBLISH")
	cal.SetRefreshInterval("P1D")
	cal.SetTimezoneId("UTC")
	cal.SetName("Upcoming releases - IMDb")
	cal.SetVersion("2.0")
	cal.SetCalscale("GREGORIAN")

	for _, movie := range movies {
		id, _ := uuid.NewUUID()

		start, err := time.Parse("Jan 2, 2006", movie.ReleaseDate)

		if err != nil {
			continue
		}

		event := cal.AddEvent(id.String())
		event.SetSummary(movie.Title)
		event.SetDescription(fmt.Sprintf("Categories: %s\nPerformers: %s", movie.Category, movie.Presentation))
		event.SetURL(movie.Link)
		event.SetCreatedTime(time.Now())
		event.SetAllDayStartAt(start)
	}

	return cal.Serialize(), nil
}

func (i IMDb) Run() (string, error) {
	return i.getCalendar()
}
