package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "gakon"
	dbname   = "library"
)

var db *sql.DB

func initDB() {
	// Формируем строку подключения
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка проверки соединения с базой данных: %v\n", err)
	}

	fmt.Println("Успешное подключение к базе данных!")
}

func main() {
	initDB()
	defer db.Close()

	// Простая HTML-страница
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>E-Library</title>
			</head>
			<body>
				<h1>Добро пожаловать в E-Library!</h1>
				<button onclick="alert('Запрос отправлен')">Пример кнопки</button>
			</body>
			</html>
		`))
	})

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
