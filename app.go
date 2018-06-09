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
)

var (
	consumerKey = secrets.API_KEY
	consumerSecret = secrets.API_SECRET
	bearerUrl = "https://api.twitter.com/oauth2/token"
)



func main() {
	encodedToken := getBearerCred(consumerKey,consumerSecret)
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
	fmt.Println(string(respBody))
}

/* function to obtain bearer token credentials
input (2): consumer_key, consumer_secret
output (1): token credential, base64 encoded
*/
func getBearerCred(key string, secret string) string {
	builder := strings.Builder{}
	builder.WriteString(url.QueryEscape(key))
	builder.WriteString(":")
	builder.WriteString(url.QueryEscape(secret))
	encodedToken := base64.StdEncoding.EncodeToString([]byte(builder.String()))
	return encodedToken
} 

