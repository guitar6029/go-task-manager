package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "tasks.db")
	if err != nil {
		return nil, err
	}

	//create tasks table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		done BOOLEAN
	)`)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the sqlite db!")
	return db, nil
}
