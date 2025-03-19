package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"time"
)

const (
	eTypeInform string = "Inform"
	eTypeError  string = "Error"

	coutMess            = 1000
	ErrorLogJson string = "ui/static/json/errorLog.json"
)

var TypeError = []string{eTypeInform, eTypeError}

type ErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
}

func logAndSaveError(errorLog error, filename string, typeError string) error {

	jsonError := []ErrorResponse{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err))
	}

	err = json.Unmarshal(data, &jsonError)
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err))
	}

	errResponse := ErrorResponse{
		Timestamp: time.Now(),
		Type:      typeError,
		Message:   errorLog.Error(),
	}

	jsonError = append(jsonError, errResponse)

	if len(jsonError) > coutMess {
		jsonError = jsonError[len(jsonError)-coutMess:]
	}

	rawDataOut, err := json.MarshalIndent(&jsonError, "", "  ")
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err))
	}

	err = ioutil.WriteFile(filename, rawDataOut, 0)
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err))
	}

	return nil
}

func SetLogError(errLog error) {
	_ = logAndSaveError(errLog, ErrorLogJson, eTypeError)

}

func SetLogInform(infLog error) {
	_ = logAndSaveError(infLog, ErrorLogJson, eTypeInform)
}

func LoadErrorFromJSONFile() ([]ErrorResponse, error) {
	val := []ErrorResponse{}
	data, err := ioutil.ReadFile(ErrorLogJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &val)
	if err != nil {
		return nil, err
	}

	// Отсортировать записи в порядке возрастания Timestamp
	//sort.Slice(val, func(i, j int) bool {
	//	return val[i].Timestamp.Before(val[j].Timestamp)
	//})

	// Отсортировать записи в порядке убывания Timestamp
	sort.Slice(val, func(i, j int) bool {
		return val[i].Timestamp.After(val[j].Timestamp)
	})

	return val, nil
}
