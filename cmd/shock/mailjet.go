package main

import (
	"log"
	"os"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

func notify(status string) string {
	mjpublic := os.Getenv("MAILJET_PUBLICKEY")
	mjsecret := os.Getenv("MAILJET_SECRETKEY")
	mjfrom := os.Getenv("MAILJET_FROM")

	//TODO store alertto with session
	alertto := os.Getenv("MAILJET_TO")

	// check for reserved flag that toggles notifications
	if mjsecret == "DISABLE_NOTIFICATION" {
		return ""
	}
	// check target mail destination
	if alertto == "" {
		return ""
	}

	var body string
	mj := mailjet.NewMailjetClient(mjpublic, mjsecret)

	if status == "" {
		//fail alert
		body = "shock API call failed"
	} else {
		//success alert
		body = "shock API status - " + status
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: mjfrom,
				Name:  "TWITCH-EXT-DEMO",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: alertto,
					Name:  alertto,
				},
			},
			Subject:  "twitch extension demo alert test.",
			TextPart: body,
			HTMLPart: "<div>**test**</div>" + body,
			CustomID: "TwitchExtDemoTest",
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mj.SendMailV31(&messages)
	if err != nil {
		log.Print("Mail failed")
		log.Print(err)
		return ""
	}

	log.Print(res)
	return status
}
