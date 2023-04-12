package mysql

import (
	"context"
	"database/sql"
	"github.com/eminetto/api-o11y/auth/user"
	"github.com/eminetto/api-o11y/internal/telemetry"
	"go.opentelemetry.io/otel/codes"
)

// UserMySQL mysql repo
type UserMySQL struct {
	db        *sql.DB
	telemetry telemetry.Telemetry
}

// NewUserMySQL create new repository
func NewUserMySQL(db *sql.DB, telemetry telemetry.Telemetry) *UserMySQL {
	return &UserMySQL{
		db:        db,
		telemetry: telemetry,
	}
}

// Get an user
func (r *UserMySQL) Get(ctx context.Context, email string) (*user.User, error) {
	ctx, span := r.telemetry.Start(ctx, "mysql")
	defer span.End()
	stmt, err := r.db.Prepare(`select id, email, password, first_name, last_name from user where email = ?`)
	if err != nil {
		return nil, err
	}
	var u user.User
	rows, err := stmt.Query(email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}
	return &u, nil
}
