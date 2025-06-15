package main

import (
	"bufio"
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
	updated_at  int
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

	//TODO
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice := scanner.Text()
		switch choice {
		case "delete":
			delete()
		case "add":
			add()
		case "change status":
			changeStatus()
		case "printing":
			printing()
		default:
			fmt.Println("Неправильный ввод")
		}

	}

}

// Добавление новой задачи
func add() {
	fmt.Println("Введите информацию для добавления задачи:")
	table := os.Getenv("DB_TABLE")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	info := scanner.Text()
	if len(info) < 3 {
		fmt.Println("Введите корректную задачу")
		return
	}

	query := fmt.Sprintf("INSERT INTO %s(info) VALUES ($1)", table)
	_, err := db.Exec(query, info)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Задача успешно добавлена")
}

// Изменение статуса задачи
func changeStatus() {
	var id int
	fmt.Print("Введите ID задачи для изменения статуса: ")
	_, err := fmt.Scan(&id)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	// Обновляем статус на противоположный и обновляем время
	result, err := db.Exec(`
        UPDATE list 
        SET status = NOT status, 
            updated_at = CURRENT_TIMESTAMP 
        WHERE id = $1`,
		id)
	if err != nil {
		fmt.Println("Ошибка при обновлении статуса:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Printf("Задача с ID %d не найдена\n", id)
		return
	}

	fmt.Printf("Статус задачи с ID %d успешно изменен\n", id)
}

// Удаление задачи
func delete() {
	fmt.Println("Введите ID задачи для удаления:")
	table := os.Getenv("DB_TABLE")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	id := scanner.Text()
	delete := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
	_, err := db.Exec(delete, id)
	if err != nil {
		log.Fatal(err)
	}
}

// Вывод всех задач на экран
func printing() {
	print, err := db.Exec("SELECT id, info, status, created_at, updated_at FROM list")
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
