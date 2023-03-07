package services

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type DataType struct {
	Id   int32  `json:"id"`
	Key  string `json:"key"`
	Unit string `json:"unit"`
}

func GetDataTypes(deviceId string) (dt []DataType, err error) {
	resp, err := http.Get(os.Getenv("API_BASE_URL") + "/devices/" + deviceId + "/data-types")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	dataTypeErr := json.Unmarshal(body, &dt)

	if dataTypeErr != nil {
		return nil, dataTypeErr
	}

	return dt, nil

}
