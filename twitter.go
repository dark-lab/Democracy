package main

import (
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/dark-lab/Democracy/shared/config"
)

type timelinesTweets map[string]anaconda.Tweet

func GetTwitter(conf *config.Configuration) *anaconda.TwitterApi {

	anaconda.SetConsumerKey(conf.TwitterConsumerKey)
	anaconda.SetConsumerSecret(conf.TwitterConsumerSecret)
	api := anaconda.NewTwitterApi(conf.TwitterAccessToken, conf.TwitterAccessTokenSecret)
	return api
}

func GetFollowers(api *anaconda.TwitterApi, account string) []int64 {
	v := url.Values{}
	v.Set("screen_name", account)
	v.Set("count", "200")
	var User anaconda.User
	var Followers []int64
	pages := api.GetFollowersListAll(v)
	for page := range pages {
		//Print the current page of followers
		for _, User = range page.Followers {
			Followers = append(Followers, User.Id)
			log.Debug("Getting another follower", User.Id)
		}
	}
	return Followers
}

func GetFollowing(api *anaconda.TwitterApi, account string) []int64 {
	v := url.Values{}
	v.Set("screen_name", account)
	v.Set("count", "5000")
	var Following []int64
	var id int64
	pages := api.GetFriendsIdsAll(v)
	for page := range pages {
		//Print the current page of "Friends"
		for _, id = range page.Ids {
			Following = append(Following, id)
			log.Debug("Getting another Following", id)

		}
	}
	return Following
}

func GetTimelines(api *anaconda.TwitterApi, account string, since int64) timelinesTweets {
	myTweets := make(timelinesTweets)
	log.Info("Searching info for: %#v\n", account)

	var max_id int64

	searchresult, _ := api.GetUsersShow(account, nil)
	v := url.Values{}
	v.Set("user_id", searchresult.IdStr)
	v.Set("count", "1") //getting twitter first tweet
	timeline, _ := api.GetUserTimeline(v)
	max_id = timeline[0].Id // putting it as max_id
	time, _ := timeline[0].CreatedAtTime()

	for time.Unix() >= since { //until we don't exceed our range of interest

		v = url.Values{}
		v.Set("user_id", searchresult.IdStr)
		v.Set("count", "200")
		v.Set("max_id", strconv.FormatInt(max_id, 10))
		timeline, _ := api.GetUserTimeline(v)
		for _, tweet := range timeline {
			if tweet.Id < max_id {
				max_id = tweet.Id
			}
			time, _ = tweet.CreatedAtTime()
			myTweets[tweet.IdStr] = tweet
		}
	}
	return myTweets

}
