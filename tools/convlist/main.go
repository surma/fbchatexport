package main

import (
	"flag"
	"fmt"
	"log"

	fb "github.com/huandu/facebook"
	"github.com/surma-dump/fbchatexport/lib"
)

var (
	token = flag.String("token", "", "Facebook OAuth token (use explorer)")
)

func main() {
	flag.Parse()

	if *token == "" {
		log.Fatalf("-token has to be defined")
	}

	session := (&fb.App{}).Session(*token)
	res, err := session.Get("/me/conversations", fb.Params{
		"fields": "message_count,id,participants",
	})
	if err != nil {
		log.Fatalf("Error requesting conversation list: %s", err)
	}
	pres, err := res.Paging(session)
	if err != nil {
		log.Fatalf("Could not do paging: %s", err)
	}
	for {
		ress := pres.Data()
		for _, res := range ress {
			conv := lib.ConversationHeader{}
			if err := res.DecodeField("", &conv); err != nil {
				log.Printf("Error decoding conversation: %s", err)
			}
			fmt.Printf("%s (%6d)\n", conv.Id, conv.MessageCount)
			for _, p := range conv.Participants.Data {
				fmt.Printf("\t%s\n", p.Name)
			}
		}
		if !pres.HasNext() {
			break
		}
		_, err := pres.Next()
		if err != nil {
			log.Fatalf("Error getting next page: %s", err)
		}
	}
}