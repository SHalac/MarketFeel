package main

import (
	"fmt"
	"net/http"
	"marketfeel/secrets"
	"net/url"
	"strings"
	"encoding/base64"
	"log"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"reflect"
)

var (
	consumerKey = secrets.API_KEY
	consumerSecret = secrets.API_SECRET
	bearerUrl = "https://api.twitter.com/oauth2/token"
	tweetSearchUrl = "https://api.twitter.com/1.1/search/tweets.json?"
	bucketName = []byte("marketFeel")
)



func main() {
	token, err := getDbToken()
	if err != nil {
		encodedToken := encodeToken(consumerKey,consumerSecret)
		bearerToken := getBearer(encodedToken)
		err2 := addTokenDb(bearerToken)
		if err2 != nil {
			log.Fatal("something went wrong")
		}
		fmt.Println("added token to db")
		token, _ = getDbToken()
	} else {
		fmt.Println("token already in DB ")
	}
	var tweets []string
	tweets = searchTweets("world cup switzerland",token)
	fmt.Println(tweets)
}

/*
Function: Get bearer token from database
In: NONE
out: token (string) and error, error is not nil 
if token isn't found 
*/
func getDbToken() (string, error){ // the issue right now is returning byte
	db,err := bolt.Open("token.db",0600,nil)
	if err != nil {
		return "error", errors.New("Could not get token from db")
	}
	defer db.Close()
	var token string
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return errors.New("can't find main bucket")
		}
		token = string(bucket.Get([]byte("token"))[:])
		if token == "" {
			return errors.New("can't find token key")
		}
		return nil
		})
	if err != nil {
		return "error", errors.New("Could not get token from db")
	}
	return token, nil

}

func addTokenDb(token string) error {
	db,err := bolt.Open("token.db",0600,nil)
	if err != nil {
		return errors.New("Could not open db")
	}
	defer db.Close()
	err2 := db.Update(func(tx *bolt.Tx) error {
		bucket,err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte("token"),[]byte(token))
		return err
		})
	if err2 != nil {
		return errors.New("Could not add token to DB")
	}
	return nil
}

/*
Function to get bearer token from twitter
In: encoded token for request
Out: access token string
*/
func getBearer(encodedToken string) string {
	body := []byte("grant_type=client_credentials")
	req,err := http.NewRequest("POST",bearerUrl,bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s",encodedToken))
	req.Header.Add("Content-Type","application/x-www-form-urlencoded;charset=UTF-8")
	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var jsonResp interface{} // turns into map[string]interface{}
	err = json.Unmarshal(respBody, &jsonResp)
	if err != nil {
		log.Fatal(err.Error())
	}
	jsonResp2 := jsonResp.(map[string]interface{})
	if jsonResp2["token_type"] == "bearer" {
		return jsonResp2["access_token"].(string)
	} else {
		return "ERR"
	}
}


// could this be improved by returning pointer insteads
// map[string]interface{}
func twitterRequest(query string, token string) []byte {
	req,err := http.NewRequest("GET",query,nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",token))
	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		fmt.Println(err2)
		return nil
	}
	defer resp.Body.Close()
	respBody, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Println("bad body")
	}
	return respBody
}


/*
var jsonResp interface{} // turns into map[string]interface{}
	err = json.Unmarshal(respBody, &jsonResp)
	if err != nil {
		log.Fatal("wrong")
	}
	return jsonResp.(map[string]interface{})

*/

func searchTweets(queries string,token string) []string {
	type Tweets struct {
		Statuses [] struct {
			Text string `json:"full_text"`
			User struct {
				Name string `json:"screen_name"`
				} `json:"user"`
		} `json:"statuses"`
	}
	var m Tweets
	builder := strings.Builder{}
	builder.WriteString(tweetSearchUrl)
	builder.WriteString("q=")
	builder.WriteString(url.QueryEscape(queries))
	builder.WriteString("&count=5")
	builder.WriteString("&result_type=popular")
	builder.WriteString("&tweet_mode=extended")
	queryUrl := builder.String()
	respbody := twitterRequest(queryUrl,token)
	err := json.Unmarshal(respbody, &m)
	if err != nil {
		log.Fatal("wrong")
	}
	var tweetSlice []string
	for _,status := range m.Statuses {
		//fmt.Println(status.User.Name)
		//fmt.Println(status.Text)
		tweetSlice = append(tweetSlice,status.Text)
	}
	return tweetSlice

}

/* function to encode credentials to be used to get bearer
input (2): consumer_key, consumer_secret
output (1): encoded token, base64 encoded
*/
func encodeToken(key string, secret string) string {
	builder := strings.Builder{}
	builder.WriteString(url.QueryEscape(key))
	builder.WriteString(":")
	builder.WriteString(url.QueryEscape(secret))
	encodedToken := base64.StdEncoding.EncodeToString([]byte(builder.String()))
	return encodedToken
} 

