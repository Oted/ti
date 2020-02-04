package main

import (
  "net/http"
	"net/url"
	"io/ioutil"
	"bytes"
  "fmt"
  "regexp"
  "encoding/base64"
)

type body struct {
  PhoneNumber string `json:"phone_number"`
  Message string `json:"message"`
  Type string `json:"message_type"`
  Sender string `json:"sender_id"`
}

func main() {
	search()
}

func search() {
  resp, err := http.Get("https://www.axs.com/venues/101921/ericsson-globe-stockholm-live-stockholm-tickets")
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

  re := regexp.MustCompile(`headliner`)
	res := re.FindAllString(string(b), -1)
	fmt.Println(res)
}

func send() {
  customerId := "4628D531-83FB-4617-834A-0BE6878D2042"
  apiKey := "wlZCJ9X67ATwZR3xom7iRELr03AOhjrsPoDcuus65pUQ1ict56cOax9mnc+x+O0znuGIQYx25j8jkwl+VXFUlA=="

  b := body{
    PhoneNumber:"4915175243414",
    Message:"haiii",
    Type:"ARN",
    Sender:"Snip",
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

  fmt.Println("response Status:", resp.Status)
  fmt.Println("response Headers:", resp.Header)
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))

	if err != nil {
		panic(err)
	}
}
