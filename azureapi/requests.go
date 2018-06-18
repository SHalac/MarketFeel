package azureapi

import (
	"strconv"
	"marketfeel/secrets"
	"encoding/json"
	"fmt"
)

var (
	accesskey = secrets.AZURE_KEY1
	endpoint = "https://westus2.api.cognitive.microsoft.com/text/analytics/v2.0/sentiment"
	lang = "en"
)

type Document struct {
	Language string `json:language`
	Id string `json:id`
	Text string `json:text`
}

type Body struct {
	Documents []Document
}

func ConstructBody(texts []string) string{
	var reqbody = Body{}
	for idx, text := range texts {
		var doc = Document{lang,strconv.Itoa(idx),text}
		reqbody.Documents = append(reqbody.Documents,doc)
	}
	bytebody, err := json.Marshal(reqbody)
	if err != nil {
		fmt.Println("body construction to json went wrong")
	}
	return string(bytebody)
}


