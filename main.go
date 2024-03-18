package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Location struct {
	City      string `json:"city"`
	Continent string `json:"continent"`
}

func compressJson(jsonData string) string {
	compressedData := &bytes.Buffer{}

	if err := json.Compact(compressedData, []byte(jsonData)); err != nil {
		panic(err)
	}

	return compressedData.String()
}
func jsonToMap(jsonData []byte) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		fmt.Print("Error")
	}
	return result
}

func sendPostRequest(url string, wg *sync.WaitGroup, ch chan map[string]interface{}) {
	defer wg.Done()

	startTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error sending request to %s: %v\n", url, err)
		ch <- map[string]interface{}{"url": url, "duration": time.Since(startTime), "error": err}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response from %s: %v\n", url, err)
		ch <- map[string]interface{}{"url": url, "duration": time.Since(startTime), "error": err}
		return
	}

	ch <- map[string]interface{}{
		"duration": time.Since(startTime),
		"body":     jsonToMap(body),
		"url":      url,
	}
}

func main() {
	urls := []string{
		// "https://cat-fact.herokuapp.com/facts",
		// "https://cataas2.com/cat",
		// "https://api.bigdatacloud.net/data/reverse-geocode-with-timezone?latitude=52.2250321&longitude=104.9633463&localityLanguage=ru",
		// "https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=43.2250321&longitude=76.9633463&localityLanguage=ru",
		"https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=51.2250321&longitude=71.9633463&localityLanguage=ru",
	}

	var wg sync.WaitGroup
	ch := make(chan map[string]interface{}, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go sendPostRequest(url, &wg, ch)
	}

	wg.Wait()
	close(ch)

	for response := range ch {
		respData, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("ERROR")
		}
		fmt.Println(string(respData))
	}
}
