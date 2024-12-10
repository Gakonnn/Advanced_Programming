package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Book struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"not null"`
	Rating    float32 `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func initGormDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=gakon dbname=library port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных с использованием Gorm: %v\n", err)
	}

	err = db.AutoMigrate(&Book{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v\n", err)
	}

	log.Println("База данных успешно инициализирована!")
	return db
}

func main() {
	db := initGormDB()
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	book := Book{Name: "Война и мир", Rating: 4.5}
	result := db.Create(&book)
	if result.Error != nil {
		log.Fatalf("Ошибка добавления книги: %v\n", result.Error)
	}

	log.Println("Книга успешно добавлена:", book)
}
