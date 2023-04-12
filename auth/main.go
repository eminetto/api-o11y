package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/eminetto/api-o11y/auth/security"
	"github.com/eminetto/api-o11y/auth/user"
	"github.com/eminetto/api-o11y/auth/user/mysql"
	"github.com/eminetto/api-o11y/internal/telemetry"
	"github.com/go-chi/httplog"
	"github.com/go-chi/telemetry"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"os"
	"time"

	"context"
	_ "github.com/go-sql-driver/mysql"
)

var (
	AppMetrics = &MyAppMetrics{telemetry.NewNamespace("auth")}
)

type MyAppMetrics struct {
	*telemetry.Namespace
}

func (m *MyAppMetrics) RecordMyAppHit() {
	m.RecordHit("my_app_hit", nil)
}

func (m *MyAppMetrics) RecordAppGuage(value float64) {
	m.RecordGauge("my_app_gauge", nil, value)
}

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

	ctx := context.Background()
	otel, err := telemetry.New(ctx, "auth")
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
	defer otel.Shutdown(ctx)

	repo := mysql.NewUserMySQL(db, otel)
	uService := user.NewService(repo, otel)

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))

	r.Use(telemetry.Collector(telemetry.Config{
		AllowAny: true,
	}, []string{"/v1"})) // path prefix filters basically records generic http request metrics

	r.Post("/v1/auth", userAuth(ctx, uService, otel))
	r.Post("/v1/validate-token", validateToken(ctx, otel))

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

func userAuth(ctx context.Context, uService user.UseCase, otel telemetry.Telemetry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		ctx, span := otel.Start(ctx, "userAuth")
		defer span.End()
		var param struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			oplog.Error().Msg(err.Error())
			return
		}
		err = uService.ValidateUser(ctx, param.Email, param.Password)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			oplog.Error().Msg(err.Error())
			return
		}
		var result struct {
			Token string `json:"token"`
		}
		result.Token, err = security.NewToken(param.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			oplog.Error().Msg(err.Error())
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return
		}

		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			oplog.Error().Msg(err.Error())
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return
		}
		return
	}
}

func validateToken(ctx context.Context, otel telemetry.Telemetry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		_, span := otel.Start(ctx, "validateToken")
		defer span.End()
		var param struct {
			Token string `json:"token"`
		}
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			oplog.Error().Msg(err.Error())
			return
		}

		t, err := security.ParseToken(param.Token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			oplog.Error().Msg(err.Error())
			return
		}
		tData, err := security.GetClaims(t)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			oplog.Error().Msg(err.Error())
			return
		}
		var result struct {
			Email string `json:"email"`
		}
		result.Email = tData["email"].(string)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			oplog.Error().Msg(err.Error())
			return
		}
		return
	}
}
