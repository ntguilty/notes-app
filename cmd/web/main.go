package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"ntguilty.me/notes-app/pkg/models"
	"ntguilty.me/notes-app/pkg/models/mysql"
	"os"
	"time"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")


// Changed application struct with notes and users with *mysql.* to interface
// because we wanted to enable mocking database.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	notes         interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Note, error)
		Latest() ([]*models.Note, error)
	}
	templateCache map[string]*template.Template
	users         interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
}

func main() {
	// Add flag -addr to give http network address
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:d8NaffvJEqQVJS5XnxmY@/notesapp?parseTime=true", "MySQL data source name")
	// Yes, I know that it shouldn't be here.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		notes:         &mysql.NoteModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Create a custom server struct to set my own errorLog as a stuff for logging errors
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)

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
