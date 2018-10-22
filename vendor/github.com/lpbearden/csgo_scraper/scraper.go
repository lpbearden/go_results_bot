package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	s "strings"
)

var monthMap = map[string]string{
	"January":   "01",
	"February":  "02",
	"March":     "03",
	"April":     "04",
	"May":       "05",
	"June":      "06",
	"July":      "07",
	"August":    "08",
	"September": "09",
	"October":   "10",
	"November":  "11",
	"December":  "12",
}

var CsMap = map[string]string{
	"inf":  "Inferno",
	"d2":   "Dust 2",
	"mrg":  "Mirage",
	"ovp":  "Overpass",
	"nuke": "Nuke",
	"cch":  "Cache",
	"trn":  "Train",
	"bo3":  "Bo3",
	"bo5":  "Bo5",
}

type Match struct {
	Date      []string
	MatchUrl  string
	Winner    string
	Loser     string
	WinScore  string
	LoseScore string
	Event     string
	num       int
	Id        string
	MapName   string
	Maps      []string
}

func (m Match) String() string {
	if len(m.Maps) > 0 {
		return fmt.Sprintf("%s %s > %s %s :: %s :: %s", m.Winner, m.WinScore, m.LoseScore, m.Loser, CsMap[m.MapName], s.Join(m.Maps, ", "))
	} else {
		return fmt.Sprintf("%s %s > %s %s :: %s", m.Winner, m.WinScore, m.LoseScore, m.Loser, CsMap[m.MapName])
	}
}

func GetMatch() Match {
	m := scrapeLastMatch()
	return m
}

func scrapeLastMatch() Match {
	c := colly.NewCollector()
	index := 0
	match := Match{}

	// Find all matches
	c.OnHTML("div.results-sublist", func(e *colly.HTMLElement) {
		if index > 0 {
			return
		}

		matchDate := e.ChildText(".standard-headline")
		if matchDate == "" {
			return
		}

		parsedDate := parseDate(s.Split(matchDate, " "))

		e.ForEach("div.result-con", func(n int, el *colly.HTMLElement) {
			if n > 0 {
				return
			}
			match = Match{
				Date:      parsedDate,
				MatchUrl:  el.ChildAttr("a", "href"),
				Winner:    el.ChildText("div.team-won"),
				Loser:     el.ChildText("div.team:not(div.team-won)"),
				WinScore:  el.ChildText("span.score-won"),
				LoseScore: el.ChildText("span.score-lost"),
				Event:     el.ChildText("span.event-name"),
				num:       n,
				Id:        el.Attr("data-zonedgrouping-entry-unix"),
				MapName:   el.ChildText("div.map"),
			}
			if s.Contains(match.MapName, "bo") {
				match.Maps = getMaps(match.MatchUrl)
			}
		})
		index++
	})

	c.Visit("https://www.hltv.org/results?stars=1")
	return match
}

func scrapeAllMatches() []Match {
	c := colly.NewCollector()
	//detailsCollector := c.Clone()
	matches := make([]Match, 0)

	// Find all matches
	c.OnHTML("div.results-sublist", func(e *colly.HTMLElement) {
		matchDate := e.ChildText(".standard-headline")

		if matchDate == "" {
			return
		}
		parsedDate := parseDate(s.Split(matchDate, " "))

		e.ForEach("div.result-con", func(n int, el *colly.HTMLElement) {

			match := Match{
				Date:      parsedDate,
				MatchUrl:  el.ChildAttr("a", "href"),
				Winner:    el.ChildText("div.team-won"),
				Loser:     el.ChildText("div.team:not(div.team-won)"),
				WinScore:  el.ChildText("span.score-won"),
				LoseScore: el.ChildText("span.score-lost"),
				Event:     el.ChildText("span.event-name"),
				num:       n,
				Id:        el.Attr("data-zonedgrouping-entry-unix"),
				MapName:   el.ChildText("div.map"),
			}
			matches = append(matches, match)

			if s.Contains(match.MapName, "bo") {
				match.Maps = getMaps(match.MatchUrl)
			}
		})
	})

	c.Visit("https://www.hltv.org/results?stars=1")
	return matches
}

func parseDate(input []string) []string {
	// date format of [dd, mm, yyyy, monthName]
	date := make([]string, 4)

	r, _ := regexp.Compile("[0-9]")
	day := r.FindAllString(input[3], -1)
	if len(day) == 2 {
		date[0] = day[0] + day[1]
	} else {
		date[0] = "0" + day[0]
	}

	date[1] = monthMap[input[2]]
	date[2] = input[len(input)-1]
	date[3] = input[2]

	return date
}

func getMaps(url string) []string {
	var maps []string
	detailsCollector := colly.NewCollector()

	detailsCollector.OnHTML("div.mapname", func(e *colly.HTMLElement) {
		if e.Text != "" {
			maps = append(maps, e.Text)
		}
	})
	detailsCollector.Visit("https://www.hltv.org/" + url)
	return maps
}
