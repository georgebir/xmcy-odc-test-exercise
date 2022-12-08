package util

import (
	"encoding/json"
	"net/http"
)

func ResponseOK(response http.ResponseWriter, content interface{}) {
	response.Header().Set("ACCESS-CONTROL-ALLOW-ORIGIN", "*")
	response.Header().Set("ACCESS-CONTROL-ALLOW-HEADERS", "*")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if content != nil {
		if contentJson, err := json.Marshal(content); err == nil {
			response.Write(contentJson)
		} //marshal
	} //if content not nil
}

func ResponseInternalServerError(response http.ResponseWriter) {
	response.Header().Set("ACCESS-CONTROL-ALLOW-ORIGIN", "*")
	response.Header().Set("ACCESS-CONTROL-ALLOW-HEADERS", "*")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusInternalServerError)
	content := map[string]string{"error": "Internal server error."}
	if contentJson, err := json.Marshal(content); err == nil {
		response.Write(contentJson)
	} //if not err
}

func ResponseNotFound(response http.ResponseWriter, message string) {
	response.Header().Set("ACCESS-CONTROL-ALLOW-ORIGIN", "*")
	response.Header().Set("ACCESS-CONTROL-ALLOW-HEADERS", "*")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusNotFound)
	response.Write([]byte(message))
}

func ResponseBadRequest(response http.ResponseWriter, message string) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "*")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusBadRequest)
	response.Write([]byte(message))
}

func ResponseForbidden(response http.ResponseWriter) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "*")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusForbidden)
	response.Write([]byte("Forbidden."))
}

func ResponseUnprocessableEntity(response http.ResponseWriter, content *map[string]interface{}) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "*")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusUnprocessableEntity)
	contentJson, err := json.Marshal(*content)
	if err == nil {
		response.Write([]byte(contentJson))
	} //if not err
}

func ResponseUnauthorized(response http.ResponseWriter) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "*")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusUnauthorized)
	content := make(map[string]string)
	content["error"] = "Unauthorized."
	contentJson, err := json.Marshal(content)
	if err == nil {
		response.Write(contentJson)
	} //if not err
}

func ResponseNotAcceptable(response http.ResponseWriter) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "*")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusNotAcceptable)
	response.Write([]byte("Not Acceptable."))
}

func ResponseOptions(response http.ResponseWriter) {
	response.Header().Set("ACCESS-CONTROL-ALLOW-ORIGIN", "*")
	response.Header().Set("ACCESS-CONTROL-ALLOW-HEADERS", "*")
	response.Header().Set("ACCESS-CONTROL-ALLOW-METHODS", "GET, POST, PUT, OPTIONS")
	response.WriteHeader(http.StatusOK)
}
