package mysql

import (
	"context"
	"database/sql"
	"github.com/eminetto/api-o11y/auth/user"
)

// UserMySQL mysql repo
type UserMySQL struct {
	db *sql.DB
}

// NewUserMySQL create new repository
func NewUserMySQL(db *sql.DB) *UserMySQL {
	return &UserMySQL{
		db: db,
	}
}

// Get an user
func (r *UserMySQL) Get(ctx context.Context, email string) (*user.User, error) {
	stmt, err := r.db.Prepare(`select id, email, password, first_name, last_name from user where email = ?`)
	if err != nil {
		return nil, err
	}
	var u user.User
	rows, err := stmt.Query(email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName)
	}
	return &u, nil
}
