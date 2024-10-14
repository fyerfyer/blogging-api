package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	Post *PostModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/blogdb?parseTime=true", "MySQL data source name")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	app := &Application{
		Post: &PostModel{DB: db},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.Routers(),
	}

	log.Println("Starting server on :4000")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error opening database:", err) // 添加日志
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Println("Error pinging database:", err) // 添加日志
		return nil, err
	}

	log.Println("Database connection successful") // 添加日志
	return db, nil
}
