package handlers

import (
	"encoding/json"
	"net/http"
)

// Добавляет ошибки в JSON и возвращает ответ в формате {"error":"ваш текст для ошибки"}.
func ErrReturn(err error, code int, w http.ResponseWriter) {
	result := make(map[string]string)
	result["error"] = err.Error()
	jsonResp, errJSON := json.Marshal(result)
	if errJSON != nil {
		http.Error(w, errJSON.Error(), http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(code)

	if _, errWrite := w.Write(jsonResp); errWrite != nil {
		http.Error(w, errWrite.Error(), http.StatusInternalServerError)
		return
	}
}
