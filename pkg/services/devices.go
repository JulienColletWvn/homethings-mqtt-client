package services

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Device struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func GetDevices() (device []Device, err error) {
	resp, err := http.Get(os.Getenv("API_BASE_URL") + "/devices")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	d := []Device{}

	deviceErr := json.Unmarshal(body, &d)

	if deviceErr != nil {
		return nil, deviceErr
	}

	return d, nil

}
