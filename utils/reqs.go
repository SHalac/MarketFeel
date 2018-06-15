package utils


import (
	"fmt"
	"net/http"
	"log"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"strings"
	"encoding/base64"
)

var (
	tweetLimit = "26"
	bearerUrl = "https://api.twitter.com/oauth2/token"
	tweetSearchUrl = "https://api.twitter.com/1.1/search/tweets.json?"
)
/*
Function to get bearer token from twitter
In: encoded token for request
Out: access token string
*/
func GetBearer(encodedToken string) string {
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


// general twitter request function
func TwitterRequest(query string, token string) []byte {
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

func ParseConfig(config map[string]string) string {
	builder := strings.Builder{}
	builder.WriteString(tweetSearchUrl)
	builder.WriteString("q=")
	builder.WriteString(url.QueryEscape(config["q"]))
	for k,v := range config {
		if k == "q" || v == "" {
			continue
		}
		builder.WriteString("&")
		builder.WriteString(fmt.Sprintf("%s=%s",k,v))
	}
	return builder.String()
}


func SearchTweets(queryUrl string,token string) []string {
	type Tweets struct {
		Statuses [] struct {
			Text string `json:"full_text"`
			User struct {
				Name string `json:"screen_name"`
				} `json:"user"`
		} `json:"statuses"`
	}
	var m Tweets
	respbody := TwitterRequest(queryUrl,token)
	err := json.Unmarshal(respbody, &m)
	if err != nil {
		log.Fatal("wrong")
	}
	var tweetSlice []string
	for _,status := range m.Statuses {
		tweetSlice = append(tweetSlice,status.Text)
	}
	return tweetSlice

}

/* function to encode credentials to be used to get bearer
input (2): consumer_key, consumer_secret
output (1): encoded token, base64 encoded
*/
func EncodeToken(key string, secret string) string {
	builder := strings.Builder{}
	builder.WriteString(url.QueryEscape(key))
	builder.WriteString(":")
	builder.WriteString(url.QueryEscape(secret))
	encodedToken := base64.StdEncoding.EncodeToString([]byte(builder.String()))
	return encodedToken
} 