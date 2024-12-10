package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestData struct {
	Message string `json:"message"`
}

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost || r.Method == http.MethodGet {
		var requestData RequestData
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil || requestData.Message == "" {
			response := ResponseData{
				Status:  "fail",
				Message: "Некорректное JSON-сообщение",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		fmt.Println("Получено сообщение:", requestData.Message)

		response := ResponseData{
			Status:  "success",
			Message: "Данные успешно приняты",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		response := ResponseData{
			Status:  "fail",
			Message: "Метод запроса не поддерживается",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	http.HandleFunc("/", handleRequest)

	fmt.Println("Сервер запущен на порту 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
