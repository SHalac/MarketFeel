package main

import (
	"fmt"
	"marketfeel/secrets"
	"marketfeel/db"
	"marketfeel/twitterapi"
	"marketfeel/azureapi"
	"log"
	"bufio"
	"os"
	"strings"
)

var (
	consumerKey = secrets.API_KEY
	consumerSecret = secrets.API_SECRET
)

var searchConfig = map[string]string{
	"lang": "en",
	//"result_type": "popular",
	"count": "28",
	"q": "-filter:retweets since:2018-06-10 ",
	"tweet_mode": "extended",
}



func main() {
	token, err := db.GetDbToken()
	if err != nil {
		encodedToken := twitterapi.EncodeToken(consumerKey,consumerSecret)
		bearerToken := twitterapi.GetBearer(encodedToken)
		err2 := db.AddTokenDb(bearerToken)
		if err2 != nil {
			log.Fatal("something went wrong")
		}
		token, _ = db.GetDbToken()
	}
	fmt.Println("Enter a stock name (using $ before): ")
	reader := bufio.NewReader(os.Stdin)
	stockname, _ := reader.ReadString('\n')
	stockname = strings.TrimSuffix(stockname,"\n")
	searchConfig["q"] = searchConfig["q"] + stockname
	var texts []string
	reqUrl := twitterapi.ParseParams(searchConfig)
	_,texts = twitterapi.SearchTweets(reqUrl,token)
	/*
	fmt.Println(len(tweets), " TOTAL TWEETS FOUND for ",reqUrl)
	for _, tweet := range tweets {
		fmt.Println(*tweet)
	}
	*/
	score := azureapi.EvalSentiment(texts)
	//ui_resp := fmt.Sprintf("Score for %s is %f",searchConfig["q"],score)
	ui_resp := fmt.Sprintf("Score for %s is %f",stockname,score)
	fmt.Println(ui_resp)
}


