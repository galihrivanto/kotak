package smtp

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/galihrivanto/kotak/config"
	"github.com/galihrivanto/kotak/db"
	"github.com/galihrivanto/kotak/log"
	"github.com/galihrivanto/kotak/module"

	"github.com/mhale/smtpd"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *config.Config
	db     *db.DB

	srv *smtpd.Server
}

func (s *Server) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	// Start the server
	go func() {
		addr := fmt.Sprintf("%s:%s", s.config.SmtpServer.Host, s.config.SmtpServer.Port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			log.Error("SMTP server listen error: %v", err)

			return
		}

		log.Info("SMTP server listening on %s", addr)

		err = s.srv.Serve(ln)
		if err != nil {
			log.Error("SMTP server error: %v", err)
		}
	}()

	return nil
}

func (s *Server) Close() error {
	log.Info("Stopping SMTP server")
	// stop server
	s.cancel()

	return nil
}

// handleMail processes incoming emails
func (s *Server) handleMail(_ net.Addr, from string, to []string, data []byte) error {
	log.Info("Received mail from %s to %v", from, to)

	// Extract email components
	body := string(data)
	subject := extractHeader(body, "Subject:")

	// Process each recipient
	for _, recipient := range to {
		// Extract account ID from email address
		parts := strings.Split(recipient, "@")
		if len(parts) != 2 {
			continue // Invalid email format
		}

		accountID := parts[0]

		// Check if account exists
		exists, err := s.db.AccountExists(accountID)
		if err != nil || !exists {
			log.Error("Account %s not found, skipping email", accountID)
			continue
		}

		// Store the email
		_, err = s.db.StoreEmail(accountID, from, recipient, subject, body)
		if err != nil {
			log.Error("Failed to store email: %v", err)
		} else {
			log.Info("Stored email for account %s", accountID)
		}
	}

	return nil
}

// extractHeader extracts a header value from email data
func extractHeader(body, header string) string {
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, header) {
			return strings.TrimSpace(strings.TrimPrefix(line, header))
		}
	}
	return ""
}

func NewServer(config *config.Config, db *db.DB) *Server {
	svc := &Server{config: config, db: db}

	// Create SMTP server
	svc.srv = &smtpd.Server{
		Addr:        fmt.Sprintf("%s:%s", config.SmtpServer.Host, config.SmtpServer.Port),
		Handler:     svc.handleMail,
		Appname:     config.SmtpServer.AppName,
		Hostname:    config.SmtpServer.Hostname,
		MaxSize:     config.SmtpServer.MaxSize,
		AuthHandler: nil, // No authentication for temporary mail
	}

	return svc
}

func init() {
	module.RegisterModule("smtp", func(config *config.Config, db *db.DB) module.Module {
		return NewServer(config, db)
	})
}
