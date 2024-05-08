package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(MysqlConnectionString string) *Repository {
	db, err := sql.Open("mysql", MysqlConnectionString)
	if err != nil {
		log.Fatal("Failed to create repository:", err)
	}
	fmt.Println(MysqlConnectionString)
	log.Println("Open Db")

	if err := db.Ping(); err != nil {
		db.Close()
		panic(err)
	}

	return &Repository{
		Db: db,
	}
}
