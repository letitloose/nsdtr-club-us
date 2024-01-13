package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/letitloose/nsdtr-club-us/internal/models"
	"github.com/letitloose/nsdtr-club-us/internal/services"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog         *log.Logger
	infoLog          *log.Logger
	members          *models.MemberModel
	userService      *services.UserService
	templateCache    map[string]*template.Template
	sessionManager   *scs.SessionManager
	useTemplateCache bool
	email            *services.Email
}

func main() {

	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "lougar:thewarrior@/nsdtrc?parseTime=true", "MySQL data source name")
	emailUser := flag.String("emailUser", "test@gmail.com", "user account to send emails from")
	emailPassword := flag.String("emailPassword", "not-real-password", "password to emailUser account")
	emailHost := flag.String("emailHost", "smtp.gmail.com", "password to emailUser account")
	useTemplateCache := flag.Bool("useTemplateCache", false, "When false, templates will render on each request.")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	var templateCache = map[string]*template.Template{}
	infoLog.Println("Using Template Cache:", *useTemplateCache)
	if *useTemplateCache {
		templateCache, err = newTemplateCache()
		if err != nil {
			errorLog.Fatal(err)
		}
	}

	email := &services.Email{
		Username: *emailUser,
		Password: *emailPassword,
		Host:     *emailHost,
	}

	members := &models.MemberModel{DB: db}
	users := &models.UserModel{DB: db}
	userService := &services.UserService{UserModel: users, Email: email}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:         errorLog,
		infoLog:          infoLog,
		members:          members,
		userService:      userService,
		templateCache:    templateCache,
		sessionManager:   sessionManager,
		useTemplateCache: *useTemplateCache,
		email:            email,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.infoLog.Print("Starting server on ", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
