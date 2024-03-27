package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	repositories "github.com/SowaOwl/goasync.git/repositories"
)

type AsyncRequestData struct {
	Url     string      `json:"url"`
	Type    string      `json:"type"`
	Headers []string    `json:"headers"`
	Data    interface{} `json:"data"`
}

type AsyncWithOptionRequestData struct {
	Options AsyncOptions  `json:"options"`
	Data    []interface{} `json:"data"`
}

type AsyncOptions struct {
	Url     string   `json:"url"`
	Type    string   `json:"post"`
	Headers []string `json:"headers"`
}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func validateAsyncData(dataList []AsyncRequestData) error {
	for i, data := range dataList {
		if data.Url == "" || data.Type == "" {
			return errors.New("fileds 'url' and 'type' reqired to fill" + ". Data index: " + fmt.Sprint(i+1))
		}
	}
	return nil
}

func AsyncHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, "This route allowed only post method", http.StatusMethodNotAllowed)
		return
	}

	var requestDataList []AsyncRequestData
	if err := json.NewDecoder(r.Body).Decode(&requestDataList); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	ch := make(chan map[string]interface{}, len(requestDataList))
	var returnData []map[string]interface{}

	if err := validateAsyncData(requestDataList); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, requestData := range requestDataList {
		wg.Add(1)

		switch requestData.Type {
		case "get":
			go repositories.GetAsync(requestData.Url, requestData.Headers, &wg, ch)
		case "post":
			go repositories.PostAsync(requestData.Url, requestData.Headers, requestData.Data, &wg, ch)
		default:
			handleError(w, "Неподдерживаемый тип запроса", http.StatusBadRequest)
			wg.Done()
		}
	}

	wg.Wait()
	close(ch)

	for response := range ch {
		returnData = append(returnData, response)
	}

	successResponse(w, "Requests completed successfully", http.StatusOK, returnData)
}

func AsyncWithOptionsHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, "This route allowed only post method", http.StatusMethodNotAllowed)
		return
	}
}

func handleError(w http.ResponseWriter, errMsg string, statusCode int) {
	errorResponse := Response{Status: false, Message: errMsg, Data: ""}
	responseJSON, _ := json.Marshal(errorResponse)

	log.Println(errMsg)
	http.Error(w, string(responseJSON), statusCode)
}

func successResponse(w http.ResponseWriter, msg string, statusCode int, data interface{}) {
	successResponse := Response{Status: true, Message: msg, Data: data}
	responseJSON, _ := json.Marshal(successResponse)

	w.WriteHeader(statusCode)
	w.Write(responseJSON)
}
