package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"ntguilty.me/notes-app/pkg/models/mysql"
	"os"
	"time"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	notes    *mysql.NoteModel
	templateCache map[string]*template.Template
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

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session: session,
		notes:    &mysql.NoteModel{DB: db},
		templateCache: templateCache,
	}

	// Create a custom server struct to set my own errorLog as a stuff for logging errors
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
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
