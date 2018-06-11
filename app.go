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
)

var (
	consumerKey = secrets.API_KEY
	consumerSecret = secrets.API_SECRET
	bearerUrl = "https://api.twitter.com/oauth2/token"
	bucketName = []byte("marketFeel")
)



func main() {
	token, err := dbToken()
	if err != nil {
		encodedToken := encodeToken(consumerKey,consumerSecret)
		bearerToken := getBearer(encodedToken)
		token = addTokenDb(bearerToken)
		fmt.Println("added token to db")
	}
	fmt.Println(token)
}

/*
Function: Get bearer token from database
In: NONE
out: token (string) and error, error is not nil 
if token isn't found 
*/
func dbToken() (s string, err error){
	db,err := bolt.Open("token.db",0600,&bolt.Options{Timeout: 1 * time.Second})
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
		token = bucket.Get([]byte("token"))
		if token == nil {
			return errors.New("can't find token key")
		}
		return nil
		})
	if err != nil {
		return "error", errors.New("Could not get token from db")
	}
	return token

}

func addTokenDb(token string) error {
	db,err := bolt.Open("token.db",0600,&bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return errors.New("Could not open db")
	}
	defer db.Close()
	err2 = db.Update(func(tx *bolt.Tx) error {
		bucket,err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		err = bucket.Put("token",token)
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
func getBearer(encodedToken) string{
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
	err := json.Unmarshal(respBody, &f)
	if jsonResp["token_type"] == "bearer" {
		return jsonResp["access_token"].(string)
	} else {
		return "ERR"
	}
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

