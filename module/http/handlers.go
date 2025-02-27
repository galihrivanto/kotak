package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/rand"
)

// createAccount handles the creation of new temporary email accounts
func (s *Server) createAccount(c echo.Context) error {
	// Generate a random account ID
	accountID := generateAccountID(8)

	// Store in database
	if err := s.db.CreateAccount(accountID); err != nil {
		fmt.Println("Failed to create account", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create account",
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"account_id": accountID,
		"email":      fmt.Sprintf("%s@%s", accountID, s.cfg.SmtpServer.Hostname),
	})
}

// getEmails retrieves all emails for an account
func (s *Server) getEmails(c echo.Context) error {
	accountID := c.Param("id")

	// Check if account exists
	exists, err := s.db.AccountExists(accountID)
	if err != nil || !exists {
		fmt.Println("Failed to check if account exists", err)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Account not found",
		})
	}

	// Get emails from database
	emails, err := s.db.GetEmails(accountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch emails",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"emails": emails,
	})
}

// getEmail retrieves a specific email
func (s *Server) getEmail(c echo.Context) error {
	accountID := c.Param("id")
	emailIDStr := c.Param("email_id")

	// Convert email ID to int64
	emailID, err := strconv.ParseInt(emailIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid email ID",
		})
	}

	// Get email from database
	email, err := s.db.GetEmail(emailID, accountID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Email not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"email": email,
	})
}

// Helper function to generate a random account ID
func generateAccountID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Initialize random seed
func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}
