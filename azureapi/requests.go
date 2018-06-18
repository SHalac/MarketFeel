package azureapi

import (
	"strconv"
	"marketfeel/secrets"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"
	"bytes"
	"io/ioutil"
)

var (
	accesskey = secrets.AZURE_KEY1
	endpoint = "https://westus2.api.cognitive.microsoft.com/text/analytics/v2.0/sentiment"
	lang = "en"
)

/* 
Dcoument represents the JSON "blueprint"
of each document that is passed to the 
azure api
*/
type Document struct {
	Language string `json:language`
	Id string `json:id`
	Text string `json:text`
}

/*
Top level representation for array of docs
in JSON structure 
*/
type Body struct {
	Documents []Document
}

/*
Representation of sentiment API response
for each document
*/
type SentimentResp struct {
	Documents []struct {
		Score string `json:score`
		Id string `json:string`
		} `json:documents`
}

/*
This function constructs a string in JSON format
to pass into the body of the request for sentiment
analysis.
In: array of each tweet string
Out: JSON body in bytes
*/
func ConstructBody(texts []string) []byte{
	var reqbody = Body{}
	for idx, text := range texts {
		var doc = Document{Language:lang,Id:strconv.Itoa(idx),Text:text}
		reqbody.Documents = append(reqbody.Documents,doc)
	}
	bytebody, err := json.Marshal(reqbody)
	if err != nil {
		fmt.Println("body construction to json went wrong")
	}
	return bytebody
}

/*
This function sends the JSON docs to azure
for sentiment processing. Returns 
*/
func FetchSentiment(reqbody []byte) []byte{
	req,err := http.NewRequest("POST",endpoint,bytes.NewBuffer(reqbody))
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Add("Ocp-Apim-Subscription-Key",accesskey)
	req.Header.Add("Content-Type","application/json")
	req.Header.Add("Accept","application/json")
	var azureclient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err2 := azureclient.Do(req)
	if err2 != nil {
		fmt.Println(err2)
	}
	defer resp.Body.Close()
	respBody, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Println("bad body")
	}
	return respBody
}

func EvalSentiment(texts []string) int {
	bytebody := ConstructBody(texts)
	byteresp := FetchSentiment(bytebody)
	var azureresp SentimentResp
	err := json.Unmarshal(byteresp, &azureresp)
	if err != nil {
		log.Fatal("azure response unmarshal error!")
	}
	var scoretotal int
	for _,doc := range azureresp.Documents {
		score_int, err2 := strconv.Atoi(doc.Score)
		if err2 != nil {
			fmt.Println("doc socre conversion error in Eval Sentiment")
		} else {
			scoretotal += score_int
		}
	}// for loop
	return scoretotal/len(azureresp.Documents)
}


