package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/cfstras/textmail/Godeps/_workspace/src/github.com/micrypt/go-plivo/plivo"
)

var (
	fromNumber string
	authID     string
	authToken  string

	mailAuthPrefix string

	client *plivo.Client

	HeaderRegex = regexp.MustCompile(`^(([a-zA-Z\-0-9]+): .+|\s.+)$`)
	toRegex     *regexp.Regexp
)

func main() {
	fromNumber = os.Getenv("FROM_NUMBER")
	authID = os.Getenv("AUTH_ID")
	authToken = os.Getenv("AUTH_TOKEN")
	mailAuthPrefix = os.Getenv("MAIL_AUTH_PREFIX")
	toRegex = regexp.MustCompile(`^` + regexp.QuoteMeta(mailAuthPrefix) +
		`-(\+?[0-9]+)@([a-zA-Z0-9\.\-]+)$`)

	if authID == "" || authToken == "" || fromNumber == "" || mailAuthPrefix == "" {
		log.Fatalln("Please set AUTH_ID, AUTH_TOKEN, MAIL_AUTH_PREFIX and FROM_NUMBER")
	}

	client = plivo.NewClient(nil, authID, authToken)
	acc, _, err := client.Account.Get()
	if err != nil {
		log.Fatalf("AccountGet failed: %v", err)
	} else {
		log.Printf("Plivo Account: %+v\n", acc)
	}

	SetCallback(sendMail)
	go mainLoop()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
}

func sendMail(msg Message) {
	toMatch := toRegex.FindStringSubmatch(msg.To)
	if toMatch == nil {
		log.Println("Invalid address/unauthenticated:", msg.To)
		return
	}
	toAddress := toMatch[1]
	toAddress = strings.TrimPrefix(toAddress, "+")

	body := strings.Replace(msg.Body, "\r\n", "\n", -1)
	lines := strings.Split(body, "\n")
	body = ""
	header := true
	for _, l := range lines {
		if !header || !HeaderRegex.MatchString(l) {
			header = false
			if body == "" {
				body = l
			} else {
				body += "\n" + l
			}
		}
	}
	msg.Body = body

	fmt.Printf("Got: %+v\n", msg)

	mp := &plivo.MessageSendParams{}
	mp.Dst = toAddress
	mp.Src = fromNumber
	mp.Text = "[" + msg.From + "] " + msg.Subject + ":\n" + msg.Body

	respBody, _, err := client.Message.Send(mp)
	if err != nil {
		log.Println("Plivo error:", err, respBody.Error)
	} else {
		log.Printf("Sent: %+v\n", respBody)
	}
}
