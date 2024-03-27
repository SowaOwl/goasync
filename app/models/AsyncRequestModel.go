package model

type AsyncRequestData struct {
	Url     string      `json:"url"`
	Type    string      `json:"type"`
	Headers []string    `json:"headers"`
	Data    interface{} `json:"data"`
}
