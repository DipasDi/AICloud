package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var (
	ConnPSQL *pgx.Conn
	ErrorDB  error
)

func InitDB() {
	ConnPSQL, ErrorDB = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if ErrorDB != nil {
		log.Fatal("failed to connect database: ", ErrorDB)
	}
}
