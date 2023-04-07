package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/eminetto/api-o11y/feedbacks/feedback"
	"github.com/eminetto/api-o11y/feedbacks/feedback/mysql"
	"github.com/eminetto/api-o11y/internal/middleware"
	"github.com/eminetto/api-o11y/internal/telemetry"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"os"
	"time"
)

func main() {
	// Logger
	logger := httplog.NewLogger("feedbacks", httplog.Options{
		JSON: true,
	})
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
	defer db.Close()
	ctx := context.Background()
	otel, err := telemetry.New(ctx, "feedbacks")
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
	defer otel.Shutdown(ctx)
	repo := mysql.NewUserMySQL(db, otel)

	fService := feedback.NewService(repo, otel)

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.IsAuthenticated)
	r.Post("/v1/feedback", storeFeedback(ctx, fService, otel))

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

func storeFeedback(ctx context.Context, fService feedback.UseCase, otel telemetry.Telemetry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		ctx, span := otel.Start(ctx, "store")
		defer span.End()
		var f feedback.Feedback
		err := json.NewDecoder(r.Body).Decode(&f)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			oplog.Error().Msg(err.Error())
			return
		}
		f.Email = r.Context().Value("email").(string)
		var result struct {
			ID uuid.UUID `json:"id"`
		}
		result.ID, err = fService.Store(ctx, &f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			oplog.Error().Msg(err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			oplog.Error().Msg(err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
		return
	}
}
