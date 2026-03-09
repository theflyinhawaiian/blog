package db

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/peterblog/blog/internal/models"
)

// UpsertUser finds or creates a user for the given provider identity.
// Returns the user and whether it was newly created.
func UpsertUser(db *sqlx.DB, provider, providerUserID, displayName string, email *string) (*models.User, error) {
	var identity models.UserIdentity
	err := db.Get(&identity,
		`SELECT id, user_id, provider, provider_user_id, email
		 FROM user_identities WHERE provider = ? AND provider_user_id = ?`,
		provider, providerUserID)

	if err == nil {
		// Existing identity - return user
		var user models.User
		err = db.Get(&user, `SELECT id, display_name, created_at FROM users WHERE id = ?`, identity.UserID)
		return &user, err
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// New identity - create user then identity
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO users (display_name) VALUES (?)`, displayName)
	if err != nil {
		return nil, err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		`INSERT INTO user_identities (user_id, provider, provider_user_id, email) VALUES (?, ?, ?, ?)`,
		userID, provider, providerUserID, email)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var user models.User
	err = db.Get(&user, `SELECT id, display_name, created_at FROM users WHERE id = ?`, userID)
	return &user, err
}

func GetUserByID(db *sqlx.DB, id uint64) (*models.User, error) {
	var user models.User
	err := db.Get(&user, `SELECT id, display_name, created_at FROM users WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
