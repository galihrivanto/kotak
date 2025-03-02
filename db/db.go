package db

import (
	"context"
	"fmt"
	"time"

	"github.com/galihrivanto/kotak/config"
	"github.com/galihrivanto/kotak/log"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type configKey struct{}

var contextKey = configKey{}

// DB is a wrapper around gorm.DB
type DB struct {
	*gorm.DB
}

// Email represents a stored email
type Email struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	AccountID  string    `gorm:"index" json:"account_id"`
	From       string    `json:"from"`
	To         string    `json:"to"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	ReceivedAt time.Time `gorm:"autoCreateTime" json:"received_at"`
}

// Account represents a temporary email account
type Account struct {
	ID        string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Emails    []Email   `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
}

// New initializes the database
func New(cfg config.Database) (*DB, error) {
	log.Info("Initializing database...")

	var dialector gorm.Dialector
	if cfg.Driver == "postgres" {
		dialector = postgres.Open(cfg.DSN())
	} else if cfg.Driver == "mysql" {
		dialector = mysql.Open(cfg.DSN())
	} else if cfg.Driver == "sqlite" {
		dialector = sqlite.Open(cfg.DSN())

	} else {
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate schema
	err = db.AutoMigrate(&Account{}, &Email{})
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// CreateAccount creates a new temporary email account
func (db *DB) CreateAccount(id string) error {
	return db.Create(&Account{ID: id}).Error
}

// GetAccount retrieves an account by ID
func (db *DB) GetAccount(id string) (*Account, error) {
	var account Account
	if err := db.First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// StoreEmail stores a new email in the database
func (db *DB) StoreEmail(accountID, from, to, subject, body string) (int64, error) {
	email := Email{
		AccountID: accountID,
		From:      from,
		To:        to,
		Subject:   subject,
		Body:      body,
	}
	if err := db.Create(&email).Error; err != nil {
		return 0, err
	}
	return email.ID, nil
}

// GetEmails retrieves all emails for an account
func (db *DB) GetEmails(accountID string) ([]Email, error) {
	var emails []Email
	if err := db.Where("account_id = ?", accountID).Order("received_at DESC").Find(&emails).Error; err != nil {
		return nil, err
	}
	return emails, nil
}

// GetEmail retrieves a specific email
func (db *DB) GetEmail(id int64, accountID string) (*Email, error) {
	var email Email
	if err := db.Where("id = ? AND account_id = ?", id, accountID).First(&email).Error; err != nil {
		return nil, err
	}
	return &email, nil
}

// AccountExists checks if an account exists
func (db *DB) AccountExists(id string) (bool, error) {
	var count int64
	if err := db.Model(&Account{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Cleanup deletes emails older than given interval
func (db *DB) Cleanup(hours int) error {
	cutoff := time.Now().Add(-time.Duration(hours) * time.Hour)
	return db.Where("created_at < ?", cutoff).Delete(&Account{}).Error
}

// Close closes the database (not needed with GORM unless using raw SQL DB)
func (db *DB) Close() error {
	// GORM handles connections automatically, but if needed:
	if sqlDB, err := db.DB.DB(); err == nil {
		return sqlDB.Close()
	}
	return nil
}

// FromContext returns the DB instance from the context
func FromContext(ctx context.Context) *DB {
	return ctx.Value(contextKey).(*DB)
}

// WithContext returns a new context with the DB instance
func WithContext(ctx context.Context, db *DB) context.Context {
	return context.WithValue(ctx, contextKey, db)
}
