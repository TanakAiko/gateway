package tools

import (
	"encoding/json"
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
