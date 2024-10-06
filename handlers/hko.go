package handlers

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

import (
	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

type HKO struct {
	UserAgent string
	URL       string
}

type moonPhase struct {
	XMLName xml.Name `xml:"MoonPhase"`
	Moons   []moon   `xml:"MOON"`
}

type moon struct {
	C      int     `xml:"C,attr"`
	Phases []phase `xml:"PHASE"`
}

type phase struct {
	P  int     `xml:"P,attr"`
	Y  int     `xml:"Y"`
	M  int     `xml:"M"`
	D  int     `xml:"D"`
	Hm string  `xml:"hm"`
	JD float64 `xml:"JD"`
}

type event struct {
	Name string
	Time string
}

var phases = map[string]string{
	"0": "🌙新（朔）月",
	"1": "🌓上弦月",
	"2": "🌕满（望）月",
	"3": "🌗下弦月",
}

func (h HKO) get() (string, error) {
	req, err := http.NewRequest("GET", h.URL, nil)

	if err != nil {
		return "", err
	}

	if len(h.UserAgent) > 0 {
		req.Header.Set("User-Agent", h.UserAgent)
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

func (h HKO) parse() ([]event, error) {
	content, err := h.get()

	var events []event

	if err != nil {
		return events, err
	}

	var m moonPhase

	err = xml.Unmarshal([]byte(content), &m)
	if err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)

		return events, err
	}

	for _, moon := range m.Moons {
		for _, phase := range moon.Phases {
			e := event{
				Name: phases[fmt.Sprintf("%d", phase.P)],
				Time: fmt.Sprintf("%d-%02d-%02d %s:00", phase.Y, phase.M, phase.D, phase.Hm),
			}

			events = append(events, e)
		}

	}

	return events, nil
}

func (h HKO) getCalendar() (string, error) {
	events, err := h.parse()

	if err != nil {
		return "", err
	}

	cal := ics.NewCalendar()

	cal.SetMethod("PUBLISH")
	cal.SetRefreshInterval("P1D")
	cal.SetTimezoneId("Asia/Shanghai")
	cal.SetName("月相")
	cal.SetVersion("2.0")
	cal.SetCalscale("GREGORIAN")

	for _, e := range events {
		id, _ := uuid.NewUUID()

		location, err := time.LoadLocation("Asia/Shanghai")

		if err != nil {
			continue
		}

		start, err := time.ParseInLocation(time.DateTime, e.Time, location)

		if err != nil {
			continue
		}

		evt := cal.AddEvent(id.String())
		evt.SetSummary(e.Name)
		evt.SetCreatedTime(time.Now())
		evt.SetStartAt(start.Local())
	}

	return cal.Serialize(), nil
}

func (h HKO) Run() (string, error) {
	return h.getCalendar()
}
