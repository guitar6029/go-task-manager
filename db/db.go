package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "tasks.db")
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
