package validator

import (
	"errors"
	"fmt"

	model "github.com/SowaOwl/goasync.git/app/models"
)

func ValidateAsyncData(dataList []model.AsyncRequestData) error {
	for i, data := range dataList {
		if data.Url == "" || data.Type == "" {
			return errors.New("fields 'url' and 'type' required to be filled" + ". Data index: " + fmt.Sprint(i+1))
		}
	}
	return nil
}

func ValidateAsyncWithOptionsData(data model.AsyncWithOptionRequestData) error {
	if data.Options.Url == "" || data.Options.Type == "" {
		return errors.New("fields 'url' and 'type' in 'options' required to be filled")
	}

	switch data.Options.Type {
	case "get":
		if data.Options.Count == 0 {
			return errors.New("the 'count' field cannot be null when 'type' is 'get'")
		}
	case "post":
		if data.Options.Count != 0 && len(data.Data) > 1 {
			return errors.New("when 'type' is 'post', 'count' must be 0 or the length of 'data' array must not be greater than 1")
		} else if data.Options.Count != 0 && len(data.Data) < 1 {
			return errors.New("when 'type' is 'post', and 'count' more 0 length of 'data' array must be 1")
		}
	default:
		return errors.New("unsupported request type")
	}

	return nil
}
