package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/dark-lab/Democracy/shared/config"
	"github.com/dark-lab/Democracy/shared/utils"
	. "github.com/mattn/go-getopt"
	"github.com/op/go-logging"
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

	myTweets := make(map[string]timelinesTweets)
	api := GetTwitter(&conf)

	retweetRegex, _ := regexp.Compile(`^RT`) // detecting retweets

	for _, account := range conf.TwitterAccounts {
		myTweets[account] = GetTimelines(api, account, conf.FetchFrom)
	}

	log.Info("Analyzing && collecting data")

	for i := range myTweets {
		var retweets int
		var mymentions int
		var mentions struct {
			Name        string
			Indices     []int
			Screen_name string
			Id          int64
			Id_str      string
		}
		var myUniqueMentions map[int64]int
		myUniqueMentions = make(map[int64]int)
		fmt.Println("-== Account: " + i + " ==-")
		fmt.Println("\tTweets: " + strconv.Itoa(len(myTweets[i])))

		for _, t := range myTweets[i] {
			// detecting retweets
			if retweetRegex.MatchString(t.Text) == true {
				retweets++
			} else {
				//detecting mentions outside retweets
				for _, mentions = range t.Entities.User_mentions {

					mymentions++
					myUniqueMentions[mentions.Id]++

				}
			}
		}
		Followers := GetFollowers(api, i)
		Following := GetFollowing(api, i)
		var Corrispective []int64
		for _, i := range Followers {
			if utils.IntInSlice(i, Following) == true {
				Corrispective = append(Corrispective, i)
			}
		}

		fmt.Println("\tof wich, there are " + strconv.Itoa(retweets) + " retweets")
		fmt.Println("\tof wich, there are " + strconv.Itoa(len(myUniqueMentions)) + " unique mentions (not in retweets)")
		fmt.Println("\tof wich, there are " + strconv.Itoa(mymentions) + " total mentions (not in retweets)")
		fmt.Println("\tFollowers: " + strconv.Itoa(len(Followers)))
		fmt.Println("\tFollowing: " + strconv.Itoa(len(Following)))
		fmt.Println("\tFollowers && Following: " + strconv.Itoa(len(Corrispective)))

	}

}
