package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type task struct {
	id          int
	info        string
	status      bool
	createdTime int
}

func main() {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Инициализация БД
	db, err = initDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось проверить соединение с БД: %v", err)
	}

	fmt.Println("Hello! Here's your tasks:")
	printing()

}

// Добавление новой задачи
func add() {

}

// Изменение статуса задачи
func changeStatus() {

}

// Удаление задачи
func delete() {

}

// Вывод всех задач на экран
func printing() {
	initDB()
	defer db.Close()

	print, err := db.Exec("insert * from list")
	if err != nil {
		fmt.Println("Не удалось записать данные в таблицу БД, попробуйте заново")
	}
	fmt.Println(print)
}

func initDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Проверяем, что все обязательные переменные установлены
	if host == "" || port == "" || user == "" || dbname == "" {
		return nil, fmt.Errorf("не все обязательные переменные окружения установлены")
	}

	// Преобразуем порт в число
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("неверный формат порта: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
		host, p, user, dbname, sslmode)

	return sql.Open("postgres", connStr)
}
