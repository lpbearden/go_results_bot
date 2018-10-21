package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	scraper "github.com/lpbearden/csgo_scraper"
	"log"
	s "strings"
	"time"
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

		// check if the match is in the last match
		if match.MatchUrl != lastMatch {
			// double check by looking at the last tweet
			if !isLastTweet(match) {
				// set lastMatch to current match and send tweet
				lastMatch = match.MatchUrl
				sendTweet(match)
			}
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
		tweetStr = fmt.Sprintf("%s beat %s %s-%s on %s \nhttps://hltv.org%s\n#%s", m.Winner, m.Loser, m.WinScore, m.LoseScore, s.Join(m.Maps, ", "), m.MatchUrl, m.Event)
	} else {
		tweetStr = fmt.Sprintf("%s beat %s %s-%s on %s \nhttps://hltv.org%s\n#%s", m.Winner, m.Loser, m.WinScore, m.LoseScore, scraper.CsMap[m.MapName], m.MatchUrl, m.Event)
	}

	tweet, _, _ := client.Statuses.Update(tweetStr, nil)
	log.Print("Posted Tweet: ", tweet.Text)
}
