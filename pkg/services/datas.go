package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type Data struct {
	DataTypeID int32   `json:"data_type_id"`
	Value      float64 `json:"value"`
}

func CreateData(data Data) (err error) {
	body, err := json.Marshal(data)

	if err != nil {
		return err
	}

	_, postErr := http.Post(os.Getenv("API_BASE_URL")+"/datas", "application/json", bytes.NewBuffer(body))

	if postErr != nil {
		return err
	}

	return nil

}
