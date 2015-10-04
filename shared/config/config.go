package config

import (
	"encoding/json"
	"github.com/op/go-logging"
	"os"
	"strconv"
)

var log = logging.MustGetLogger("democracy")

type Configuration struct {
	TwitterConsumerKey       string   `json:"twitter_consumer_key"`
	TwitterConsumerSecret    string   `json:"twitter_consumer_secret"`
	TwitterAccessToken       string   `json:"twitter_access_token"`
	TwitterAccessTokenSecret string   `json:"twitter_access_token_secret"`
	TwitterAccounts          []string `json:"twitter_accounts"`
	Date                     string   `json:"fetch_from"`
	FetchFrom                int64
	FetchFollow              bool `json:"fetch_follow"`
}

func LoadConfig(f string) (Configuration, error) {

	jsonFile, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	var conf Configuration
	conf.FetchFollow = false
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatal("Couldn't parse json file")
	}
	if conf.TwitterConsumerKey == "" {
		log.Fatal("You need to specify 'twitter_consumer_key'")
	}
	if conf.TwitterConsumerSecret == "" {
		log.Fatal("You need to specify 'twitter_consumer_secret'")
	}
	if conf.TwitterAccessToken == "" {
		log.Fatal("You need to specify 'twitter_access_token'")
	}
	if conf.TwitterAccessTokenSecret == "" {
		log.Fatal("You need to specify 'twitter_access_token_secret'")
	}
	if len(conf.TwitterAccounts) == 0 {
		log.Fatal("You need to specify at least one account in 'twitter_accounts'")
	}
	if conf.Date == "" {
		log.Fatal("You need to specify 'fetch_from'")
	} else {
		i, _ := strconv.ParseInt(conf.Date, 10, 64)
		conf.FetchFrom = i
	}

	log.Info("TwitterConsumerKey: %#v\n", conf.TwitterConsumerKey)
	log.Info("TwitterConsumerSecret: %#t\n", conf.TwitterConsumerSecret)

	log.Info("TwitterAccessToken: %#v\n", conf.TwitterAccessToken)
	log.Info("TwitterAccessTokenSecret: %#v\n", conf.TwitterAccessTokenSecret)
	log.Info("FetchFrom: %#t\n", conf.FetchFrom)

	for _, i := range conf.TwitterAccounts {
		log.Info("Account to inspect: %#v\n", i)
	}

	return conf, err
}
