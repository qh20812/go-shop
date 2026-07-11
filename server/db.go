package main

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // driver PostgreSQL (import để đăng ký, không gọi trực tiếp)
)

var db *sql.DB

func connectDB() error {
	// Ví dụ: postgres://goshop:matkhau@127.0.0.1:5432/goshop?sslmode=disable
	dsn := os.Getenv("DATABASE_URL")

	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	return db.Ping() // sql.Open chưa thật sự kết nối — Ping mới kiểm tra thật
}