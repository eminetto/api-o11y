package mysql

import (
	"context"
	"database/sql"
	"github.com/eminetto/api-o11y/feedbacks/feedback"
	"github.com/eminetto/api-o11y/internal/telemetry"
	"go.opentelemetry.io/otel/codes"
	"time"
)

// FeedbackMySQL mysql repo
type FeedbackMySQL struct {
	db        *sql.DB
	telemetry telemetry.Telemetry
}

// NewFeedbackMySQL create new repository
func NewUserMySQL(db *sql.DB, telemetry telemetry.Telemetry) *FeedbackMySQL {
	return &FeedbackMySQL{
		db:        db,
		telemetry: telemetry,
	}
}

// Store a feedback
func (r *FeedbackMySQL) Store(ctx context.Context, f *feedback.Feedback) error {
	ctx, span := r.telemetry.Start(ctx, "mysql")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	defer tx.Commit()
	stmt, err := tx.Prepare(`
		insert into feedback (id, email, title, body, created_at) 
		values(?,?,?,?,?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		f.ID,
		f.Email,
		f.Title,
		f.Body,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	return nil
}
