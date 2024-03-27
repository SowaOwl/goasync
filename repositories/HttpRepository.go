package repositories

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	client = &http.Client{}
)

func GetAsync(url string, headers []string, wg *sync.WaitGroup, ch chan map[string]interface{}) {
	defer wg.Done()

	startTime := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		sendError(url, startTime, err.Error(), ch)
		return
	}
	addHeaders(req, headers)

	resp, err := client.Do(req)
	if err != nil {
		sendError(url, startTime, err.Error(), ch)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sendError(url, startTime, err.Error(), ch)
		return
	}

	ch <- map[string]interface{}{
		"duration": time.Since(startTime) / 1e6,
		"body":     jsonToMap(body),
		"url":      url,
	}
}

func PostAsync(url string, headers []string, data interface{}, wg *sync.WaitGroup, ch chan map[string]interface{}) {
	defer wg.Done()

	startTime := time.Now()

	postBody, err := json.Marshal(data)
	if err != nil {
		sendError(url, startTime, err.Error(), ch)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		sendError(url, startTime, err.Error(), ch)
		return
	}
	addHeaders(req, headers)

	resp, err := client.Do(req)
	if err != nil {
		sendError(url, startTime, err.Error(), ch)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sendError(url, startTime, err.Error(), ch)
		return
	}

	ch <- map[string]interface{}{
		"duration": time.Since(startTime) / 1e6,
		"body":     jsonToMap(body),
		"url":      url,
	}
}

func sendError(url string, startTime time.Time, errorMsg string, ch chan<- map[string]interface{}) {
	ch <- map[string]interface{}{
		"url":      url,
		"duration": time.Since(startTime),
		"error":    errorMsg,
	}
}

func addHeaders(req *http.Request, headers []string) {
	for _, header := range headers {
		headerParts := strings.SplitN(header, ": ", 2)
		if len(headerParts) == 2 {
			req.Header.Set(headerParts[0], headerParts[1])
		}
	}
}

func jsonToMap(jsonData []byte) interface{} {
	var result interface{}
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		return string(jsonData)
	}
	return result
}
