package main

import (
	"log"

	"github.com/resend/resend-go/v2"
)

func main() {

	apiKey := "re_7zJ1XWch_Df9MSNkcK4Pn27GydncqhEX9"

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{"hoangvule100@gmail.com"},
		Subject: "Hello World",
		Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
	}

	sent, err := client.Emails.Send(params)

	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Print(sent)
}
