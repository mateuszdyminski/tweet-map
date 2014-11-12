package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/golang/glog"
)

const (
	TweetUrl      = "http://search-twitter-proxy.herokuapp.com/search/tweets?q=%s"
	GoogleMapsUrl = "http://maps.googleapis.com/maps/api/geocode/json?sensor=false&address=%s"
)

type Tweets struct {
	Tweets []Tweet `json:"statuses"`
}

type Tweet struct {
	Text        string    `json:"text"`
	User        TweetUser `json:"user"`
	Coordinates Location  `json:"coordinates"`
}

type TweetUser struct {
	Location string `json:"location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type GmLocations struct {
	GmGeometry []GmGeometry `json:"results"`
}

type GmGeometry struct {
	GmLocation GmLocation `json:"geometry"`
}

type GmLocation struct {
	Location Location `json:"location"`
}

func Search(searchText string) <-chan []byte {
	return encodeJson(findLocationForTweets(parseJson(searchTweets(searchText))))
}

func searchTweets(searchText string) <-chan []byte {
	out := make(chan []byte)
	go func() {
		resp, err := http.Get(fmt.Sprintf(TweetUrl, searchText))
		if err != nil {
			glog.Errorf("%+v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		glog.Infof("Request for tweets with keyword: %s done", searchText)
		out <- body
		close(out)
	}()
	return out
}

func parseJson(in <-chan []byte) <-chan Tweets {
	out := make(chan Tweets)
	go func() {
		tweets := Tweets{}
		err := json.Unmarshal(<-in, &tweets)
		if err != nil {
			glog.Errorf("%+v", err)
		}
		glog.Infof("Unmarshal of tweets done")
		out <- tweets
		close(out)
	}()
	return out
}

func encodeJson(in <-chan Tweet) <-chan []byte {
	out := make(chan []byte)
	go func() {
		for tw := range in {
			data, err := json.Marshal(tw)
			if err != nil {
				glog.Errorf("%+v", err)
			}
			out <- data
		}
		close(out)
	}()
	return out
}

func findLocationForTweets(in <-chan Tweets) <-chan Tweet {
	out := make(chan Tweet)
	go func() {
		tweets := (<-in).Tweets
		done := make(chan bool)
		glog.Infof("Number of tweets: %d", len(tweets))

		for i := 0; i < len(tweets); i++ {
			tweetWithoutLocation := tweets[i]
			go func() {
				out <- findLocation(tweetWithoutLocation)
				done <- true
			}()
		}

		for i := 0; i < len(tweets); i++ {
			<-done
		}

		glog.Infof("All tweets with location sent")
		close(out)
	}()
	return out
}

func findLocation(tweet Tweet) Tweet {
	glog.Infof("Getting location for tweet: %+v", tweet)
	if tweet.Coordinates.Lat == 0 && tweet.Coordinates.Lng == 0 {
		if tweet.User.Location != "" {
			resp, err := http.Get(fmt.Sprintf(GoogleMapsUrl, tweet.User.Location))
			if err != nil {
				glog.Errorf("%+v", err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			glog.Infof("Request for location: %s done", tweet.User.Location)
			if err != nil {
				glog.Errorf("Request for location: %s failed", tweet.User.Location)
			}

			loc := GmLocations{}
			json.Unmarshal(body, &loc)

			glog.Infof("%+v", loc)

			if len(loc.GmGeometry) == 0 {
				tweet.Coordinates = getRandCoordinates()
			} else {
				tweet.Coordinates = loc.GmGeometry[0].GmLocation.Location
			}

			return tweet
		} else {
			tweet.Coordinates = getRandCoordinates()
			return tweet
		}
	} else {
		return tweet
	}
}

func getRandCoordinates() Location {
	glog.Infof("Generating random location!")
	return Location{(rand.Float64() * 180) - 90, (rand.Float64() * 360) - 180}
}
