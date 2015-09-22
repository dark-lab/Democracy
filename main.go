package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/dark-lab/Democracy/shared/config"
	. "github.com/mattn/go-getopt"

	"github.com/op/go-logging"
	"net/url"
	"os"
	_ "regexp"
	"strconv"
)

var log = logging.MustGetLogger("democracy")
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortpkg}.%{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

func main() {
	var c int
	var configurationFile string
	OptErr = 0
	for {
		if c = Getopt("c:h"); c == EOF {
			break
		}
		switch c {
		case 'c':
			configurationFile = OptArg
		case 'h':
			println("usage: " + os.Args[0] + " [-c config.yaml -h]")
			os.Exit(1)
		}
	}
	if configurationFile == "" {
		panic("I can't work without a configuration file")
	}
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)
	log.Info("Loading config")
	conf, err := config.LoadConfig(configurationFile)
	if err != nil {
		panic(err)
	}

	api := GetTwitter(&conf)

	type timelinesTweets map[string]anaconda.Tweet

	myTweets := make(map[string]timelinesTweets)

	for _, account := range conf.TwitterAccounts {

		myTweets[account] = make(timelinesTweets)
		log.Info("Searching info for: %#v\n", account)

		var max_id int64

		searchresult, _ := api.GetUsersShow(account, nil)
		fmt.Println("URL:" + searchresult.URL)

		v := url.Values{}
		v.Set("user_id", searchresult.IdStr)
		v.Set("count", "1") //getting twitter first tweet
		timeline, _ := api.GetUserTimeline(v)
		max_id = timeline[0].Id // putting it as max_id
		time, _ := timeline[0].CreatedAtTime()

		//retweetRegex, _ := regexp.Compile(`^RT`)
		for time.Unix() >= conf.FetchFrom { //until we don't exceed our range of interest

			v = url.Values{}
			v.Set("user_id", searchresult.IdStr)
			v.Set("count", "200")
			v.Set("max_id", strconv.FormatInt(max_id, 10))
			timeline, _ := api.GetUserTimeline(v)
			for _, tweet := range timeline {
				//fmt.Println("Tweet time:" + tweet.CreatedAt + " Tweet: " + tweet.Text)
				if tweet.Id < max_id {
					max_id = tweet.Id
				}
				//	if retweetRegex.MatchString(tweet.Text) == false {
				time, _ = tweet.CreatedAtTime()

				myTweets[account][tweet.IdStr] = tweet

				//for _, mentions := range tweet.Entities.User_mentions {
				//	fmt.Println("HA MENZIONATO:" + mentions.Name)
				//}

				//	}
			}
		}

	}

	for i := range myTweets {
		fmt.Println("Account: " + i)
		fmt.Println("Tweets: " + strconv.Itoa(len(myTweets[i])))
	}

}
