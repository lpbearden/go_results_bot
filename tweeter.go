package main

import (
	"fmt"
	"log"
	s "strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	scraper "github.com/lpbearden/csgo_scraper"
)

var client = GetTwitterClient()

func main() {
	// sentinel value to check against
	lastMatch := ""
	duration := time.Duration(60) * time.Second

	for {
		fmt.Println("-------------- checking for the latest match")
		// go get the latest match
		match := scraper.GetMatch()

		//trim match winner and check if its an empty string (some results are that way)
		if len(s.TrimSpace(match.Winner)) != 0 {
			// check if the match is in the last match
			if match.MatchUrl != lastMatch {
				// double check by looking at the last tweet
				if !isLastTweet(match) {
					// set lastMatch to current match and send tweet
					lastMatch = match.MatchUrl
					sendTweet(match)
				}
			}
		} else {
			fmt.Println("empty match: ")
			fmt.Println(match)
		}
		fmt.Println("-------------- sleeping")
		time.Sleep(duration)
	}
}

func isLastTweet(m scraper.Match) bool {
	isLast := false
	lastMatch := ""
	homeTimelineParams := &twitter.HomeTimelineParams{
		Count:     1,
		TweetMode: "extended",
	}
	tweets, _, _ := client.Timelines.HomeTimeline(homeTimelineParams)

	if len(m.Maps) > 0 {
		lastMatch = fmt.Sprintf("%s beat %s %s-%s on %s", m.Winner, m.Loser, m.WinScore, m.LoseScore, s.Join(m.Maps, ", "))
	} else {
		lastMatch = fmt.Sprintf("%s beat %s %s-%s on %s", m.Winner, m.Loser, m.WinScore, m.LoseScore, scraper.CsMap[m.MapName])
	}

	if s.Contains(tweets[0].FullText, lastMatch) {
		isLast = true
		log.Print("current match is same as last match")
	}
	return isLast
}

func sendTweet(m scraper.Match) {
	tweetStr := ""
	if len(m.Maps) > 0 {
		tweetStr = fmt.Sprintf("%s beat %s %s-%s on %s \nhttps://hltv.org%s\n#%s", m.Winner, m.Loser, m.WinScore, m.LoseScore, s.Join(m.Maps, ", "), m.MatchUrl, s.Split(m.Event, " ")[0])
	} else {
		tweetStr = fmt.Sprintf("%s beat %s %s-%s on %s \nhttps://hltv.org%s\n#%s", m.Winner, m.Loser, m.WinScore, m.LoseScore, scraper.CsMap[m.MapName], m.MatchUrl, s.Split(m.Event, " ")[0])
	}

	tweet, _, _ := client.Statuses.Update(tweetStr, nil)
	log.Print("Posted Tweet: ", tweet.Text)
}
