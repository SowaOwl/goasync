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

	body, err := sendRequest("GET", url, headers, nil)
	if err != nil {
		sendError(url, startTime, err.Error(), ch)
		return
	}

	sendSuccess(url, startTime, body, ch)
}

func PostAsync(url string, headers []string, data interface{}, wg *sync.WaitGroup, ch chan map[string]interface{}) {
	defer wg.Done()
	startTime := time.Now()

	body, err := sendRequest("POST", url, headers, data)
	if err != nil {
		sendError(url, startTime, err.Error(), ch)
		return
	}

	sendSuccess(url, startTime, body, ch)
}

func sendRequest(method string, url string, headers []string, data interface{}) ([]byte, error) {
	var postBody []byte
	var err error
	if data != nil {
		postBody, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	var req *http.Request
	if data != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(postBody))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	addHeaders(req, headers)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func sendSuccess(url string, startTime time.Time, body []byte, ch chan<- map[string]interface{}) {
	ch <- map[string]interface{}{
		"duration": time.Since(startTime) / time.Millisecond,
		"body":     jsonToMap(body),
		"url":      url,
	}
}

func sendError(url string, startTime time.Time, errorMsg string, ch chan<- map[string]interface{}) {
	ch <- map[string]interface{}{
		"url":      url,
		"duration": time.Since(startTime) / time.Millisecond,
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
