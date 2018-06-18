package azureapi

import (
	"strconv"
	"marketfeel/secrets"
	"encoding/json"
)

var (
	accesskey = secrets.AZURE_KEY1
	endpoint = "https://westus2.api.cognitive.microsoft.com/text/analytics/v2.0/sentiment"
	lang = "en"
)

type Document struct {
	language string
	id string
	text string
}

type Body struct {
	documents []Document
}

func ConstructBody(texts []string) string{
	var reqbody = &Body{}
	for idx, text := range texts {
		var doc = &Document{lang,strconv.Itoa(idx),text}
		reqbody.documents = append(reqbody.documents,doc)
	}
	return json.Marshal(reqbody).string()
}


