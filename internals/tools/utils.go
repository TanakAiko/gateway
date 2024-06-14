package tools

import (
	"bytes"
	"encoding/json"
	md "gateway/model"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, data any, status int) {
	w.WriteHeader(status)
	dataMarshaled, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error : Marshal data to send", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(dataMarshaled)
	if err != nil {
		http.Error(w, "Error : Writing the data to the response", http.StatusInternalServerError)
		return
	}
}

func SendRequest(w http.ResponseWriter, bodyData md.RequestBody, methode string, URL string) (*http.Response, error) {
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(methode, URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
