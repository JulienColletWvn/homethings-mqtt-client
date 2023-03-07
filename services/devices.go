package services

import "net/http"

func getDevices() {
	resp, err := http.Get("http://example.com/")
}
