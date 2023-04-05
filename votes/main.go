package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/eminetto/api-o11y/pkg/middleware"
	"github.com/eminetto/api-o11y/votes/vote"
	"github.com/eminetto/api-o11y/votes/vote/mysql"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"net/http"
	"os"
	"time"
)

func main() {
	// Logger
	logger := httplog.NewLogger("auth", httplog.Options{
		JSON: true,
	})
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
	defer db.Close()
	repo := mysql.NewVoteMySQL(db)

	vService := vote.NewService(repo)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(middleware.IsAuthenticated)
	r.Post("/v1/vote", storeVote(vService))

	http.Handle("/", r)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      http.DefaultServeMux,
	}
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func storeVote(vService vote.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		var v vote.Vote
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			oplog.Error().Msg(err.Error())
			return
		}
		v.Email = r.Context().Value("email").(string)
		var result struct {
			ID uuid.UUID `json:"id"`
		}
		result.ID, err = vService.Store(r.Context(), &v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			oplog.Error().Msg(err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			oplog.Error().Msg(err.Error())
			return
		}
		return
	}
}
