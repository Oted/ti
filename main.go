package main

import (
  "net/http"
	"strings"
	"errors"
	"net/url"
	"bytes"
  "fmt"
  "os"
	"github.com/PuerkitoBio/goquery"
  "encoding/base64"
)

type body struct {
  PhoneNumber string `json:"phone_number"`
  Message string `json:"message"`
  Type string `json:"message_type"`
  Sender string `json:"sender_id"`
}

func main() {
	if (search()) {
		fmt.Println("found " + os.Getenv("TARGET"))

		send()
	}
}

func search() bool {
  req, err := http.NewRequest("GET", os.Getenv("URL"), nil)
	if err != nil {
		panic(err)
	}

  req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:72.0) Gecko/20100101 Firefox/72.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "sv-SE,sv;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	found := false

	doc.Find(os.Getenv("CSS_PATH")).Each(func(i int, s *goquery.Selection) {
		t := strings.ToLower(s.Text())
		if strings.Contains(t, os.Getenv("TARGET")) {
			found = true
		}
	})

	return found
}

func send() {
  customerId := ""
  apiKey := ""

  b := body{
    PhoneNumber:"4915175243414",
    Message: os.Getenv("MESSAGE"),
    Type:"ARN",
    Sender:"informer",
  }

  u := "https://rest-api.telesign.com/v1/messaging"

  data := url.Values{}
	data.Set("phone_number", b.PhoneNumber)
	data.Set("message", b.Message)
	data.Set("message_type", b.Type)
	data.Set("sender_id", "gnugglet")

	encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", customerId, apiKey)))

  req, err := http.NewRequest("POST", u, bytes.NewBufferString(data.Encode()))
  if err != nil {
      panic(err)
  }

  req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

  client := &http.Client{}
	resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }

  defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(errors.New("invalid statuscode"))
	}
}
