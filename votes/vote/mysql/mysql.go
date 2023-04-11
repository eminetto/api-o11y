package mysql

import (
	"context"
	"database/sql"
	"github.com/eminetto/api-o11y/internal/telemetry"
	"github.com/eminetto/api-o11y/votes/vote"
	"go.opentelemetry.io/otel/codes"
	"time"
)

// VoteMySQL mysql repo
type VoteMySQL struct {
	db        *sql.DB
	telemetry telemetry.Telemetry
}

// NewVoteMySQL create new repository
func NewVoteMySQL(db *sql.DB, telemetry telemetry.Telemetry) *VoteMySQL {
	return &VoteMySQL{
		db:        db,
		telemetry: telemetry,
	}
}

// Store a feedback
func (r *VoteMySQL) Store(ctx context.Context, v *vote.Vote) error {
	ctx, span := r.telemetry.Start(ctx, "mysql")
	defer span.End()
	stmt, err := r.db.Prepare(`
		insert into vote (id, email, talk_name, score, created_at) 
		values(?,?,?,?,?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		v.ID,
		v.Email,
		v.TalkName,
		v.Score,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		span.RecordError(err)
		return err
	}
	err = stmt.Close()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
