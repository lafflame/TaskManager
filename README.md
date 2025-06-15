# Task Manager CLI Application

Простое консольное приложение для управления задачами с использованием PostgreSQL в качестве базы данных.

## Функциональность

- Добавление новых задач
- Изменение статуса задач (выполнено/не выполнено)
- Удаление задач
- Просмотр списка всех задач в удобном табличном формате
- Подсчет оставшихся невыполненных задач

## Установка и настройка

1. Убедитесь, что у вас установлен Go (версии 1.16 или выше)
2. Установите PostgreSQL и создайте базу данных для приложения
3. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/yourusername/task-manager-cli.git
   cd task-manager-cli
   ```
4. Установите зависимости:
   ```bash
   go mod download
   ```
5. Создайте файл `.env` в корне проекта со следующими переменными:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database_name
   DB_SSLMODE=disable
   DB_TABLE=list
   ```
6. Создайте таблицу в базе данных:
   ```sql
   CREATE TABLE list (
       id SERIAL PRIMARY KEY,
       info TEXT NOT NULL,
       status BOOLEAN DEFAULT FALSE,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP
   );
   ```

## Использование

Запустите приложение:
```bash
go run main.go
```

Доступные команды:
- `add` - добавить новую задачу
- `change status` - изменить статус задачи (по ID)
- `delete` - удалить задачу (по ID)
- `printing` - показать список всех задач

## Зависимости

- [github.com/alexeyco/simpletable](https://github.com/alexeyco/simpletable) - для красивого вывода таблиц
- [github.com/fatih/color](https://github.com/fatih/color) - для цветного вывода в консоли
- [github.com/lib/pq](https://github.com/lib/pq) - драйвер PostgreSQL для Go
- [github.com/joho/godotenv](https://github.com/joho/godotenv) - для работы с .env файлами
