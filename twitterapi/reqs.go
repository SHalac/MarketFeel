package twitterapi


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
	bearerUrl = "https://api.twitter.com/oauth2/token"
	tweetSearchUrl = "https://api.twitter.com/1.1/search/tweets.json?"
)

type Tweet struct {
	Text string
	Author string
	Retweets int 
	Favorites int
}

func (t Tweet) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Username: %s \n",t.Author))
	builder.WriteString(fmt.Sprintf("Retweets: %v \n",t.Retweets))
	builder.WriteString(fmt.Sprintf("Favorites: %v \n",t.Favorites))
	builder.WriteString(fmt.Sprintf("Text: %s \n",t.Text))
	return builder.String()
}

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


/*
General twitter API request function, only 
supports "GET" for now, might add further functionality 
if app requires new resources from APi
In (2): Request url, auth token
Out (1): byte slice that is the returned json
*/
func TwitterRequest(query string, token string) []byte {
	req,err := http.NewRequest("GET",query,nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",token))
	twitterclient := &http.Client{}
	resp, err2 := twitterclient.Do(req)
	if err2 != nil {
		fmt.Println(err2)
		return nil
	}
	defer resp.Body.Close()
	respBody, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Println("bad body")
		return nil
	}
	return respBody
}


/*
Function to parse search preferences into request url
In (1): Configuration map (holds search preferences)
Out (1): Url string ready for request
*/
func ParseParams(config map[string]string) string {
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
/*
Specialized function for searching tweets
In (2): request url, bearer token for auth
Out (1): Slice of tweet array
*/
func SearchTweets(queryUrl string,token string) ([]*Tweet, []string) {
	type Tweets struct {
		Statuses [] struct {
			Text string `json:"full_text"`
			User struct {
				Name string `json:"screen_name"`
				} `json:"user"`
			Favorites int `json:"favorite_count"`
			Retweets int `json:"retweet_count"`
		} `json:"statuses"`
	}
	var responsestruct Tweets
	respbody := TwitterRequest(queryUrl,token)
	err := json.Unmarshal(respbody, &responsestruct)
	if err != nil {
		log.Fatal("twitter unmarshal error!")
	}
	var tweetSlice []*Tweet
	var texts []string
	for _,status := range responsestruct.Statuses {
		tweetInfo := &Tweet{
			Text:status.Text,
			Author:status.User.Name,
			Retweets:status.Retweets,
			Favorites:status.Favorites,
		}
		tweetSlice = append(tweetSlice,tweetInfo)
		texts = append(texts,status.Text)
	}
	return tweetSlice,texts

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