package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/micrypt/go-plivo/plivo"
)

var (
	fromNumber string
	authID     string
	authToken  string

	mailAuthPrefix string

	client *plivo.Client

	mailPlusRegex = regexp.MustCompile(`^\s*([^\+@]*)(\+(?:[^@]*))?(@.*)\s*$`)
	toRegex       *regexp.Regexp
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

	fmt.Printf("Got: %+v\n", msg)
	clearBody := getClearBody(msg)
	fmt.Printf("Clear Body:", clearBody)

	mp := &plivo.MessageSendParams{}
	mp.Dst = toAddress
	mp.Src = fromNumber

	from := msg.From
	for _, h := range msg.Headers {
		if h.K == "From" {
			from = h.V
		}
	}

	mp.Text = "[" + from + "] " + msg.Subject + ":\n" + clearBody

	respBody, _, err := client.Message.Send(mp)
	if err != nil {
		log.Println("Plivo error:", err, respBody.Error)
	} else {
		log.Printf("Sent: %+v\n", respBody)
	}
}

func getClearBody(msg Message) string {
	for _, h := range msg.Headers {
		if h.K != "Content-Type" {
			continue
		}
		split := strings.Split(h.V, ";")
		if len(split) < 2 {
			continue
		}
		t := split[0]
		if !strings.HasPrefix(t, "multipart") {
			continue
		}
		boundary := ""
		for _, kv := range split[1:] {
			splitkv := strings.Split(kv, "=")
			if len(splitkv) == 2 && splitkv[0] == "boundary" {
				boundary = strings.Trim(splitkv[1], `"`)
				break
			}
		}
		if boundary == "" {
			continue
		}
		split = strings.Split(msg.Body, boundary)
		clearBody := ""
	parts:
		for _, part := range split {
			part = strings.TrimPrefix(part, "--")
			part = strings.TrimSuffix(part, "--")
			part = strings.TrimSpace(part)
			headers, body := splitHeaders(part)
			for _, h2 := range headers {
				if h2.K != "Content-Type" {
					continue
				}
				if strings.Contains(h2.V, "/plain") {
					clearBody = body
					break parts
				}
			}
		}
		if clearBody != "" {
			return strings.TrimSpace(clearBody)
		}
	}
	return msg.Body
}

func getClearAddress(addr string) string {
	match := mailPlusRegex.FindStringSubmatch(addr)
	if len(match) > 1 {
		return match[1] + match[3]
	}
	return addr
}
