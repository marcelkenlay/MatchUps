package request_utils

import (
	. "../utils"
	"encoding/json"
	"net/http"
)

func EncodeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := json.NewEncoder(w).Encode(data)
	CheckErr(err)
}