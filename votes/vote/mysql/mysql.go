package mysql

import (
	"context"
	"database/sql"
	"github.com/eminetto/api-o11y/votes/vote"
	"time"
)

// VoteMySQL mysql repo
type VoteMySQL struct {
	db *sql.DB
}

// NewVoteMySQL create new repository
func NewVoteMySQL(db *sql.DB) *VoteMySQL {
	return &VoteMySQL{
		db: db,
	}
}

// Store a feedback
func (r *VoteMySQL) Store(ctx context.Context, v *vote.Vote) error {
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
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}
