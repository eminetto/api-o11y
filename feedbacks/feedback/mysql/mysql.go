package mysql

import (
	"context"
	"database/sql"
	"github.com/eminetto/api-o11y/feedbacks/feedback"
	"time"
)

// FeedbackMySQL mysql repo
type FeedbackMySQL struct {
	db *sql.DB
}

// NewFeedbackMySQL create new repository
func NewUserMySQL(db *sql.DB) *FeedbackMySQL {
	return &FeedbackMySQL{
		db: db,
	}
}

// Store a feedback
func (r *FeedbackMySQL) Store(ctx context.Context, f *feedback.Feedback) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
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
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	return nil
}
