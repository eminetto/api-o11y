package main

import (
	"database/sql"
	"encoding/json"
	"feedbacks/feedback"
	"feedbacks/feedback/mysql"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/eminetto/talk-microservices-go/pkg/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

// @todo get this variables from env or config
const (
	DB_USER     = "feedbacks_user"
	DB_PASSWORD = "feedbacks_pwd"
	DB_HOST     = "localhost"
	DB_DATABASE = "feedbacks_db"
	DB_PORT     = "3307"
)

func main() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_DATABASE)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Panic(err.Error())
	}
	defer db.Close()
	repo := mysql.NewUserMySQL(db)

	fService := feedback.NewService(repo)
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.NewLogger(),
	)
	r.Handle("/v1/feedback", n.With(
		negroni.HandlerFunc(middleware.IsAuthenticated()),
		negroni.Wrap(storeFeedback(fService)),
	)).Methods("POST", "OPTIONS")

	http.Handle("/", r)
	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":8082", //@TODO usar variável de ambiente
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func storeFeedback(fService feedback.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var f feedback.Feedback
		err := json.NewDecoder(r.Body).Decode(&f)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		f.Email = r.Header.Get("email")
		var result struct {
			ID uuid.UUID `json:"id"`
		}
		result.ID, err = fService.Store(r.Context(), &f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		return
	})
}
