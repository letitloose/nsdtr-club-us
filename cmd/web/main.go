package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/letitloose/nsdtr-club-us/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	members       *models.MemberModel
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "lougar:thewarrior@/nsdtrc?parseTime=true", "MySQL data source name")

	flag.Parse()

	app := &application{}
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()

	app.templateCache, err = newTemplateCache()
	if err != nil {
		app.errorLog.Fatal(err)
	}

	app.members = &models.MemberModel{DB: db}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Print("Starting server on ", *addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
