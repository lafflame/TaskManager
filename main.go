package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type task struct {
	id          int
	info        string
	status      bool
	createdTime time.Time
	updated_at  time.Time
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
	doneColor := color.New(color.FgGreen).SprintFunc()
	pendingColor := color.New(color.FgRed).SprintFunc()
	headerColor := color.New(color.FgCyan, color.Bold).SprintFunc()
	//footerColor := color.New(color.FgYellow).SprintFunc()

	// Создаем таблицу
	table := simpletable.New()

	// Устанавливаем заголовки
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: headerColor("ID")},
			{Align: simpletable.AlignCenter, Text: headerColor("Task")},
			{Align: simpletable.AlignCenter, Text: headerColor("Done?")},
			{Align: simpletable.AlignCenter, Text: headerColor("Created At")},
			{Align: simpletable.AlignCenter, Text: headerColor("Updated At")},
		},
	}

	// Получаем данные из БД
	rows, err := db.Query("SELECT id, info, status, created_at, updated_at FROM list ORDER BY id")
	if err != nil {
		log.Fatal("Ошибка при выполнении запроса:", err)
	}
	defer rows.Close()

	var remainTasks int

	// Обрабатываем результаты
	for rows.Next() {
		var (
			id        int
			info      string
			status    bool
			createdAt time.Time
			updatedAt sql.NullTime
		)

		if err := rows.Scan(&id, &info, &status, &createdAt, &updatedAt); err != nil {
			log.Fatal("Ошибка при сканировании данных:", err)
		}

		if !status {
			remainTasks++
		}

		statusStr := pendingColor("✗")
		if status {
			statusStr = doneColor("✓")
		}

		updatedAtStr := "—"
		if updatedAt.Valid {
			updatedAtStr = updatedAt.Time.Format("2006-01-02 15:04:05")
		}

		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{Text: headerColor(id)},
			{Text: info},
			{Align: simpletable.AlignCenter, Text: statusStr},
			{Text: createdAt.Format("2006-01-02 15:04:05")},
			{Text: updatedAtStr},
		})
	}

	// Проверяем ошибки после итерации
	if err = rows.Err(); err != nil {
		log.Fatal("Ошибка при чтении строк:", err)
	}

	// Добавляем футер с статистикой
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Span: 5, Align: simpletable.AlignCenter, Text: fmt.Sprintf(
				"Remain tasks: %d",
				remainTasks,
			)},
		},
	}

	// Выводим таблицу
	table.Println()
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
