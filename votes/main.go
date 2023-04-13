package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/eminetto/api-o11y/internal/middleware"
	"github.com/eminetto/api-o11y/internal/telemetry"
	"github.com/eminetto/api-o11y/votes/vote"
	"github.com/eminetto/api-o11y/votes/vote/mysql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	telemetrymiddleware "github.com/go-chi/telemetry"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"os"
	"time"
)

func main() {
	// Logger
	logger := httplog.NewLogger("votes", httplog.Options{
		JSON: true,
	})
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
	defer db.Close()

	ctx := context.Background()
	otel, err := telemetry.New(ctx, "votes")
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
	defer otel.Shutdown(ctx)

	repo := mysql.NewVoteMySQL(db, otel)

	vService := vote.NewService(repo, otel)

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(telemetrymiddleware.Collector(telemetrymiddleware.Config{
		AllowAny: true,
	}, []string{"/v1"})) // path prefix filters basically records generic http request metrics
	r.Use(middleware.IsAuthenticated(ctx, otel))
	r.Post("/v1/vote", storeVote(ctx, vService, otel))

	http.Handle("/", r)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      http.DefaultServeMux,
	}
	err = srv.ListenAndServe()
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
}

func storeVote(ctx context.Context, vService vote.UseCase, otel telemetry.Telemetry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		ctx, span := otel.Start(ctx, "store")
		defer span.End()
		var v vote.Vote
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			oplog.Error().Msg(err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
		v.Email = r.Context().Value("email").(string)
		var result struct {
			ID uuid.UUID `json:"id"`
		}
		result.ID, err = vService.Store(ctx, &v)
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
