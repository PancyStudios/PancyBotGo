package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type CraiyonRequest struct {
	Prompt         string `json:"prompt"`
	Token          string `json:"token"`
	Model          string `json:"model"`
	NegativePrompt string `json:"negative_prompt"`
	Size           string `json:"size"`
}

func main() {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_120),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Println(err)
		return
	}

	reqBody := CraiyonRequest{
		Prompt:         "A beautiful red dog",
		Token:          "turbis",
		Model:          "auto",
		NegativePrompt: "",
		Size:           "256x256",
	}
	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(http.MethodPost, "https://api.craiyon.com/v4", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("origin", "https://www.craiyon.com")
	req.Header.Set("referer", "https://www.craiyon.com/")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response:", string(body)[:min(200, len(body))]) // Print first 200 chars
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
