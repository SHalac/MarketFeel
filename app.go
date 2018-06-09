package main

import (
	"fmt"
	_ "net/http"
	"marketfeel/secrets"
)

var consumer_key = secrets.API_KEY
var consumer_secret = secrets.API_SECRET


func main() {
	fmt.Println(consumer_key)
/*	messages := make(chan string)
	go func() { messages <- "hello" 
	fmt.Println("done hello!")
	}()
    go func() { messages <- "ping" 
    fmt.Println("done ping!")
    }()
    msg := <-messages
    msg2 := <-messages
    fmt.Println(msg)
    fmt.Println(msg2)*/
    
}

