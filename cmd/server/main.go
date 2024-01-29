package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"storage/internal/handler"
	"storage/internal/repository"
	"storage/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
)

func main() {
	// config

	cfg := mysql.Config{
		User:      "root",
		Passwd:    "localhost",
		Addr:      "localhost:3306",
		Net:       "tcp",
		DBName:    "my_db",
		ParseTime: true,
	}

	// open connection to db
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	router := chi.NewRouter()

	rp := repository.NewProductMysql(db)

	sv := service.NewProductDefault(rp)

	hd := handler.NewProductDefault(sv)

	router.Route("/api/v1/products", func(r chi.Router) {
		// - GET
		r.Get("/{id}", hd.GetByID())

	})

	err = http.ListenAndServe(":8080", router)

	if err != nil {
		fmt.Println(err)
		return
	}

}
