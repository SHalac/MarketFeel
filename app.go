package main

import (
	"fmt"
	"net/http"
	"marketfeel/secrets"
	"net/url"
	"strings"
	"encoding/base64"
)

var (
	consumerKey = secrets.API_KEY
	consumerSecret = secrets.API_SECRET
)



func main() {
	encodedToken := getBearerCred(consumerKey,consumerSecret)

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

