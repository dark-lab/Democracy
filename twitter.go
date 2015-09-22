package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/dark-lab/Democracy/shared/config"
	_ "time"
)

func GetTwitter(conf *config.Configuration) *anaconda.TwitterApi {

	anaconda.SetConsumerKey(conf.TwitterConsumerKey)
	anaconda.SetConsumerSecret(conf.TwitterConsumerSecret)
	api := anaconda.NewTwitterApi(conf.TwitterAccessToken, conf.TwitterAccessTokenSecret)
	//api.SetDelay(10 * time.Second)
	return api
}
