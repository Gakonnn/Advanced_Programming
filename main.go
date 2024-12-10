package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Book struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `gorm:"not null" json:"name"`
	Rating    float32 `gorm:"default:0" json:"rating"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type RequestData struct {
	Message string `json:"message"`
}

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var gormDB *gorm.DB
var tmpl *template.Template

func initGormDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=gakon dbname=library port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
	}
	db.AutoMigrate(&Book{})
	return db
}

// Новый обработчик для поиска книг по имени
func searchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("name")
	if query == "" {
		http.Error(w, "Отсутствует параметр поиска", http.StatusBadRequest)
		return
	}

	var books []Book
	// Поиск книг, имя которых содержит query
	gormDB.Where("name ILIKE ?", "%"+query+"%").Find(&books)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func handleBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		var books []Book
		gormDB.Find(&books)
		json.NewEncoder(w).Encode(books)
	case "POST":
		var book Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, "Некорректный JSON", http.StatusBadRequest)
			return
		}
		gormDB.Create(&book)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	case "PUT":
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Некорректный ID", http.StatusBadRequest)
			return
		}

		var book Book
		err = json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, "Некорректный JSON", http.StatusBadRequest)
			return
		}

		if err := gormDB.Model(&Book{}).Where("id = ?", id).Updates(book).Error; err != nil {
			http.Error(w, "Ошибка обновления книги", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(book)

	case "DELETE":
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Некорректный ID", http.StatusBadRequest)
			return
		}

		if err := gormDB.Delete(&Book{}, id).Error; err != nil {
			http.Error(w, "Ошибка удаления книги", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResponseData{
			Status:  "success",
			Message: "Книга удалена",
		})

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func handleTemplate(w http.ResponseWriter, r *http.Request) {
	var books []Book
	gormDB.Find(&books)
	tmpl.Execute(w, books)
}
func handleBookByID(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из запроса
	idStr := r.URL.Path[len("/books/"):]

	// Преобразуем ID из строки в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	// Ищем книгу по ID
	var book Book
	result := gormDB.First(&book, id)
	if result.Error != nil {
		http.Error(w, "Книга не найдена", http.StatusNotFound)
		return
	}

	// Рендерим шаблон с данными о книге
	tmpl, err := template.ParseFiles("templates/book.html")
	if err != nil {
		http.Error(w, "Ошибка при парсинге шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, book)
	if err != nil {
		http.Error(w, "Ошибка при рендеринге шаблона", http.StatusInternalServerError)
		return
	}
}

func main() {
	gormDB = initGormDB()

	tmpl = template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", handleTemplate)
	http.HandleFunc("/books", handleBooks)
	http.HandleFunc("/books/", handleBookByID)

	fmt.Println("Сервер запущен на порту 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
