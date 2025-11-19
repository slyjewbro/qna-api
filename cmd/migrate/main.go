package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"qna-api/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	// Для тестов используем другую БД если указано
	dbName := "qna_db"
	if len(os.Args) > 1 && os.Args[1] == "test" {
		dbName = "qna_test"
	}

	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "password",
		DBName:     dbName,
	}

	if dbName == "qna_test" {
		cfg.DBPort = "5433" // Тестовая БД на другом порту
	}

	// Открываем соединение с БД
	db, err := sql.Open("postgres", cfg.GetDBConnectionString())
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	// Выполняем миграции
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS questions (
            id SERIAL PRIMARY KEY,
            text TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS answers (
            id SERIAL PRIMARY KEY,
            question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
            user_id VARCHAR(36) NOT NULL,
            text TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE INDEX IF NOT EXISTS idx_answers_question_id ON answers(question_id)`,
		`CREATE INDEX IF NOT EXISTS idx_answers_user_id ON answers(user_id)`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			log.Fatalf("Migration %d failed: %v", i+1, err)
		}
		fmt.Printf("Migration %d applied successfully\n", i+1)
	}

	fmt.Println("All migrations completed successfully for database:", dbName)
}
