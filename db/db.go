package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/galihrivanto/kotak/config"
	"github.com/galihrivanto/kotak/log"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type configKey struct{}

var contextKey = configKey{}

// DB is a wrapper around sql.DB
type DB struct {
	*sql.DB
}

// Email represents a stored email
type Email struct {
	ID         int64     `json:"id"`
	AccountID  string    `json:"account_id"`
	From       string    `json:"from"`
	To         string    `json:"to"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	ReceivedAt time.Time `json:"received_at"`
}

// Account represents a temporary email account
type Account struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// New initializes the database
func New(cfg config.Database) (*DB, error) {
	log.Info("Initializing database...")
	db, err := sql.Open(cfg.Driver, cfg.DSN())
	if err != nil {
		return nil, err
	}

	// enquire database version
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		return nil, err
	}
	log.Info("Database version: %s", version)

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id TEXT PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS emails (
			id SERIAL PRIMARY KEY,
			account_id TEXT,
			from_email TEXT,
			to_email TEXT,
			subject TEXT,
			body TEXT,
			received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db}, nil
}

// CreateAccount creates a new temporary email account
func (db *DB) CreateAccount(id string) error {
	_, err := db.DB.Exec("INSERT INTO accounts (id) VALUES ($1)", id)
	return err
}

// GetAccount retrieves an account by ID
func (db *DB) GetAccount(id string) (*Account, error) {
	var account Account
	err := db.DB.QueryRow("SELECT id, created_at FROM accounts WHERE id = $1", id).Scan(&account.ID, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// StoreEmail stores a new email in the database
func (db *DB) StoreEmail(accountID, from, to, subject, body string) (int64, error) {
	res, err := db.DB.Exec(
		"INSERT INTO emails (account_id, from_email, to_email, subject, body) VALUES ($1, $2, $3, $4, $5)",
		accountID, from, to, subject, body,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetEmails retrieves all emails for an account
func (db *DB) GetEmails(accountID string) ([]Email, error) {
	rows, err := db.DB.Query(`
		SELECT id, account_id, from_email, to_email, subject, body, received_at 
		FROM emails 
		WHERE account_id = $1 
		ORDER BY received_at DESC`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []Email
	for rows.Next() {
		var email Email
		if err := rows.Scan(&email.ID, &email.AccountID, &email.From, &email.To, &email.Subject, &email.Body, &email.ReceivedAt); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, nil
}

// GetEmail retrieves a specific email
func (db *DB) GetEmail(id int64, accountID string) (*Email, error) {
	var email Email
	err := db.DB.QueryRow(`
		SELECT id, account_id, from_email, to_email, subject, body, received_at 
		FROM emails 
		WHERE id = $1 AND account_id = $2`, id, accountID).
		Scan(&email.ID, &email.AccountID, &email.From, &email.To, &email.Subject, &email.Body, &email.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return &email, nil
}

// AccountExists checks if an account exists
func (db *DB) AccountExists(id string) (bool, error) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM accounts WHERE id = $1", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Cleanup deletes emails older than given interval
func (db *DB) Cleanup(hours int) error {
	_, err := db.DB.Exec(fmt.Sprintf("DELETE FROM accounts WHERE created_at < NOW() - INTERVAL '%d hour'", hours))
	return err
}

// Close closes the database
func (db *DB) Close() error {
	return db.DB.Close()
}

// FromContext returns the DB instance from the context
func FromContext(ctx context.Context) *DB {
	return ctx.Value(contextKey).(*DB)
}

// WithContext returns a new context with the DB instance
func WithContext(ctx context.Context, db *DB) context.Context {
	return context.WithValue(ctx, contextKey, db)
}
