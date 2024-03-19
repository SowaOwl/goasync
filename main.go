package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func jsonToMap(jsonData []byte, url string) interface{} {
	var result interface{}
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		println(url)
		println(err)
		return string(jsonData)
	}
	return result
}

func sendGetRequest(url string, wg *sync.WaitGroup, ch chan map[string]interface{}) {
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

	println(resp.S)

	ch <- map[string]interface{}{
		"duration": time.Since(startTime) / 1e6,
		"body":     jsonToMap(body, url),
		"url":      url,
	}
}

func main() {
	urls := []string{
		// "https://cat-fact.herokuapp.com/facts",
		"https://cataas.com/cat",
		// "https://api.bigdatacloud.net/data/reverse-geocode-with-timezone?latitude=52.2250321&longitude=104.9633463&localityLanguage=ru",
		// "https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=43.2250321&longitude=76.9633463&localityLanguage=ru",
		// "https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=51.2250321&longitude=71.9633463&localityLanguage=ru",
	}

	var wg sync.WaitGroup
	var returnData []map[string]interface{}
	ch := make(chan map[string]interface{}, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go sendGetRequest(url, &wg, ch)
	}

	wg.Wait()
	close(ch)

	for response := range ch {
		returnData = append(returnData, response)
	}

	// respData, err := json.Marshal(returnData)
	// if err != nil {
	// 	fmt.Printf("ERROR")
	// }
	// fmt.Println(string(respData))
}
