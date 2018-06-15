package main

import (
	"fmt"
	"marketfeel/secrets"
	"marketfeel/db"
	"marketfeel/utils"
	"log"
	_ "reflect"
)

var (
	consumerKey = secrets.API_KEY
	consumerSecret = secrets.API_SECRET
)

var searchConfig = map[string]string{
	"lang": "en",
	"result_type": "popular",
	"count": "18",
	"q": "donald trump -filter:retweets since:2018-06-01",
	"tweet_mode": "extended",
}



func main() {
	token, err := db.GetDbToken()
	if err != nil {
		encodedToken := utils.EncodeToken(consumerKey,consumerSecret)
		bearerToken := utils.GetBearer(encodedToken)
		err2 := db.AddTokenDb(bearerToken)
		if err2 != nil {
			log.Fatal("something went wrong")
		}
		fmt.Println("added token to db")
		token, _ = db.GetDbToken()
	} else {
		fmt.Println("token already in DB ")
	}
	var tweets []string
	reqUrl := utils.ParseConfig(searchConfig)
	tweets = utils.SearchTweets(reqUrl,token)
	fmt.Println(len(tweets), " TOTAL TWEETS FOUND for ",tweetQueries)
	for _, tweet := range tweets {
		//fmt.Println(idx)
		fmt.Println(tweet)
		//fmt.Println("\n")
	}
}


