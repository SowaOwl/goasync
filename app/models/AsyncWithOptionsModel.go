package model

type AsyncWithOptionRequestData struct {
	Options AsyncOptions  `json:"options"`
	Data    []interface{} `json:"data"`
}

type AsyncOptions struct {
	Url     string   `json:"url"`
	Type    string   `json:"type"`
	Count   int      `json:"count"`
	Headers []string `json:"headers"`
}
