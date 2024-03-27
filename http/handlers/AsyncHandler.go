package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	model "github.com/SowaOwl/goasync.git/app/models"
	repositories "github.com/SowaOwl/goasync.git/app/repositories"
	validator "github.com/SowaOwl/goasync.git/http/validators"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func AsyncHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, "This route allowed only post method", http.StatusMethodNotAllowed)
		return
	}

	var requestDataList []model.AsyncRequestData
	if err := json.NewDecoder(r.Body).Decode(&requestDataList); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.ValidateAsyncData(requestDataList); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	ch := make(chan map[string]interface{}, len(requestDataList))
	var returnData []map[string]interface{}

	for _, requestData := range requestDataList {
		wg.Add(1)

		switch requestData.Type {
		case "get":
			go repositories.GetAsync(requestData.Url, requestData.Headers, &wg, ch)
		case "post":
			go repositories.PostAsync(requestData.Url, requestData.Headers, requestData.Data, &wg, ch)
		default:
			handleError(w, "Unsupported request type", http.StatusBadRequest)
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

	var requestData model.AsyncWithOptionRequestData

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.ValidateAsyncWithOptionsData(requestData); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	var returnData []map[string]interface{}

	if requestData.Options.Count != 0 {
		ch := make(chan map[string]interface{}, requestData.Options.Count)
		if requestData.Options.Type == "get" {
			for i := 0; i < requestData.Options.Count; i++ {
				wg.Add(1)

				go repositories.GetAsync(requestData.Options.Url, requestData.Options.Headers, &wg, ch)
			}
		} else if requestData.Options.Type == "post" {
			for i := 0; i < requestData.Options.Count; i++ {
				wg.Add(1)

				go repositories.PostAsync(requestData.Options.Url, requestData.Options.Headers, requestData.Data[0], &wg, ch)
			}
		}
		wg.Wait()
		close(ch)

		for response := range ch {
			returnData = append(returnData, response)
		}
	} else {
		ch := make(chan map[string]interface{}, len(requestData.Data))
		for _, data := range requestData.Data {
			wg.Add(1)

			go repositories.PostAsync(requestData.Options.Url, requestData.Options.Headers, data, &wg, ch)
		}

		wg.Wait()
		close(ch)

		for response := range ch {
			returnData = append(returnData, response)
		}
	}

	successResponse(w, "Requests completed successfully", http.StatusOK, returnData)
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
