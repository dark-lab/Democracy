package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/dark-lab/Democracy/shared/config"
	"github.com/dark-lab/Democracy/shared/utils"
	"github.com/gernest/nutz"
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

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)
	OptErr = 0
	for {
		if c = Getopt("g:c:h"); c == EOF {
			break
		}
		switch c {
		case 'g':
			configurationFile = OptArg
			GenerateData(configurationFile)
		case 'c':
			configurationFile = OptArg
			GatherData(configurationFile)
		case 'h':
			println("usage: " + os.Args[0] + " [ -r ] [-c config.json -h]")
			os.Exit(1)
		}
	}

}

func GenerateData(configurationFile string) {
	if configurationFile == "" {
		panic("I can't work without a configuration file")
	}

	log.Info("Loading config")
	conf, err := config.LoadConfig(configurationFile)
	if err != nil {
		panic(err)
	}
	//api := GetTwitter(&conf)
	db := nutz.NewStorage(configurationFile+".db", 0600, nil)
	mygraph := Graph{Nodes: []Node{}, Links: []Link{}}
	count := 0
	innercount := 0
	group := 0
	for _, account := range conf.TwitterAccounts {
		tweets := db.Get(account, "tweets")
		from := db.Get(account, "from")
		retweets := db.Get(account, "retweets")
		unique_mentions := db.Get(account, "unique_mentions")
		total_mentions := db.Get(account, "total_mentions")
		followers := db.Get(account, "followers")
		following := db.Get(account, "following")
		followers_followed := db.Get(account, "followers_followed")
		mentions_to_followed := db.Get(account, "mentions_to_followed")

		log.Info("Account: " + account)
		log.Info("from: " + string(from.Data))

		log.Info("Tweets: " + string(tweets.Data))

		log.Info("retweets: " + string(retweets.Data))
		log.Info("unique_mentions: " + string(unique_mentions.Data))

		log.Info("total_mentions: " + string(total_mentions.Data))
		log.Info("followers: " + string(followers.Data))
		log.Info("following: " + string(following.Data))
		log.Info("followers_followed: " + string(followers_followed.Data))
		log.Info("mentions_to_followed: " + string(mentions_to_followed.Data))
		myUniqueMentions := db.GetAll(account, "map_unique_mentions").DataList

		mygraph.Nodes = append(mygraph.Nodes, Node{Name: account, Group: group})
		innercount += count
		for k, v := range myUniqueMentions {
			innercount++

			//id, _ := strconv.ParseInt(k, 10, 64)
			//User, _ := api.GetUsersShowById(id, nil)

			//log.Info("[" + User.ScreenName + "]:" + string(v))
			// in name del nodo poi mettici User.ScreeName
			weight, _ := strconv.Atoi(string(v))
			mygraph.Nodes = append(mygraph.Nodes, Node{Name: string(k), Group: group})
			mygraph.Links = append(mygraph.Links, Link{Source: innercount, Target: count, Value: weight})

		}
		group++
		count++

	}
	fileJson, _ := json.MarshalIndent(mygraph, "", "  ")
	err = ioutil.WriteFile("Democracy.json", fileJson, 0644)
	if err != nil {
		log.Info("WriteFileJson ERROR: " + err.Error())
	}

}

func GatherData(configurationFile string) {
	if configurationFile == "" {
		panic("I can't work without a configuration file")
	}

	log.Info("Loading config")
	conf, err := config.LoadConfig(configurationFile)
	if err != nil {
		panic(err)
	}

	myTweets := make(map[string]timelinesTweets)
	api := GetTwitter(&conf)

	db := nutz.NewStorage(configurationFile+".db", 0600, nil)

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
		var MentionsWithCorrispective []int64
		for _, i := range Following {
			if _, ok := myUniqueMentions[i]; ok {
				MentionsWithCorrispective = append(MentionsWithCorrispective, i)
			}
			if utils.IntInSlice(i, Followers) == true {
				Corrispective = append(Corrispective, i)
			}
		}

		fmt.Println("\tof wich, there are " + strconv.Itoa(retweets) + " retweets")
		fmt.Println("\tof wich, there are " + strconv.Itoa(len(myUniqueMentions)) + " unique mentions (not in retweets)")
		fmt.Println("\tof wich, there are " + strconv.Itoa(mymentions) + " total mentions (not in retweets)")
		fmt.Println("\tFollowers: " + strconv.Itoa(len(Followers)))
		fmt.Println("\tFollowing: " + strconv.Itoa(len(Following)))
		fmt.Println("\tFollowers && Following: " + strconv.Itoa(len(Corrispective)))
		fmt.Println("\tBetween mentions, those are whom the user is following: " + strconv.Itoa(len(MentionsWithCorrispective)))

		//di := DemocracyIndex(len(myUniqueMentions), len(MentionsWithCorrispective), len(myTweets[i]), retweets)
		om := OutsideMentions(len(myUniqueMentions), len(MentionsWithCorrispective))
		apt := AnswerPeopleTax(len(myUniqueMentions), len(MentionsWithCorrispective), len(myTweets[i]), retweets)

		//	fmt.Println("\tDemocracy tax: " + FloatToString(di))
		fmt.Println("\tOutside of circle mentions: " + FloatToString(om))
		fmt.Println("\t of answering to external people: " + FloatToString(apt))

		db.Create(i, "from", []byte(conf.Date))
		db.Create(i, "tweets", []byte(strconv.Itoa(len(myTweets[i]))))
		db.Create(i, "retweets", []byte(strconv.Itoa(retweets)))
		db.Create(i, "unique_mentions", []byte(strconv.Itoa(len(myUniqueMentions))))
		db.Create(i, "total_mentions", []byte(strconv.Itoa(mymentions)))
		db.Create(i, "followers", []byte(strconv.Itoa(len(Followers))))
		db.Create(i, "following", []byte(strconv.Itoa(len(Following))))
		db.Create(i, "followers_followed", []byte(strconv.Itoa(len(Corrispective))))
		db.Create(i, "mentions_to_followed", []byte(strconv.Itoa(len(MentionsWithCorrispective))))

		for k, v := range myUniqueMentions {
			db.Create(i, strconv.FormatInt(k, 10), []byte(strconv.Itoa(v)), "map_unique_mentions")
		}

		// Visualize example:
		//http://bl.ocks.org/mbostock/4062045
		// Circle size is defined by it's radius (r) : .attr("r", 5)
		// TOOlTIP: http://bl.ocks.org/Caged/6476579
	}
}

func OutsideMentions(UniqueMentions int, MensionsToFollowing int) float32 {
	return (float32(UniqueMentions) - float32(MensionsToFollowing)) / float32(UniqueMentions)
}

func AnswerPeopleTax(UniqueMentions int, MensionsToFollowing int, Tweets int, Retweets int) float32 {
	return (float32(UniqueMentions) - float32(MensionsToFollowing)) / (float32(Tweets) - float32(Retweets))
}

func FloatToString(input_num float32) string {
	// to convert a float number to a string
	return strconv.FormatFloat(float64(input_num), 'f', 6, 64)
}
